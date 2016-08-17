package mbox

import (
	"bytes"
	"net/mail"
	"strings"
	"testing"
)

func testWriter(t *testing.T, messages []*mail.Message) string {
	b := &bytes.Buffer{}
	w := NewWriter(b)

	for _, m := range messages {
		if _, err := w.WriteMessage(m); err != nil {
			t.Fatal(err)
		}
	}

	return b.String()
}

func TestWriter(t *testing.T) {
	messages := []*mail.Message{
		&mail.Message{
			Header: map[string][]string{
				"Date": {"Thu, 01 Jan 2015 00:00:01 +0100"},
			},
			Body: strings.NewReader(`This is a simple test.

And, by the way, this is how a "From" line is escaped in mboxo format:

From Herp Derp with love.

Bye.`),
		},
		&mail.Message{
			Header: map[string][]string{
				"Date": {"Thu, 02 Jan 2015 00:00:01 +0100"},
			},
			Body: strings.NewReader(`This is another simple test.

Another line.

Bye.`),
		},
	}

	expected := `From ???@??? Thu Jan  1 00:00:01 2015`+"\r"+`
Date: Thu, 01 Jan 2015 00:00:01 +0100`+"\r"+`
`+"\r"+`
This is a simple test.

And, by the way, this is how a "From" line is escaped in mboxo format:

>From Herp Derp with love.

Bye.`+"\r"+`
`+"\r"+`
From ???@??? Fri Jan  2 00:00:01 2015`+"\r"+`
Date: Thu, 02 Jan 2015 00:00:01 +0100`+"\r"+`
`+"\r"+`
This is another simple test.

Another line.

Bye.`+"\r"+`
`+"\r"+`
`

	s := testWriter(t, messages)
	if s != expected {
		t.Error("Invalid mbox output:", s)
	}
}
