package mbox

import (
	"bytes"
	"errors"
	"io"
	"time"
)

type messageWriter struct {
	w   io.Writer
	buf bytes.Buffer
}

func (mw *messageWriter) writeLine(l []byte) (int, error) {
	if bytes.HasPrefix(l, header) {
		if _, err := mw.w.Write([]byte{'>'}); err != nil {
			return 0, err
		}
	}

	return mw.w.Write(l)
}

func (mw *messageWriter) Write(p []byte) (int, error) {
	mw.buf.Write(p)
	b := mw.buf.Bytes()
	mw.buf.Reset()

	N := 0
	for {
		i := bytes.IndexByte(b, '\n')
		if i < 0 {
			n, err := mw.buf.Write(b)
			N += n
			return N, err
		}

		var l []byte
		l, b = b[:i+1], b[i+1:]

		n := len(l)
		// Replace CRLF with LF
		if len(l) > 1 && l[len(l)-2] == '\r' {
			l = l[:len(l)-2]
			l = append(l, '\n')
		}

		_, err := mw.writeLine(l)
		N += n
		if err != nil {
			return N, err
		}
	}
}

func (mw *messageWriter) Close() error {
	b := mw.buf.Bytes()
	mw.buf.Reset()
	if _, err := mw.writeLine(b); err != nil {
		return err
	}

	_, err := mw.w.Write([]byte("\n\n"))
	return err
}

// Writer writes messages to a mbox stream. The Close method must be called to
// end the stream.
type Writer struct {
	w      io.Writer
	last   *messageWriter
	closed bool
}

// NewWriter creates a new Writer that writes messages to w.
func NewWriter(w io.Writer) *Writer {
	return &Writer{w: w}
}

// CreateMessage appends a message to the mbox stream. The message text
// (including both the header and the body) should be written to the returned
// io.Writer.
func (w *Writer) CreateMessage(from string, t time.Time) (io.Writer, error) {
	if w.closed {
		return nil, errors.New("mbox: Writer.CreateMessage called after Close")
	}
	if w.last != nil {
		if err := w.last.Close(); err != nil {
			return nil, err
		}
		w.last = nil
	}

	if from == "" {
		from = "???@???"
	}

	if t.IsZero() {
		t = time.Now()
	}
	date := t.UTC().Format(time.ANSIC)

	line := "From " + from + " " + date + "\n"
	if _, err := io.WriteString(w.w, line); err != nil {
		return nil, err
	}

	w.last = &messageWriter{w: w.w}
	return w.last, nil
}

func (w *Writer) Close() error {
	if w.closed {
		return errors.New("mbox: Writer already closed")
	}
	w.closed = true
	if w.last != nil {
		return w.last.Close()
	}
	return nil
}
