//go:build windows
// +build windows

package cancelreader

import (
	"io"
)

// NewReader returns a reader and a cancel function. If the input reader is a
// File with the same file descriptor as os.Stdin, the cancel function can
// be used to interrupt a blocking read call. In this case, the cancel function
// returns true if the call was canceled successfully. If the input reader is
// not a File with the same file descriptor as os.Stdin, the cancel
// function does nothing and always returns false. The Windows implementation
// is based on WaitForMultipleObject with overlapping reads from CONIN$.
func NewReader(reader io.Reader) (CancelReader, error) {
	return newFallbackCancelReader(reader)
}
