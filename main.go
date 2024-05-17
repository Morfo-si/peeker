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

// Constants used to convert units to MB and GB
const megabyteDiv uint64 = 1024 * 1024
const gigabyteDiv uint64 = megabyteDiv * 1024

var (
	// The width of the visible terminal.
	terminalWidth, _, _ = term.GetSize(int(os.Stdout.Fd()))

	// Style used for the disk information.
	diskStyle = lipgloss.NewStyle().
			Inherit(statusBarStyle).
			Align(lipgloss.Right).
			Padding(0, 1)

	// Style used for the memory information.
	memoryStyle = lipgloss.NewStyle().
			Inherit(statusBarStyle).
			Align(lipgloss.Left).
			Padding(0, 1)

	// Base style used for the status bar.
	statusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#8CABFF"}).
			Background(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#512B81"})

	// Style used for the upper left corner of the status bar.
	highlightLeftStyle = lipgloss.NewStyle().
				Inherit(statusBarStyle).
				Background(lipgloss.Color("#35155D")).
				Bold(true).
				Padding(0, 1).
				MarginRight(1)

	// Style used for the upper right corner of the status bar.
	highlightRightStyle = lipgloss.NewStyle().
				Inherit(statusBarStyle).
				Background(lipgloss.Color("#35155D")).
				Bold(true).
				Padding(0, 1).
				MarginLeft(1)

	// Style used for generic text.
	generalTextStyle = lipgloss.NewStyle().
				Inherit(statusBarStyle).
				Foreground(lipgloss.Color("#FFFFFF")).
				Bold(true)
)

// DisplayHostMemory displays information about the host's memory.
func DisplayHostMemory(sb StatusBar, width int) string {
	// Memory information
	var (
		memoryTotal, memoryAvailable, memoryInUse uint64
		memoryUsedPercentenge                     float64
	)

	if sb.mem != nil {
		memoryTotal = sb.mem.Total / megabyteDiv
		memoryAvailable = sb.mem.Available / megabyteDiv
		memoryInUse = memoryTotal - memoryAvailable
		memoryUsedPercentenge = sb.mem.UsedPercent
	}
	memoryInformation := fmt.Sprintf("Memory: %d of %d MB used (%2.f%%)", memoryInUse, memoryTotal, memoryUsedPercentenge)

	return memoryStyle.
		Width(width).
		Render(memoryInformation)
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
func DisplayHostInformation(sb StatusBar, width int) string {
	// Host information
	var hostName, hostArch string

	if sb.host != nil {
		hostName = sb.host.Hostname
		hostArch = sb.host.KernelArch
	}
	// e.g., localhost.local arm64
	hostInformation := fmt.Sprintf("%s %s", hostName, hostArch)
	return generalTextStyle.
		Width(width).
		Align(lipgloss.Center).
		Render(hostInformation)
}

// DisplayPlatformInformation displays information about a host.
func DisplayPlatformInformation(sb StatusBar) string {
	// Platform information
	var platformName, platformVersion string

	if sb.host != nil {
		platformName = cases.Title(language.AmericanEnglish).String(sb.host.Platform)
		platformVersion = sb.host.PlatformVersion
	}
	// e.g., Darwin 14.4.1
	platformInformation := fmt.Sprintf("%s %s", platformName, platformVersion)

	return highlightLeftStyle.Render(platformInformation)
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
func DisplayDiskInformation(sb StatusBar) string {
	// Disk information
	var (
		diskTotal, diskUsed uint64
		diskUsedPercent     float64
	)

	if sb.disk != nil {
		diskTotal = sb.disk.Total / gigabyteDiv
		diskUsed = sb.disk.Used / gigabyteDiv
		diskUsedPercent = sb.disk.UsedPercent
	}
	diskInformation := fmt.Sprintf("Disk: %d of %d GB used (%2.f%%)", diskUsed, diskTotal, diskUsedPercent)
	return diskStyle.Align(lipgloss.Right).Render(diskInformation)
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

// DisplayCPUInformation displays the CPU information.
func DisplayCPUInformation(sb StatusBar) string {
	// CPU information
	var cpuInformation string

	if len(sb.cpu) != 0 {
		cpuInformation = fmt.Sprintf("%s %s MHz", sb.cpu[0].ModelName, strconv.FormatFloat(sb.cpu[0].Mhz, 'f', 2, 64))
	}
	return highlightRightStyle.Render(cpuInformation)
}

// GetCPUStat returns only one CPUInfoStat on FreeBSD
func GetCPUStat() ([]cpu.InfoStat, error) {
	cpuStat, err := cpu.Info()
	if err != nil {
		return nil, err
	}

	return cpuStat, nil
}

// Struct used to represent a StatusBar
type StatusBar struct {
	cpu  []cpu.InfoStat
	disk *disk.UsageStat
	host *host.InfoStat
	mem  *mem.VirtualMemoryStat
}

// New StatusBar with no features.
func NewStatusBar() *StatusBar {
	return &StatusBar{}
}

// New StatusBar with feature.
func (sb *StatusBar) WithHostInformation() *StatusBar {
	if host, err := GetHostInformation(); err == nil {
		sb.host = host
	}
	return sb
}

// New StatusBar with CPU feature.
func (sb *StatusBar) WithCPUInformation() *StatusBar {
	if cpu, err := GetCPUStat(); err == nil {
		sb.cpu = cpu
	}
	return sb
}

// New StatusBar with Memory feature.
func (sb *StatusBar) WithMemoryInformation() *StatusBar {
	if mem, err := GetHostMemory(); err == nil {
		sb.mem = mem
	}
	return sb
}

// New StatusBar with Disk feature.
func (sb *StatusBar) WithDiskInformation() *StatusBar {
	if disk, err := GetDiskInformation(); err == nil {
		sb.disk = disk
	}
	return sb
}

// Renders the StatusBar with all features.
func (sb StatusBar) Render() {
	// Shortcut to get accurate width from a given string.
	w := lipgloss.Width

	// Platform information
	platformCell := DisplayPlatformInformation(sb)
	// CPU information
	cpuCell := DisplayCPUInformation(sb)
	// Host information
	hostCell := DisplayHostInformation(sb, terminalWidth-w(platformCell)-w(cpuCell))
	// Disk information
	diskCell := DisplayDiskInformation(sb)
	// Memory information
	memoryCell := DisplayHostMemory(sb, terminalWidth-w(diskCell))

	// Top line for status bar.
	firstLine := lipgloss.JoinHorizontal(lipgloss.Top,
		platformCell,
		hostCell,
		cpuCell,
	)
	// Bottom line for status bar.
	secondLine := lipgloss.JoinHorizontal(lipgloss.Top,
		memoryCell,
		diskCell,
	)

	bar := lipgloss.JoinVertical(lipgloss.Top,
		firstLine, secondLine,
	)
	fmt.Println(bar)

}

func main() {
	// Display status bar with system information.
	bar := NewStatusBar().
		WithHostInformation().
		WithCPUInformation().
		WithMemoryInformation().
		WithDiskInformation()
	bar.Render()
}
