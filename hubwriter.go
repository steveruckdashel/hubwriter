package hubwriter

import (
	"io"
	"os"
)

type HubWriter struct {
	hub []io.WriteCloser
}

func NewHubWriter() *HubWriter {
	null, err := os.OpenFile(os.DevNull, os.O_WRONLY, os.ModeAppend)
	if err != nil {
		panic(err)
	}
	return &HubWriter{
		hub: []io.WriteCloser{null},
	}
}

func (hw *HubWriter) Write(p []byte) (int, error) {
	nhub := []io.WriteCloser{}
	n := len(p)

	for i := range hw.hub {
		nn, e := hw.hub[i].Write(p)
		if !(e != nil || nn != n) {
			nhub = append(nhub, hw.hub[i])
		}
	}
	
	hw.hub = nhub
	return n, nil
}

func (hw *HubWriter) Close() error {
	var e error
	for i := range hw.hub {
		e = hw.hub[i].Close()
	}
	return e
}

func (hw *HubWriter) Subscribe(s io.WriteCloser) {
	hw.hub = append(hw.hub, s)
}

func (hw *HubWriter) Unsubscribe(s io.WriteCloser) {
	nhub := []io.WriteCloser{}
	for i := range hw.hub {
		if hw.hub[i] != s {
			nhub = append(nhub, hw.hub[i])
		}
	}
	hw.hub = nhub
}
