// "THE BEER-WARE LICENSE" (Revision 42):
// <contact@emersion.fr> wrote this file. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you
// think this stuff is worth it, you can buy me a beer in return.
//                                                                 emersion

package mbox

import (
	"io"
	"net/mail"
	"net/textproto"
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

// Write writes a message to the mbox stream. It returns the number of bytes
// written.
func (w *Writer) Write(m *mail.Message) (N int, err error) {
	from := "???@???"
	if fromList, err := m.Header.AddressList("From"); err == nil && len(fromList) > 0 {
		from = fromList[0].Address
	}

	date := ""
	if t, err := m.Header.Date(); err == nil {
		date = t.Format(time.ANSIC)
	}

	line := "From " + from + " " + date + "\n"
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

	nn, err := io.Copy(w.w, m.Body)
	N += int(nn)
	if err != nil {
		return
	}

	n, err = io.WriteString(w.w, "\n")
	N += n
	return
}
