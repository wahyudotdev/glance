package client

import (
	"bufio"
	"bytes"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
)

// JavaProcess represents a running Java application.
type JavaProcess struct {
	PID  string `json:"pid"`
	Name string `json:"name"`
}

// ListJavaProcesses returns a list of active Java processes.
func ListJavaProcesses() ([]JavaProcess, error) {
	processes := make(map[string]JavaProcess)

	// 1. Try jps first (best for accurate names)
	cmd := exec.Command("jps", "-l")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err == nil {
		scanner := bufio.NewScanner(&out)
		for scanner.Scan() {
			line := scanner.Text()
			parts := strings.SplitN(line, " ", 2)
			if len(parts) < 2 || parts[0] == "" {
				continue
			}
			pid := parts[0]
			name := parts[1]

			if name == "sun.tools.jps.Jps" || name == "jps" {
				continue
			}

			// If it's a path to a jar, just take the filename
			if strings.HasSuffix(name, ".jar") {
				name = filepath.Base(name)
			}

			processes[pid] = JavaProcess{PID: pid, Name: name}
		}
		log.Printf("Found %d processes using jps", len(processes))
	} else {
		log.Printf("jps failed: %v (moving to ps fallback)", err)
	}

	// 2. Fallback to ps aux on Unix systems to find processes jps missed (e.g. different users)
	psCmd := exec.Command("ps", "aux")
	var psOut bytes.Buffer
	psCmd.Stdout = &psOut
	if err := psCmd.Run(); err == nil {
		scanner := bufio.NewScanner(&psOut)
		foundViaPs := 0
		for scanner.Scan() {
			line := scanner.Text()
			if !strings.Contains(line, "java") || strings.Contains(line, "grep") {
				continue
			}

			// Format: USER PID %CPU %MEM VSZ RSS TTY STAT START TIME COMMAND
			fields := strings.Fields(line)
			if len(fields) < 11 {
				continue
			}

			pid := fields[1]
			// If we already found it with jps, don't overwrite
			if _, ok := processes[pid]; ok {
				continue
			}

			// Extract command line
			commandLine := strings.Join(fields[10:], " ")
			name := "Java Process"

			// Look for specific apps or JARs in the command line
			if strings.Contains(commandLine, "org.elasticsearch.bootstrap.Elasticsearch") {
				name = "Elasticsearch"
			} else if strings.Contains(commandLine, "org.apache.catalina.startup.Bootstrap") {
				name = "Tomcat"
			} else if idx := strings.Index(commandLine, "-jar "); idx != -1 {
				// Try to get the jar name after -jar flag
				parts := strings.Fields(commandLine[idx+5:])
				if len(parts) > 0 {
					name = filepath.Base(parts[0])
				}
			} else {
				// Fallback: use the first word of the command (usually 'java' or path to it)
				name = filepath.Base(fields[10])
			}

			processes[pid] = JavaProcess{PID: pid, Name: name}
			foundViaPs++
		}
		log.Printf("Found %d additional processes using ps", foundViaPs)
	} else {
		log.Printf("ps aux failed: %v", err)
	}

	result := make([]JavaProcess, 0, len(processes))
	for _, p := range processes {
		result = append(result, p)
	}
	return result, nil
}
