package services

import (
	"fmt"
	"os/exec"
	"strings"
)

const NTPConfigFile = "/etc/ntp.conf"

func GetNtpEnable() (bool, error) {
	command := exec.Command("/bin/sh", "-c", `service ntp status | grep "Active:"`)
	if out, err := command.Output(); err != nil {
		return false, fmt.Errorf("exec service ntp status failed: %w", err)
	} else if len(out) == 0 {
		return false, fmt.Errorf("exec service ntp status failed: no output")
	} else {
		if strings.Contains(string(out), "active (running)") {
			return true, nil
		} else {
			return false, nil
		}
	}
}

func SetNtpEnable(en bool) error {
	var command *exec.Cmd
	if en {
		command = exec.Command("/bin/sh", "-c", `service ntp open`)
		if _, err := command.Output(); err != nil {
			return fmt.Errorf("exec service ntp open failed: %w", err)
		}
		command = exec.Command("/bin/sh", "-c", `service ntp restart`)
		if _, err := command.Output(); err != nil {
			return fmt.Errorf("exec service ntp restart failed: %w", err)
		}
	} else {
		command = exec.Command("/bin/sh", "-c", `service ntp close`)
		if _, err := command.Output(); err != nil {
			return fmt.Errorf("exec service ntp close failed: %w", err)
		}
	}
	if res, err := GetNtpEnable(); err != nil {
		return err
	} else {
		if res != en {
			return fmt.Errorf("SetNtpEnable failed")
		}
	}
	return nil
}

func GetNtpServer() (string, error) {
	// sed -n 's/server \(.*\) iburst/\1/p' /etc/ntp.conf
	command := exec.Command("sed", "-n", `s/server \(.*\) iburst/\1/p`, NTPConfigFile)
	if out, err := command.Output(); err != nil {
		return "", fmt.Errorf("sed %q failed: %w", NTPConfigFile, err)
	} else if len(out) == 0 {
		return "", fmt.Errorf("sed %q failed: no output", NTPConfigFile)
	} else {
		return strings.TrimSpace(string(out)), nil
	}
}

func SetNtpServer(server string) error {
	// sed -i "/Autoedit section start/,/Autoedit section end/s/server \(.*\) iburst/server {{server}} iburst/" ./ntp.conf
	sedcmd := `/Autoedit section start/,/Autoedit section end/s/server \(.*\) iburst/server `
	sedcmd += server
	sedcmd += ` iburst/`
	command := exec.Command("sed", "-i", sedcmd, NTPConfigFile)
	if out, err := command.Output(); err != nil {
		return fmt.Errorf("sed %q failed: %w", NTPConfigFile, err)
	} else if len(out) != 0 {
		return fmt.Errorf("sed %q failed: output: %q", NTPConfigFile, string(out))
	} else {
		if res, err := GetNtpServer(); err != nil {
			return err
		} else if res != strings.TrimSpace(server) {
			return fmt.Errorf("SetNtpServer failed")
		}
		return nil
	}
}
