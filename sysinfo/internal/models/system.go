package models

type SystemInfo struct {

	OS        string `json:"os"`        
	Platform  string `json:"platform"`  
	Arch      string `json:"arch"`  
	Kernel    string `json:"kernel"`

	CPUModel  string  `json:"cpu_model"`
	CPUCores  int     `json:"cpu_cores"`
	TotalRAM  uint64  `json:"total_ram"`
	
	Hostname  string `json:"hostname"`
	MachineID string `json:"machine_id"` 
	Username  string `json:"username"`
}