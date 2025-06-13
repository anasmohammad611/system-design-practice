package main

import (
	"fmt"
	"log"
	"net/http"
)

func startBackendServer(port int) {
	// Start the backend server
	// This function creates a simple HTTP server that listens on the specified port.
	// It responds with a message indicating which port it is running on.
	// The server is started in a separate goroutine so that multiple servers can run concurrently.
	// The server will log a message when it starts and when it receives a request.
	// The server will also log an error if it fails to start.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		msg := "Hello from backend server on port " + fmt.Sprint(port)
		fmt.Println(msg)
		w.WriteHeader(http.StatusOK)
	})

	// Create a new HTTP server with the specified port and handler
	// The server will listen for incoming HTTP requests on the specified port.
	// The handler will respond to requests with a message indicating the port.
	server := &http.Server{
		Addr: fmt.Sprintf(":%d", port),
		Handler: handler,
	}

	// Start the server in a new goroutine
	// This allows the server to run concurrently with other servers.
	// The server will log a message when it starts and will handle requests
	// in the handler function defined above.
	go func() {
		log.Println("Starting backend server on port", port)
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("Backend server on port %d failed: %v", port, err)
		}
	}()
}

func main() {
	for i := 8000; i <= 8002; i++ {
		startBackendServer(i)
	}
	log.Println("Backend servers started on ports 8000, 8001, and 8002")
	// Keep the main goroutine running
	// This is necessary to prevent the program from exiting immediately
	// since the backend servers are running in goroutines.
	// In a real application, you might use a more sophisticated mechanism
	// to keep the application running, such as a signal handler or a wait group.
	// Here, we simply block forever.		
	log.Println("Load balancer is starting...")
	startLoadBalancer(9000)
	log.Println("Load balancer started on port 9000")
	select {} // Block forever to keep the main goroutine running
}

var (
	backendServers = []string{"localhost:8000", "localhost:8001", "localhost:8002"}
	currentIdx     = -1 // Index for round-robin selection
	currentServer = 0
	weights = []int{2, 1, 2} // Example weights for each backend server
	totalWeight = 5 // Total weight of all backend servers
)
/*
// getNextBackendServer returns the next backend server in a round-robin fashion.
func getNextBackendServer() string {
	// This function returns the next backend server in a round-robin fashion.
	// It uses a global slice of backend server addresses and an index to keep track of the current server.
	server := backendServers[currentIdx]
	currentIdx = (currentIdx + 1) % len(backendServers)
	return server
}
*/

// getNextBackendServer returns the next backend server in a weighted round-robin fashion.
func getNextBackendServerWeighted() string {
	currentIdx = (currentIdx + 1) % totalWeight
	cnt := 0
	for i := range weights {
		cnt += weights[i]
		if currentIdx < cnt {
			return backendServers[i]
		}
	}
	return backendServers[0] // Fallback to the first server if something goes wrong
}

func startLoadBalancer(port int) {
	// This function starts a simple load balancer that listens on the specified port.
	// The load balancer will forward requests to backend servers in a round-robin manner.
	// It uses the getNextBackendServer function to determine which backend server to forward the request to.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		// target := getNextBackendServer()
		target := getNextBackendServerWeighted()
		resp, err := http.Get("http://" + target)
		if err != nil {
			http.Error(w, "Failed to reach backend server "+err.Error(), http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		log.Printf("Forwarding request to %s\n", target)
		// Write the response from the backend server to the client
		// The response status code and body are written to the client.
		// The status code is set to the same as the backend server's response.
		w.WriteHeader(resp.StatusCode)
		_, err = fmt.Fprint(w, target+" responded\n")
		if err != nil {
			log.Println("Error writing response to client:", err)
		}
	})

	// Create a new HTTP server with the specified port and handler
	// The server will listen for incoming HTTP requests on the specified port.
	// The handler will forward requests to the backend servers using the round-robin method.
	server := &http.Server{
		Addr:   fmt.Sprintf(":%d", port),
		Handler: handler,
	}

	log.Println("Starting load balancer on port", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Load balancer failed: %v", err)
	}
	log.Println("Load balancer started on port", port)
	log.Println("Load balancer is running, waiting for requests...")
}