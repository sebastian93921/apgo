package util

import (
	"bytes"
	"encoding/json"
	"log"
	"unicode/utf8"
)

func ConvertToReadableString(message []byte, body []byte) string {
	replyBody := "---Cannot convert it into string---"

	// Check if the request body is valid UTF-8
	if utf8.Valid(body) {
		if json.Valid(body) {
			var prettyJSON bytes.Buffer
			err := json.Indent(&prettyJSON, body, "", "  ")
			if err != nil {
				log.Println("JSON parse error: ", err)
				replyBody = string(body)
			}
			replyBody = string(prettyJSON.Bytes())
		} else {
			replyBody = string(body) // Convert the request body to a string and update the request panel
		}
	}
	// If not, display an error message

	return string(message) + replyBody
}

func ConvertToSendableString(bodyBuffer *bytes.Buffer) {
	// Convert the json string to one line
	if bytePayload := bodyBuffer.Bytes(); json.Valid(bytePayload) {
		dst := &bytes.Buffer{}
		if err := json.Compact(dst, []byte(bytePayload)); err != nil {
			log.Panicln("Error compacting json string:", err)
		}
		bodyBuffer.Reset()
		bodyBuffer.Write(dst.Bytes())
	}
}
