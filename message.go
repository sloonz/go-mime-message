// The package mime/message can be used to procduce MIME messages which can be sent
// as mails to a SMTP server, or put in a local mailbox.
package message

import (
	"bytes"
	"encoding/base64"
	"github.com/sloonz/go-qprintable"
	"io"
	"net/http"
)

type Message struct {
	// The transfer encoding for this message. General rules are:
	//  - message/rfc822, message/partial and message/external-body only accept 7bit
	//  - multipart/* only accept 7bit, 8bit and binary
	//  - you should not use "binary" and "8bit", since such messages will not
	//    be conform with SMTP
	//  - for encodings other than base64 and quoted-printable, it is your responsibility
	//    to ensure that the given data conforms to the encoding, and that data does not
	//    contain the multipart boundary in multipart parts
	// If you use NewTextMessage, NewBinaryMessage and NewMultipartMessage, you shouldn't
	// have to worry about this. It is wise not to modify it yourself, since defautlts
	// are standard compliants and works well with multipart messages
	TE TransferEncoding

	// For quoted-printable transfer-encoding, define the canonical form of the body for
	// ends of line encoding. Use BinaryEncoding to avoid any end of line conversion, but
	// please note that this is invalid for text/* entities if you want to be pedantic
	// (in practice, few MUA are perturbated by bad end of lines)
	QPEncoding *qprintable.Encoding

	// Headers of the message. They are stored in the http.CanonicalHeaderKey format.
	// Don't put content-transfer-encoding nor mime-version into this, it will be handled
	// internally.
	Headers map[string]string

	// The body of the message
	Body io.Reader

	isMultipartPart bool
	buf             *bytes.Buffer
	bodyReader      io.Reader
}

// New message containing text data. It will be encoded with quoted-printable encoding.
// You should use this for text/* media types.
func NewTextMessage(qpEncoding *qprintable.Encoding, body io.Reader) *Message {
	m := new(Message)
	m.TE = TE_qprintable
	m.QPEncoding = qpEncoding
	m.Body = body
	m.Headers = make(map[string]string)
	return m
}

// New message containing binary data. It will be encoded with base64 encoding.
// You should use this for all media types but text/* and multipart/*
func NewBinaryMessage(body io.Reader) *Message {
	m := new(Message)
	m.TE = TE_base64
	m.Body = body
	m.Headers = make(map[string]string)
	return m
}

// Set an header. val will be directly written ; to escape it, see EncodeWord.
// Returns self.
func (m *Message) SetHeader(name, val string) *Message {
	m.Headers[http.CanonicalHeaderKey(name)] = val
	return m
}

// Read the MIME representation of the message (headers + body). You can do this 
// only once, since after the first representation this will always return os.EOF.
// For base64 and quoted-printable encodings, also take care of encoding the body.
func (m *Message) Read(p []byte) (n int, err error) {
	// Write message header to buffer on first call
	// TODO: wrap headers ?
	if m.buf == nil {
		m.buf = bytes.NewBuffer(nil)
		if !m.isMultipartPart {
			m.buf.WriteString("MIME-Version: 1.0\r\n")
		}
		if m.TE != TE_7bit {
			if (m.TE == TE_8bit || m.TE == TE_binary) && m.isMultipartPart {
				return n, PartInvalidTransferEncoding
			} else {
				m.buf.WriteString("Content-Transfer-Encoding: " + string(m.TE) + "\r\n")
			}
		}
		for name, val := range m.Headers {
			m.buf.WriteString(name + ": " + val + "\r\n")
		}
		m.buf.WriteString("\r\n")
	}

	// Create body transform (transfer encoding)
	if m.bodyReader == nil {
		buf := bytes.NewBuffer(nil)
		if m.TE == TE_qprintable {
			m.bodyReader = &qprintableReader{m.Body, buf, qprintable.NewEncoder(m.QPEncoding, buf)}
		} else if m.TE == TE_base64 {
			m.bodyReader = &base64Reader{m.Body, buf, base64.NewEncoder(base64.StdEncoding, buf), 0, false}
		} else {
			m.bodyReader = m.Body
		}
	}

	// Main loop
	for len(p) > n && err != io.EOF {
		if m.buf.Len() > 0 {
			nn, _ := m.buf.Read(p[n:])
			n += nn
		} else {
			nn, merr := m.bodyReader.Read(p[n:])
			err = merr
			n += nn
		}
	}

	return n, err
}
