package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

// Proxy function creates a reverse proxy for the given target host.
func Proxy(targetHost string) (*httputil.ReverseProxy, error) {
	// Parse the target host URL.
	url, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}
	return httputil.NewSingleHostReverseProxy(url), nil
}

// ProxyRequestHandler returns a function that handles HTTP requests using the provided proxy.
func ProxyRequestHandler(proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Serve the HTTP request using the proxy.
		proxy.ServeHTTP(w, r)
	}
}

func main() {
	fmt.Println("----------------------------------------")
	fmt.Println("Example Server Url: https://10.255.100.80:1234")
	fmt.Println("Example Listening Port: 8088")
	fmt.Println("----------------------------------------")

	// Reader for user input.
	reader := bufio.NewReader(os.Stdin)

	// Get the server URL from the user.
	fmt.Print("Enter Server Url: ")
	serverUrl, err := reader.ReadString('\n')
	serverUrl = strings.Replace(serverUrl, "\r\n", "", -1)
	if err != nil {
		fmt.Println("Error reading server URL:", err)
		fmt.Println("Press Enter to exit.")
		fmt.Scanln()
		os.Exit(1)
	}

	// Create the reverse proxy.
	proxy, err := Proxy(serverUrl)
	if err != nil {
		fmt.Println("Error creating proxy:", err)
		fmt.Println("Press Enter to exit.")
		fmt.Scanln()
		os.Exit(1)
	}

	// Get the listening port from the user.
	fmt.Print("Enter Listening Port: ")
	listeningPort, err := reader.ReadString('\n')
	listeningPort = strings.Replace(listeningPort, "\r\n", "", -1)
	if err != nil {
		fmt.Println("Error reading listening port:", err)
		fmt.Println("Press Enter to exit.")
		fmt.Scanln()
		os.Exit(1)
	}

	// Clear the console screen.
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error clearing screen:", err)
		fmt.Println("Press Enter to exit.")
		fmt.Scanln()
		os.Exit(1)
	}

	// Display information about the server and port.
	fmt.Println("Server Url: " + serverUrl)
	fmt.Println("Listening Port: " + listeningPort)
	fmt.Println("Now Listening...")

	// Start the HTTP server.
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	http.HandleFunc("/", ProxyRequestHandler(proxy))
	if err := http.ListenAndServe(":"+listeningPort, nil); err != nil {
		fmt.Println("Error starting server:", err)
		fmt.Println("Press Enter to exit.")
		fmt.Scanln()
		os.Exit(1)
	}
}
