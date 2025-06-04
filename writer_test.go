package mbox

import (
	"bytes"
	"io"
	"strings"
	"testing"
	"time"
)

type testMessage struct {
	date string
	text string
}

func testWriter(t *testing.T, messages []testMessage) string {
	var b bytes.Buffer
	wc := NewWriter(&b)

	for _, m := range messages {
		r := strings.NewReader(m.text)
		date, _ := time.Parse(time.RFC1123Z, m.date)

		mw, err := wc.CreateMessage("", date)
		if err != nil {
			t.Fatal(err)
		}

		_, err = io.Copy(mw, r)
		if err != nil {
			t.Fatal(err)
		}
	}

	if err := wc.Close(); err != nil {
		t.Fatal(err)
	}

	return b.String()
}

func TestWriter(t *testing.T) {
	messages := []testMessage{
		{
			"Thu, 01 Jan 2015 00:00:01 +0100",
			`Date: Thu, 01 Jan 2015 00:00:01 +0100

This is a simple test.

And, by the way, this is how a "From" line is escaped in mboxo format:

From Herp Derp with love.

Bye.`,
		},
		{
			"Thu, 02 Jan 2015 00:00:01 +0100",
			`Date: Thu, 02 Jan 2015 00:00:01 +0100` + "\r" + `
` + "\r" + `
This is another simple test.` + "\r" + `
` + "\r" + `
Another line.` + "\r" + `
` + "\r" + `
Bye.`,
		},
	}

	expected := `From ???@??? Wed Dec 31 23:00:01 2014
Date: Thu, 01 Jan 2015 00:00:01 +0100

This is a simple test.

And, by the way, this is how a "From" line is escaped in mboxo format:

>From Herp Derp with love.

Bye.

From ???@??? Thu Jan  1 23:00:01 2015
Date: Thu, 02 Jan 2015 00:00:01 +0100

This is another simple test.

Another line.

Bye.

`

	s := testWriter(t, messages)
	if s != expected {
		t.Error("Invalid mbox output:", s)
	}
}

// Test for the fix for https://github.com/emersion/go-mbox/issues/21
func TestWriter_multipleWritesCount(t *testing.T) {
	var buffer bytes.Buffer
	mboxWriter := NewWriter(&buffer)
	w, err := mboxWriter.CreateMessage("-", time.Time{})
	if err != nil {
		panic(err)
	}
	b := []byte("some text   ")
	n, err := w.Write(b)
	if n != len(b) {
		t.Errorf("unexpected return value for write: %d (expected %d)", n, len(b))
	}
	b = []byte("end of line\n")
	n, err = w.Write(b)
	if n != len(b) {
		t.Errorf("unexpected return value for write: %d (expected %d)", n, len(b))
	}
}
