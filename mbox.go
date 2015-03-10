// "THE BEER-WARE LICENSE" (Revision 42):
// <tobias.rehbein@web.de> wrote this file. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you
// think this stuff is worth it, you can buy me a beer in return.
//                                                             Tobias Rehbein

// Package mbox parses the mbox file format into messages.
//
// As the mbox file format is not standardized this package expects the least
// common denominator, the so called mboxo format.
package mbox

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"net/mail"
)

// ErrInvalidMboxFormat is the error returned by the Next method of type Mbox if
// its content is malformed in a way that it is not possible to extract a
// message.
var ErrInvalidMboxFormat = errors.New("invalid mbox format")

// scanMessage is a split function for a bufio.Scanner that returns a message in
// RFC 822 format or an error.
func scanMessage(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if len(data) == 0 && atEOF {
		return 0, nil, nil
	}

	var n int
	e := bytes.Index(data, []byte("\nFrom "))
	if e == -1 && !atEOF {
		// request more data
		return 0, nil, nil
	}

	if !bytes.HasPrefix(data, []byte("From ")) {
		return 0, nil, ErrInvalidMboxFormat
	}

	n = bytes.IndexByte(data, '\n')
	if n == -1 {
		return 0, nil, ErrInvalidMboxFormat
	}

	if atEOF {
		if data[len(data)-1] != '\n' {
			return 0, nil, ErrInvalidMboxFormat
		}

		return len(data), data[n+1:], nil
	}

	if data[e-1] != '\n' {
		return e + 1, data[n+1 : e+1], nil
	}

	return e + 1, data[n+1 : e], nil
}

// Mbox provides an interface to read a sequence of messages from an mbox.
// Calling the Next method steps through the messages. The current message can
// then be accessed by calling the Message method.
//
// The Next method returns true while there are messages to skip to and no error
// occurs. When Next returns false, you can call the Err method to check for an
// error.
//
// The Message method returns the current message as mail.Message, or nil if an
// error occurs or if you have skipped past the last message using Next. You
// cann call the Err method to check for an error.
type Mbox struct {
	s   *bufio.Scanner
	m   *mail.Message
	err error
}

// Mbox provides an interface to read a sequence of messages from an mbox.
// Calling the Next method steps through the messages. The current message can
// then be accessed by calling the Message method.
//
// The Next method returns true while there are messages to skip to and no error
// occurs. When Next returns false, you can call the Err method to check for an
// error.
//
// The Message method returns the current message as mail.Message, or nil if an
// error occured while calling Next or if you have skipped past the last message
// using Next. If Next returned true, you can expect Message to return a valid
// *mail.Message.
func New(r io.Reader) *Mbox {
	s := bufio.NewScanner(r)
	s.Split(scanMessage)

	return &Mbox{s, nil, nil}
}

// Next skips to the next message and returns true. It will return false if
// there are no messages left or an error occurs. You can call the Err method to
// check if an error occured. If Next returns false and Err returns nil there
// are no messages left.
func (m *Mbox) Next() bool {
	m.m = nil
	if m.err != nil {
		return false
	}

	if !m.s.Scan() {
		m.err = m.s.Err()
		return false
	}

	m.m, m.err = mail.ReadMessage(bytes.NewBuffer(m.s.Bytes()))
	if m.err != nil {
		return false
	}

	return true
}

// Err returns the first error that occured while calling Next.
func (m *Mbox) Err() error {
	return m.err
}

// Message returns the current message. It returns nil if you never called Next,
// skipped past the last message or if an error occured during a call to Next.
//
// If Next returned true, you can expect Message to return a valid
// *mail.Message.
func (m *Mbox) Message() *mail.Message {
	if m.err != nil {
		return nil
	}

	return m.m
}
