package fp

import "github.com/npillmayer/gotype/syntax/terex"

type ListSeq struct {
	list *terex.GCons
	seq  ListGenerator
}

func (seq *ListSeq) Break() {
	seq.seq = nil
}

func (seq *ListSeq) Done() bool {
	return seq.seq == nil
}

func (seq ListSeq) First() (terex.Atom, ListSeq) {
	return seq.list.Car, seq
}

func (seq *ListSeq) Next() terex.Atom {
	if seq.Done() {
		return terex.NilAtom
	}
	next := seq.seq()
	seq.list = next.list
	seq.seq = next.seq
	return seq.list.Car
}

type ListGenerator func() ListSeq

// --- Trees -----------------------------------------------------------------

type TreeSeq struct {
	list terex.GCons
	seq  TreeGenerator
}

func (seq *TreeSeq) Break() {
	seq.seq = nil
}

func (seq *TreeSeq) Done() bool {
	return seq.seq == nil
}

func (seq TreeSeq) First() (terex.Atom, TreeSeq) {
	return seq.list.Car, seq
}

func (seq *TreeSeq) Next() terex.Atom {
	if seq.Done() {
		return terex.NilAtom
	}
	next := seq.seq()
	seq.list = next.list
	seq.seq = next.seq
	return seq.list.Car
}

type TreeGenerator func() TreeSeq
