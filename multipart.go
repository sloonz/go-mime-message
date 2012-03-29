package message

import (
	"bytes"
	"fmt"
	"io"
	"sync"
)

// A multipart message is a messsage containing other messages
type MultipartMessage struct {
	Message
	Parts    []*Message
	Boundary string
}

var boundaryGenerator chan string
var boundaryGeneratorInit sync.Once

// Create a new multipart message. If you supply the boundary yourself, you must
// ensure that it is valid and not taken anywhere else.
// You should not modify Body field of the returned structure.
func NewMultipartMessage(subtype, boundary string) *MultipartMessage {
	boundaryGeneratorInit.Do(func() {
		boundaryGenerator = make(chan string)
		go (func() {
			for i := uint(0); true; i++ {
				// ==G Can't appear in quoted-printable and base64 streams
				boundaryGenerator <- fmt.Sprintf("==GoMultipartBoundary:%d.", i)
			}
		})()
	})
	if boundary == "" {
		boundary = <-boundaryGenerator
	}
	m := new(MultipartMessage)
	m.TE = TE_7bit
	m.Headers = make(map[string]string)
	m.Body = &multipartReader{m, -1, bytes.NewBuffer(nil)}
	m.SetHeader("Content-Type", "multipart/"+subtype+"; boundary=\""+boundary+"\"")
	m.Boundary = boundary
	return m
}

// Add a message to the multipart message.
// Returns self.
func (m *MultipartMessage) AddPart(c *Message) *MultipartMessage {
	m.Parts = append(m.Parts, c)
	c.isMultipartPart = true
	return m
}

type multipartReader struct {
	m   *MultipartMessage
	cur int
	buf *bytes.Buffer
}

func (r *multipartReader) Read(p []byte) (n int, err error) {
	if r.m.TE != TE_7bit && r.m.TE != TE_8bit && r.m.TE != TE_binary {
		return 0, MultipartInvalidTransferEncoding
	}

	if r.cur == -1 {
		r.cur = 0
		r.buf.WriteString("--")
		r.buf.WriteString(r.m.Boundary)
		r.buf.WriteString("\r\n")
	}

	for len(p) > n {
		if r.buf.Len() > 0 {
			nn, _ := r.buf.Read(p[n:])
			n += nn
		}
		if r.cur < len(r.m.Parts) {
			nn, merr := r.m.Parts[r.cur].Read(p[n:])
			n += nn
			if merr != nil && merr != io.EOF {
				return n, merr
			}
			if merr == io.EOF {
				r.cur++
				r.buf.WriteString("\r\n--")
				r.buf.WriteString(r.m.Boundary)
				if r.cur < len(r.m.Parts) {
					r.buf.WriteString("\r\n")
				} else {
					r.buf.WriteString("--\r\n")
				}
			}
		} else {
			return n, io.EOF
		}
	}

	return n, err
}
