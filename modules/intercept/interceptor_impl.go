package intercept

import (
	"apgo/util"
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/elazarl/goproxy"
)

func (i *Interceptor) OnInterception(messageByte []byte, bodyBytes []byte, req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	var newRequest *http.Request
	var newResponse *http.Response

	// Check if interception is enabled
	if !i.isIntercepting {
		return req, newResponse
	}

	log.Printf("[%d] Starting intercept...\n", ctx.Session)

	i.requestPanel.SetText(util.ConvertToReadableString(messageByte, bodyBytes))
	i.targetEntry.SetText(req.URL.Scheme + "://" + req.URL.Host)

	select {
	// press forward
	case <-i.forwardChan:
		requestBody, _ := i.requestBinding.Get()

		lines := strings.Split(requestBody, "\n")
		var bodyBuffer bytes.Buffer

		// Add the headers to the request
		var headerBuffer bytes.Buffer
		for i, line := range lines {
			line = strings.TrimRight(line, "\r\n")

			if line == "" {
				// Start of the body
				bodyBuffer.WriteString(strings.Join(lines[i+1:], "\r\n"))
				break
			}

			// Add headers including the potocal headers, not including the content length
			if !strings.HasPrefix(line, "Content-Length") {
				headerBuffer.WriteString(line + "\r\n")
			}
		}

		// Covertion
		util.ConvertToSendableString(&bodyBuffer)

		headerBuffer.WriteString(fmt.Sprintf("Content-Length: %d\r\n", bodyBuffer.Len()))

		fullRequestBody := headerBuffer.String() + "\r\n" + bodyBuffer.String()
		reader := strings.NewReader(fullRequestBody)
		buf := bufio.NewReader(reader)

		// Only support HTTP/1.x
		r, err := http.ReadRequest(buf)
		if err != nil {
			log.Printf("[%d] Error reading forward request: %v", ctx.Session, err)
			newRequest = req
			log.Printf("[%d] Revert to the original request", ctx.Session)
		} else {
			r.URL.Scheme = req.URL.Scheme
			r.URL.Host = req.URL.Host
			r.RequestURI = ""
			newRequest = r
			log.Printf("[%d] Forward new request", ctx.Session)

			// Use for debugging
			// test, _ := httputil.DumpRequest(newRequest, true)
			// log.Println("Dump intercepted request>>", string(test))
		}

		i.interceptSessionId = ctx.Session
		i.requestPanel.SetText("")
		i.targetEntry.SetText("")

	// press drop
	case <-i.dropChan:
		newRequest = req
		newResponse = goproxy.NewResponse(req,
			goproxy.ContentTypeText, http.StatusForbidden, "Dropped")

		i.requestPanel.SetText("")
		i.targetEntry.SetText("")

	}

	return newRequest, newResponse
}

func (i *Interceptor) OnReceive(messageByte []byte, bodyBytes []byte, r *http.Response, ctx *goproxy.ProxyCtx) {
	// Check if interception is enabled
	if !i.isIntercepting {
		return
	}

	log.Printf("[%d] Received...\n", ctx.Session)
	if i.interceptSessionId == ctx.Session {
		i.responsePanel.SetText(util.ConvertToReadableString(messageByte, bodyBytes))
	} else {
		log.Printf("[%d] Received wrong session intercept event, expected: %d\n", ctx.Session, i.interceptSessionId)
	}
}
