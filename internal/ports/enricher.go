package ports

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// ProcessInfo holds metadata about the process owning a port.
type ProcessInfo struct {
	PID  int
	Name string
	Exe  string
}

// String returns a human-readable representation of ProcessInfo.
func (p ProcessInfo) String() string {
	if p.PID == 0 {
		return "unknown"
	}
	return fmt.Sprintf("%s(pid=%d)", p.Name, p.PID)
}

// Enricher resolves process information for a given inode.
type Enricher struct {
	procRoot string
}

// NewEnricher creates an Enricher that reads from the given proc root
// (typically "/proc").
func NewEnricher(procRoot string) *Enricher {
	return &Enricher{procRoot: procRoot}
}

// Lookup finds the process that owns the given socket inode.
func (e *Enricher) Lookup(inode uint64) (ProcessInfo, error) {
	pids, err := e.listPIDs()
	if err != nil {
		return ProcessInfo{}, err
	}
	target := fmt.Sprintf("socket:[%d]", inode)
	for _, pid := range pids {
		if owns, _ := e.pidOwnsSocket(pid, target); owns {
			return e.readProcessInfo(pid)
		}
	}
	return ProcessInfo{}, nil
}

func (e *Enricher) listPIDs() ([]int, error) {
	entries, err := os.ReadDir(e.procRoot)
	if err != nil {
		return nil, err
	}
	var pids []int
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		if pid, err := strconv.Atoi(entry.Name()); err == nil {
			pids = append(pids, pid)
		}
	}
	return pids, nil
}

func (e *Enricher) pidOwnsSocket(pid int, target string) (bool, error) {
	fdDir := filepath.Join(e.procRoot, strconv.Itoa(pid), "fd")
	entries, err := os.ReadDir(fdDir)
	if err != nil {
		return false, err
	}
	for _, entry := range entries {
		link, err := os.Readlink(filepath.Join(fdDir, entry.Name()))
		if err == nil && link == target {
			return true, nil
		}
	}
	return false, nil
}

func (e *Enricher) readProcessInfo(pid int) (ProcessInfo, error) {
	commPath := filepath.Join(e.procRoot, strconv.Itoa(pid), "comm")
	data, err := os.ReadFile(commPath)
	if err != nil {
		return ProcessInfo{PID: pid}, nil
	}
	exePath := filepath.Join(e.procRoot, strconv.Itoa(pid), "exe")
	exe, _ := os.Readlink(exePath)
	return ProcessInfo{
		PID:  pid,
		Name: strings.TrimSpace(string(data)),
		Exe:  exe,
	}, nil
}
