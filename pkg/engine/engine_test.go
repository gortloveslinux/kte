package engine

import (
	"io"
	"os"
	"testing"
)

func TestNewEngine(t *testing.T) {
	e := NewEngine()
	bs := e.GetBuffers()
	if len(bs) != 1 {
		t.Error("New engine should have 1 buffer")
	}
}

func TestNewBuffer(t *testing.T) {
	e := NewEngine()
	e.NewBuffer()
	bs := e.GetBuffers()
	if len(bs) != 2 {
		t.Error("Engine should have 2 buffers")
	}
}

func TestRemoveBuffer(t *testing.T) {
	e := NewEngine()
	b1 := e.NewBuffer()
	b2 := e.NewBuffer()
	b3 := e.NewBuffer()
	total := 4 // Engine has one new buffer then we added 3 buffers to it

	for _, v := range []Buffer{b2, b1, b3} {
		if len(e.GetBuffers()) != total {
			t.Errorf("Engine should have %d buffers", total)
		}
		e.RemoveBuffer(v.GetID())
		total = total - 1
	}
}

func TestDetermineLineEndings(t *testing.T) {
	f1, err := os.CreateTemp("", "")
	if err != nil {
		t.Skip("Couldn't create temp file, skipping")
	}
	f2, err := os.CreateTemp("", "")
	if err != nil {
		t.Skip("Couldn't create temp file, skipping")
	}
	f3, err := os.CreateTemp("", "")
	if err != nil {
		t.Skip("Couldn't create temp file, skipping")
	}
	f4, err := os.CreateTemp("", "")
	if err != nil {
		t.Skip("Couldn't create temp file, skipping")
	}
	f5, err := os.CreateTemp("", "")
	if err != nil {
		t.Skip("Couldn't create temp file, skipping")
	}

	defer func() {
		os.Remove(f1.Name())
		os.Remove(f2.Name())
		os.Remove(f3.Name())
		os.Remove(f4.Name())
		os.Remove(f5.Name())
	}()

	file1 := "This is a test\nThis file contains line feeds\n"
	file2 := "This is a test\r\nThis file contains carriage returns and line feeds\r\n"
	file3 := "This is a test\rThis file contains carriage returns only\r"
	file4 := "This is a test This file contains no line endings"

	f1.WriteString(file1)
	f2.WriteString(file2)
	f3.WriteString(file3)
	f4.WriteString(file4)

	f1.Sync()
	f2.Sync()
	f3.Sync()
	f4.Sync()

	f1.Seek(0, io.SeekStart)
	f2.Seek(0, io.SeekStart)
	f3.Seek(0, io.SeekStart)
	f4.Seek(0, io.SeekStart)
	f5.Seek(0, io.SeekStart)

	le1 := detirmineLineEndings(f1)
	le2 := detirmineLineEndings(f2)
	le3 := detirmineLineEndings(f3)
	le4 := detirmineLineEndings(f4)
	le5 := detirmineLineEndings(f5)

	if le1 != LF {
		t.Error("File should contain line feeds")
	}

	if le2 != CRLF {
		t.Error("File should contain carriage returns and line feeds")
	}

	if le3 != INVALID {
		t.Error("File should contain invalid line endings", le3)
	}

	if le4 != NONE {
		t.Error("File should not contain line endings. Got:", le4)
	}

	if le5 != UNKNOWN {
		t.Error("File should have unknow line endings. Got:", le5)
	}

}
