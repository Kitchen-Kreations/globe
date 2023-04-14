package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

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

	fmt.Println("Scanning " + fmt.Sprint(len(port)) + " ports on " + target + "("+target_ip+")")
}

func scanner_worker() {
	
}