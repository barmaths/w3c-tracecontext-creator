package w3c_traceparent_creator

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

const defaultHeader = "traceparent"

var src = rand.New(rand.NewSource(time.Now().UnixNano()))

// Config the plugin configuration.
type Config struct {
	HeaderName string
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		HeaderName: defaultHeader,
	}
}

// W3CTraceParentCreator Plugin Struct.
type W3CTraceParentCreator struct {
	next       http.Handler
	headerName string
	name       string
}

// New creates a new plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if len(config.HeaderName) == 0 {
		return nil, fmt.Errorf("no header name provided")
	}
	return &W3CTraceParentCreator{
		next:       next,
		name:       name,
		headerName: config.HeaderName,
	}, nil
}

// ServeHTTP defines the behaviour of the plugin
func (tpc *W3CTraceParentCreator) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if len(req.Header.Get(tpc.headerName)) > 0 {
		rw.Header().Add(tpc.headerName, req.Header.Get(tpc.headerName))
		tpc.next.ServeHTTP(rw, req)
	} else {
		traceParent := "00-" + RandomHexaDecimalStringOfLength(32) + "-" + RandomHexaDecimalStringOfLength(16) + "-" + "00"
		req.Header.Add(tpc.headerName, traceParent)
		rw.Header().Add(tpc.headerName, traceParent)
		tpc.next.ServeHTTP(rw, req)
	}
}

// RandomHexaDecimalStringOfLength returns a random hexadecimal string of length n.
func RandomHexaDecimalStringOfLength(n int) string {
	b := make([]byte, n/2)

	if _, err := src.Read(b); err != nil {
		panic(err)
	}

	return hex.EncodeToString(b)[:n]
}
