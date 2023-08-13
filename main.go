package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	portPtr := flag.String("p", "30000", "Listening port")
	logPtr := flag.String("log", "null", "Log output directory")

	flag.Parse()

	serverPort := *portPtr
	logFilePath := *logPtr

	if logFilePath != "null" {
		logfile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			fmt.Printf("Error opening file: %v\n", err)
			log.Fatal(err)
		}
		defer logfile.Close()
		log.SetOutput(logfile)
	}

	http.HandleFunc("/", getIPAddress)
	fmt.Printf("INFO: Server listening on 0.0.0.0 port %s.\n", serverPort)
	http.ListenAndServe("0.0.0.0:"+serverPort, nil)
}

func getIPAddress(w http.ResponseWriter, r *http.Request) {
	IPAddress := strings.Split(r.RemoteAddr, ":")[0]
	if r.Header.Get("X-Forwarded-For") != "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	} else if r.Header.Get("X-Real-IP") != "" {
		IPAddress = r.Header.Get("X-Real-IP")
	}

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		logRequest(r, IPAddress, http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "%s\n", IPAddress)
	logRequest(r, IPAddress, http.StatusOK)
}

func logRequest(r *http.Request, IPAddress string, status int) {
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")
	referer := r.Referer()
	if referer == "" {
		referer = "-"
	}

	logEntry := fmt.Sprintf("%s - %s \"%s %s %s\" %d \"%s\" \"%s\"",
		formattedTime, IPAddress, r.Method, r.URL, r.Proto, status, referer, r.UserAgent())

	log.SetFlags(0)
	log.Println(logEntry)
}
