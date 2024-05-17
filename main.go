package main

import (
	"fmt"
	"strconv"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
)

const megabyteDiv uint64 = 1024 * 1024
const gigabyteDiv uint64 = megabyteDiv * 1024

// DisplayHostMemory displays information about the host's memory.
func DisplayHostMemory(vmStat *mem.VirtualMemoryStat) {
	fmt.Println("Total memory:", strconv.FormatUint(vmStat.Total/megabyteDiv, 10), "MB")
	fmt.Println("Free memory:", strconv.FormatUint(vmStat.Available/megabyteDiv, 10), "MB")
	fmt.Println("Percentage used memory:", strconv.FormatFloat(vmStat.UsedPercent, 'f', 2, 64), "%")
}

// GetHostMemory fetches information about the host's memory.
func GetHostMemory() (*mem.VirtualMemoryStat, error) {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}
	return vmStat, nil
}

// DisplayHostInformation displays information about a host.
func DisplayHostInformation(hostStat *host.InfoStat) {
	fmt.Printf("Hostname: %v\n", hostStat.Hostname)
	fmt.Printf("Operating System: %v %v (%v)\n", hostStat.Platform, hostStat.PlatformVersion, hostStat.KernelArch)
}

// GetHostInformation fetches information for the host.
func GetHostInformation() (*host.InfoStat, error) {
	hostStat, err := host.Info()
	if err != nil {
		return nil, err
	}
	return hostStat, nil
}

// DisplayDiskInformation displays disk information on the console
func DisplayDiskInformation(diskStat *disk.UsageStat) {
	fmt.Printf("Total disk space: %v GB\n", strconv.FormatUint(diskStat.Total/gigabyteDiv, 10))
	fmt.Printf("Used disk space: %v GB\n", strconv.FormatUint(diskStat.Used/gigabyteDiv, 10))
	fmt.Printf("Free disk space: %v GB\n", strconv.FormatUint(diskStat.Free/gigabyteDiv, 10))
	fmt.Printf("Percentage disk space usage: %v %%\n", strconv.FormatFloat(diskStat.UsedPercent, 'f', 2, 64))
}

// GetDiskInformation returns the file system usage.
func GetDiskInformation() (*disk.UsageStat, error) {
	diskStat, err := disk.Usage("/")
	if err != nil {
		return nil, err
	}
	return diskStat, nil
}

// DisplayCPUPercentage displays CPU usage percentage.
func DisplayCPUPercentage(percentage []float64) {
	firstCpus := percentage[:len(percentage)/2]
	secondCpus := percentage[len(percentage)/2:]

	fmt.Println("Cores:")
	for idx, cpupercent := range firstCpus {
		fmt.Printf("\tCPU [%v]: %v %%\n", strconv.Itoa(idx), strconv.FormatFloat(cpupercent, 'f', 2, 64))
	}

	oftset := len(firstCpus)
	for idx, cpupercent := range secondCpus {
		fmt.Printf("\tCPU [%v]: %v %%\n", strconv.Itoa(idx+oftset), strconv.FormatFloat(cpupercent, 'f', 2, 64))
	}
}

// GetCPUPercentage calculates the percentage of cpu used either per CPU or combined.
func GetCPUPercentage() ([]float64, error) {
	percentage, err := cpu.Percent(0, true)
	if err != nil {
		return nil, err
	}
	return percentage, nil
}

// DisplayCPUStat displays the CPU information.
func DisplayCPUStat(cpuStat []cpu.InfoStat) {
	if len(cpuStat) != 0 {
		fmt.Printf("Model Name: %v ", cpuStat[0].ModelName)
		fmt.Printf("Family: %v ", cpuStat[0].Family)
		fmt.Printf("Speed: %v MHz\n", strconv.FormatFloat(cpuStat[0].Mhz, 'f', 2, 64))
	}
}

// GetCPUStat returns only one CPUInfoStat on FreeBSD
func GetCPUStat() ([]cpu.InfoStat, error) {
	cpuStat, err := cpu.Info()
	if err != nil {
		return nil, err
	}

	return cpuStat, nil
}

func main() {
	if host, err := GetHostInformation(); err == nil {
		DisplayHostInformation(host)
	}

	if stat, err := GetCPUStat(); err == nil {
		DisplayCPUStat(stat)
	}

	if mem, err := GetHostMemory(); err == nil {
		DisplayHostMemory(mem)
	}

	if disk, err := GetDiskInformation(); err == nil {
		DisplayDiskInformation(disk)
	}
	if percentage, err := GetCPUPercentage(); err == nil {
		DisplayCPUPercentage(percentage)
	}
}
