package dto

// NetInfo 设置网络的结构体
type NetInfo struct {
	Nettype string `json:"nettype" binding:"required"`
	Address string `json:"address" binding:"required"`
	Netmask string `json:"netmask" binding:"required"`
	Gateway string `json:"gateway" binding:"required"`
}

// NTPInfo 设置时区的结构体
type NTPInfo struct {
	Date     int64  `json:"date" binding:"required"`
	TimeType string `json:"timetype" binding:"required"`
	NtpURL   string `json:"ntpurl" binding:"required"`
}

// VersionInfo 微服务版本信息
type VersionInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}
