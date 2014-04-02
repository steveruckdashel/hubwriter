// Copyright 2014

// Package hubwriter provies a means to write to multiple io.WriteClosers which are subscribed. 
// I wrote this for sending rpc output over websockets (code.google.com/p/go.net/websocket).
//
// Here's an example of how I used this with websockets.
//
//		func ConsoleServer(ws *websocket.Conn) {
//			pr, pw := io.Pipe()
//			hubWriter.Subscribe(pw)
//			io.Copy(ws, pr)
//		}
//
package hubwriter

import (
	"io"
)

// HubWriter represents the array of io.WriteClosers to write to.
type HubWriter struct {
	hub []io.WriteCloser
}

// New returns a new HubWriter
func New() *HubWriter {
	return &HubWriter{hub: []io.WriteCloser{}}
}

// Write writes len(b) bytes to the HubWriter. It currently suppresses all errors and returns the length of b, not length of what was writen.
// If there is an error writing to an io.Writer in the hub, that io.Writer is Closed and removed.
func (hw *HubWriter) Write(b []byte) (int, error) {
	nhub := []io.WriteCloser{}
	n := len(b)

	for i := range hw.hub {
		nn, e := hw.hub[i].Write(b)
		if !(e != nil || nn != n) {
			nhub = append(nhub, hw.hub[i])
			hw.hub[i].Close()
		}
	}
	
	hw.hub = nhub
	return n, nil
}

// Close closes all io.WriteClosers in the hub.
func (hw *HubWriter) Close() error {
	var e error
	for i := range hw.hub {
		e = hw.hub[i].Close()
	}
	return e
}

// Subscribe adds an io.WriteCloser to the hub.
func (hw *HubWriter) Subscribe(s io.WriteCloser) {
	hw.hub = append(hw.hub, s)
}

// Unsubscribe removes an io.WriteCloser from the hub.
func (hw *HubWriter) Unsubscribe(s io.WriteCloser) {
	nhub := []io.WriteCloser{}
	for i := range hw.hub {
		if hw.hub[i] != s {
			nhub = append(nhub, hw.hub[i])
		}
	}
	hw.hub = nhub
}
