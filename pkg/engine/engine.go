package engine

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"path"
)

var NotImplementedErr = errors.New("Not Implemented")

type lineEnding int

const (
	LF lineEnding = iota
	CRLF
)

var lineEndingString = map[lineEnding]string{LF: "\n", CRLF: "\r\n"}

type guts struct {
	buffers []Buffer
	genID   func() int
}

type buffer struct {
	name       string
	content    [][]byte
	cursorPos  struct{ ln, col int }
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
		return nil, err
	}
	b := &buffer{
		id: g.genID(),
	}

	err = readFileIntoBuffer(b, f)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func readFileIntoBuffer(b *buffer, f *os.File) error {
	b.name = f.Name()
	b.file = f
	defer f.Seek(0, io.SeekStart)
	r := bufio.NewReader(f)

	// read first line and detrmine line endings
	l, err := r.ReadBytes('\n')
	if err != nil && err != io.EOF {
		return err
	}

	if len(l) > 2 {
		if bytes.Compare(l[len(l)-2:], []byte("\r\n")) == 0 {
			b.lineEnding = CRLF
		} else {
			b.lineEnding = LF
		}
	}
	b.content = append(b.content, l)

	for {
		l, err := r.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return err
		}
		if len(l) > 0 {
			b.content = append(b.content, l)
		}
		if err == io.EOF {
			break
		}
	}
	return nil
}

func (b *buffer) Save() error {
	return NotImplementedErr
}

func (b *buffer) SaveAs(string) error                          { return NotImplementedErr }
func (b *buffer) Read([]byte) (int, error)                     { return 0, NotImplementedErr }
func (b *buffer) Write([]byte) (int, error)                    { return 0, NotImplementedErr }
func (b *buffer) Seek(offset int64, whence int) (int64, error) { return 0, NotImplementedErr }
func (b *buffer) GetID() int                                   { return b.id }
