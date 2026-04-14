package errors


import "errors"

var (
	ErrHostInfo    = errors.New("SYS_HOST_INFO_FAIL")
	ErrCPUInfo     = errors.New("SYS_CPU_INFO_FAIL")
	ErrMemoryInfo  = errors.New("SYS_MEM_INFO_FAIL")
	ErrMachineID   = errors.New("SYS_MACHINE_ID_FAIL")
)