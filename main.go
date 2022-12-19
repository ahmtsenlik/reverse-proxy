package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

func Proxy(targetHost string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(targetHost)
	if err != nil {
		panic(err)
	}
	return httputil.NewSingleHostReverseProxy(url), nil
}

func ProxyRequestHandler(proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	}
}

func main() {
	fmt.Println("----------------------------------------")
	fmt.Println("Example Server Url: https://10.255.100.80")
	fmt.Println("Example Listening Port: 8088")
	fmt.Println("----------------------------------------")
	fmt.Print("Enter Server Url: ")
	reader := bufio.NewReader(os.Stdin)
	serverUrl, err := reader.ReadString('\n')
	serverUrl = strings.Replace(serverUrl, "\r\n", "", -1)
	if err != nil {
		panic(err)
	} else {
		proxy, err := Proxy(serverUrl)
		if err != nil {
			panic(err)
		}
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		http.HandleFunc("/", ProxyRequestHandler(proxy))
	}

	fmt.Print("Enter Listening Port: ")
	listeningPort, err := reader.ReadString('\n')
	listeningPort = strings.Replace(listeningPort, "\r\n", "", -1)
	if err != nil {
		panic(err)
	} else {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		err = cmd.Run()
		if err != nil {
			panic(err)

		}
		fmt.Println("Server Url: " + serverUrl)
		fmt.Println("Listening Port: " + listeningPort)
		fmt.Println("Now Listening...")
		log.Fatal(http.ListenAndServe(":"+listeningPort, nil))
	}
}
