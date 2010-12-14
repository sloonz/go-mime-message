# PACKAGE

The package [mime/message](http://github.com/sloonz/go-mime-message) can be used to procduce MIME messages which can be sent
as mails to a SMTP server, or put in a local mailbox.

It requires the [qprintable](http://github.com/sloonz/go-qprintable) package.

## VARIABLES

	var (
	    TE_7bit       = TransferEncoding("7bit")
	    TE_8bit       = TransferEncoding("8bit")
	    TE_binary     = TransferEncoding("binary")
	    TE_base64     = TransferEncoding("base64")
	    TE_qprintable = TransferEncoding("quoted-printable")
	)
	
	var (
	    MultipartInvalidTransferEncoding = Error("multipart messages only support 7bit transfer encoding")
	    PartInvalidTransferEncoding      = Error("parts of a multipart message may not use binary or 8bit transfer encoding")
	)


## FUNCTIONS

`func EncodeWord(w string) string`

Encode a word according to RFC 2047.
You must use this when you set headers values that does contain non-ascii characters.
You must encode only "phrases" (in the sense of RFC 2822) and not full headers (unless
the header is a phrase, like the Subject: header). For example, if you want to write a
mail to Tanaka (田中), you must use:

	SetHeader("To", EncodeWord("田中") + " <tanaka@example.com>")

and not

	SetHeader("To", EncodeWord("田中 <tanaka@example.com>"))

The phrase is assumed to be valid UTF-8.


## TYPES

	type Error string

	func (e Error) String() string

### Message

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
	    // contains unexported fields
	}

`func NewBinaryMessage(body io.Reader) *Message`

New message containing binary data. It will be encoded with base64 encoding.
You should use this for all media types but text/\* and multipart/\*

`func NewTextMessage(qpEncoding *qprintable.Encoding, body io.Reader) *Message`

New message containing text data. It will be encoded with quoted-printable encoding.
You should use this for text/\* media types.

`func (m *Message) Read(p []byte) (n int, err os.Error)`

Read the MIME representation of the message (headers + body). You can do this
only once, since after the first representation this will always return os.EOF.
For base64 and quoted-printable encodings, also take care of encoding the body.

`func (m *Message) SetHeader(name, val string) *Message`

Set an header. val will be directly written ; to escape it, see EncodeWord.
Returns self.

### MultipartMessage

	type MultipartMessage struct {
	    Message
	    Parts    []*Message
	    Boundary string
	}

A multipart message is a messsage containing other messages

`func NewMultipartMessage(subtype, boundary string) *MultipartMessage`

Create a new multipart message. If you supply the boundary yourself, you must
ensure that it is valid and not taken anywhere else.
You should not modify Body field of the returned structure.

`func (m *MultipartMessage) AddPart(c *Message) *MultipartMessage`

Add a message to the multipart message.
Returns self.

type TransferEncoding string

# BUGS, MISSING FEATURES

No bug known. If you find one, you can fill bug reports on the Github
project home, or submit patches.

Headers handling is somewhat austere. But I don’t see how to improve
this without parsing the whole RFC 2822/2047 BNF, which seems rather
painful and not worthing it.
