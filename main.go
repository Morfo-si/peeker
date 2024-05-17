package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/charmbracelet/lipgloss"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"

	"golang.org/x/term"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const megabyteDiv uint64 = 1024 * 1024
const gigabyteDiv uint64 = megabyteDiv * 1024

var (
	diskStyle = lipgloss.NewStyle().
			Inherit(statusBarStyle).
			Align(lipgloss.Right).
			Padding(0, 1)

	memoryStyle = lipgloss.NewStyle().
			Inherit(statusBarStyle).
			Align(lipgloss.Left).
			Padding(0, 1)

	statusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#8CABFF"}).
			Background(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#512B81"})

	cpuStyle = lipgloss.NewStyle().
			Inherit(statusBarStyle).
			Foreground(lipgloss.Color("8CABFF")).
			Background(lipgloss.Color("#35155D")).
			Bold(true).
			Padding(0).
			MarginRight(1)

	highlightLeftStyle = lipgloss.NewStyle().
				Inherit(statusBarStyle).
		// Foreground(lipgloss.Color("#5AB2FF")).
		Background(lipgloss.Color("#35155D")).
		Bold(true).
		Padding(0, 1).
		MarginRight(1)

	highlightRightStyle = lipgloss.NewStyle().
				Inherit(statusBarStyle).
		// Foreground(lipgloss.Color("#5AB2FF")).
		Background(lipgloss.Color("#35155D")).
		Bold(true).
		Padding(0, 1).
		MarginLeft(1)

	generalTextStyle = lipgloss.NewStyle().
				Inherit(statusBarStyle).
				Foreground(lipgloss.Color("#FFFFFF")).
				Bold(true)
)

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

type StatusBar struct {
	cpu  []cpu.InfoStat
	disk *disk.UsageStat
	host *host.InfoStat
	mem  *mem.VirtualMemoryStat
}

func NewStatusBar() *StatusBar {
	return &StatusBar{}
}

func (sb *StatusBar) WithHostInformation() *StatusBar {
	if host, err := GetHostInformation(); err == nil {
		sb.host = host
	}
	return sb
}

func (sb *StatusBar) WithCPUInformation() *StatusBar {
	if cpu, err := GetCPUStat(); err == nil {
		sb.cpu = cpu
	}
	return sb
}

func (sb *StatusBar) WithMemoryInformation() *StatusBar {
	if mem, err := GetHostMemory(); err == nil {
		sb.mem = mem
	}
	return sb
}

func (sb *StatusBar) WithDiskInformation() *StatusBar {
	if disk, err := GetDiskInformation(); err == nil {
		sb.disk = disk
	}
	return sb
}

func (sb StatusBar) Render() {
	terminalWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))
	w := lipgloss.Width

	// Platform information
	platformName := cases.Title(language.AmericanEnglish).String(sb.host.Platform)
	platformVersion := sb.host.PlatformVersion
	// e.g., Darwin 14.4.1
	platformInformation := fmt.Sprintf("%s %s", platformName, platformVersion)

	// Host information
	hostName := sb.host.Hostname
	hostArch := sb.host.KernelArch
	// e.g., localhost.local arm64
	hostInformation := fmt.Sprintf("%s %s", hostName, hostArch)

	// CPU information
	var cpuInformation string
	if len(sb.cpu) != 0 {
		cpuInformation = fmt.Sprintf("%s %s MHz", sb.cpu[0].ModelName, strconv.FormatFloat(sb.cpu[0].Mhz, 'f', 2, 64))
	}

	platformCell := highlightLeftStyle.Render(platformInformation)
	cpuCell := highlightRightStyle.Render(cpuInformation)
	hostVal := generalTextStyle.
		Width(terminalWidth - w(platformCell) - w(cpuCell)).
		Align(lipgloss.Center).
		Render(hostInformation)

	barLine1 := lipgloss.JoinHorizontal(lipgloss.Top,
		platformCell,
		hostVal,
		cpuCell,
	)

	// Disk information
	diskTotal := sb.disk.Total / gigabyteDiv
	diskUsed := sb.disk.Used / gigabyteDiv
	diskUsedPercent := sb.disk.UsedPercent
	diskInformation := fmt.Sprintf("Disk: %d of %d used (%2.f %%)", diskUsed, diskTotal, diskUsedPercent)
	diskCell := diskStyle.Align(lipgloss.Right).Render(diskInformation)

	// Memory information
	memoryTotal := sb.mem.Total / megabyteDiv
	memoryAvailable := sb.mem.Available / megabyteDiv
	memoryInUse := memoryTotal - memoryAvailable
	memoryUsedPercentenge := sb.mem.UsedPercent
	memoryInformation := fmt.Sprintf("Memory used: %d of %d MB (%2.f%% used)", memoryInUse, memoryTotal, memoryUsedPercentenge)

	memoryCell := memoryStyle.
		Width(terminalWidth - w(diskCell)).
		Render(memoryInformation)

	barLine2 := lipgloss.JoinHorizontal(lipgloss.Top,
		memoryCell,
		diskCell,
	)

	bar := lipgloss.JoinVertical(lipgloss.Top,
		barLine1, barLine2,
	)
	fmt.Println(bar)

}

func main() {
	bar := NewStatusBar().
		WithHostInformation().
		WithCPUInformation().
		WithMemoryInformation().
		WithDiskInformation()
	bar.Render()
}
