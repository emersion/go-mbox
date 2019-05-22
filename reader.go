// "THE BEER-WARE LICENSE" (Revision 42):
// <tobias.rehbein@web.de> wrote this file. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you
// think this stuff is worth it, you can buy me a beer in return.
//                                                             Tobias Rehbein

package mbox

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"net/textproto"
	"strings"
)

// ErrInvalidMboxFormat is the error returned by the Next method of type Mbox if
// its content is malformed in a way that it is not possible to extract a
// message.
var ErrInvalidMboxFormat = errors.New("invalid mbox format")

// scanMessage is a split function for a bufio.Reader that returns a message in
// RFC 822 format or an error.
func scanMessage(data []byte, atEOF bool) (int, []byte, error) {
	if len(data) == 0 && atEOF {
		return 0, nil, nil
	}

	var n int
	e := bytes.Index(data, []byte("\nFrom "))
	advanceExtra := 0
	if e == 0 {
		data = data[1:] // advance past the leading LF
		advanceExtra = 1
		e = bytes.Index(data, []byte("\nFrom "))
	}
	if e == -1 && !atEOF {
		// request more data
		return advanceExtra, nil, nil
	}

	if !bytes.HasPrefix(data, []byte("From ")) {
		return advanceExtra, nil, ErrInvalidMboxFormat
	}
	n = bytes.IndexByte(data, '\n')
	if n == -1 {
		return advanceExtra, nil, ErrInvalidMboxFormat
	}

	if atEOF {
		if data[len(data)-1] != '\n' {
			return advanceExtra, nil, ErrInvalidMboxFormat
		}
		return len(data) + advanceExtra, data[n+1:], nil
	}
	tpr := textproto.NewReader(bufio.NewReader(bytes.NewReader(data[n+1:])))
	header, err := tpr.ReadMIMEHeader()
	if err != nil {
		return 0, nil, err
	}
	cth := header.Get("Content-Type")
	boundaryEnd := ""
	if strings.Contains(cth, "multipart") {
		splt := strings.Split(cth, "; ")
		for _, v := range splt {
			if strings.HasPrefix(v, "boundary=") {
				c := strings.Index(v, "=") + 1
				semi := strings.Index(v[c:], ";")
				if semi == -1 {
					semi = len(v[c:])
				}
				boundaryEnd = "--" + strings.Trim(v[c:c+semi], `"'`) + "--"
				break
			}
		}
	}
	if boundaryEnd != "" {
		b := bytes.Index(data, []byte(boundaryEnd))
		if b == -1 {
			return 0, nil, nil // need more data!
		}
		if e < b {
			e = bytes.Index(data[b:], []byte("\nFrom "))
			e += b
		}
	}
	if data[e-1] != '\n' {
		return e + 1 + advanceExtra, data[n+1 : e+1], nil
	}
	return e + 1 + advanceExtra, data[n+1 : e], nil
}

type Reader struct {
	s   *bufio.Scanner
	err error
}

// NewReader returns a new *Reader to read messages from mbox file format data
// provided by io.Reader r.
func NewReader(r io.Reader) *Reader {
	s := bufio.NewScanner(r)
	s.Split(scanMessage)

	return &Reader{s, nil}
}

// Next skips to the next message and returns true. It will return false if
// there are no messages left or an error occurs. You can call the Err method to
// check if an error occured. If Next returns false and Err returns nil there
// are no messages left.
func (m *Reader) NextMessage() (io.Reader, error) {
	if !m.s.Scan() {
		if err := m.s.Err(); err != nil {
			return nil, err
		}
		return nil, io.EOF
	}

	return bytes.NewReader(m.s.Bytes()), nil
}

// Buffer sets the initial buffer to use when scanning and the maximum size of
// buffer that may be allocated during scanning.
//
// Buffer panics if it is called after scanning has started.
func (m *Reader) Buffer(buf []byte, max int) {
	m.s.Buffer(buf, max)
}
