package dto

// NetInfo 设置网络的结构体
type NetInfo struct {
	DHCP    bool   `json:"dhcp" `
	Address string `json:"address" `
	Netmask string `json:"netmask" `
	Gateway string `json:"gateway" `
}

// NTPInfo 设置时区的结构体
type NTPInfo struct {
	Type string `json:"type" `
	Date int64  `json:"date" `
	URL  string `json:"url" `
}

// SubVersionInfo 微服务版本信息
type SubVersionInfo struct {
	Name      string `json:"name"`
	Version   string `json:"version"`
	BuildTime string `json:"buildtime"`
}

// VersionInfo 整体版本信息
type VersionInfo struct {
	Version     string           `json:"version"`
	Description string           `json:"description"`
	SubVersion  []SubVersionInfo `json:"subversion"`
}

// MacInfo Mac 信息
type MacInfo struct {
	Mac string `json:"mac"`
}

//TimeInfo 时间信息
type TimeInfo struct {
	Time      int64  `json:"time"`
	Ntpstatus bool   `json:"ntpstatus"`
	NtpURL    string `json:"ntpurl"`
}

// FileInfo 文件信息
type FileInfo struct {
	Downloading bool   `json:"downloading"`
	Message     string `json:"message"`
}
