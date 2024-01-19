package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"github.com/joho/godotenv"
)

// HTTPSDPServer starts a HTTP Server that consumes SDPs
func HTTPSDPServer(port int) chan string {
	sdpChan := make(chan string)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		fmt.Fprintf(w, "done")
		sdpChan <- string(body)
	})

	go func() {
		// nolint: gosec
		err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
		if err != nil {
			panic(err)
		}
	}()

	return sdpChan
}

func main() {
	HTTPSDPServer(8080)
	// another()

}

func simple_server() {
	// Load environment variables
	if err := loadEnv(); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Get the values of HOST and PORT from environment variables
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	address := fmt.Sprintf("%s:%s", host, port)

	// Handle requests with a simple function
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, this is a simple server!")
	})

	// Start the server
	fmt.Printf("Server is running on http://%s\n", address)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}

func loadEnv() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("error loading .env file")
	}
	return nil
}

func another() {
	// Hello world, the web server

	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "Hello, world!\n")
	}

	http.HandleFunc("/hello", helloHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
