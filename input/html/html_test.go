package html

import (
	"os"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestParseFile1(t *testing.T) {
	f, err := os.Open("small.html")
	defer f.Close()
	if err != nil {
		t.Errorf("Error opening HTML file: %s", err)
	} else {
		dom, htmlerr := ReadHTMLBook(f)
		if htmlerr != nil {
			t.Errorf("Error parsing HTML file: %s", htmlerr)
		} else {
			txt := ""
			h1 := dom.Find("h1")
			h1.Each(func(i int, s *goquery.Selection) {
				txt += s.First().Text()
			})
			if strings.Compare(txt, "My First HeadingMy Second Heading") != 0 {
				t.Error("Did not parse H1 elements correctly")
			}
		}
	}
}

func TestParseFile2(t *testing.T) {
	f, err := os.Open("htmlbook.html")
	defer f.Close()
	if err != nil {
		t.Errorf("Error opening HTML file: %s", err)
	} else {
		dom, htmlerr := ReadHTMLBook(f)
		if htmlerr != nil {
			t.Errorf("Error parsing HTML file: %s", htmlerr)
		} else {
			h1 := dom.Find("h1")
			h1.Each(func(i int, s *goquery.Selection) {
				txt := s.First().Text()
				t.Logf("H1 = %s", txt)
			})
		}
	}
}
