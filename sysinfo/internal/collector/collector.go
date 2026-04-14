package collector

import (
	"runtime"
	"techsupport/sysinfo/internal/models"

	"github.com/denisbrodbeck/machineid"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

type RealSystemCollector struct{}

func (c *RealSystemCollector) Collect() (models.SystemInfo, error) {
	hInfo, _ := host.Info()
	cInfo, _ := cpu.Info()
	mInfo, _ := mem.VirtualMemory()
	mID, _ := machineid.ID()

	cpuName := "Unknown"
	if len(cInfo) > 0 {
		cpuName = cInfo[0].ModelName
	}

	return models.SystemInfo{
		OS:        runtime.GOOS,
		Platform:  hInfo.Platform,
		Arch:      runtime.GOARCH,
		Kernel:    hInfo.KernelVersion,
		CPUModel:  cpuName,
		CPUCores:  runtime.NumCPU(),
		TotalRAM:  mInfo.Total,
		Hostname:  hInfo.Hostname,
		MachineID: mID,
	}, nil
}