package system

import (
	"bufio"
	"bytes"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type Metrics struct {
	Hostname      string
	OSInfo        string
	KernelVersion string
	CPUUsage      float64
	MemoryUsage   float64
	DiskUsage     float64
	LoadAverage1  float64
	LoadAverage5  float64
	LoadAverage15 float64
	UptimeSeconds int64
}

type ProcessInfo struct {
	PID      int
	Name     string
	User     string
	CPUUsage float64
	MemUsage float64
	Command  string
}

type PortInfo struct {
	Protocol   string
	LocalAddr  string
	Port       int
	ProcessName string
	PID        int
	State      string
}

type ServiceInfo struct {
	Name        string
	Status      string
	ActiveState string
	SubState    string
}

// CollectMetrics is the entry point for metrics collection
func CollectMetrics() (*Metrics, error) {
	return Collect()
}

func Collect() (*Metrics, error) {
	m := &Metrics{}

	if hostname, err := getHostname(); err == nil {
		m.Hostname = hostname
	}
	if osInfo, err := getOSInfo(); err == nil {
		m.OSInfo = osInfo
	}
	if cpu, err := getCPUUsage(); err == nil {
		m.CPUUsage = cpu
	}
	if mem, err := getMemoryUsage(); err == nil {
		m.MemoryUsage = mem
	}
	if disk, err := getDiskUsage(); err == nil {
		m.DiskUsage = disk
	}
	if la1, la5, la15, err := getLoadAverages(); err == nil {
		m.LoadAverage1 = la1
		m.LoadAverage5 = la5
		m.LoadAverage15 = la15
	}
	if uptime, err := getUptime(); err == nil {
		m.UptimeSeconds = uptime
	}
	if kernel, err := readKernelVersion(); err == nil {
		m.KernelVersion = kernel
	}

	return m, nil
}

func getHostname() (string, error) {
	cmd := exec.Command("hostname")
	output, err := cmd.Output()
	if err != nil {
		return "unknown", nil
	}
	return strings.TrimSpace(string(output)), nil
}

func getOSInfo() (string, error) {
	cmd := exec.Command("sw_vers", "-productVersion")
	output, err := cmd.Output()
	if err != nil {
		// Try uname -a as fallback
		cmd = exec.Command("uname", "-a")
		output, err = cmd.Output()
		if err != nil {
			return "unknown", nil
		}
		return strings.TrimSpace(string(output)), nil
	}
	return strings.TrimSpace(string(output)), nil
}

func readKernelVersion() (string, error) {
	// Try sysctl first (macOS), fall back to uname (Linux)
	cmd := exec.Command("sysctl", "-n", "kern.osrelease")
	output, err := cmd.Output()
	if err != nil {
		// Try Linux uname
		cmd = exec.Command("uname", "-r")
		output, err = cmd.Output()
		if err != nil {
			return "unknown", nil
		}
	}
	return strings.TrimSpace(string(output)), nil
}

func getCPUUsage() (float64, error) {
	cmd := exec.Command("sh", "-c", `top -l 1 -n 1 | grep "CPU usage" | awk '{print $3}' | tr -d '%'`)
	output, err := cmd.Output()
	if err != nil {
		// Try Linux alternative
		cmd = exec.Command("sh", "-c", `cat /proc/stat | head -1 | awk '{usage=($2+$4)*100/($2+$4+$5)} END {print usage}'`)
		output, err = cmd.Output()
		if err != nil {
			return 0, err
		}
	}
	usageStr := strings.TrimSpace(string(output))
	return strconv.ParseFloat(usageStr, 64)
}

func getMemoryUsage() (float64, error) {
	cmd := exec.Command("sh", "-c", `vm_stat | grep "Pages active" | awk '{print $3}' | tr -d '.'`)
	output, err := cmd.Output()
	if err != nil {
		// Try Linux
		cmd = exec.Command("sh", "-c", `free | grep Mem | awk '{print ($3/$2) * 100.0}'`)
		output, err = cmd.Output()
		if err != nil {
			return 0, err
		}
	}
	memPages := strings.TrimSpace(string(output))
	memPct, err := strconv.ParseFloat(memPages, 64)
	if err != nil {
		return 0, err
	}
	// vm_stat pages to GB-ish approximation
	return memPct * 0.01, nil
}

func getDiskUsage() (float64, error) {
	cmd := exec.Command("sh", "-c", `df -h / | tail -1 | awk '{print $5}' | tr -d '%'`)
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}
	diskStr := strings.TrimSpace(string(output))
	return strconv.ParseFloat(diskStr, 64)
}

func getLoadAverages() (float64, float64, float64, error) {
	cmd := exec.Command("sh", "-c", `sysctl -n vm.loadavg | tr -d '{' | tr -d '}'`)
	output, err := cmd.Output()
	if err != nil {
		// Try Linux
		cmd = exec.Command("sh", "-c", `uptime | awk -F'load averages:' '{print $2}' | awk '{print $1}' | tr -d ','`)
		output, err = cmd.Output()
		if err != nil {
			return 0, 0, 0, err
		}
	}
	parts := strings.Fields(string(output))
	if len(parts) < 3 {
		return 0, 0, 0, nil
	}
	la1, _ := strconv.ParseFloat(strings.Trim(parts[0], "{,"), 64)
	la5, _ := strconv.ParseFloat(strings.Trim(parts[1], ","), 64)
	la15, _ := strconv.ParseFloat(strings.Trim(parts[2], "}"), 64)
	return la1, la5, la15, nil
}

func getUptime() (int64, error) {
	cmd := exec.Command("sysctl", "-n", "kern.uptime")
	output, err := cmd.Output()
	if err != nil {
		// Try reading /proc/uptime on Linux
		raw, err := os.ReadFile("/proc/uptime")
		if err != nil {
			return 0, err
		}
		parts := strings.Fields(string(raw))
		if len(parts) == 0 {
			return 0, nil
		}
		uptimeFloat, _ := strconv.ParseFloat(parts[0], 64)
		return int64(uptimeFloat), nil
	}

	parts := strings.Fields(string(output))
	if len(parts) == 0 {
		return 0, nil
	}

	uptimeFloat, _ := strconv.ParseFloat(parts[0], 64)
	return int64(uptimeFloat), nil
}

func GetProcessList() ([]ProcessInfo, error) {
	cmd := exec.Command("ps", "aux")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var processes []ProcessInfo
	scanner := bufio.NewScanner(bytes.NewReader(output))
	if scanner.Scan() {
		scanner.Text() // skip header
	}

	re := regexp.MustCompile(`\s+`)
	for scanner.Scan() {
		fields := re.Split(strings.TrimSpace(scanner.Text()), 11)
		if len(fields) < 11 {
			continue
		}
		pid, _ := strconv.Atoi(fields[1])
		user := fields[0]
		cpu, _ := strconv.ParseFloat(fields[2], 64)
		mem, _ := strconv.ParseFloat(fields[3], 64)
		comm := fields[10]

		processes = append(processes, ProcessInfo{
			PID:      pid,
			User:     user,
			CPUUsage: cpu,
			MemUsage: mem,
			Command:  comm,
		})
	}
	return processes, nil
}

func GetPortList() ([]PortInfo, error) {
	cmd := exec.Command("lsof", "-i", "-P", "-n")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var ports []PortInfo
	scanner := bufio.NewScanner(bytes.NewReader(output))
	if scanner.Scan() {
		scanner.Text() // skip header
	}

	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 9 {
			continue
		}
		procName := fields[0]
		proto := fields[7]
		localAddr := fields[8]
		pid, _ := strconv.Atoi(fields[1])

		// Parse address and port
		parts := strings.Split(localAddr, "->")
		local := parts[0]
		var port int
		if len(parts) > 1 {
			portParts := strings.Split(parts[1], ":")
			port, _ = strconv.Atoi(portParts[len(portParts)-1])
		} else {
			portParts := strings.Split(local, ":")
			port, _ = strconv.Atoi(portParts[len(portParts)-1])
		}

		ports = append(ports, PortInfo{
			Protocol:   proto,
			LocalAddr:  local,
			Port:       port,
			ProcessName: procName,
			PID:        pid,
			State:      "LISTEN",
		})
	}
	return ports, nil
}

func GetServiceList() ([]ServiceInfo, error) {
	// Use launchctl on macOS, systemctl on Linux
	var services []ServiceInfo

	// Try macOS launchctl
	cmd := exec.Command("launchctl", "list")
	output, err := cmd.Output()
	if err == nil {
		scanner := bufio.NewScanner(bytes.NewReader(output))
		if scanner.Scan() {
			scanner.Text() // skip header
		}
		for scanner.Scan() {
			fields := strings.Fields(scanner.Text())
			if len(fields) < 3 {
				continue
			}
			pid, _ := strconv.Atoi(fields[0])
			status := fields[1]
			name := fields[2]

			var state string
			if pid > 0 {
				state = "running"
			} else if status == "-" {
				state = "stopped"
			} else {
				state = status
			}

			services = append(services, ServiceInfo{
				Name:        name,
				Status:      status,
				ActiveState: state,
				SubState:    state,
			})
		}
	}

	return services, nil
}

// DeepScanResult holds the results of a deep scan
type DeepScanResult struct {
	TopCPUProcesses []ProcessInfo
	OpenPorts       []PortInfo
	Services        []ServiceInfo
}

// RunDeepScan performs a comprehensive scan of the system
func RunDeepScan() (*DeepScanResult, error) {
	result := &DeepScanResult{}

	// Get top CPU processes
	procs, err := GetProcessList()
	if err == nil && len(procs) > 0 {
		// Get top 10 by CPU
		for _, p := range procs {
			if len(result.TopCPUProcesses) >= 10 {
				break
			}
			result.TopCPUProcesses = append(result.TopCPUProcesses, p)
		}
	}

	// Get open ports
	ports, err := GetPortList()
	if err == nil {
		result.OpenPorts = ports
	}

	// Get services
	svcs, err := GetServiceList()
	if err == nil {
		result.Services = svcs
	}

	return result, nil
}
