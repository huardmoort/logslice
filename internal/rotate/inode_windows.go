//go:build windows

package rotate

import "os"

// inode is not meaningful on Windows; always returns 0.
// Rotation detection on Windows relies solely on size decrease.
func inode(info os.FileInfo) uint64 {
	return 0
}
