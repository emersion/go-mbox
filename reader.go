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
)

// ErrInvalidFormat is the error returned by the NextMessage method of Reader if
// its content is malformed in a way that it is not possible to extract a
// message.
var ErrInvalidFormat = errors.New("invalid mbox format")

type messageReader struct {
	r                  *bufio.Reader
	next               bytes.Buffer
	atEOF, atSeparator bool
	atMiddleOfLine     bool
}

func (mr *messageReader) Read(p []byte) (int, error) {
	if mr.atEOF || mr.atSeparator {
		return 0, io.EOF
	}

	if mr.next.Len() == 0 {
		b, isPrefix, err := mr.r.ReadLine()
		if err != nil {
			mr.atEOF = true
			return 0, err
		}

		if !mr.atMiddleOfLine {
			if bytes.HasPrefix(b, header) {
				mr.atSeparator = true
				return 0, io.EOF
			} else if len(b) == 0 {
				// Check if the next line is separator. In such case the new
				// line should not be written to not have double new line.
				b, isPrefix, err = mr.r.ReadLine()
				if err != nil {
					mr.atEOF = true
					return 0, err
				}

				if bytes.HasPrefix(b, header) {
					mr.atSeparator = true
					return 0, io.EOF
				}

				mr.next.Write([]byte("\r\n"))
			}

			if bytes.HasPrefix(b, escapedHeader) {
				b = b[1:]
			}
		}

		mr.next.Write(b)
		if !isPrefix {
			mr.next.Write([]byte("\r\n"))
		}
		mr.atMiddleOfLine = isPrefix
	}

	return mr.next.Read(p)
}

// Reader reads an mbox archive.
type Reader struct {
	r  *bufio.Reader
	mr *messageReader
}

// NewReader returns a new Reader to read messages from mbox file format data
// provided by io.Reader r.
func NewReader(r io.Reader) *Reader {
	return &Reader{r: bufio.NewReader(r)}
}

// NextMessage returns the next message text (containing both the header and the
// body). It will return io.EOF if there are no messages left.
func (r *Reader) NextMessage() (io.Reader, error) {
	if r.mr == nil {
		for {
			b, isPrefix, err := r.r.ReadLine()
			if err != nil {
				return nil, err
			}
			// Discard the rest of the line.
			for isPrefix {
				_, isPrefix, err = r.r.ReadLine()
				if err != nil {
					return nil, err
				}
			}
			if len(b) == 0 {
				continue
			}
			if bytes.HasPrefix(b, header) {
				break
			} else {
				return nil, ErrInvalidFormat
			}
		}
	} else {
		if _, err := io.Copy(ioutil.Discard, r.mr); err != nil {
			return nil, err
		}
		if r.mr.atEOF {
			return nil, io.EOF
		}
	}
	r.mr = &messageReader{r: r.r}
	return r.mr, nil
}
