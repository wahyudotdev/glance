package client

import (
	"bufio"
	"bytes"
	"os/exec"
	"strings"
)

type JavaProcess struct {
	PID  string `json:"pid"`
	Name string `json:"name"`
}

func ListJavaProcesses() ([]JavaProcess, error) {
	// Use jps (Java Process Status tool) which is part of the JDK
	cmd := exec.Command("jps", "-l")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	var processes []JavaProcess
	scanner := bufio.NewScanner(&out)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, " ", 2)
		if len(parts) < 2 {
			continue
		}
		pid := parts[0]
		name := parts[1]

		// Ignore jps itself
		if name == "sun.tools.jps.Jps" || name == "jps" {
			continue
		}

		processes = append(processes, JavaProcess{
			PID:  pid,
			Name: name,
		})
	}

	return processes, nil
}
