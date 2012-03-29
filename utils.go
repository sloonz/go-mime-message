package message

import (
	"bytes"
)

// Enforced only for base64 and quoted-printable. No limit for binary.
const maxLineSize = 76

/**
 * Transfer encodings
 */

type TransferEncoding string

var (
	TE_7bit       = TransferEncoding("7bit")
	TE_8bit       = TransferEncoding("8bit")
	TE_binary     = TransferEncoding("binary")
	TE_base64     = TransferEncoding("base64")
	TE_qprintable = TransferEncoding("quoted-printable")
)

/**
 * Errors
 */

type Error string

func (e Error) Error() string {
	return string(e)
}

var (
	MultipartInvalidTransferEncoding = Error("multipart messages only support 7bit transfer encoding")
	PartInvalidTransferEncoding      = Error("parts of a multipart message may not use binary or 8bit transfer encoding")
)

/**
 * Encoded-word
 */
const hexTable = "0123456789ABCDEF"
const acceptableSpecialChars = "!*+-/="

func isAcceptable(b byte) bool {
	return (b >= '0' && b <= '9') ||
		(b >= 'a' && b <= 'z') ||
		(b >= 'A' && b <= 'Z') ||
		bytes.IndexByte([]byte(acceptableSpecialChars), b) != -1
}

// Encode a word according to RFC 2047.
// You must use this when you set headers values that does contain non-ascii characters.
// You must encode only "phrases" (in the sense of RFC 2822) and not full headers (unless
// the header is a phrase, like the Subject: header). For example, if you want to write a
// mail to Tanaka (田中), you must use:
//   SetHeader("To", EncodeWord("田中") + " <tanaka@example.com>")
// and not
//   SetHeader("To", EncodeWord("田中 <tanaka@example.com>"))
// The phrase is assumed to be valid UTF-8.
func EncodeWord(w string) string {
	// If it's ascii, no need to encode it (more readable)
	ascii := true
	for i := 0; i < len(w) && ascii; i++ {
		if !isAcceptable(w[i]) {
			ascii = false
		}
	}
	if ascii {
		return w
	}

	// Else, encode the word
	buf := bytes.NewBufferString("=?UTF-8?Q?")
	for i := 0; i < len(w); i++ {
		if isAcceptable(w[i]) {
			buf.WriteByte(w[i])
		} else if w[i] == byte(' ') {
			buf.WriteByte('_')
		} else {
			buf.Write([]byte{'=', hexTable[w[i]>>4], hexTable[w[i]&0xf]})
		}
	}
	buf.WriteString("?=")
	return buf.String()
}
