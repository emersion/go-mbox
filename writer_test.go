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
			`Date: Thu, 02 Jan 2015 00:00:01 +0100

This is another simple test.

Another line.

Bye.`,
		},
	}

	expected := `From ???@??? Thu Jan  1 00:00:01 2015` + "\r" + `
Date: Thu, 01 Jan 2015 00:00:01 +0100

This is a simple test.

And, by the way, this is how a "From" line is escaped in mboxo format:

>From Herp Derp with love.

Bye.` + "\r" + `
` + "\r" + `
From ???@??? Fri Jan  2 00:00:01 2015` + "\r" + `
Date: Thu, 02 Jan 2015 00:00:01 +0100

This is another simple test.

Another line.

Bye.` + "\r" + `
` + "\r" + `
`

	s := testWriter(t, messages)
	if s != expected {
		t.Error("Invalid mbox output:", s)
	}
}
