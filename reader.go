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
	"io/ioutil"

	"log"
)

// ErrInvalidFormat is the error returned by the NextMessage method of Reader if
// its content is malformed in a way that it is not possible to extract a
// message.
var ErrInvalidFormat = errors.New("invalid mbox format")

type messageReader struct {
	s                  *bufio.Scanner
	next               bytes.Buffer
	atEOF, atSeparator bool
}

func (mr *messageReader) readLine() ([]byte, error) {
	if !mr.s.Scan() {
		log.Println("EOF")
		mr.atEOF = true
		if err := mr.s.Err(); err != nil {
			return nil, err
		}
		return nil, io.EOF
	}
	return mr.s.Bytes(), nil
}

func (mr *messageReader) Read(p []byte) (int, error) {
	if mr.atEOF || mr.atSeparator {
		return 0, io.EOF
	}

	if mr.next.Len() == 0 {
		b, err := mr.readLine()
		if err != nil {
			return 0, err
		}

		if bytes.HasPrefix(b, header) {
			mr.atSeparator = true
			return 0, io.EOF
		} else if len(b) == 0 {
			// Empty line, maybe the next one will contain a header
			if b, err = mr.readLine(); err != nil {
				return 0, err
			}

			if bytes.HasPrefix(b, header) {
				mr.atSeparator = true
				return 0, io.EOF
			}

			mr.next.Write([]byte("\r\n"))
		}

		escapedHeader := append([]byte{'>'}, header...)
		if len(b) > 0 && bytes.HasPrefix(b, escapedHeader) {
			b = b[1:]
		}

		mr.next.Write(b)
		mr.next.Write([]byte("\r\n"))
	}

	return mr.next.Read(p)
}

// Reader reads an mbox archive.
type Reader struct {
	s  *bufio.Scanner
	mr *messageReader
}

// NewReader returns a new Reader to read messages from mbox file format data
// provided by io.Reader r.
func NewReader(r io.Reader) *Reader {
	return &Reader{s: bufio.NewScanner(r)}
}

// NextMessage returns the next message text (containing both the header and the
// body). It will return io.EOF if there are no messages left.
func (r *Reader) NextMessage() (io.Reader, error) {
	if r.mr == nil {
		for r.s.Scan() {
			b := r.s.Bytes()
			if len(b) == 0 {
				continue
			}
			if bytes.HasPrefix(b, header) {
				break
			} else {
				return nil, ErrInvalidFormat
			}
		}
		if err := r.s.Err(); err != nil {
			return nil, err
		}
	} else {
		if _, err := io.Copy(ioutil.Discard, r.mr); err != nil {
			return nil, err
		}
		if r.mr.atEOF {
			return nil, io.EOF
		}
	}
	r.mr = &messageReader{s: r.s}
	return r.mr, nil
}
