package syntax

import (
	"fmt"
)

// --- Memory Frames (Call Stack) ---------------------------------------

type DynamicMemoryFrame struct {
	Name        string
	Scope       *Scope
	SymbolTable *SymbolTable
	Parent      *DynamicMemoryFrame
}

func NewDynamicMemoryFrame(nm string, scope *Scope) *DynamicMemoryFrame {
	mf := &DynamicMemoryFrame{
		Name:  nm,
		Scope: scope,
	}
	return mf
}

/* Prettyfied Stringer.
 */
func (mf *DynamicMemoryFrame) String() string {
	return fmt.Sprintf("<mem %s -> %v>", mf.Name, mf.Scope)
}

func (mf *DynamicMemoryFrame) GetName() string {
	return mf.Name
}

func (mf *DynamicMemoryFrame) GetParent() *DynamicMemoryFrame {
	return mf.Parent
}

func (mf *DynamicMemoryFrame) Symbols() *SymbolTable {
	return mf.SymbolTable
}

func (mf *DynamicMemoryFrame) GetScope() *Scope {
	return mf.Scope
}

func (mf *DynamicMemoryFrame) IsRoot() bool {
	return (mf.Parent == nil)
}

// ---------------------------------------------------------------------------

type MemoryFrameStack struct {
	memoryFrameBase *DynamicMemoryFrame
	memoryFrameTOS  *DynamicMemoryFrame
}

/* Get the current memory frame of a stack (TOS).
 */
func (mfst *MemoryFrameStack) Current() *DynamicMemoryFrame {
	if mfst.memoryFrameTOS == nil {
		panic("attempt to access memory frame from empty stack")
	}
	return mfst.memoryFrameTOS
}

/* Get the outermost memory frame, containing global symbols.
 */
func (mfst *MemoryFrameStack) Globals() *DynamicMemoryFrame {
	if mfst.memoryFrameBase == nil {
		panic("attempt to access global memory frame from empty stack")
	}
	return mfst.memoryFrameBase
}

/* Push a memory frame. A frame is constructed, having the recent TOS as its
 * parent. If the new frame is not the bottommost frame, it will copy the
 * symbol-creator from the parent frame. Otherwise callers will have to provide
 * a scope (if needed) in a separate step.
 */
func (mfst *MemoryFrameStack) PushNewMemoryFrame(nm string, scope *Scope) *DynamicMemoryFrame {
	mfp := mfst.memoryFrameTOS
	newmf := NewDynamicMemoryFrame(nm, scope)
	newmf.Parent = mfp
	if mfp == nil { // the new frame is the global frame
		mfst.memoryFrameBase = newmf // make new mf anchor
	} else {
		symcreator := mfp.SymbolTable.GetSymbolCreator()
		symtab := NewSymbolTable(symcreator)
		newmf.SymbolTable = symtab
	}
	mfst.memoryFrameTOS = newmf // new frame now TOS
	T.P("mem", newmf.Name).Debug("pushing new memory frame")
	return newmf
}

/* Pop a memory frame. Returns the popped frame.
 */
func (mfst *MemoryFrameStack) PopMemoryFrame() *DynamicMemoryFrame {
	if mfst.memoryFrameTOS == nil {
		panic("attempt to pop memory frame from empty call stack")
	}
	mf := mfst.memoryFrameTOS
	T.Debugf("popping memory frame [%s]", mf.Name)
	mfst.memoryFrameTOS = mfst.memoryFrameTOS.Parent
	return mf
}

/* Find the top-most memory frame pointing to scope.
 */
func (mfst *MemoryFrameStack) FindMemoryFrameWithScope(scope *Scope) *DynamicMemoryFrame {
	mf := mfst.memoryFrameTOS
	for mf != nil {
		if mf.Scope == scope {
			return mf
		}
		mf = mf.Parent
	}
	return nil
}
