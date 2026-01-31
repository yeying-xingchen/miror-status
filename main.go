package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

// 数据结构
type SystemInfo struct {
	CPU    []CPUInfo  `json:"cpu"`
	Memory MemoryInfo `json:"memory"`
	Disk   DiskInfo   `json:"disk"`
	Host   HostInfo   `json:"host"`
}

type CPUInfo struct {
	ModelName string  `json:"model_name"`
	Cores     int32   `json:"cores"`
	Usage     float64 `json:"usage"`
}

type MemoryInfo struct {
	Total       uint64  `json:"total"`
	Available   uint64  `json:"available"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"used_percent"`
}

type DiskInfo struct {
	Total       uint64  `json:"total"`
	Free        uint64  `json:"free"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"used_percent"`
}

type HostInfo struct {
	Hostname   string    `json:"hostname"`
	OS         string    `json:"os"`
	Platform   string    `json:"platform"`
	BootTime   uint64    `json:"boot_time"`
	Uptime     uint64    `json:"uptime"`
	Procs      int       `json:"procs"`
	CreateTime time.Time `json:"create_time"`
}

func main() {
	router := gin.Default()

	// 获取总信息
	router.GET("/", func(c *gin.Context) {
		info := GetSystemInfo()
		c.JSON(200, info)
	})

	// 获取CPU信息
	router.GET("/cpu", func(c *gin.Context) {
		cpuInfo := GetCPUInfo()
		c.JSON(200, cpuInfo)
	})

	// 获取内存信息
	router.GET("/memory", func(c *gin.Context) {
		memInfo := GetMemoryInfo()
		c.JSON(200, memInfo)
	})

	// 获取磁盘信息
	router.GET("/disk", func(c *gin.Context) {
		diskInfo := GetDiskInfo()
		c.JSON(200, diskInfo)
	})

	// 获取主机信息
	router.GET("/host", func(c *gin.Context) {
		hostInfo := GetHostInfo()
		c.JSON(200, hostInfo)
	})

	router.Run(":8080")
}

// 完整的系统信息
func GetSystemInfo() SystemInfo {
	return SystemInfo{
		CPU:    GetCPUInfo(),
		Memory: GetMemoryInfo(),
		Disk:   GetDiskInfo(),
		Host:   GetHostInfo(),
	}
}

// CPU信息
func GetCPUInfo() []CPUInfo {
	var cpuInfos []CPUInfo

	// 型号
	infos, _ := cpu.Info()
	for _, info := range infos {
		cpuInfo := CPUInfo{
			ModelName: info.ModelName,
			Cores:     info.Cores,
		}
		cpuInfos = append(cpuInfos, cpuInfo)
	}

	// 使用率
	if len(cpuInfos) > 0 {
		percents, _ := cpu.Percent(time.Second, false)
		if len(percents) > 0 {
			cpuInfos[0].Usage = percents[0]
		}
	}

	return cpuInfos
}

// 内存信息
func GetMemoryInfo() MemoryInfo {
	v, _ := mem.VirtualMemory()

	return MemoryInfo{
		Total:       v.Total,
		Available:   v.Available,
		Used:        v.Used,
		UsedPercent: v.UsedPercent,
	}
}

// 磁盘信息
func GetDiskInfo() DiskInfo {
	parts, _ := disk.Partitions(false)
	var total, free, used uint64

	for _, part := range parts {
		usage, _ := disk.Usage(part.Mountpoint)
		total += usage.Total
		free += usage.Free
		used += usage.Used
	}

	var usedPercent float64
	if total > 0 {
		usedPercent = float64(used) / float64(total) * 100
	}

	return DiskInfo{
		Total:       total,
		Free:        free,
		Used:        used,
		UsedPercent: usedPercent,
	}
}

// 设备信息
func GetHostInfo() HostInfo {
	h, _ := host.Info()

	return HostInfo{
		Hostname:   h.Hostname,
		OS:         h.OS,
		Platform:   h.Platform,
		BootTime:   h.BootTime,
		Uptime:     h.Uptime,
		Procs:      int(h.Procs),
		CreateTime: time.Now(),
	}
}
