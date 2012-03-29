package message

import (
	"bytes"
	"io"
)

type qprintableReader struct {
	body    io.Reader
	buf     *bytes.Buffer
	encoder io.Writer
}

func (r *qprintableReader) Read(p []byte) (n int, err error) {
	for len(p) > n {
		// Take unencoded data from body and put it to encoder (which will put encoded data in buffer)
		if (len(p) - n) > r.buf.Len() {
			_, cerr := io.CopyN(r.encoder, r.body, int64(len(p)-n-r.buf.Len()))
			if (cerr != nil && cerr != io.EOF) || (cerr == io.EOF && r.buf.Len() == 0) {
				return n, cerr
			}
		}

		// Take data from buffer
		nn, _ := r.buf.Read(p[n:])
		n += nn
	}
	return n, nil
}

type base64Reader struct {
	body          io.Reader
	buf           *bytes.Buffer
	encoder       io.WriteCloser
	lineSize      int
	shouldWriteLF bool
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (r *base64Reader) Read(p []byte) (n int, err error) {
	if r.shouldWriteLF && len(p) > 0 {
		n = 1
		p[0] = byte('\n')
		r.shouldWriteLF = false
	}

	for len(p) > n {
		// Manage wraping
		if r.lineSize == maxLineSize {
			r.lineSize = 0
			p[n] = byte('\r')
			n++
			if len(p) > n {
				p[n] = byte('\n')
				n++
			} else {
				r.shouldWriteLF = true
				return n, err
			}
		}

		// Take unencoded data from body and put it to encoder (which will put encoded data in buffer)
		if (len(p) - n) > r.buf.Len() {
			_, cerr := io.CopyN(r.encoder, r.body, int64(len(p)-n-r.buf.Len()))
			if cerr == io.EOF {
				r.encoder.Close()
			}
			if (cerr != nil && cerr != io.EOF) || (cerr == io.EOF && r.buf.Len() == 0) {
				return n, cerr
			}
		}

		// Take data from buffer
		toread := min(maxLineSize-r.lineSize, len(p)-n)
		nn, _ := r.buf.Read(p[n : n+toread])
		r.lineSize += nn
		n += nn
	}

	return n, nil
}
