/*
Package pdfapi will implement a low-level API for the PDF format.

The original package is by Ross Light (zombiezen). I copied it, rather than
importing it, because I intend to write a low level PDF API from scratch and
will use Ross's code as a loose reference.

There are other PDF packages around (most notably gofpdf
(/github.com/jung-kurt/gofpdf) and gopdf (/github.com/signintech/gopdf).
However, for my purpose these a often too high-level.
Ross's package is the most concise of all and has a Go-like API (in contrast
to, e.g., gofpdf).

My focus for the PDF API will be on concurrency, Unicode typesetting
and the usage of PDF fragment templates.

Status

This is just a bag of ideas. Nothing useful for anyone else yet.

License

The original code has a BSD 2-clause license.
Changes are published under a BSD 3-clause licencse. Please
refer to the license file for details.

PDF Format

The PDF format can roughly be broken into four parts:

	- Header
	- List of Objects (with Streams)
	- XRef
	- Trailer

This package will stick closely to the nomenclature of PDF where possible.
*/
package pdfapi
