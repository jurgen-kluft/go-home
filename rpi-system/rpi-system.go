package main

// https://github.com/meinside/rpi-tools

// Raspberry PI system monitor
// - Uptime
// - Disk
// - Memory
// - CPU
// - CPU temperature

// Log metrics to StatsD (Influx.Telegraf)

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"os"
	"runtime"
	"strconv"
)

func dealwithErr(err error) {
	if err != nil {
		fmt.Println(err)
		//os.Exit(-1)
	}
}

func getHardwareData() string {
	runtimeOS := runtime.GOOS
	// memory
	vmStat, err := mem.VirtualMemory()
	dealwithErr(err)

	// disk - start from "/" mount point for Linux
	// might have to change for Windows!!
	// don't have a Window to test this out, if detect OS == windows
	// then use "\" instead of "/"
	diskStat, err := disk.Usage(string(rune(os.PathSeparator)))
	dealwithErr(err)

	// cpu - get CPU number of cores and speed
	cpuStat, err := cpu.Info()
	dealwithErr(err)
	percentage, err := cpu.Percent(0, true)
	dealwithErr(err)

	// host or machine kernel, uptime, platform Info
	hostStat, err := host.Info()
	dealwithErr(err)

	// get interfaces MAC/hardware address
	interfStat, err := net.Interfaces()
	dealwithErr(err)

	html := "<html>OS : " + runtimeOS + "<br>"
	html = html + "Total memory: " + strconv.FormatUint(vmStat.Total, 10) + " bytes <br>"
	html = html + "Free memory: " + strconv.FormatUint(vmStat.Free, 10) + " bytes<br>"
	html = html + "Percentage used memory: " + strconv.FormatFloat(vmStat.UsedPercent, 'f', 2, 64) + "%<br>"

	// get disk serial number.... strange... not available from disk package at compile time
	// undefined: disk.GetDiskSerialNumber
	//serial := disk.GetDiskSerialNumber("/dev/sda")

	//html = html + "Disk serial number: " + serial + "<br>"

	html = html + "Total disk space: " + strconv.FormatUint(diskStat.Total, 10) + " bytes <br>"
	html = html + "Used disk space: " + strconv.FormatUint(diskStat.Used, 10) + " bytes<br>"
	html = html + "Free disk space: " + strconv.FormatUint(diskStat.Free, 10) + " bytes<br>"
	html = html + "Percentage disk space usage: " + strconv.FormatFloat(diskStat.UsedPercent, 'f', 2, 64) + "%<br>"

	// since my machine has one CPU, I'll use the 0 index
	// if your machine has more than 1 CPU, use the correct index
	// to get the proper data
	html = html + "CPU index number: " + strconv.FormatInt(int64(cpuStat[0].CPU), 10) + "<br>"
	html = html + "VendorID: " + cpuStat[0].VendorID + "<br>"
	html = html + "Family: " + cpuStat[0].Family + "<br>"
	html = html + "Number of cores: " + strconv.FormatInt(int64(cpuStat[0].Cores), 10) + "<br>"
	html = html + "Model Name: " + cpuStat[0].ModelName + "<br>"
	html = html + "Speed: " + strconv.FormatFloat(cpuStat[0].Mhz, 'f', 2, 64) + " MHz <br>"

	for idx, cpupercent := range percentage {
		html = html + "Current CPU utilization: [" + strconv.Itoa(idx) + "] " + strconv.FormatFloat(cpupercent, 'f', 2, 64) + "%<br>"
	}

	html = html + "Hostname: " + hostStat.Hostname + "<br>"
	html = html + "Uptime: " + strconv.FormatUint(hostStat.Uptime, 10) + "<br>"
	html = html + "Number of processes running: " + strconv.FormatUint(hostStat.Procs, 10) + "<br>"

	// another way to get the operating system name
	// both darwin for Mac OSX, For Linux, can be ubuntu as platform
	// and linux for OS

	html = html + "OS: " + hostStat.OS + "<br>"
	html = html + "Platform: " + hostStat.Platform + "<br>"

	// the unique hardware id for this machine
	html = html + "Host ID(uuid): " + hostStat.HostID + "<br>"

	for _, interf := range interfStat {
		html = html + "------------------------------------------------------<br>"
		html = html + "Interface Name: " + interf.Name + "<br>"

		if interf.HardwareAddr != "" {
			html = html + "Hardware(MAC) Address: " + interf.HardwareAddr + "<br>"
		}

		for _, flag := range interf.Flags {
			html = html + "Interface behavior or flags: " + flag + "<br>"
		}

		for _, addr := range interf.Addrs {
			html = html + "IPv6 or IPv4 addresses: " + addr.String() + "<br>"

		}

	}

	html = html + "</html>"
	return html
}

func main() {

	fmt.Println(getHardwareData())
}
