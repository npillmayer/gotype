/*
Package font is for typeface and font handling.

There is a certain confusion in the nomenclature of typesetting. We will
stick to the following definitions:

* A "typeface" is a family of fonts. An example is "Helvetica".
This corresponds to a TrueType "collection" (*.ttc).

* A "scalable font" is a font, i.e. a variant of a typeface with a
certain weight, slant, etc.  An example is "Helvetica regular".

* A "typecase" is a scaled font, i.e. a font in a certain size for
a certain script and language. The name is reminiscend on the wooden
boxes of typesetters in the aera of metal type.
An example is "Helvetica regular 11pt, Latin, en_US".

Please not that Go (Golang) does use the terms "font" and "face"
differently–actually more or less in an opposite manner.

TODO: font collections (*.ttc), e.g., /System/Library/Fonts/Helvetica.ttc

Utility to view a character map of a font: http://torinak.com/font/lsfont.html

Website for fonts:
https://www.fontsquirrel.com/fonts/list/popular

OpenType explained:
https://docs.microsoft.com/en-us/typography/opentype/

----------------------------------------------------------------------

BSD License

Copyright (c) 2017-20, Norbert Pillmayer

All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions
are met:

1. Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright
notice, this list of conditions and the following disclaimer in the
documentation and/or other materials provided with the distribution.

3. Neither the name of this software nor the names of its contributors
may be used to endorse or promote products derived from this software
without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE. */
package font

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
)

type ScalableFont struct {
	Fontname string
	Filepath string     // file path
	Binary   []byte     // raw data
	SFNT     *sfnt.Font // the font's container // TODO: not threadsafe???
}

type TypeCase struct {
	scalableFontParent *ScalableFont
	font               font.Face // Go uses 'face' and 'font' in an inverse manner
	size               float64
	// script
	// language
}

func LoadOpenTypeFont(fontfile string) (*ScalableFont, error) {
	ff := &ScalableFont{}
	ff.Fontname = filepath.Base(fontfile)
	ff.Filepath = fontfile
	bytez, err := ioutil.ReadFile(fontfile)
	ff.Binary = bytez
	if err != nil {
		panic(fmt.Sprintf("cannot read font file: %s", fontfile))
	} else {
		reader := bytes.NewReader(ff.Binary)
		ff.SFNT, err = sfnt.ParseReaderAt(reader)
		if err != nil {
			ff.SFNT = nil
			panic(fmt.Sprintf("cannot parse font file: %s", fontfile))
		} else {
			ff.Fontname, _ = ff.SFNT.Name(nil, sfnt.NameIDFull)
		}
	}
	return ff, err
}

// TODO: check if language fits to script
// TODO: check if font suports script
func (sf *ScalableFont) PrepareCase(fontsize float64) (*TypeCase, error) {
	typecase := &TypeCase{}
	typecase.scalableFontParent = sf
	if fontsize < 5.0 || fontsize > 500.0 {
		fmt.Printf("*** font size must be 5pt < size < 500pt, is %g (set to 10pt)\n", fontsize)
		fontsize = 10.0
	}
	options := &opentype.FaceOptions{
		Size: fontsize,
		DPI:  600,
	}
	f, err := opentype.NewFace(sf.SFNT, options)
	if err == nil {
		typecase.font = f
		typecase.size = fontsize
	}
	return typecase, err
}

func (tc *TypeCase) ScalableFontParent() *ScalableFont {
	return tc.scalableFontParent
}

func (tc *TypeCase) PtSize() float64 {
	return tc.size
}
