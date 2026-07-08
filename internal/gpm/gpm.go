// Copyright (c) 2026 Robin <robin413@protonmail.com>
// SPDX-License-Identifier: MIT

//go:build linux

package gpm

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/gdamore/tcell/v2"
)

const (
	gpmSocket = "/dev/gpmctl"

	// Event type bits (int32 bitmask)
	gpmMove = 1
	gpmUp   = 8

	// Button bits (uint8)
	gpmBRight  = 1
	gpmBMiddle = 2
	gpmBLeft   = 4
	gpmBUp     = 16
	gpmBDown   = 32
)

// gpmConnect is the 16-byte handshake struct sent to /dev/gpmctl.
type gpmConnect struct {
	EventMask   uint16
	DefaultMask uint16
	MinMod      uint16
	MaxMod      uint16
	Pid         int32
	Vc          int32
}

// gpmEvent is the 28-byte event struct read back from /dev/gpmctl.
type gpmEvent struct {
	Buttons   uint8
	Modifiers uint8
	Vc        uint16
	Dx        int16
	Dy        int16
	X         int16
	Y         int16
	Type      int32
	Clicks    int32
	Margin    int32
	Wdx       int16
	Wdy       int16
}

// Client manages the GPM connection lifecycle.
type Client struct {
	conn      net.Conn
	postEvent func(tcell.Event)
	done      chan struct{}
	wg        sync.WaitGroup
}

// IsGPMAvailable returns true when the process is running in a Linux TTY
// (not inside a graphical terminal emulator) and the GPM daemon is running.
func IsGPMAvailable() bool {
	if os.Getenv("DISPLAY") != "" {
		return false
	}
	if os.Getenv("WAYLAND_DISPLAY") != "" {
		return false
	}
	if os.Getenv("TERM") != "linux" {
		return false
	}
	_, err := os.Stat(gpmSocket)
	return err == nil
}

// currentVC returns the active virtual console number by reading
// /sys/class/tty/tty0/active. Returns 0 on any error (safe GPM fallback).
func currentVC() int32 {
	data, err := os.ReadFile("/sys/class/tty/tty0/active")
	if err != nil {
		return 0
	}
	name := strings.TrimSpace(string(data))
	name = strings.TrimPrefix(name, "tty")
	n, err := strconv.Atoi(name)
	if err != nil || n < 0 || n > 63 {
		return 0
	}
	return int32(n)
}

// Connect establishes a connection to the GPM daemon and returns a Client.
// postEvent is called for each mouse event — pass app.QueueEvent.
// Returns (nil, nil) when GPM is unavailable. A non-nil error means the
// daemon socket exists but the connection or handshake failed.
func Connect(postEvent func(tcell.Event)) (*Client, error) {
	if !IsGPMAvailable() {
		return nil, nil
	}

	conn, err := net.Dial("unix", gpmSocket)
	if err != nil {
		return nil, fmt.Errorf("gpm: dial %s: %w", gpmSocket, err)
	}

	handshake := gpmConnect{
		EventMask:   0xFF,
		DefaultMask: 0,
		MinMod:      0,
		MaxMod:      ^uint16(0),
		Pid:         int32(os.Getpid()),
		Vc:          currentVC(),
	}
	if err := binary.Write(conn, binary.LittleEndian, handshake); err != nil {
		conn.Close()
		return nil, fmt.Errorf("gpm: write handshake: %w", err)
	}

	return &Client{
		conn:      conn,
		postEvent: postEvent,
		done:      make(chan struct{}),
	}, nil
}

// Start launches the GPM event reading goroutine.
func (c *Client) Start() {
	c.wg.Add(1)
	go c.readLoop()
}

// Stop signals the event reading goroutine to exit and waits for it.
// Safe to call on a nil receiver or multiple times.
func (c *Client) Stop() {
	if c == nil {
		return
	}
	select {
	case <-c.done:
		// already stopped
	default:
		close(c.done)
	}
	if c.conn != nil {
		c.conn.Close()
	}
	c.wg.Wait()
}

func (c *Client) readLoop() {
	defer c.wg.Done()

	for {
		select {
		case <-c.done:
			return
		default:
		}

		var ev gpmEvent
		if err := binary.Read(c.conn, binary.LittleEndian, &ev); err != nil {
			select {
			case <-c.done:
				// Normal shutdown — conn.Close() was called by Stop().
			default:
				log.Printf("gpm: read error: %v", err)
			}
			return
		}

		if tcellEv := convertEvent(&ev); tcellEv != nil {
			c.postEvent(tcellEv)
		}
	}
}

// convertEvent maps a gpmEvent to a tcell.EventMouse.
// Returns nil for events with no useful tcell equivalent.
func convertEvent(ev *gpmEvent) *tcell.EventMouse {
	var btn tcell.ButtonMask
	var mod tcell.ModMask

	// GPM modifier bits: 0=Shift, 2=Ctrl, 3=Alt
	if ev.Modifiers&0x01 != 0 {
		mod |= tcell.ModShift
	}
	if ev.Modifiers&0x04 != 0 {
		mod |= tcell.ModCtrl
	}
	if ev.Modifiers&0x08 != 0 {
		mod |= tcell.ModAlt
	}

	// Wheel events: prefer the explicit wheel delta fields, fall back to button
	// flags. Vertical takes priority over horizontal when both are present.
	if ev.Wdy > 0 {
		btn = tcell.WheelUp
	} else if ev.Wdy < 0 {
		btn = tcell.WheelDown
	} else if ev.Wdx > 0 {
		btn = tcell.WheelRight
	} else if ev.Wdx < 0 {
		btn = tcell.WheelLeft
	} else if ev.Buttons&gpmBDown != 0 {
		btn = tcell.WheelDown
	} else if ev.Buttons&gpmBUp != 0 {
		btn = tcell.WheelUp
	} else if ev.Type&gpmUp != 0 {
		// Button release: send ButtonNone so tview sees the transition from
		// pressed → released and fires MouseLeftUp / MouseRightUp / etc.
		btn = tcell.ButtonNone
	} else {
		if ev.Buttons&gpmBLeft != 0 {
			btn |= tcell.Button1
		}
		if ev.Buttons&gpmBMiddle != 0 {
			btn |= tcell.Button2
		}
		if ev.Buttons&gpmBRight != 0 {
			btn |= tcell.Button3
		}
	}

	// GPM coordinates are 1-based; tcell uses 0-based.
	x := int(ev.X) - 1
	y := int(ev.Y) - 1
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}

	return tcell.NewEventMouse(x, y, btn, mod)
}
