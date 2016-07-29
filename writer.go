// "THE BEER-WARE LICENSE" (Revision 42):
// <contact@emersion.fr> wrote this file. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you
// think this stuff is worth it, you can buy me a beer in return.
//                                                                 emersion

package mbox

import (
	"io"
	"io/ioutil"
	"net/mail"
	"net/textproto"
	"strings"
	"time"
)

// Write a MIME header.
func writeMIMEHeader(w io.Writer, header textproto.MIMEHeader) (N int, err error) {
	var n int

	for name, values := range header {
		for _, value := range values {
			n, err = io.WriteString(w, name + ": " + value + "\r\n")
			N += n
			if err != nil {
				return
			}
		}
	}

	n, err = io.WriteString(w, "\r\n")
	N += n
	return
}

// Writer writes messages to a mbox stream.
type Writer struct {
	w io.Writer
}

// NewWriter creates a new *Writer that writes messages to w.
func NewWriter(w io.Writer) *Writer {
	return &Writer{w}
}

// WriteMessage writes a message to the mbox stream. It returns the number of
// bytes written.
func (w *Writer) WriteMessage(m *mail.Message) (N int, err error) {
	from := "???@???"
	if fromList, err := m.Header.AddressList("From"); err == nil && len(fromList) > 0 {
		from = fromList[0].Address
	}

	date := ""
	if t, err := m.Header.Date(); err == nil {
		date = t.Format(time.ANSIC)
	}

	line := "From " + from + " " + date + "\r\n"
	n, err := io.WriteString(w.w, line)
	N += n
	if err != nil {
		return
	}

	n, err = writeMIMEHeader(w.w, textproto.MIMEHeader(m.Header))
	N += n
	if err != nil {
		return
	}

	// Escape lines begining with "From "
	// TODO: use golang.org/x/text/transform
	b, err := ioutil.ReadAll(m.Body)
	if err != nil {
		return
	}

	r := strings.NewReplacer("\nFrom ", "\n>From ")
	n, err = r.WriteString(w.w, string(b))
	N += n
	if err != nil {
		return
	}

	n, err = io.WriteString(w.w, "\r\n\r\n")
	N += n
	return
}
