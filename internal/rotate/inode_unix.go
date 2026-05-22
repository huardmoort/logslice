//go:build !windows

package rotate

import (
	"os"
	"syscall"
)

// inode returns the inode number for the given FileInfo on Unix systems.
func inode(info os.FileInfo) uint64 {
	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		return stat.Ino
	}
	return 0
}
