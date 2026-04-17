package models

// SystemInfo represents a snapshot of the host machine state.
// It is used by system collectors to standardize hardware and OS information
// across different environments and platforms.
type SystemInfo struct {
	// OS is the operating system name (e.g. linux, windows, darwin)
	OS string `json:"os"`

	// Platform is a more detailed OS platform identifier (e.g. ubuntu, debian)
	Platform string `json:"platform"`

	// Arch defines CPU architecture (e.g. amd64, arm64)
	Arch string `json:"arch"`

	// Kernel represents OS kernel version
	Kernel string `json:"kernel"`

	// CPUModel describes the processor model name
	CPUModel string `json:"cpu_model"`

	// CPUCores is the number of logical CPU cores available
	CPUCores int `json:"cpu_cores"`

	// TotalRAM represents total physical memory in bytes
	TotalRAM uint64 `json:"total_ram"`

	// Hostname is the system network name
	Hostname string `json:"hostname"`

	// MachineID is a unique identifier for the machine instance
	// used for fingerprinting and tracking across runs
	MachineID string `json:"machine_id"`

	// Username is the current logged-in system user
	Username string `json:"username"`
}