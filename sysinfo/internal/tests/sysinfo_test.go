package tests

import (
	"fmt"
	"strings"
	"testing"
	"techsupport/sysinfo/internal/collector"
	"techsupport/sysinfo/pkg"
)

func TestRealSystemCollector(t *testing.T) {
	var col pkg.SystemCollector = &collector.RealSystemCollector{}

	info, err := col.Collect()
	if err != nil {
		t.Skipf("Skipping: %v", err)
	}

	if info.OS == "" {
		t.Error("OS field is empty")
	}
	if info.CPUCores <= 0 {
		t.Errorf("Invalid cores: %d", info.CPUCores)
	}

	divider := strings.Repeat("=", 45)
	fmt.Println(divider)
	fmt.Printf(" [DEVICE AGENT REPORT]\n")
	fmt.Println(divider)
	fmt.Printf(" %-15s: %s (%s)\n", "OPERATING SYS", strings.ToUpper(info.OS), info.Platform)
	fmt.Printf(" %-15s: %s\n", "ARCHITECTURE", info.Arch)
	fmt.Printf(" %-15s: %s\n", "HOSTNAME", info.Hostname)
	fmt.Println(strings.Repeat("-", 45))
	fmt.Printf(" %-15s: %s\n", "CPU MODEL", info.CPUModel)
	fmt.Printf(" %-15s: %d Cores\n", "LOGICAL CORES", info.CPUCores)
	fmt.Printf(" %-15s: %.2f GB\n", "TOTAL RAM", float64(info.TotalRAM)/1024/1024/1024)
	fmt.Println(strings.Repeat("-", 45))
	fmt.Printf(" %-15s: %s\n", "MACHINE ID", info.MachineID)
	fmt.Println(divider)
}