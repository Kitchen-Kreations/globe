package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/akamensky/argparse"
)

func scanner(ip string, portsToScanChan chan int, results chan int) {
	for p := range portsToScanChan {
		address := fmt.Sprintf("%s:%d", ip, p)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			results <- 0
			continue
		}
		conn.Close()
		results <- p
	}
}

func main() {
	startTime := time.Now()

	// init parser
	parser := argparse.NewParser("globe", "port scanner")

	// arguments
	var ip *string = parser.String("i", "ip", &argparse.Options{Required: true, Help: "What address to scan"})
	var ports *string = parser.String("p", "ports", &argparse.Options{Required: false, Help: "What ports to scan"})

	err := parser.Parse(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	var portsToScan []int

	if *ports == "" {
		var defaultPorts []int
		for i := 0; i <= 1000; i++ {
			defaultPorts = append(defaultPorts, i)
		}

		portsToScan = defaultPorts
	} else if *ports == "a" || *ports == "all" {
		var allPorts []int
		for i := 0; i <= 65535; i++ {
			allPorts = append(allPorts, i)
		}

		portsToScan = allPorts
	} else {
		// parse through ports

		if strings.Contains(*ports, ",") {
			portsToScanString := strings.Split(*ports, ",")

			for _, i := range portsToScanString {
				port, err := strconv.Atoi(i)
				if err != nil {
					log.Fatal(err)
				}

				portsToScan = append(portsToScan, port)
			}
		} else if strings.Contains(*ports, "-") {
			portRangeString := strings.Split(*ports, "-")
			var portRange []int

			for _, i := range portRangeString {
				port, err := strconv.Atoi(i)
				if err != nil {
					log.Fatal(err)
				}

				portRange = append(portRange, port)
			}

			for i := portRange[0]; i <= portRange[1]; i++ {
				portsToScan = append(portsToScan, i)
			}
		} else {
			port, err := strconv.Atoi(*ports)
			if err != nil {
				log.Fatal(err)
			}

			portsToScan = append(portsToScan, port)
		}
	}

	target := *ip

	fmt.Println("Target: ", target)
	fmt.Println("Start Time: ", startTime.String())

	var wg = &sync.WaitGroup{}

	portsToScanChan := make(chan int, 500)
	results := make(chan int)
	var openPorts []int

	for i := 0; i <= cap(portsToScanChan); i++ {
		go scanner(target, portsToScanChan, results)
	}

	go func() {
		for _, element := range portsToScan {
			portsToScanChan <- element
		}
	}()

	for i := range portsToScan {
		wg.Add(1)
		i++
		go func(port int) {
			if port != 0 {
				openPorts = append(openPorts, port)
			}
		}(<-results)
	}

	close(portsToScanChan)
	close(results)

	fmt.Printf("\n PORT \t STATE \n======\t=======\n")
	sort.Ints(openPorts)
	for _, port := range openPorts {
		fmt.Printf(" %d \t open \n", port)
	}
	fmt.Print("\n")
	endTime := time.Now()
	fmt.Println("End Time: " + endTime.String())
	timeDiff := endTime.Sub(startTime)
	fmt.Println("Globe took " + timeDiff.String() + " to run")
}
