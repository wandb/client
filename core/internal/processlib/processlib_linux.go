package processlib

import (
	"os"
	"syscall"

	"golang.org/x/sys/unix"
)

func ShutdownOnParentExit(parentPid int) bool {
	// SIGTERM allows the installed handler to flush within its grace period.
	// The server also polls parent lifetime if prctl is unavailable.
	if err := unix.Prctl(unix.PR_SET_PDEATHSIG, uintptr(syscall.SIGTERM), 0, 0, 0); err != nil {
		return false
	}
	if os.Getppid() != parentPid {
		_ = syscall.Kill(os.Getpid(), syscall.SIGTERM)
	}
	return true
}
