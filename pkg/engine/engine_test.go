package engine

import (
	"fmt"
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

func TestReadFileIntoBuffer(t *testing.T) {
	testData := []struct {
		b  *buffer
		fn string
		le lineEnding
		lc int
	}{
		{
			&buffer{},
			"testLineFeedFile1.txt",
			LF,
			4,
		},
		{
			&buffer{},
			"testCarriageLineFeedFile1.txt",
			CRLF,
			5,
		},
	}

	for _, v := range testData {
		wd, err := os.Getwd()
		if err != nil {
			t.Error("Could not get working directory: ", err)
		}
		filepath := wd + "/testdata/" + v.fn
		t.Log("Test file: ", filepath)
		f, err := os.Open(filepath)
		if err != nil {
			t.Error("Could not open file: ", err)
		}
		err = readFileIntoBuffer(v.b, f)
		if err != nil {
			t.Error("Should not be an error", err)
		}
		if v.b.lineEnding != v.le {
			t.Error(fmt.Sprintf("Line ending should be %v. Got: %v", []byte(lineEndingString[v.le]), []byte(lineEndingString[v.b.lineEnding])))
		}
		if len(v.b.content) != v.lc {
			t.Error(fmt.Sprintf("Buffer should only contain %d lines. Has: %d", v.lc, len(v.b.content)))
		}
	}
}

func TestNewBufferFromFileExisting(t *testing.T) {
	f, err := os.CreateTemp("", "")
	if err != nil {
		t.Skip("Couldn't create temp file, skipping")
	}
	filename := f.Name()
	defer os.Remove(filename)

	f.WriteString("This is a test.\nThere will be 4 lines.\nThis is line 3.\nGoodbye.")
	f.Close()

	e := NewEngine()
	b, err := e.NewBufferFromFile(filename)
	if err != nil {
		t.Error("Could not create error: ", err)
	}
	buf := b.(*buffer)

	if buf.name != filename {
		t.Error(fmt.Sprintf("File name mismatch. Got: %s Expected: %s", buf.name, filename))
	}

	if buf.lineEnding != LF {
		t.Error("Buffer has wrong line ending")
	}

	if len(buf.content) != 4 {
		t.Error("Buffer should only contain 4 lines. Has: ", len(buf.content), fmt.Sprintf("%v", buf.content))
	}
}

func TestNewBufferFromFileCreated(t *testing.T) {
	dir := os.TempDir()
	filename := dir + "testfile"
	defer os.Remove(filename)
	e := NewEngine()
	_, err := e.NewBufferFromFile(filename)
	if err != nil {
		t.Error("Could not create error: ", err)
	}
	_, err = os.Stat(filename)
	if err != nil {
		t.Error("There should be a file created")
	}
}

func TestBufferSave(t *testing.T) {

}

func TestBufferSaveAs(t *testing.T) {

}
func TestBufferRead(t *testing.T) {

}
func TestBufferWrite(t *testing.T) {

}
func TestBufferSeek(t *testing.T) {

}
