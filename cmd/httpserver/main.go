package main

import (
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"httpfromtcp/internal/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const port = 42069

func main() {
	server, err := server.Serve(port, mainHandler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func mainHandler(w *response.Writer, req *request.Request) {
	if req.RequestLine.RequestTarget == "/yourproblem" {
		body := []byte("<html>\n<head>\n<title>400 Bad Request</title>\n</head>\n<body>\n" +
			"<h1>Bad Request</h1>\n<p>Your request honestly kinda sucked.</p>\n" +
			"</body>\n</html>")

		err := w.WriteStatusLine(response.StatusBadRequest)
		if err != nil {
			log.Printf("Error writing status line: %v", err)
		}
		err = w.WriteHeaders(response.GetDefaultHeaders(len(body)))
		if err != nil {
			log.Printf("Error writing header: %v", err)
		}
		_, err = w.WriteBody(body)
		if err != nil {
			log.Printf("Error writing body: %v", err)
		}
	}

	if req.RequestLine.RequestTarget == "/myproblem" {
		body := []byte("<html>\n<head>\n<title>500 Internal Server Error</title>\n</head>\n<body>\n" +
			"<h1>Internal Server Error</h1>\n<p>Okay, you know what? This one is on me.</p>\n" +
			"</body>\n</html>")

		err := w.WriteStatusLine(response.StatusInternalServerError)
		if err != nil {
			log.Printf("Error writing status line: %v", err)
		}
		err = w.WriteHeaders(response.GetDefaultHeaders(len(body)))
		if err != nil {
			log.Printf("Error writing header: %v", err)
		}
		_, err = w.WriteBody(body)
		if err != nil {
			log.Printf("Error writing body: %v", err)
		}
	}

	body := []byte("<html>\n<head>\n<title>200 OK</title>\n</head>\n<body>\n" +
		"<h1>Success!</h1>\n<p>Your request was an absolute banger.</p>\n" +
		"</body>\n</html>")

	err := w.WriteStatusLine(response.StatusOK)
	if err != nil {
		log.Printf("Error writing status line: %v", err)
	}
	err = w.WriteHeaders(response.GetDefaultHeaders(len(body)))
	if err != nil {
		log.Printf("Error writing header: %v", err)
	}
	_, err = w.WriteBody(body)
	if err != nil {
		log.Printf("Error writing body: %v", err)
	}
}
