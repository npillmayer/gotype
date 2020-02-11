package fn

type IntSeq struct {
	n   int64
	seq IntGenerator
}

func (iseq *IntSeq) Break() {
	iseq.seq = nil
}

func (iseq *IntSeq) Done() bool {
	return iseq.seq == nil
}

// TODO for testing only
func (iseq *IntSeq) N() int64 {
	return iseq.n
}

func (iseq IntSeq) First() (int64, IntSeq) {
	//n := iseq.n
	//seq := iseq.seq()
	//seq := iseq
	return iseq.n, iseq
}

func (iseq *IntSeq) Next() int64 {
	//n := iseq.n
	if iseq.Done() {
		return iseq.n // no possibility to return in-band error
	}
	next := iseq.seq()
	iseq.n = next.n
	iseq.seq = next.seq
	return iseq.n
}

//type IntGenerator func() (int64, IntGenerator)
type IntGenerator func() IntSeq

func N() IntSeq {
	var n int64
	var N IntGenerator
	N = func() IntSeq {
		n++
		return IntSeq{n, N}
	}
	return IntSeq{n, N}
}

type FloatSeq struct {
	n   float64
	seq FloatGenerator
}

func (rseq FloatSeq) First() (float64, FloatSeq) {
	n := rseq.n
	seq := rseq.seq()
	return n, seq
}

func (rseq *FloatSeq) Next() float64 {
	n := rseq.n
	next := rseq.seq()
	rseq.n = next.n
	rseq.seq = next.seq
	return n
}

type FloatGenerator func() FloatSeq

func R() FloatSeq {
	var x float64
	var R FloatGenerator
	R = func() FloatSeq {
		x += 1.0
		return FloatSeq{x, R}
	}
	return FloatSeq{x, R}
}

type IntFilter func(n int64) bool

func LessThanN(b int64) IntFilter {
	return func(n int64) bool {
		return n < b
	}
}

func EvenN() IntFilter {
	return func(n int64) bool {
		return n%2 == 0
	}
}

func (seq IntSeq) Where(filt IntFilter) IntSeq {
	var F IntGenerator
	//inner := seq
	//n, inner := seq.First()
	n, inner := seq.n, seq
	F = func() IntSeq {
		//fmt.Printf("F  called, n=%d\n", n)
		n = inner.Next()
		for !filt(n) {
			//fmt.Printf("   skip n=%d\n", n)
			n = inner.Next()
		}
		//fmt.Printf("F' n=%d\n", n)
		return IntSeq{n, F}
	}
	return IntSeq{n, F}
}

type IntMapper func(n int64) int64

func SquareN() IntMapper {
	return func(n int64) int64 {
		return n * n
	}
}

func (seq IntSeq) Map(mapper IntMapper) IntSeq {
	var F IntGenerator
	//inner := seq
	n, inner := seq.n, seq
	//n, inner := seq.First()
	v := mapper(n)
	F = func() IntSeq {
		//fmt.Printf("F  called, n=%d\n", n)
		n = inner.Next()
		v = mapper(n)
		//fmt.Printf("F' n=%d, v=%d\n", n, v)
		return IntSeq{v, F}
	}
	return IntSeq{v, F}
}

func (seq IntSeq) Vec() []int64 {
	return nil
}
