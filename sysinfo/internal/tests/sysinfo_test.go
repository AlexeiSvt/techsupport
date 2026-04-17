package tests

import (
	"fmt"
	"strings"
	"testing"

	"techsupport/sysinfo/internal/collector"
	"techsupport/sysinfo/pkg"
)

// TestRealSystemCollector verifies that the real system collector
// can successfully gather basic system information from the host machine.
// It is an integration-style test and may be skipped in restricted environments.
func TestRealSystemCollector(t *testing.T) {
	var col pkg.SystemCollector = &collector.RealSystemCollector{}

	info, err := col.Collect()
	if err != nil {
		// In CI or restricted environments system access may fail,
		// so the test is skipped instead of marked as failed.
		t.Skipf("Skipping: %v", err)
	}

	// Validate essential system fields
	if info.OS == "" {
		t.Error("OS field is empty")
	}
	if info.CPUCores <= 0 {
		t.Errorf("Invalid cores: %d", info.CPUCores)
	}

	//Console report formatting section

	divider := strings.Repeat("=", 45)

	fmt.Println(divider)
	fmt.Printf(" [DEVICE AGENT REPORT]\n")
	fmt.Println(divider)

	// Operating system information
	fmt.Printf(" %-15s: %s (%s)\n", "OPERATING SYS", strings.ToUpper(info.OS), info.Platform)

	// Architecture details
	fmt.Printf(" %-15s: %s\n", "ARCHITECTURE", info.Arch)

	// Host identification
	fmt.Printf(" %-15s: %s\n", "HOSTNAME", info.Hostname)

	fmt.Println(strings.Repeat("-", 45))

	// CPU information
	fmt.Printf(" %-15s: %s\n", "CPU MODEL", info.CPUModel)
	fmt.Printf(" %-15s: %d Cores\n", "LOGICAL CORES", info.CPUCores)

	// Memory information (converted to GB for readability)
	fmt.Printf(" %-15s: %.2f GB\n", "TOTAL RAM", float64(info.TotalRAM)/1024/1024/1024)

	fmt.Println(strings.Repeat("-", 45))

	// Machine identification
	fmt.Printf(" %-15s: %s\n", "MACHINE ID", info.MachineID)

	fmt.Println(divider)
}
