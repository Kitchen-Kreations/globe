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
	"github.com/fatih/color"
)

var (
	banner string = ` 
	 ██████╗ ██╗      ██████╗ ██████╗ ███████╗
	██╔════╝ ██║     ██╔═══██╗██╔══██╗██╔════╝
	██║  ███╗██║     ██║   ██║██████╔╝█████╗  
	██║   ██║██║     ██║   ██║██╔══██╗██╔══╝  
	╚██████╔╝███████╗╚██████╔╝██████╔╝███████╗
	 ╚═════╝ ╚══════╝ ╚═════╝ ╚═════╝ ╚══════╝
	`
	startTime = time.Now()
)

func main() {
	parser := argparse.NewParser("globe", "Port Scanner")

	var port *string = parser.String("p", "port", &argparse.Options{Required: false, Help: "Ports to scan"})
	var target *string = parser.String("t", "target", &argparse.Options{Required: true, Help: "IP/Domain to target"})

	err := parser.Parse(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	ports_to_scan := get_ports(*port)
	target_ip := conv_domain_to_ip(*target).String()
	print_start_banner(*target, ports_to_scan, target_ip)

	scan(target_ip, ports_to_scan)
}

func scan(ip string, ports []int) {
	var wg = &sync.WaitGroup{}

	ports_to_scan_channel := make(chan int, len(ports))
	results_channel := make(chan int)
	var open_ports []int

	for i := 0; i <= cap(ports_to_scan_channel); i++ {
		go scanner(ip, ports_to_scan_channel, results_channel)
	}

	go func() {
		for _, element := range ports {
			ports_to_scan_channel <- element
		}
	}()

	for i := range ports {
		wg.Add(1)
		i++
		go func(port int) {
			if port != 0 {
				open_ports = append(open_ports, port)
			}
		}(<-results_channel)
	}

	close(ports_to_scan_channel)
	close(results_channel)

	print_out_banner(open_ports, ip, ports)
}

func print_out_banner(open_ports []int, ip string, ports []int) {
	fmt.Printf("\n PORT \t STATE \n======\t=======\n")
	sort.Ints(open_ports)
	for _, port := range open_ports {
		fmt.Printf(" %d \t open \n", port)
	}

	endTime := time.Now()
	timeDiff := endTime.Sub(startTime)
	fmt.Println("\nFound " + fmt.Sprint(len(open_ports)) + " ports OPEN out of " + fmt.Sprint(len(ports)) + " ports scanned in " + timeDiff.String())
}

func conv_domain_to_ip(target string) net.IP {
	ips, _ := net.LookupIP(target)
	return ips[0]
}

func get_ports(port string) []int {
	if port == "" {
		return makeRange(1, 1000)
	} else if strings.Contains(port, "-") {
		min_max := strings.Split(port, "-")
		min, err := strconv.Atoi(min_max[0])
		if err != nil {
			log.Fatal(err)
		}
		max, err := strconv.Atoi(min_max[1])
		if err != nil {
			log.Fatal(err)
		}

		if max < min {
			log.Fatal("Lower bound must be less than upper bound")
		}

		return makeRange(min, max)
	} else if strings.Contains(port, ",") {
		split_port := strings.Split(port, ",")
		var ports []int
		for _, port := range split_port {
			port, err := strconv.Atoi(port)
			if err != nil {
				log.Fatal(err)
			}
			ports = append(ports, port)
		}
		return ports
	} else if strings.Contains(port, "all") {
		return makeRange(1, 65535)
	}

	return makeRange(1, 1000)
}

func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}

func print_start_banner(target string, port []int, target_ip string) {
	blue := color.New(color.FgCyan)

	blue.Println(banner)
	fmt.Println("Created By: BlessedToastr")
	fmt.Println("github.com/Kitchen-Kreations/globe")
	fmt.Println("")

	fmt.Println("Scanning " + fmt.Sprint(len(port)) + " ports on " + target + "(" + target_ip + ")")
}

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
