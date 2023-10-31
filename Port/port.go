package port // Initializes the file as a package.

import (
	"net" // Golang portable interface for network input/output.
	"strconv" // Implements conversions to and from string representations of basic data types.
	"time" // Provides functionality for measuring and displaying time.
	"sort" // Provides functionality for sorting the resulting port scan
	"sync" // Provides locking and unlocking
	log "github.com/sirupsen/logrus" // Adds advanced logging functionality using the logrus package.
)

type ScanResult struct { // Creates a structure to be called later for implimentation into an array. Allows the array to easily store both the port and state of the port in each element
	Port int
	Protocol string
	State string
}

func worker(id int, jobs <-chan int, resultC chan<- ScanResult, hostname string) {
	for i := range jobs {
		func(port int) {
			result := ScanPort("tcp", hostname, port)
			if result.State == "Open" {
				mutex.Lock()
				resultC <- ScanResult{result}
				mutex.Unlock()
			} else {
				mutex.Lock()
				ClosedCounter++
				mutex.Unlock()
			}
		}(i)
	}
}

var (
	mutex sync.Mutex
	ClosedCounter int
)

func ScanPort(protocol, hostname string, port int) ScanResult { // Function that takes a protocol, hostname, and port. Returns as the ScanResult structure.
	result := ScanResult{Port: port} // Sets the Port element to the port number taken in by this function.
	result.Protocol = protocol // Sets the Protocol element to the protocol type taken in by this function.
	address := hostname + ":" + strconv.Itoa(port) // Takes the hostname (represented as an ip), concatinates a ':' (signifies a socket) to it, and then turns the port number into a string so that it can be concatinated to the rest of the address. Stores in address.
	conn, err := net.DialTimeout(protocol, address, 1*time.Second) /* DialTimeout is a function in the net package. It takes a protocol, address, and a timeout duration. It attempts to connect to the address
 									using the given protocol. Returns conn (which reads data from the connection), and error. If there is no error, then it returns nil for error.
	  								conn and error are then assigned to conn and err in this function to be used later.*/

	if err != nil { // If there is an error, run this:
		result.State = "Closed" // Sets the state in this port's element to "Closed".
		return result // Returns the element.
	}
	defer func(conn net.Conn) { // Waits for surrounding functions to return before this function executes. Tries to close the connection that was established previously.
		err := conn.Close() // Stores any produced errors from trying to close the connection.
		if err != nil { // If there is an error, print that an error was encountered.
			log.Error("Connection close error") // If the connection fails to close, log an error.
		}
	}(conn) // Immediately Invoked Function. Right after the above function is declared, we invoke it with conn as an argument.
	result.State = "Open" // Should no errors be encountered, set the state atribute of the element in go to Open.
	return result // Returns the element of the array.
}

func InitialScan(hostname string) []ScanResult { // Takes an IP address as an argument, and returns an array
	const totalTask = 60000
	jobs := make(chan int, totalTask) // Creates a jobs channel with a buffer size of numJobs
	resultsC := make(chan ScanResult, numJobs) // Creates a resultsC channel with a buffer size of numJobs
	for i := 1; i <= 100; i++ {
		go worker(i, jobs, resultsC, hostname)
	}
	for i := 1; i <= totalTask; i++ { // As long as i is less than or equal to totalTask, send i to jobs channel.
		jobs <- i
	}
	close(jobs)
	for i := 1; i <= totalTask; i++ { // Recieves the results of the workers
		<-resultsC
	}
	close(results)
	sort.SliceStable(results, func(i, j int) bool {
		return results[i].Port < results[j].Port
	})
	return results // Return the results array.
}
