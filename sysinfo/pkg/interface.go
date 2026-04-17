package pkg

import "techsupport/sysinfo/pkg/models"

// SystemCollector defines an abstraction for collecting system-level information.
// It allows different implementations (real, mocked, test, remote agents)
// to provide system data in a unified format.
type SystemCollector interface {
	// Collect retrieves current system information from the host machine.
	// It returns a SystemInfo snapshot or an error if collection fails.
	Collect() (models.SystemInfo, error)
}