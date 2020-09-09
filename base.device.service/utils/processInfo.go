package utils

import "github.com/shirou/gopsutil/process"

//ProcessInfo struct for current process information
type ProcessInfo struct {
	process.Process
	CPUPercent float64
}
