package gtcore

type Kern struct {
	w int
}

type Glue struct {
	w       int
	shrink  int
	stretch int
}

type BreakOpportunity struct {
	nobreak int
	before  int
	after   int
}

type Box struct {
	w int
	h int
	d int
}
