package collector

import (
	"runtime"
	"sync"

	"techsupport/log/pkg"
	"techsupport/sysinfo/internal/errors"
	"techsupport/sysinfo/pkg/models"

	"github.com/denisbrodbeck/machineid"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

// RealSystemCollector is a thread-safe implementation of system information collector.
// It gathers runtime and hardware information from the host machine.
type RealSystemCollector struct {
	Log pkg.Logger

	mu sync.Mutex // protects logger usage in concurrent scenarios
}

// handleErr processes errors from system calls and logs them if logger is available.
// It returns true if an error occurred.
func (c *RealSystemCollector) handleErr(err error, sysErr error) bool {
	if err != nil {
		c.mu.Lock()
		if c.Log != nil {
			c.Log.Errorw(sysErr.Error(), "raw_error", err)
		}
		c.mu.Unlock()
		return true
	}
	return false
}

// Collect gathers system information from OS and hardware sources.
// It is safe for concurrent usage due to internal synchronization.
func (c *RealSystemCollector) Collect() (models.SystemInfo, error) {
	c.mu.Lock()
	if c.Log != nil {
		c.Log.Debugw("system info collection started")
	}
	c.mu.Unlock()

	hInfo, errH := host.Info()
	c.handleErr(errH, errors.ErrHostInfo)

	mInfo, errM := mem.VirtualMemory()
	c.handleErr(errM, errors.ErrMemoryInfo)

	cInfo, errC := cpu.Info()
	c.handleErr(errC, errors.ErrCPUInfo)

	mID, errID := machineid.ID()
	c.handleErr(errID, errors.ErrMachineID)

	cpuName := "Unknown"
	if len(cInfo) > 0 {
		cpuName = cInfo[0].ModelName
	}

	platform := "Unknown"
	if hInfo != nil {
		platform = hInfo.Platform
	}

	var totalRAM uint64
	if mInfo != nil {
		totalRAM = mInfo.Total
	}

	result := models.SystemInfo{
		OS:        runtime.GOOS,
		Platform:  platform,
		Arch:      runtime.GOARCH,
		CPUModel:  cpuName,
		CPUCores:  runtime.NumCPU(),
		TotalRAM:  totalRAM,
		MachineID: mID,
	}

	return result, nil
}