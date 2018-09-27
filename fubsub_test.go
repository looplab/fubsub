package fubsub

import (
	"os"
	"testing"
)

func TestFubsub(t *testing.T) {
	os.Setenv("PUBSUB_EMULATOR_HOST", "localhost:8793")

	f, err := New("looplab-playground", "foo")
	if err != nil {
		t.Fatal("there should be no error:", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			t.Error("there should be no error:", err)
		}
	}()

	written := "Hello world!"
	n, err := f.Write([]byte(written))
	if err != nil {
		t.Error("there should be no error:", err)
	}
	if n != 12 {
		t.Error("the length should be correct:", n)
	}

	read := make([]byte, 100)
	n, err = f.Read(read)
	if err != nil {
		t.Error("there should be no error:", err)
	}
	if n != 12 {
		t.Error("the length should be correct:", n)
	}

	if string(read[:n]) != written {
		t.Error("the read data should be correct:", string(read))
	}
}
