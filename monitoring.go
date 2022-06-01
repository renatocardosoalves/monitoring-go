package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const monitoramentos = 3
const delay = 5

func showMenu() {
	fmt.Println("1 - start monitoring")
	fmt.Println("2 - show logs")
	fmt.Println("0 - exit")
	fmt.Print("Enter one of the options: ")
}

func readCommand() int {
	var command int
	fmt.Scan(&command)
	fmt.Println("")
	return command
}

func readFile(filename string) []string {
	var sites []string
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("An error has occurred:", err)
	}
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		sites = append(sites, strings.TrimSpace(line))
		if err == io.EOF {
			break
		}
	}
	file.Close()
	return sites
}

func writeLog(site string, status bool) {
	file, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("An error has occurred:", err)
	}
	file.WriteString(time.Now().Format("02/01/2006 15:04:05") + " " + site + " - online: " + strconv.FormatBool(status) + "\n")
	file.Close()
}

func printLogs() {
	logs := readFile("log.txt")

	for _, log := range logs {
		fmt.Println(log)
	}
}

func sendRequest(url string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("An error has occurred:", err)
	}
	url = resp.Request.URL.String()
	statusCode := resp.StatusCode
	fmt.Printf("(%d)\turl: %s\n", statusCode, url)
	if statusCode == 200 {
		writeLog(url, true)
	} else {
		writeLog(url, false)
	}
}

func startsMonitoring() {
	fmt.Println("Monitoring...")
	urls := readFile("sites.txt")
	for i := 0; i < monitoramentos; i++ {
		fmt.Println("Status\tSite")
		for _, url := range urls {
			sendRequest(url)
		}
		time.Sleep(delay * time.Second)
		fmt.Println("")
	}
}

func exitOnSuccess() {
	fmt.Println("Exiting the program...")
	os.Exit(0)
}

func exitOnError() {
	fmt.Println("I don't know this command")
	os.Exit(-1)
}

func executeCommand() {
	switch readCommand() {
	case 1:
		startsMonitoring()
	case 2:
		printLogs()
	case 0:
		exitOnSuccess()
	default:
		exitOnError()
	}
}

func main() {
	for {
		showMenu()
		executeCommand()
	}
}
