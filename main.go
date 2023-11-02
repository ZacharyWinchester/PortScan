package main

import (
	"flag"                                            // Allows for command line arguments
	"fmt"                                             // Formatted input/output
	port "github.com/ZacharyWinchester/PortScan/Port" // Package that contains the functions needed to scan a host.
	"github.com/thediveo/netdb"                       // Package for basic network connectivity.
	"time"                                            // Package used to get how long the program took to run.
)

func main() {
	host := flag.String("host", "127.0.0.1", "What host is being scanned?")
	help := flag.Bool("help", false, "Tells the program to print out all of the flags and their usages.")
	ports := flag.Int("ports", 1024, "How many ports should be scanned?") // Modify later to allow for specific ports and ranges.
	concurrent := flag.Int("con", 100, "How many connections should be attempted at once?")
	flag.Parse()
	if *help == false {
		fmt.Printf("Port Scanning...")                          // Prints out text to alert the user that the program is running
		startTime := time.Now()                                 // Gets the time before the program runs.
		results := port.InitialScan(*host, *ports, *concurrent) // Passes an IP address to the InitialScan function within the port package. Stores in "results"
		endTime := time.Now()                                   // Gets the time after the program runs.
		executionTime := endTime.Sub(startTime)                 // Subtracts the start time from the end time to get the time the program took to run.
		fmt.Printf(" Finished.\nScan finished in %v\n", executionTime)
		fmt.Printf("There are %d closed ports\n", port.ClosedCounter)
		for _, res := range results { // Formats the []ScanResult structure and prints it out.
			if res.State != "Closed" {
				service := netdb.ServiceByPort(res.Port, res.Protocol)
				if service != nil {
					fmt.Printf("%5d/%s %s %s\n", res.Port, res.Protocol, res.State, service.Name)
				} else {
					fmt.Printf("%5d/%s %s\n", res.Port, res.Protocol, res.State)
				}
			}
		}
	} else {
		flag.PrintDefaults()
	}
}
