package message

import (
	"bufio"
	"bytes"
	"github.com/sloonz/go-qprintable"
	"net/textproto"
	"testing"
)

const MESSAGE = "Lorem ipsum dolor sit amet, consectetur adipiscing " +
	"elit. Ut elit eros, viverra in laoreet nec, iaculis in libero. Duis " +
	"id velit quis enim lobortis bibendum. Sed purus nisl, luctus et " +
	"pharetra vel, condimentum quis ipsum. Ut scelerisque molestie ligula id " +
	"aliquet. Curabitur fringilla metus eu dui consectetur aliquet. Aenean " +
	"euismod eros tempor libero facilisis vitae rutrum arcu ultricies. Fusce " +
	"vel elit sit amet velit lobortis scelerisque vel nec orci. Ut eu sapien " +
	"quis magna imperdiet mattis sed sit amet diam. Maecenas id ipsum erat, " +
	"eu tristique dolor. Pellentesque aliquet mi eu quam sodales cursus. Nulla " +
	"erat risus, interdum vitae commodo sed, viverra in risus. Nam molestie " +
	"cursus neque, ut venenatis nibh fringilla quis.\n\n" +
	"Aliquam in sem neque. Nullam scelerisque ligula porttitor nunc semper " +
	"scelerisque. Proin urna diam, consequat quis accumsan in, suscipit a " +
	"diam. In laoreet interdum nunc, et fringilla arcu volutpat varius. Sed " +
	"lorem odio, sagittis vel iaculis congue, convallis non tellus. Suspendisse " +
	"consectetur aliquam feugiat. Quisque luctus sollicitudin eros in " +
	"tempor. Suspendisse sit amet risus urna, fringilla tempus nibh. Praesent " +
	"aliquam euismod erat ac congue. Phasellus neque nibh, sodales vitae " +
	"tincidunt et, blandit a dui."

const MESSAGE_QENCODED = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Ut elit eros, vive=\r\n" +
	"rra in laoreet nec, iaculis in libero. Duis id velit quis enim lobortis bib=\r\n" +
	"endum. Sed purus nisl, luctus et pharetra vel, condimentum quis ipsum. Ut s=\r\n" +
	"celerisque molestie ligula id aliquet. Curabitur fringilla metus eu dui con=\r\n" +
	"sectetur aliquet. Aenean euismod eros tempor libero facilisis vitae rutrum =\r\n" +
	"arcu ultricies. Fusce vel elit sit amet velit lobortis scelerisque vel nec =\r\n" +
	"orci. Ut eu sapien quis magna imperdiet mattis sed sit amet diam. Maecenas =\r\n" +
	"id ipsum erat, eu tristique dolor. Pellentesque aliquet mi eu quam sodales =\r\n" +
	"cursus. Nulla erat risus, interdum vitae commodo sed, viverra in risus. Nam=\r\n" +
	" molestie cursus neque, ut venenatis nibh fringilla quis.\r\n\r\n" +
	"Aliquam in sem neque. Nullam scelerisque ligula porttitor nunc semper scele=\r\n" +
	"risque. Proin urna diam, consequat quis accumsan in, suscipit a diam. In la=\r\n" +
	"oreet interdum nunc, et fringilla arcu volutpat varius. Sed lorem odio, sag=\r\n" +
	"ittis vel iaculis congue, convallis non tellus. Suspendisse consectetur ali=\r\n" +
	"quam feugiat. Quisque luctus sollicitudin eros in tempor. Suspendisse sit a=\r\n" +
	"met risus urna, fringilla tempus nibh. Praesent aliquam euismod erat ac con=\r\n" +
	"gue. Phasellus neque nibh, sodales vitae tincidunt et, blandit a dui."

const MESSAGE_B64ENCODED = "TG9yZW0gaXBzdW0gZG9sb3Igc2l0IGFtZXQsIGNvbnNlY3RldHVyIGFkaXBpc2NpbmcgZWxpdC4g\r\n" +
	"VXQgZWxpdCBlcm9zLCB2aXZlcnJhIGluIGxhb3JlZXQgbmVjLCBpYWN1bGlzIGluIGxpYmVyby4g\r\n" +
	"RHVpcyBpZCB2ZWxpdCBxdWlzIGVuaW0gbG9ib3J0aXMgYmliZW5kdW0uIFNlZCBwdXJ1cyBuaXNs\r\n" +
	"LCBsdWN0dXMgZXQgcGhhcmV0cmEgdmVsLCBjb25kaW1lbnR1bSBxdWlzIGlwc3VtLiBVdCBzY2Vs\r\n" +
	"ZXJpc3F1ZSBtb2xlc3RpZSBsaWd1bGEgaWQgYWxpcXVldC4gQ3VyYWJpdHVyIGZyaW5naWxsYSBt\r\n" +
	"ZXR1cyBldSBkdWkgY29uc2VjdGV0dXIgYWxpcXVldC4gQWVuZWFuIGV1aXNtb2QgZXJvcyB0ZW1w\r\n" +
	"b3IgbGliZXJvIGZhY2lsaXNpcyB2aXRhZSBydXRydW0gYXJjdSB1bHRyaWNpZXMuIEZ1c2NlIHZl\r\n" +
	"bCBlbGl0IHNpdCBhbWV0IHZlbGl0IGxvYm9ydGlzIHNjZWxlcmlzcXVlIHZlbCBuZWMgb3JjaS4g\r\n" +
	"VXQgZXUgc2FwaWVuIHF1aXMgbWFnbmEgaW1wZXJkaWV0IG1hdHRpcyBzZWQgc2l0IGFtZXQgZGlh\r\n" +
	"bS4gTWFlY2VuYXMgaWQgaXBzdW0gZXJhdCwgZXUgdHJpc3RpcXVlIGRvbG9yLiBQZWxsZW50ZXNx\r\n" +
	"dWUgYWxpcXVldCBtaSBldSBxdWFtIHNvZGFsZXMgY3Vyc3VzLiBOdWxsYSBlcmF0IHJpc3VzLCBp\r\n" +
	"bnRlcmR1bSB2aXRhZSBjb21tb2RvIHNlZCwgdml2ZXJyYSBpbiByaXN1cy4gTmFtIG1vbGVzdGll\r\n" +
	"IGN1cnN1cyBuZXF1ZSwgdXQgdmVuZW5hdGlzIG5pYmggZnJpbmdpbGxhIHF1aXMuCgpBbGlxdWFt\r\n" +
	"IGluIHNlbSBuZXF1ZS4gTnVsbGFtIHNjZWxlcmlzcXVlIGxpZ3VsYSBwb3J0dGl0b3IgbnVuYyBz\r\n" +
	"ZW1wZXIgc2NlbGVyaXNxdWUuIFByb2luIHVybmEgZGlhbSwgY29uc2VxdWF0IHF1aXMgYWNjdW1z\r\n" +
	"YW4gaW4sIHN1c2NpcGl0IGEgZGlhbS4gSW4gbGFvcmVldCBpbnRlcmR1bSBudW5jLCBldCBmcmlu\r\n" +
	"Z2lsbGEgYXJjdSB2b2x1dHBhdCB2YXJpdXMuIFNlZCBsb3JlbSBvZGlvLCBzYWdpdHRpcyB2ZWwg\r\n" +
	"aWFjdWxpcyBjb25ndWUsIGNvbnZhbGxpcyBub24gdGVsbHVzLiBTdXNwZW5kaXNzZSBjb25zZWN0\r\n" +
	"ZXR1ciBhbGlxdWFtIGZldWdpYXQuIFF1aXNxdWUgbHVjdHVzIHNvbGxpY2l0dWRpbiBlcm9zIGlu\r\n" +
	"IHRlbXBvci4gU3VzcGVuZGlzc2Ugc2l0IGFtZXQgcmlzdXMgdXJuYSwgZnJpbmdpbGxhIHRlbXB1\r\n" +
	"cyBuaWJoLiBQcmFlc2VudCBhbGlxdWFtIGV1aXNtb2QgZXJhdCBhYyBjb25ndWUuIFBoYXNlbGx1\r\n" +
	"cyBuZXF1ZSBuaWJoLCBzb2RhbGVzIHZpdGFlIHRpbmNpZHVudCBldCwgYmxhbmRpdCBhIGR1aS4=\r\n"

func TestMessage(t *testing.T) {
	m := NewMultipartMessage("alternative", "")
	m.SetHeader("Subject", EncodeWord("昨日の会議"))
	m.SetHeader("From", EncodeWord("Miller")+" <miller@example.com>")
	m.SetHeader("To", EncodeWord("田中")+" <tanaka@example.com>")
	m1 := NewTextMessage(qprintable.UnixTextEncoding, bytes.NewBufferString(MESSAGE))
	m1.SetHeader("Content-Type", "text/plain")
	m2 := NewBinaryMessage(bytes.NewBufferString(MESSAGE))
	m2.SetHeader("Content-Type", "application/octet-stream")
	m.AddPart(m1)
	m.AddPart(m2)

	expected_headers := map[string]string{
		"Mime-Version": "1.0",
		"Subject":      "=?UTF-8?Q?=E6=98=A8=E6=97=A5=E3=81=AE=E4=BC=9A=E8=AD=B0?=",
		"Content-Type": "multipart/alternative; boundary=\"==GoMultipartBoundary:0.\"",
		"From":         "Miller <miller@example.com>",
		"To":           "=?UTF-8?Q?=E7=94=B0=E4=B8=AD?= <tanaka@example.com>"}
	expected_data := "--==GoMultipartBoundary:0.\r\n" +
		"Content-Transfer-Encoding: quoted-printable\r\n" +
		"Content-Type: text/plain\r\n\r\n" +
		MESSAGE_QENCODED + "\r\n" +
		"--==GoMultipartBoundary:0.\r\n" +
		"Content-Transfer-Encoding: base64\r\n" +
		"Content-Type: application/octet-stream\r\n\r\n" +
		MESSAGE_B64ENCODED + "\r\n" +
		"--==GoMultipartBoundary:0.--\r\n"

	mbuf := bufio.NewReader(m)
	tp := textproto.NewReader(mbuf)
	headers, err := tp.ReadMIMEHeader()
	if err != nil {
		t.Errorf("Can't parse resulting message: %v", err)
	}
	for k, v := range headers {
		expected_v, ok := expected_headers[k]
		if !ok {
			t.Errorf("Unexpected header %s", k)
			return
		}
		if len(v) != 1 {
			t.Errorf("Unexpected multiple-values header %s", k)
			return
		}
		if v[0] != expected_v {
			t.Errorf("%s is %s, expected %s", k, v[0], expected_v)
			return
		}
		delete(expected_headers, k)
	}

	buf := bytes.NewBuffer(nil)
	_, err = buf.ReadFrom(tp.R)
	if err != nil {
		t.Errorf("Can't read body: %v", err)
		return
	}
	data := buf.String()
	if data != expected_data {
		t.Logf("Message is not what was expected !")
		t.Logf("Expected:")
		t.Logf("%#v", expected_data)
		t.Logf("Message:")
		t.Logf("%#v", data)
		t.Fail()
	}
}
