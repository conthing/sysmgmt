package services

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/conthing/sysmgmt/dto"

	"github.com/conthing/utils/common"
)

const networkSettingFile = "/etc/network/interfaces"
const dhcpcdConfFile = "/etc/dhcpcd.conf"

// 从 /etc/network/interfaces 文件中读取网络设置，此文件必须
//    从"auto eth0开始"到"\n\n"结束是这样的设置：
//      DHCP时
//auto eth0
//iface eth0 inet dhcp
//
//      静态IP时
//auto eth0
//iface eth0 inet static
//	address 192.168.0.101
//	netmask	255.255.255.0
//	gateway 192.168.0.1
//
// #####兼容dhcpcd的设计
// 先判断是否存在 /etc/dhcpcd.conf 如果不存在就用 /etc/network/interfaces 如果存在就查找 "\ninterface eth0"
//      静态IP时
//# from "interface eth0" to "\n\n" will be read by sysmgmt
//interface eth0
//static ip_address=192.168.0.188/24
//static routers=192.168.0.1
//static domain_name_servers=192.168.0.1
//      DHCP时
// 没有这些行

func isDhcpcd() bool {
	_, err := os.Lstat(dhcpcdConfFile)
	return !os.IsNotExist(err)
}

func GetDhcpcdSetting(info *dto.NetInfo) error {
	data, err := ioutil.ReadFile(dhcpcdConfFile)
	if err != nil {
		return fmt.Errorf("network interfaces file read failed: %v", err)
	}
	ethx := common.GetMajorInterface() // 确定接口名称

	contents := strings.ReplaceAll(string(data), "\r", "") // 去除\r
	contents = strings.ReplaceAll(contents, "\t", " ")     // 去除Tab
	begin := strings.Index(contents, "\ninterface "+ethx)  // 找到 interface eth0 的位置
	if begin == -1 {
		info.DHCP = true
		info.Netmask = ""
		info.Gateway = ""
		info.Address = ""
		return nil
	}
	end := strings.Index(contents[begin:], "\n\n")
	if end == -1 {
		return fmt.Errorf("dhcpcd.conf file don't contain \"\\n\\n\"")
	}

	info.DHCP = false
	info.Netmask = ""
	info.Gateway = ""
	info.Address = ""

	setting := contents[begin : begin+end+1]
	//common.Log.Debugf("setting:%v", setting)
	setting = strings.ReplaceAll(setting, "static", "") // 去除"static"
	reader := bufio.NewReader(strings.NewReader(setting))
	for {
		line, isPrefix, err := reader.ReadLine()
		if err != nil {
			if info.Address != "" && info.Gateway != "" && info.Netmask != "" {
				return nil
			} else {
				return fmt.Errorf("dhcpcd.conf setting parse failed:\n%s", setting)
			}
		}
		if !isPrefix {
			strline := strings.TrimSpace((string(line)))
			index := strings.Index(strline, "=")
			if index != -1 && index+1 < len(strline) {
				name := strings.TrimSpace(strline[:index])
				value := strings.TrimSpace(strline[index+1:])
				if name == "ip_address" {
					index = strings.Index(value, "/")
					if index != -1 && index+1 < len(value) {
						ip := value[:index]
						if net.ParseIP(ip) != nil { // 前面是合法的IP字符串
							mask := value[index+1:]
							if n, err := strconv.ParseUint(mask, 10, 8); err == nil {
								if n >= 1 && n <= 31 {
									var netmaskbytes [4]byte
									for i := 0; i < int(n); i++ {
										netmaskbytes[i/8] |= (byte(1) << (7 - i%8))
									}
									netmask := fmt.Sprintf("%d.%d.%d.%d", netmaskbytes[0], netmaskbytes[1], netmaskbytes[2], netmaskbytes[3])
									info.Address = ip
									info.Netmask = netmask
								} else {
									return fmt.Errorf("invaild ip_address(4):%q", value)
								}
							} else {
								return fmt.Errorf("invaild ip_address(3):%q", value)
							}
						} else {
							return fmt.Errorf("invaild ip_address(2):%q", value)
						}
					} else {
						return fmt.Errorf("invaild ip_address(1):%q", value)
					}
				} else if name == "routers" {
					if net.ParseIP(value) != nil { // 前面是合法的IP字符串
						info.Gateway = value
					}
				}
			}
		}
	}
}

func netmaskToBits(netmask string) int {
	ip := net.ParseIP(netmask)
	if ip != nil {
		ipv4 := ip.To4()
		if ipv4 != nil {
			ones, _ := net.IPv4Mask(ipv4[0], ipv4[1], ipv4[2], ipv4[3]).Size()
			return ones
		}
	}
	return 0
}

func SetDhcpcdSetting(info *dto.NetInfo) error {
	data, err := ioutil.ReadFile(dhcpcdConfFile)
	if err != nil {
		return fmt.Errorf("network interfaces file read failed: %v", err)
	}
	ethx := common.GetMajorInterface() // 确定接口名称

	contents := strings.ReplaceAll(string(data), "\r", "") // 去除\r
	contents = strings.ReplaceAll(contents, "\t", " ")     // 去除Tab

	prefix := contents
	postfix := ""
	setting := ""

	begin := strings.Index(contents, "\ninterface "+ethx) // 找到 interface eth0 的位置
	if begin != -1 {
		prefix = contents[:begin]
		end := strings.Index(contents[begin:], "\n\n")
		if end != -1 {
			postfix = contents[begin+end:]
		}

	}
	if !info.DHCP {
		setting = fmt.Sprintf("\ninterface %s\nstatic ip_address=%s/%d\nstatic routers=%s\nstatic domain_name_servers=%s\n\n", ethx, info.Address, netmaskToBits(info.Netmask), info.Gateway, info.Gateway)
	}

	//common.Log.Debugf("setting:%v", setting)
	err = ioutil.WriteFile(dhcpcdConfFile, []byte(prefix+setting+postfix), 0777)
	if err != nil {
		return fmt.Errorf("dhcpcd conf file write failed: %v", err)
	}

	// 使网络设置立即生效
	go flushNetSetting(ethx)

	return nil

}

func GetCurrentNetInfo(info *dto.NetInfo) error {
	if isDhcpcd() {
		return GetDhcpcdSetting(info)
	}

	data, err := ioutil.ReadFile(networkSettingFile)
	if err != nil {
		return fmt.Errorf("network interfaces file read failed: %v", err)
	}
	ethx := common.GetMajorInterface() // 确定接口名称

	contents := strings.ReplaceAll(string(data), "\r", "") // 去除\r
	contents = strings.ReplaceAll(contents, "\t", " ")     // 去除Tab
	begin := strings.Index(contents, "auto "+ethx)
	if begin == -1 {
		return fmt.Errorf("network interfaces file don't contain \"auto %s\"", ethx)
	}
	end := strings.Index(contents[begin:], "\n\n")
	if end == -1 {
		return fmt.Errorf("network interfaces file don't contain \"\\n\\n\"")
	}
	//common.Log.Debugf("begin:%v,end:%v", begin, end)

	setting := contents[begin : begin+end]
	//common.Log.Debugf("setting:%v", setting)
	setting = strings.ReplaceAll(setting, "auto "+ethx, "")           // 去除"auto eth0"
	setting = strings.ReplaceAll(setting, "iface "+ethx+" inet ", "") // 去除"iface eth0 inet"
	//common.Log.Debugf("setting:%v", setting)

	strs := strings.Fields(setting)
	//common.Log.Debugf("strs:%v", strs)
	length := len(strs)
	for i, s := range strs {
		if s == "address" && i < length-1 {
			info.Address = strs[i+1]
			info.DHCP = false
		}
		if s == "netmask" && i < length-1 {
			info.Netmask = strs[i+1]
			info.DHCP = false
		}
		if s == "gateway" && i < length-1 {
			info.Gateway = strs[i+1]
			info.DHCP = false
		}
		if s == "dhcp" {
			info.DHCP = true
		}
	}

	return nil
}

func SetNetInfo(info *dto.NetInfo) error {
	if isDhcpcd() {
		return SetDhcpcdSetting(info)
	}

	data, err := ioutil.ReadFile(networkSettingFile)
	if err != nil {
		return fmt.Errorf("network interfaces file read failed: %v", err)
	}
	ethx := common.GetMajorInterface() // 确定接口名称

	contents := strings.ReplaceAll(string(data), "\r", "") // 去除\r
	begin := strings.Index(contents, "auto "+ethx)
	end := begin + strings.Index(contents[begin:], "\n\n")
	prefix := contents[:begin]
	postfix := contents[end:]

	setting := "auto " + ethx + "\niface " + ethx + " inet "
	if info.DHCP {
		setting = setting + "dhcp"
	} else {
		setting = setting + "static\n\taddress " + info.Address + "\n\tnetmask " + info.Netmask + "\n\tgateway " + info.Gateway
	}

	//common.Log.Debugf("setting:%v", setting)
	err = ioutil.WriteFile(networkSettingFile, []byte(prefix+setting+postfix), 0777)
	if err != nil {
		return fmt.Errorf("network interfaces file write failed: %v", err)
	}

	// 使网络设置立即生效
	go flushNetSetting(ethx)

	return nil
}

func flushNetSetting(ethx string) {
	if isDhcpcd() {
		flushDhcpcdSetting()
		return
	}

	command := exec.Command("ip", "addr", "flush", "dev", ethx) //初始化Cmd
	_, _ = command.Output()
	//common.Log.Debugf("ip addr flush output:%s", string(out))

	command = exec.Command("/etc/init.d/networking", "restart") //初始化Cmd
	_, _ = command.Output()
	//common.Log.Debugf("modified ip output:%s", string(out))
}

func flushDhcpcdSetting() {
	command := exec.Command("systemctl", "restart", "dhcpcd") //初始化Cmd
	_, _ = command.Output()
	//common.Log.Debugf("modified ip output:%s", string(out))
}
