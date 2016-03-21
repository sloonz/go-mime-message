# message
--
    import "github.com/sloonz/go-mime-message"

The package mime/message can be used to procduce MIME messages which can be sent
as mails to a SMTP server, or put in a local mailbox.

## Usage

```go
var (
	TE_7bit       = TransferEncoding("7bit")
	TE_8bit       = TransferEncoding("8bit")
	TE_binary     = TransferEncoding("binary")
	TE_base64     = TransferEncoding("base64")
	TE_qprintable = TransferEncoding("quoted-printable")
)
```

```go
var (
	MultipartInvalidTransferEncoding = Error("multipart messages only support 7bit transfer encoding")
	PartInvalidTransferEncoding      = Error("parts of a multipart message may not use binary or 8bit transfer encoding")
)
```

#### func  EncodeWord

```go
func EncodeWord(w string) string
```
Encode a word according to RFC 2047. You must use this when you set headers
values that does contain non-ascii characters. You must encode only "phrases"
(in the sense of RFC 2822) and not full headers (unless the header is a phrase,
like the Subject: header). For example, if you want to write a mail to Tanaka
(田中), you must use:

    SetHeader("To", EncodeWord("田中") + " <tanaka@example.com>")

and not

    SetHeader("To", EncodeWord("田中 <tanaka@example.com>"))

The phrase is assumed to be valid UTF-8.

#### type Error

```go
type Error string
```


#### func (Error) Error

```go
func (e Error) Error() string
```

#### type Message

```go
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

	// End of line characters. Defaults to CRLF (as required by most standards), but you may
	// want change this to "\n" if you intend to write in a Maildir, which requires LF line
	// endings.
	EOL string

	// The body of the message
	Body io.Reader
}
```


#### func  NewBinaryMessage

```go
func NewBinaryMessage(body io.Reader) *Message
```
New message containing binary data. It will be encoded with base64 encoding. You
should use this for all media types but text/* and multipart/*

#### func  NewTextMessage

```go
func NewTextMessage(qpEncoding *qprintable.Encoding, body io.Reader) *Message
```
New message containing text data. It will be encoded with quoted-printable
encoding. You should use this for text/* media types.

#### func (*Message) Read

```go
func (m *Message) Read(p []byte) (n int, err error)
```
Read the MIME representation of the message (headers + body). You can do this
only once, since after the first representation this will always return os.EOF.
For base64 and quoted-printable encodings, also take care of encoding the body.

#### func (*Message) SetHeader

```go
func (m *Message) SetHeader(name, val string) *Message
```
Set an header. val will be directly written ; to escape it, see EncodeWord.
Returns self.

#### type MultipartMessage

```go
type MultipartMessage struct {
	Message
	Parts    []*Message
	Boundary string
}
```

A multipart message is a messsage containing other messages

#### func  NewMultipartMessage

```go
func NewMultipartMessage(subtype, boundary string) *MultipartMessage
```
Create a new multipart message.

If boundary is empty, a new one will be automatically generated. If you supply
one, you must ensure that it is valid and not taken anywhere else.

You should not modify Body field of the returned structure.

#### func  NewMultipartMessageParams

```go
func NewMultipartMessageParams(subtype, boundary string, params map[string]string) *MultipartMessage
```
Create a new multipart message with additional parameters.

If boundary is empty, a new one will be automatically generated. If you supply
one, you must ensure that it is valid and not taken anywhere else.

Additional parameters (e.g. type for multipart/related) can be supplied. It is
the responsibility of the caller to encode them (atom / quoted-string according
to RFC 2822)

You should not modify Body field of the returned structure.

#### func (*MultipartMessage) AddPart

```go
func (m *MultipartMessage) AddPart(c *Message) *MultipartMessage
```
Add a message to the multipart message. EOL for the part will be inherited from
the multipart message. Returns self.

#### type TransferEncoding

```go
type TransferEncoding string
```
