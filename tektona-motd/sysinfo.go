package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type Session struct {
	User string
	Line string
	From string
	When string
}

type Info struct {
	Hostname    string
	OS          string
	Kernel      string
	Uptime      time.Duration
	MemTotalKB  uint64
	MemAvailKB  uint64
	DiskTotalB  uint64
	DiskFreeB   uint64
	Load1       float64
	NCPU        int
	IPs         []string
	User        string
	SessionType string // "SSH from x.x.x.x", "VNC :0", or "interactive"
	SandboxID   string
	Sessions    []Session
}

func gather() Info {
	i := Info{
		NCPU: runtime.NumCPU(),
	}
	if h, err := os.Hostname(); err == nil {
		i.Hostname = h
	}
	if u, err := user.Current(); err == nil {
		i.User = u.Username
	}
	i.OS = readOSPrettyName()
	i.Kernel = readKernel()
	i.Uptime = readUptime()
	i.MemTotalKB, i.MemAvailKB = readMeminfo()
	i.DiskTotalB, i.DiskFreeB = readDisk("/")
	i.Load1 = readLoad1()
	i.IPs = readIPs()
	i.SessionType = readSessionType()
	i.SandboxID = readSandboxID()
	i.Sessions = readSessions()
	return i
}

func readOSPrettyName() string {
	f, err := os.Open("/etc/os-release")
	if err != nil {
		return runtime.GOOS
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()
		if strings.HasPrefix(line, "PRETTY_NAME=") {
			return strings.Trim(strings.TrimPrefix(line, "PRETTY_NAME="), `"`)
		}
	}
	return runtime.GOOS
}

func readKernel() string {
	b, err := os.ReadFile("/proc/sys/kernel/osrelease")
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(b))
}

func readUptime() time.Duration {
	b, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return 0
	}
	parts := strings.Fields(string(b))
	if len(parts) == 0 {
		return 0
	}
	secs, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return 0
	}
	return time.Duration(secs * float64(time.Second))
}

func readMeminfo() (total, avail uint64) {
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		return 0, 0
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		fields := strings.Fields(s.Text())
		if len(fields) < 2 {
			continue
		}
		switch fields[0] {
		case "MemTotal:":
			total, _ = strconv.ParseUint(fields[1], 10, 64)
		case "MemAvailable:":
			avail, _ = strconv.ParseUint(fields[1], 10, 64)
		}
	}
	return total, avail
}

func readDisk(path string) (total, free uint64) {
	var st syscall.Statfs_t
	if err := syscall.Statfs(path, &st); err != nil {
		return 0, 0
	}
	total = st.Blocks * uint64(st.Bsize)
	free = st.Bavail * uint64(st.Bsize)
	return total, free
}

func readLoad1() float64 {
	b, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		return 0
	}
	parts := strings.Fields(string(b))
	if len(parts) == 0 {
		return 0
	}
	v, _ := strconv.ParseFloat(parts[0], 64)
	return v
}

func readIPs() []string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil
	}
	var out []string
	for _, ifc := range ifaces {
		if ifc.Flags&net.FlagLoopback != 0 || ifc.Flags&net.FlagUp == 0 {
			continue
		}
		addrs, err := ifc.Addrs()
		if err != nil {
			continue
		}
		for _, a := range addrs {
			ipn, ok := a.(*net.IPNet)
			if !ok || ipn.IP.IsLinkLocalUnicast() {
				continue
			}
			if v4 := ipn.IP.To4(); v4 != nil {
				out = append(out, v4.String())
			}
		}
	}
	return out
}

func readSessionType() string {
	if c := os.Getenv("SSH_CONNECTION"); c != "" {
		fields := strings.Fields(c)
		if len(fields) >= 1 {
			return "SSH from " + fields[0]
		}
		return "SSH"
	}
	if d := os.Getenv("DISPLAY"); d != "" {
		return "VNC " + d
	}
	return "interactive"
}

func readSandboxID() string {
	if v := os.Getenv("TEKTONA_SANDBOX_ID"); v != "" {
		return v
	}
	for _, p := range []string{"/etc/tektona/sandbox-id", "/run/tektona/sandbox-id"} {
		if b, err := os.ReadFile(p); err == nil {
			if s := strings.TrimSpace(string(b)); s != "" {
				return s
			}
		}
	}
	return ""
}

func readSessions() []Session {
	out, err := exec.Command("who").Output()
	if err != nil {
		return nil
	}
	var sessions []Session
	s := bufio.NewScanner(strings.NewReader(string(out)))
	for s.Scan() {
		fields := strings.Fields(s.Text())
		if len(fields) < 2 {
			continue
		}
		sess := Session{User: fields[0], Line: fields[1]}
		if len(fields) >= 4 {
			sess.When = fields[2] + " " + fields[3]
		}
		// Remaining fields may contain "(host)" for remote logins.
		for _, f := range fields[4:] {
			if strings.HasPrefix(f, "(") && strings.HasSuffix(f, ")") {
				sess.From = strings.Trim(f, "()")
			}
		}
		sessions = append(sessions, sess)
	}
	return sessions
}

func formatUptime(d time.Duration) string {
	if d <= 0 {
		return "unknown"
	}
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	mins := int(d.Minutes()) % 60
	switch {
	case days > 0:
		return fmt.Sprintf("%dd %dh %dm", days, hours, mins)
	case hours > 0:
		return fmt.Sprintf("%dh %dm", hours, mins)
	default:
		return fmt.Sprintf("%dm", mins)
	}
}

func formatBytes(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}
