package decoder

import (
	"bytes"
	"io"
	"log"
	"net/http"

	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/flate"
	"github.com/klauspost/compress/gzip"
)

func chunked(te []string) bool { return len(te) > 0 && te[0] == "chunked" }

func DecodeBodyRequest(r *http.Request) ([]byte, error) {
	var bodyBytes []byte
	bodyBytes, _ = io.ReadAll(r.Body)
	// Restore the io.ReadCloser to its original state
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// // Decode bodyBytes
	// var newBodyBytes []byte
	// buf := bytes.NewBuffer(bodyBytes)
	// var reader io.Reader
	// if chunked(r.TransferEncoding) {
	// 	reader = httputil.NewChunkedReader(buf)
	// } else {
	// 	reader = buf
	// }
	// if r.Uncompressed {
	// 	newBodyBytes, _ = io.ReadAll(reader)
	// } else {
	// 	// Check if the response is gzip
	// 	if r.Header.Get("Content-Encoding") == "gzip" {
	// 		greader, err := gzip.NewReader(reader)
	// 		if err != nil {
	// 			log.Println("Error unzip response:", err)
	// 			return nil, err
	// 		} else {
	// 			newBodyBytes, _ = io.ReadAll(greader)
	// 		}
	// 	} else {
	// 		// Default is flate
	// 		freader := flate.NewReader(reader)
	// 		newBodyBytes, _ = io.ReadAll(freader)
	// 	}
	// }

	// log.Println(r.Request.URL, r.Uncompressed, "?", chunked(r.TransferEncoding), "-", string(bodyBytes[0:100]), "->>", string(newBodyBytes[0:100]))

	return bodyBytes, nil
}

func DecodeBodyResponse(r *http.Response) ([]byte, error) {
	var bodyBytes []byte
	// bodyBytes, _ = io.ReadAll(r.Body)
	// // Restore the io.ReadCloser to its original state
	// r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// // Decode bodyBytes
	// var newBodyBytes []byte
	// buf := bytes.NewBuffer(bodyBytes)
	// var reader io.Reader
	// if chunked(r.TransferEncoding) {
	// 	reader = httputil.NewChunkedReader(buf)
	// } else {
	// 	reader = buf
	// }

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	if !r.Uncompressed {
		switch r.Header.Get("Content-Encoding") {
		case "gzip":
			reader, err := gzip.NewReader(io.NopCloser(bytes.NewBuffer(bodyBytes)))
			defer reader.Close()
			if err != nil {
				log.Println("Error unzip response:", err)
				return nil, err
			}

			bodyBytes, err = io.ReadAll(reader)
			if err != nil {
				log.Println("error decoding gzip response", err)
				return nil, err
			}
		case "flate":
			// Check if the response is flate
			reader := flate.NewReader(io.NopCloser(bytes.NewBuffer(bodyBytes)))
			defer reader.Close()
			var err error
			bodyBytes, err = io.ReadAll(reader)
			if err != nil {
				log.Println("error decoding flate response", err)
				return nil, err
			}
		case "br":
			bodyBytes, err = io.ReadAll(brotli.NewReader(io.NopCloser(bytes.NewBuffer(bodyBytes))))
			if err != nil {
				log.Println("error decoding gzip response", err)
				return nil, err
			}
		}
	}

	// log.Println(r.Request.URL, r.Uncompressed, "?", chunked(r.TransferEncoding), "-", string(bodyBytes[0:100]), "->>", string(newBodyBytes[0:100]))

	return bodyBytes, nil
}
