package collector

import (
	"runtime"
	"techsupport/log/pkg"
	"techsupport/sysinfo/internal/errors"
	"techsupport/sysinfo/internal/models"

	"github.com/denisbrodbeck/machineid"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

type RealSystemCollector struct {
	Log pkg.Logger
}

func (c *RealSystemCollector) handleErr(err error, sysErr error) bool {
	if err != nil {
		if c.Log != nil {
			c.Log.Errorw(sysErr.Error(), "raw_error", err)
		}
		return true
	}
	return false
}

func (c *RealSystemCollector) Collect() (models.SystemInfo, error) {
	if c.Log != nil {
		c.Log.Debugw("collection started")
	}

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

	return models.SystemInfo{
		OS:        runtime.GOOS,
		Platform:  platform,
		Arch:      runtime.GOARCH,
		CPUModel:  cpuName,
		CPUCores:  runtime.NumCPU(),
		TotalRAM:  totalRAM,
		MachineID: mID,
	}, nil
}