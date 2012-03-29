package message

import (
	"testing"
)

var testData = []struct{ decoded, encoded string }{
	{"test", "test"},
	{"Bonjour à tous!", "=?UTF-8?Q?Bonjour_=C3=A0_tous!?="},
	{"Right?", "=?UTF-8?Q?Right=3F?="},
	{"田中", "=?UTF-8?Q?=E7=94=B0=E4=B8=AD?="},
}

func TestEncodeWord(t *testing.T) {
	for _, data := range testData {
		if EncodeWord(data.decoded) != data.encoded {
			t.Errorf("EncodeWord(%#v) should be %#v, got %#v", data.decoded, data.encoded, EncodeWord(data.decoded))
		}
	}
}
