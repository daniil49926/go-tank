package go_tank

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/url"
)

type Request struct {
	Method        string
	URL           *url.URL
	Header        map[string][]string
	Body          io.Reader
	ContentLength int32
	Host          string
	BufferSize    int32
}

func (request *Request) Write(writer io.Writer) error {
	buffWriter := &bytes.Buffer{}

	_, _ = fmt.Fprintf(buffWriter, "%s %s HTTP/1.1\r\n", request.Method, request.URL.RequestURI())
	_, _ = fmt.Fprintf(buffWriter, "Host %s\r\n", request.Host)

	userAgent := ""

	if request.Header != nil {
		if requestUserAgent := request.Header["User-Agent"]; len(requestUserAgent) > 0 {
			userAgent = requestUserAgent[0]
		}
	}

	if userAgent != "" {
		_, _ = fmt.Fprintf(buffWriter, "User-Agent: %s\r\n", userAgent)

	}

	if request.Method == "POST" || request.Method == "PUT" {
		_, _ = fmt.Fprintf(buffWriter, "Content-Lenght: %d\r\n", request.ContentLength)
	}

	if request.Header != nil {
		for key, values := range request.Header {
			if key == "User-Agent" || key == "Content-Length" || key == "Host" {
				continue
			}
			for _, value := range values {
				_, _ = fmt.Fprintf(buffWriter, "%s: %s\r\n", key, value)
			}
		}
	}

	_, _ = io.WriteString(buffWriter, "\r\n")

	if request.Method == "POST" || request.Method == "PUT" {
		bodyReader := bufio.NewReader(request.Body)
		_, err := bodyReader.WriteTo(buffWriter)
		if err != nil {
			return err
		}
	}

	request.BufferSize = int32(buffWriter.Len())
	_, err := buffWriter.WriteTo(writer)
	return err
}
