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

const MONITORING_TIMES = 5
const MONITORING_DELAY = 30

func check(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func main() {
	showIntroduction()
	for {
		showMenuOptions()
		command := getCommand()

		selectOption(command)
	}
}

func showIntroduction() {
	name := "Pedro"
	version := 0.1

	fmt.Println("Olá,", name)
	fmt.Println("Este programa está na versão ", version)
}

func getCommand() int {
	var inputedCommand int
	var _, err = fmt.Scan(&inputedCommand)
	check(err)

	return inputedCommand
}

func showMenuOptions() {
	fmt.Println("01 - Iniciar monitoramento dos sites")
	fmt.Println("02 - Exibir logs")
	fmt.Println("03 - Encerrar programa")
}

func selectOption(command int) {
	switch command {
	case 1:
		initMonitoring()
	case 2:
		fmt.Println("Logs:")
		printLogs()
	case 3:
		fmt.Println("Saindo...")
		os.Exit(0)
	default:
		fmt.Println("Comando inexistente. Reinicie e tente novamente!")
		os.Exit(0)
	}
}

func initMonitoring() {
	fmt.Println("Monitorando...")
	sites := extractWebsitesFromFile()

	for i := 0; i < MONITORING_TIMES; i++ {
		for i, site := range sites {
			testSiteStatus(site, i)
		}
		time.Sleep(MONITORING_DELAY * time.Minute)
	}

}

func testSiteStatus(site string, index int) {
	resp, err := http.Get(site)
	check(err)

	if resp.StatusCode == 200 {
		fmt.Printf("Site %02d: %v foi carregado com sucesso!\n", index+1, site)
		fmt.Println(strings.Repeat("-", 50))
		recordLog(site, true)
	} else {
		fmt.Printf("Site %02d: %v está com problemas! status_code: %d \n", index+1, site, resp.StatusCode)
		fmt.Println(strings.Repeat("-", 50))
		recordLog(site, false)
	}
}

func extractWebsitesFromFile() []string {
	var sites []string
	file, err := os.Open("sites.txt")
	check(err)

	reader := bufio.NewReader(file)

	for {
		readString, err := reader.ReadString('\n')
		readString = strings.TrimSpace(readString)

		sites = append(sites, readString)
		if err == io.EOF {
			break
		}
	}

	errOnCloseFile := file.Close()
	check(errOnCloseFile)

	return sites
}

func recordLog(site string, status bool) {
	file, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	check(err)

	_, err = file.WriteString(time.Now().Format("02/01/2006 15:04:05") + " - " + site + " - online: " + strconv.FormatBool(status) + "\n")
	check(err)
}

func printLogs() {
	file, err := os.ReadFile("log.txt")
	check(err)

	fmt.Println(string(file))
	return
}
