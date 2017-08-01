package mbox

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

const mboxWithOneMessage = `From herp.derp at example.com  Thu Jan  1 00:00:01 2015
From: herp.derp at example.com (Herp Derp)
Date: Thu, 01 Jan 2015 00:00:01 +0100
Subject: Test

This is a simple test.

And, by the way, this is how a "From" line is escaped in mboxo format:

>From Herp Derp with love.

Bye.
`

const mboxWithOneMessageMissingHeaders = `From herp.derp at example.com  Thu Jan  1 00:00:01 2015
This is a simple test.

And, by the way, this is how a "From" line is escaped in mboxo format:

>From Herp Derp with love.

Bye.
`

const mboxWithThreeMessages = `From herp.derp at example.com  Thu Jan  1 00:00:01 2015
From: herp.derp at example.com (Herp Derp)
Date: Thu, 01 Jan 2015 00:00:01 +0100
Subject: Test

This is a simple test.

And, by the way, this is how a "From" line is escaped in mboxo format:

>From Herp Derp with love.

Bye.

From derp.herp at example.com  Thu Jan  1 00:00:01 2015
From: derp.herp at example.com (Derp Herp)
Date: Thu, 02 Jan 2015 00:00:01 +0100
Subject: Another test

This is another simple test.

Another line.

Bye.

From bernd.lauert at example.com  Thu Jan  3 00:00:01 2015
From: bernd.lauert at example.com (Bernd Lauert)
Date: Thu, 03 Jan 2015 00:00:01 +0100
Subject: A last test

This is the last simple test.

Bye.
`

const mboxWithStartingLF = `
From herp.derp at example.com  Thu Jan  1 00:00:01 2015
From: herp.derp at example.com (Herp Derp)
Date: Thu, 01 Jan 2015 00:00:01 +0100
Subject: Test

This is a simple test.

And, by the way, this is how a "From" line is escaped in mboxo format:

>From Herp Derp with love.

Bye.

From derp.herp at example.com  Thu Jan  1 00:00:01 2015
From: derp.herp at example.com (Derp Herp)
Date: Thu, 02 Jan 2015 00:00:01 +0100
Subject: Another test

This is another simple test.

Another line.

Bye.

From bernd.lauert at example.com  Thu Jan  3 00:00:01 2015
From: bernd.lauert at example.com (Bernd Lauert)
Date: Thu, 03 Jan 2015 00:00:01 +0100
Subject: A last test

This is the last simple test.

Bye.
`

const mboxWithThreeMessagesMalformedButValid = `From herp.derp at example.com  Thu Jan  1 00:00:01 2015
From: herp.derp at example.com (Herp Derp)
Date: Thu, 01 Jan 2015 00:00:01 +0100
Subject: Test

This is a simple test.

And, by the way, this is how a "From" line is escaped in mboxo format:

>From Herp Derp with love.

Bye.
From derp.herp at example.com  Thu Jan  1 00:00:01 2015
From: derp.herp at example.com (Derp Herp)
Date: Thu, 02 Jan 2015 00:00:01 +0100
Subject: Another test

This is another simple test.

Another line.

Bye.

From bernd.lauert at example.com  Thu Jan  3 00:00:01 2015
From: bernd.lauert at example.com (Bernd Lauert)
Date: Thu, 03 Jan 2015 00:00:01 +0100
Subject: A last test

This is the last simple test.

Bye.
`

const mboxWithOneMessageMissingSeparator = `From: herp.derp at example.com (Herp Derp)
Date: Thu, 01 Jan 2015 00:00:01 +0100
Subject: Test

This is a simple test.

And, by the way, this is how a "From" line is escaped in mboxo format:

>From Herp Derp with love.

Bye.
`

const mboxFirstMessage = `From: herp.derp at example.com (Herp Derp)
Date: Thu, 01 Jan 2015 00:00:01 +0100
Subject: Test

This is a simple test.

And, by the way, this is how a "From" line is escaped in mboxo format:

>From Herp Derp with love.

Bye.
`

const mboxFirstMessageBody = `This is a simple test.

And, by the way, this is how a "From" line is escaped in mboxo format:

>From Herp Derp with love.

Bye.
`

const mboxSecondMessageSubjectHeader = "Another test"

type tsmInput struct {
	data  string
	atEOF bool
}

type tsmExpected struct {
	advance     int
	token       string
	yieldsError bool
}

func testScanMessage(t *testing.T, input *tsmInput, expected *tsmExpected) {
	advance, token, err := scanMessage([]byte(input.data), input.atEOF)

	if err == nil && expected.yieldsError {
		t.Errorf("unexpected success")
	}
	if err != nil && !expected.yieldsError {
		t.Errorf("unexpected error: %v", err)
	}
	if advance != expected.advance {
		t.Errorf("unexpected advance: %d", advance)
	}
	if string(token) != expected.token {
		t.Errorf("unexpected token: %q", token)
	}
}

func TestScanMessageMboxEmptyAtEOF(t *testing.T) {
	input := &tsmInput{
		atEOF: true,
		data:  "",
	}

	expected := &tsmExpected{
		yieldsError: false,
		advance:     0,
		token:       "",
	}

	testScanMessage(t, input, expected)
}

func TestScanMessageMboxWithOneMessageAtEOF(t *testing.T) {
	input := &tsmInput{
		atEOF: true,
		data:  mboxWithOneMessage,
	}

	expected := &tsmExpected{
		yieldsError: false,
		advance:     281,
		token:       mboxFirstMessage,
	}

	testScanMessage(t, input, expected)
}

func TestScanMessageMboxWithOneMessageMissingSeparatorAtEOF(t *testing.T) {
	input := &tsmInput{
		atEOF: true,
		data:  mboxWithOneMessageMissingSeparator,
	}

	expected := &tsmExpected{
		yieldsError: true,
		advance:     0,
		token:       "",
	}

	testScanMessage(t, input, expected)
}

func TestScanMessageMboxWithThreeMessages(t *testing.T) {
	input := &tsmInput{
		atEOF: false,
		data:  mboxWithThreeMessages,
	}

	expected := &tsmExpected{
		yieldsError: false,
		advance:     282,
		token:       mboxFirstMessage,
	}

	testScanMessage(t, input, expected)
}

func TestScanMessageWithThreeMessagesMalformedButValid(t *testing.T) {
	input := &tsmInput{
		atEOF: false,
		data:  mboxWithThreeMessagesMalformedButValid,
	}

	expected := &tsmExpected{
		yieldsError: false,
		advance:     281,
		token:       mboxFirstMessage,
	}

	testScanMessage(t, input, expected)
}

func TestScanMessageIncompleteRecord(t *testing.T) {
	input := &tsmInput{
		atEOF: false,
		data:  mboxWithOneMessage[:100],
	}

	expected := &tsmExpected{
		yieldsError: false,
		advance:     0,
		token:       "",
	}

	testScanMessage(t, input, expected)
}

func TestScanMessageVeryShortIncompleteRecord(t *testing.T) {
	input := &tsmInput{
		atEOF: false,
		data:  "From",
	}

	expected := &tsmExpected{
		yieldsError: false,
		advance:     0,
		token:       "",
	}

	testScanMessage(t, input, expected)
}

func TestScanMessageOnlySeperatorAtEOF(t *testing.T) {
	input := &tsmInput{
		atEOF: true,
		data:  mboxWithOneMessage[:55],
	}

	expected := &tsmExpected{
		yieldsError: true,
		advance:     0,
		token:       "",
	}

	testScanMessage(t, input, expected)
}

func TestScanMessageMboxWithOneMessageWithoutNewlineAtEOF(t *testing.T) {
	input := &tsmInput{
		atEOF: true,
		data:  mboxWithOneMessage[:len(mboxWithOneMessage)-1],
	}

	expected := &tsmExpected{
		yieldsError: true,
		advance:     0,
		token:       "",
	}

	testScanMessage(t, input, expected)
}

func testMboxMessage(t *testing.T, mbox string, count int) {
	b := bytes.NewBufferString(mbox)
	m := NewScanner(b)

	for i := 0; i < count; i++ {
		if !m.Next() {
			t.Errorf("Next() failed; pass %d", i)
		}
		if m.Err() != nil {
			t.Errorf("Unexpected error after Next(): %v", m.Err())
		}

		msg := m.Message()
		if msg == nil {
			t.Errorf("message is nil; pass %d", i)
			continue
		}
		body := new(bytes.Buffer)
		_, err := body.ReadFrom(msg.Body)
		if err != nil {
			t.Errorf("Unexpected error reading message body: %v", err)
		}
		if i == 0 && body.String() != mboxFirstMessageBody {
			t.Errorf("Expected:\n %q\ngot\n%q", mboxFirstMessageBody, body.String())
		}
		if i == 1 && msg.Header.Get("Subject") != mboxSecondMessageSubjectHeader {
			t.Errorf("Unexpected subject header: %q", msg.Header.Get("Subject"))
		}
		if m.Err() != nil {
			t.Errorf("Unexpected error after Message(): %v", m.Err())
		}
	}

	if m.Next() {
		t.Errorf("Next() succeeded")
	}
	if m.Err() != nil {
		t.Errorf("Unexpected error after Next(): %v", m.Err())
	}
	if msg := m.Message(); msg != nil {
		t.Errorf("message is not nil")
	}
	if m.Err() != nil {
		t.Errorf("Unexpected error after Message(): %v", m.Err())
	}
}

func TestMboxMessageWithOneMessage(t *testing.T) {
	testMboxMessage(t, mboxWithOneMessage, 1)
}

func TestMboxMessageWithThreeMessages(t *testing.T) {
	testMboxMessage(t, mboxWithThreeMessages, 3)
}

func TestMboxMessageWithStartingLF(t *testing.T) {
	testMboxMessage(t, mboxWithStartingLF, 3)
}

func TestMboxMessageWithThreeMessagesMalformedButValid(t *testing.T) {
	testMboxMessage(t, mboxWithThreeMessagesMalformedButValid, 3)
}

func testMboxMessageInvalid(t *testing.T, mbox string) {
	b := bytes.NewBufferString(mbox)
	m := NewScanner(b)

	if m.Next() {
		t.Errorf("Next() succeeded")
	}
	if m.Err() == nil {
		t.Errorf("Missing error after Next(): %v", m.Err())
	}
	if msg := m.Message(); msg != nil {
		t.Errorf("message is not nil")
	}
	if m.Err() == nil {
		t.Errorf("Missing error after Message(): %v", m.Err())
	}
	if m.Next() {
		t.Errorf("Next() after error succeeded")
	}
}

func TestMboxMessageWithOneMessageMissingSeparator(t *testing.T) {
	testMboxMessageInvalid(t, mboxWithOneMessageMissingSeparator)
}

func TestMboxMessageWithOneMessageMissingHeaders(t *testing.T) {
	testMboxMessageInvalid(t, mboxWithOneMessageMissingHeaders)
}

func TestScanMessageWithBoundaries(t *testing.T) {
	sourceData := `
From one place.
From: herp.derp at example.com (Herp Derp)
Date: Thu, 01 Jan 2015 00:00:01 +0100
Subject: Test
Content-Type: multipart/alternative; boundary="_----------=bestboundaryever"

This is a test of boundaries.  Don't accept a new email via \nFrom until the boundary is done!'

And, by the way, this is how a "From" line is escaped in mboxo format:
From Herp Derp with love.

From Herp Derp with love.

Bye.
--_----------=bestboundaryever--

From another!
From: herp.derp at example.com (Herp Derp)
Date: Thu, 01 Jan 2015 00:00:01 +0100
Subject: Test

This is the second email in a test of boundaries.
`
	expected := []string{
		"This is a test of boundaries.  Don't accept a new email via \\nFrom until the boundary is done!'\n\nAnd, by the way, this is how a \"From\" line is escaped in mboxo format:\nFrom Herp Derp with love.\n\nFrom Herp Derp with love.\n\nBye.\n--_----------=bestboundaryever--\n",
		"This is the second email in a test of boundaries.\n",
	}
	b := bytes.NewBufferString(sourceData)
	m := NewScanner(b)

	for i := range expected {
		if !m.Next() {
			t.Errorf("Next() failed; pass %d", i)
		}
		if m.Err() != nil {
			t.Errorf("Unexpected error after Next(): %v", m.Err())
		}

		msg := m.Message()
		if msg == nil {
			t.Errorf("message is nil; pass %d", i)
			continue
		}
		body := new(bytes.Buffer)
		_, err := body.ReadFrom(msg.Body)
		if err != nil {
			t.Errorf("%d - Unexpected error reading message body: %v", i, err)
			continue
		}
		if body.String() != expected[i] {
			t.Errorf("%d - Expected:\n %q\ngot\n%q", i, expected[i], body.String())
		}
		if m.Err() != nil {
			t.Errorf("%d - Unexpected error after Message(): %v", i, m.Err())
		}
	}

	if m.Next() {
		t.Errorf("Next() succeeded")
	}
	if m.Err() != nil {
		t.Errorf("Unexpected error after Next(): %v", m.Err())
	}
	if msg := m.Message(); msg != nil {
		t.Errorf("message is not nil")
	}
	if m.Err() != nil {
		t.Errorf("Unexpected error after Message(): %v", m.Err())
	}
}

func ExampleScanner() {
	r := strings.NewReader(`From herp.derp at example.com  Thu Jan  1 00:00:01 2015
From: herp.derp at example.com (Herp Derp)
Date: Thu, 01 Jan 2015 00:00:01 +0100
Subject: Test

This is a simple test.

CU.

From derp.herp at example.com  Thu Jan  1 00:00:01 2015
From: derp.herp at example.com (Derp Herp)
Date: Thu, 02 Jan 2015 00:00:01 +0100
Subject: Another test

This is another simple test.

Bye.
`)

	mbox := NewScanner(r)
	for mbox.Next() {
		// If Next() returns true, you can expect Message() to return a
		// valid *mail.Message.
		fmt.Printf("Message from %v\n", mbox.Message().Header.Get("from"))
	}
	if mbox.Err() != nil {
		fmt.Print("Oops, something went wrong!")
	}
	// Output:
	// Message from herp.derp at example.com (Herp Derp)
	// Message from derp.herp at example.com (Derp Herp)
}
