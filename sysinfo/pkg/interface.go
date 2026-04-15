package pkg

import "techsupport/sysinfo/pkg/models"


type SystemCollector interface {
	Collect() (models.SystemInfo, error)
}