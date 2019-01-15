package pdfapi

import (
	"compress/lzw"
	"compress/zlib"
	"io/ioutil"
	"testing"
)

const streamTestString = "Hello, this is a PDF API!\n"

func TestUnfilteredStream(t *testing.T) {
	st := newStream(streamNoFilter)
	st.WriteString(streamTestString)
	st.Close()

	if st.String() != streamTestString {
		t.Errorf("Stream is %q, wanted %q", st.String(), streamTestString)
	}
}

func TestLZWStream(t *testing.T) {
	st := newStream(streamLZWDecode)
	st.WriteString(streamTestString)
	st.Close()

	output, _ := ioutil.ReadAll(lzw.NewReader(st, lzw.MSB, 8))
	if string(output) != streamTestString {
		t.Errorf("Stream is %q, wanted %q", output, streamTestString)
	}
}

func TestFlateStream(t *testing.T) {
	st := newStream(streamFlateDecode)
	st.WriteString(streamTestString)
	st.Close()

	r, _ := zlib.NewReader(st)
	output, _ := ioutil.ReadAll(r)
	if string(output) != streamTestString {
		t.Errorf("Stream is %q, wanted %q", output, streamTestString)
	}
}

const expectedMarshalStreamOutput = "<< /Length 26 >> stream" + Newline + streamTestString + Newline + "endstream"

func TestMarshalStream(t *testing.T) {
	t.Logf("marshal(stream): length of data is %d", len(streamTestString))
	b, err := marshalStream(nil, streamInfo{Length: len(streamTestString)}, []byte(streamTestString))
	if err == nil {
		if string(b) != expectedMarshalStreamOutput {
			t.Errorf("marshalStream(...) != %q (got %q)", expectedMarshalStreamOutput, b)
		}
	} else {
		t.Errorf("Error: %v", err)
	}
}
