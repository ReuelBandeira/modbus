package utils

import (
	"os"
	"path/filepath"

	process "github.com/shirou/gopsutil/v3/process"
)

func Running() bool {

	processes := []string{}

	myself, _ := os.Executable()
	myself = filepath.Base(myself)

	v, _ := process.Processes()

	for _, p := range v {
		proc, err := p.Name()
		if err == nil && proc == myself {
			processes = append(processes, proc)
		}
	}

	return len(processes) > 1
}
