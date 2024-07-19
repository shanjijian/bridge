package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
)

var clientConn net.Conn

func main() {
	go startTCPServer()
	http.HandleFunc("/", handleHTTPRequest)
	fmt.Println("HTTP Server is listening on port 6602...")
	err := http.ListenAndServe(":6602", nil)
	if err != nil {
		fmt.Printf("Failed to start HTTP server: %v\n", err)
	}
}

func startTCPServer() {
	listener, err := net.Listen("tcp", ":6601")
	if err != nil {
		fmt.Printf("Failed to start TCP server: %v\n", err)
		return
	}
	defer listener.Close()
	fmt.Println("TCP Server is listening on port 6601...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Failed to accept connection: %v\n", err)
			continue
		}
		clientConn = conn
		fmt.Println("Client connected:", conn.RemoteAddr().String())
	}
}

func handleHTTPRequest(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Query().Get("file")
	cmd := r.URL.Query().Get("cmd")
	if filePath == "" && cmd == "" {
		http.Error(w, "file or cmd parameter is required", http.StatusBadRequest)
		return
	}

	if clientConn == nil {
		http.Error(w, "no client connected", http.StatusInternalServerError)
		return
	}

	request := ""
	if filePath != "" {
		request = "FILE:" + filePath
	} else if cmd != "" {
		request = "CMD:" + cmd
	}

	_, err := clientConn.Write([]byte(request + "\n"))
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to send request to client: %v", err), http.StatusInternalServerError)
		return
	}

	reader := bufio.NewReader(clientConn)
	resp, err := reader.ReadString('\n')
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to read response from client: %v", err), http.StatusInternalServerError)
		return
	}
	resp = strings.TrimSpace(resp)

	if resp == "ERROR" {
		http.Error(w, "failed to process request from client", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, err = io.WriteString(w, resp)
	if err != nil {
		fmt.Printf("Failed to write response: %v\n", err)
	}
}
