package traefikretryplugin

import (
	"bytes"
	"context"
	"fmt"
	. "github.com/niki-timofe/traefikretryplugin/internal"
	. "github.com/atidev/golib/pkg/structuredheaders"
	"io"
	"net/http"
	"sync"
)

type Config struct {
}

type retryPlugin struct {
	next http.Handler
	name string
	ctx  context.Context
}

//goland:noinspection GoUnusedExportedFunction
func CreateConfig() *Config {
	return &Config{}
}

//goland:noinspection GoUnusedExportedFunction
func New(_ context.Context, next http.Handler, _ *Config, name string) (http.Handler, error) {
	return &retryPlugin{
		next: next,
		name: name,
	}, nil
}

var bbPool = sync.Pool{New: func() interface{} { return make([]byte, 0, 512) }}

func (p *retryPlugin) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if bypass(req.Header) {
		p.next.ServeHTTP(rw, req)
		return
	}

	bb := bbPool.Get().([]byte)
	defer func() {
		bb = bb[:0]
		bbPool.Put(bb)
	}()

	rdr, err := bufferBody(bb, req)
	if err != nil {
		fmt.Printf("ServeHTTP: %s\n", err)

		internalServerError(rw)
		return
	}

	pl := policyFrom(req)

	fmt.Printf("ServeHTTP: %s\n", pl.String())

	var rrw *RetryResponseWriter

	for attempt := 0; rrw == nil || rrw.Retrying; attempt++ {
		if err = copyBody(rw, req, rdr); err != nil {
			fmt.Printf("ServeHTTP: %s\n", err)

			internalServerError(rw)
			return
		}

		rrw = NewRetryResponseWriter(rw, pl, attempt)

		p.next.ServeHTTP(rrw, req)
	}
}

func policyFrom(req *http.Request) *RetryPolicy {
	ph, err := NewStructuredHeader(req.Header).Dictionary("Retry-Policy")
	if err != nil {
		fmt.Printf("traefikretryplugin.policyFrom: error reading policy header as dictionary: %s\n", err)
		return nil
	}

	pl, err := ParsePolicy(ph)
	if err != nil {
		fmt.Printf("traefikretryplugin.policyFrom: error parsing policy: %s\n", err)
		return nil
	}

	return pl
}

func copyBody(rw http.ResponseWriter, req *http.Request, reader *bytes.Reader) error {
	_, err := reader.Seek(0, 0)
	if err != nil {
		internalServerError(rw)
		return err
	}

	req.Body = io.NopCloser(reader)

	return nil
}

func internalServerError(rw http.ResponseWriter) {
	http.Error(rw, "Internal Server Error", 500)
}

func bufferBody(buffer []byte, req *http.Request) (*bytes.Reader, error) {
L:
	for {
		if len(buffer) == cap(buffer) {
			buffer = append(buffer, 0)[:len(buffer)]
		}

		n, err := req.Body.Read(buffer[len(buffer):cap(buffer)])
		buffer = buffer[:len(buffer)+n]

		switch {
		case err == io.EOF:
			break L
		case err != nil:
			return nil, fmt.Errorf("traefikretryplugin.bufferBody: can't buffer body: %w", err)
		}
	}

	return bytes.NewReader(buffer), nil
}

func bypass(header http.Header) bool {
	return header.Get("Connection") == "Upgrade" && header.Get("Upgrade") == "websocket" ||
		header.Get("Transfer-Encoding") == "chunked" ||
		header.Get("Retry-Policy") == ""
}
