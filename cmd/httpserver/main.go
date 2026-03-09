package main

import (
	"crypto/sha256"
	"fmt"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"httpfromtcp/internal/server"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

const port = 42069
const chunkedBufferSize = 1024

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
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/stream") {
		streamHandler(w, req)
		return
	}

	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/html") {
		htmlHandler(w, req)
		return
	}

	if req.RequestLine.RequestTarget == "/yourproblem" {
		mainHandler400(w, req)
		return
	}

	if req.RequestLine.RequestTarget == "/myproblem" {
		mainHandler500(w, req)
		return
	}

	if req.RequestLine.RequestTarget == "/video" {
		videoHandler(w, req)
		return
	}

	mainHandler200(w, req)
}

func mainHandler400(w *response.Writer, _ *request.Request) {
	body := []byte("<html>\n<head>\n<title>400 Bad Request</title>\n</head>\n<body>\n" +
		"<h1>Bad Request</h1>\n<p>Your request honestly kinda sucked.</p>\n" +
		"</body>\n</html>")

	headers := response.GetDefaultHeaders(len(body))
	headers.Replace("Content-Type", "text/html")

	w.WriteStatusLine(response.StatusBadRequest)
	w.WriteHeaders(headers)
	w.WriteBody(body)
}

func mainHandler500(w *response.Writer, _ *request.Request) {
	body := []byte("<html>\n<head>\n<title>500 Internal Server Error</title>\n</head>\n<body>\n" +
		"<h1>Internal Server Error</h1>\n<p>Okay, you know what? This one is on me.</p>\n" +
		"</body>\n</html>")

	headers := response.GetDefaultHeaders(len(body))
	headers.Replace("Content-Type", "text/html")

	w.WriteStatusLine(response.StatusInternalServerError)
	w.WriteHeaders(headers)
	w.WriteBody(body)
}

func mainHandler200(w *response.Writer, _ *request.Request) {
	body := []byte("<html>\n<head>\n<title>200 OK</title>\n</head>\n<body>\n" +
		"<h1>Success!</h1>\n<p>Your request was an absolute banger.</p>\n" +
		"</body>\n</html>")

	headers := response.GetDefaultHeaders(len(body))
	headers.Replace("Content-Type", "text/html")

	w.WriteStatusLine(response.StatusOK)
	w.WriteHeaders(headers)
	w.WriteBody(body)
}

func streamHandler(w *response.Writer, req *request.Request) {
	// Trim the request target prefix to get the route parameter
	param := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")

	// Put together the response headers
	headers := response.GetDefaultHeaders(0)
	headers.Remove("Content-Length")
	headers.Set("Transfer-Encoding", "chunked")

	w.WriteStatusLine(response.StatusOK)
	w.WriteHeaders(headers)

	// Call httpbin.org api with route parameter
	res, err := http.Get("https://httpbin.org/" + param)
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()

	buf := make([]byte, chunkedBufferSize)
	for {
		n, err := res.Body.Read(buf)
		fmt.Printf("httpbin.org response body bytes read: %d\n", n)

		if err != nil && err != io.EOF {
			log.Println(err)
			return
		}

		if n == 0 && err == io.EOF {
			w.WriteChunkedBodyDone()
			return
		}

		w.WriteChunkedBody(buf[:n])
	}
}

func htmlHandler(w *response.Writer, req *request.Request) {
	param := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")

	// Put together the response headers
	headers := response.GetDefaultHeaders(0)
	headers.Remove("Content-Length")
	headers.Set("Transfer-Encoding", "chunked")
	headers.Set("Trailer", "X-Content-SHA256")
	headers.Set("Trailer", "X-Content-Length")

	w.WriteStatusLine(response.StatusOK)
	w.WriteHeaders(headers)

	// Call httpbin.org api with route parameter
	res, err := http.Get("https://httpbin.org/" + param)
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()

	buf := make([]byte, chunkedBufferSize)
	resBody := make([]byte, 0)
	for {
		n, err := res.Body.Read(buf)

		if err != nil && err != io.EOF {
			log.Println(err)
			return
		}

		if n == 0 && err == io.EOF {
			hash := fmt.Sprintf("%x", sha256.Sum256(resBody))
			length := strconv.Itoa(len(resBody))
			headers.Set("X-Content-SHA256", hash)
			headers.Set("X-Content-Length", length)

			w.WriteTrailers(headers)
			return
		}

		w.WriteChunkedBody(buf[:n])
		resBody = append(resBody, buf[:n]...)
	}
}

func videoHandler(w *response.Writer, _ *request.Request) {
	// Read the video file into the response body
	body, err := os.ReadFile("assets/vim.mp4")
	if err != nil {
		log.Println(err)
	}

	// Put together the response headers
	headers := response.GetDefaultHeaders(len(body))
	headers.Replace("Content-Type", "video/mp4")

	w.WriteStatusLine(response.StatusOK)
	w.WriteHeaders(headers)
	w.WriteBody(body)
}
