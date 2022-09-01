package engine

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
)

var NotImplementedErr = errors.New("Not Implemented")

type lineEnding int

const (
	LF lineEnding = iota
	CRLF
	INVALID
	UNKNOWN
	NONE
)

type guts struct {
	buffers []Buffer
	genID   func() int
}

type buffer struct {
	name       string
	body       []bytes.Buffer
	cursorPos  struct{ x, y int }
	file       *os.File
	dirty      bool
	id         int
	lineEnding lineEnding
}

type Engine interface {
	// Buffer Related
	GetBuffers() []Buffer
	NewBuffer() Buffer
	NewBufferFromFile(string) (Buffer, error)
	RemoveBuffer(int)
}

type Buffer interface {
	io.ReadWriteSeeker
	Save() error
	SaveAs(string) error
	GetID() int
}

func NewEngine() Engine {
	genID := func() func() int {
		id := 0
		return func() int {
			nid := id
			id = id + 1
			return nid
		}
	}()
	g := &guts{genID: genID}

	g.NewBuffer()
	return g
}

func (g *guts) GetBuffers() []Buffer {
	return g.buffers
}

func (g *guts) NewBuffer() Buffer {
	b := &buffer{
		name: "new",
		id:   g.genID(),
	}
	g.buffers = append(g.buffers, b)
	return b
}

func (g *guts) RemoveBuffer(ID int) {
	for i, v := range g.buffers {
		if v.GetID() == ID {
			g.buffers = append(g.buffers[:i], g.buffers[i+1:]...)
		}
	}
}

func (g *guts) NewBufferFromFile(filename string) (Buffer, error) {
	filename = path.Clean(filename)
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, fmt.Errorf("Could not create buffer: %w", err)
	}
	b := &buffer{id: g.genID(), file: f, name: filename}
	b.lineEnding = detirmineLineEndings(f)
	return b, nil
}

func (b *buffer) Save() error                                  { return NotImplementedErr }
func (b *buffer) SaveAs(string) error                          { return NotImplementedErr }
func (b *buffer) Read([]byte) (int, error)                     { return 0, NotImplementedErr }
func (b *buffer) Write([]byte) (int, error)                    { return 0, NotImplementedErr }
func (b *buffer) Seek(offset int64, whence int) (int64, error) { return 0, NotImplementedErr }
func (b *buffer) GetID() int                                   { return b.id }

func detirmineLineEndings(f *os.File) lineEnding {
	r := bufio.NewReader(f)
	checked := false
	for {
		b, err := r.ReadByte()
		if err != nil {
			if !checked {
				return UNKNOWN
			}
			return NONE
		}
		checked = true
		switch b {
		case '\r':
			bs, err := r.Peek(1)
			if err != nil {
				return INVALID
			}
			if bs[0] == '\n' {
				return CRLF
			}
		case '\n':
			return LF
		}
	}
	return NONE
}
