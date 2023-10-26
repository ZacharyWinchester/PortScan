package port // Initializes the file as a package.

import (
	"net" // Golang portable interface for network input/output.
	"strconv" // Implements conversions to and from string representations of basic data types.
	"time" // Provides functionality for measuring and displaying time.
	"sync"
	log "github.com/sirupsen/logrus" // Adds advanced logging functionality using the logrus package.
)

type ScanResult struct { // Creates a structure to be called later for implimentation into an array. Allows the array to easily store both the port and state of the port in each element
	Port  int
	Protocol string
	State string
}

var (
	mutex sync.Mutex
	wg sync.WaitGroup
	results []ScanResult // Initalizes a zero-filled array in results
)

func ScanPort(protocol, hostname string, port int, wg *sync.WaitGroup) ScanResult { // Function that takes a protocol, hostname, and port. Returns as the ScanResult structure.
	mutex.Lock()
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
	func(conn net.Conn) { // Waits for surrounding functions to return before this function executes. Tries to close the connection that was established previously.
		err := conn.Close() // Stores any produced errors from trying to close the connection.
		if err != nil { // If there is an error, print that an error was encountered.
			log.Error("Connection close error") // If the connection fails to close, log an error.
		}
	}(conn) // Immediately Invoked Function. Right after the above function is declared, we invoke it with conn as an argument.
	result.State = "Open" // Should no errors be encountered, set the state atribute of the element in go to Open.
	mutex.Unlock()
	defer wg.Done()
	results = append(results, result)
	return result // Returns the element of the array.
}

func InitialScan(hostname string) []ScanResult { // Takes an IP address as an argument, and returns an array
	for i := 0; i <= 3; i++ { // As long as i is less than or equal to 1024, run the following and increase i by one.
		wg.Add(1)
		go ScanPort("tcp", hostname, i, &wg)
		wg.Wait()
	}
	return results // Return the results array.
}
