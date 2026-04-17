package errors

import "errors"

// Package errors defines system-level error constants used by the sysinfo collector.
// These errors represent failures in low-level system information retrieval operations.

// ErrHostInfo indicates failure while retrieving host system information
// (e.g. OS version, platform, uptime).
var ErrHostInfo = errors.New("SYS_HOST_INFO_FAIL")

// ErrCPUInfo indicates failure while retrieving CPU information
// (e.g. model name, architecture details, core data).
var ErrCPUInfo = errors.New("SYS_CPU_INFO_FAIL")

// ErrMemoryInfo indicates failure while retrieving memory information
// (e.g. total RAM, available memory).
var ErrMemoryInfo = errors.New("SYS_MEM_INFO_FAIL")

// ErrMachineID indicates failure while retrieving machine unique identifier
// used for host fingerprinting.
var ErrMachineID = errors.New("SYS_MACHINE_ID_FAIL")