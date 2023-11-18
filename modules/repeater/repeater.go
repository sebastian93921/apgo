package repeater

import (
	"apgo/util"
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func SendRawRequest(target string, rawRequest string) (*http.Response, error) {
	// Split the raw request into lines
	lines := strings.Split(rawRequest, "\n")

	// Split the request line into method, URL, and protocol
	requestLineParts := strings.Split(lines[0], " ")
	method := requestLineParts[0]
	url := requestLineParts[1]

	// Create a buffer to hold the request body
	var bodyBuffer bytes.Buffer

	// Add the headers to the request
	headers := make(http.Header)
	for i, line := range lines[1:] {
		line = strings.TrimRight(line, "\r\n")

		if line == "" {
			// Start of the body
			bodyBuffer.WriteString(strings.Join(lines[i+2:], "\n"))
			break
		}
		parts := strings.SplitN(line, ": ", 2)
		if len(parts) == 2 {
			if !strings.HasPrefix(parts[0], "Content-Length") {
				headers.Set(parts[0], parts[1])
			}
		}
	}

	// Covertion
	util.ConvertToSendableString(&bodyBuffer)
	headers.Set("Content-Length", fmt.Sprintf("Content-Length: %d", bodyBuffer.Len()))

	// Create a new request
	targetUrl := target
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		targetUrl += url
	}
	log.Println("Sending request to", targetUrl, "\n headers: ", headers)
	req, err := http.NewRequest(method, targetUrl, &bodyBuffer)
	if err != nil {
		return nil, err
	}
	req.Header = headers

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
