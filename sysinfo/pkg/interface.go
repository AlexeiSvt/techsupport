package pkg

import "techsupport/sysinfo/internal/models"


type SystemCollector interface {
	Collect() (models.SystemInfo, error)
}