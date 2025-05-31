package filelist

import (
	"fmt"
	"os"
	"os/user"
	"strings"
	"syscall"
)

func ownerName(entry os.DirEntry) string {
	info, err := entry.Info()
	if err != nil {
		return "?"
	}
	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		uid := stat.Uid
		user, err := user.LookupId(fmt.Sprintf("%d", uid))
		if err == nil {
			return user.Username
		}
		return fmt.Sprintf("%d", uid)
	}
	return "?"
}

func groupName(entry os.DirEntry) string {
	info, err := entry.Info()
	if err != nil {
		return "?"
	}
	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		if group, err := user.LookupGroupId(fmt.Sprintf("%d", stat.Gid)); err == nil {
			return group.Name
		}
	}
	return "?"
}

func permissions(entry os.DirEntry) string {
	info, err := entry.Info()
	if err != nil {
		return "?"
	}
	return info.Mode().String()
}

func sortCaseInsensitiveFunc(a string, b string) int {
	aLower := strings.ToLower(a)
	bLower := strings.ToLower(b)
	if aLower < bLower {
		return -1
	} else if aLower > bLower {
		return 1
	}
	return 0
}

func sortNameFunc(a os.DirEntry, b os.DirEntry) int {
	aLower := strings.ToLower(a.Name())
	bLower := strings.ToLower(b.Name())

	// Cheesy hack to ensure directories are sorted before files
	if a.IsDir() {
		aLower = string(rune(0)) + aLower
	}
	if b.IsDir() {
		bLower = string(rune(0)) + bLower
	}

	if aLower < bLower {
		return -1
	} else if aLower > bLower {
		return 1
	}
	return 0
}

func sortSizeFunc(a os.DirEntry, b os.DirEntry) int {
	var aSize int64 = 0
	if aInfo, err := a.Info(); err == nil {
		aSize = aInfo.Size()
	}

	var bSize int64 = 0
	if bInfo, err := b.Info(); err == nil {
		bSize = bInfo.Size()
	}

	if aSize < bSize {
		return -1
	} else if aSize > bSize {
		return 1
	}
	return 0
}

func sortModifiedFunc(a os.DirEntry, b os.DirEntry) int {
	var aModTime, bModTime int64 = 0, 0
	if aInfo, err := a.Info(); err == nil {
		aModTime = aInfo.ModTime().Unix()
	}
	if bInfo, err := b.Info(); err == nil {
		bModTime = bInfo.ModTime().Unix()
	}
	if aModTime < bModTime {
		return -1
	} else if aModTime > bModTime {
		return 1
	}
	return 0
}

func sortPermissionsFunc(a os.DirEntry, b os.DirEntry) int {
	var aPerm, bPerm string
	if aInfo, err := a.Info(); err == nil {
		aPerm = aInfo.Mode().String()
	}
	if bInfo, err := b.Info(); err == nil {
		bPerm = bInfo.Mode().String()
	}
	if aPerm < bPerm {
		return -1
	} else if aPerm > bPerm {
		return 1
	}
	return 0
}

func sortOwnerFunc(a os.DirEntry, b os.DirEntry) int {
	return sortCaseInsensitiveFunc(ownerName(a), ownerName(b))
}

func sortGroupFunc(a os.DirEntry, b os.DirEntry) int {
	return sortCaseInsensitiveFunc(groupName(a), groupName(b))
}

func formatSize(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%.1fB", float64(size))
	} else if size < 1024*1024 {
		return fmt.Sprintf("%.1fKiB", float64(size)/1024)
	} else if size < 1024*1024*1024 {
		return fmt.Sprintf("%.1fMiB", float64(size)/(1024*1024))
	} else {
		return fmt.Sprintf("%.1fGiB", float64(size)/(1024*1024*1024))
	}
}
