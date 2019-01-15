// Copyright (C) 2011, Ross Light

package pdfapi

import (
	"bytes"
	"reflect"
	"testing"
)

const encodingTestData = "%PDF-1.7\n" +
	"%\x93\x8c\x8b\x9e\n" +
	"1 0 obj\n" +
	"(Hello, World!)\n" +
	"endobj\n" +
	"2 0 obj\n" +
	"42\n" +
	"endobj\n" +
	"xref\n" +
	"0 3\n" +
	"0000000000 65535 f\n" +
	"0000000015 00000 n\n" +
	"0000000046 00000 n\n" +
	"trailer\n" +
	"<< /Size 3 /Root 0 0 R >>\n" +
	"startxref\n" +
	"64\n" +
	"%%EOF\n"

func TestEncoder(t *testing.T) {
	var e encoder
	if ref := e.add("Hello, World!"); !reflect.DeepEqual(ref, Reference{1, 0}) {
		t.Errorf("After adding first object, reference is %#v", ref)
	}
	if ref := e.add(42); !reflect.DeepEqual(ref, Reference{2, 0}) {
		t.Errorf("After adding second object, reference is %#v", ref)
	}

	var b bytes.Buffer
	if err := e.encode(&b); err != nil {
		t.Fatalf("Encoding error: %v", err)
	}
	if b.String() != encodingTestData {
		t.Errorf("Encoding result\n%q, want\n%q", b.String(), encodingTestData)
	}
}
