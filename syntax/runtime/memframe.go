package runtime

import (
	"fmt"
)

/*
----------------------------------------------------------------------

BSD License
Copyright (c) 2017, Norbert Pillmayer

All rights reserved.
Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions
are met:
1. Redistributions of source code must retain the above copyright
   notice, this list of conditions and the following disclaimer.
2. Redistributions in binary form must reproduce the above copyright
   notice, this list of conditions and the following disclaimer in the
   documentation and/or other materials provided with the distribution.
3. Neither the name of Norbert Pillmayer or the names of its contributors
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
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

----------------------------------------------------------------------

 * This module implements a stack of memory frames.
 * Memory frames are used by an interpreter to allocate local storage
 * for active scopes.

*/

// A memory frame, representing a piece of memory for a scope
type DynamicMemoryFrame struct {
	Name        string
	Scope       *Scope
	SymbolTable *SymbolTable
	Parent      *DynamicMemoryFrame
}

// Create a new memory frame
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

/* Return the name of the memory frame.
 */
func (mf *DynamicMemoryFrame) GetName() string {
	return mf.Name
}

/* Access to the memory frame's symbol table.
 */
func (mf *DynamicMemoryFrame) Symbols() *SymbolTable {
	return mf.SymbolTable
}

/* Return the corresponding scope for this frame.
 */
func (mf *DynamicMemoryFrame) GetScope() *Scope {
	return mf.Scope
}

/* Is this a root frame?
 */
func (mf *DynamicMemoryFrame) IsRoot() bool {
	return (mf.Parent == nil)
}

// ---------------------------------------------------------------------------

// A (call-)stack of memory frames
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
	T().P("mem", newmf.Name).Debugf("pushing new memory frame")
	return newmf
}

/* Pop a memory frame. Returns the popped frame.
 */
func (mfst *MemoryFrameStack) PopMemoryFrame() *DynamicMemoryFrame {
	if mfst.memoryFrameTOS == nil {
		panic("attempt to pop memory frame from empty call stack")
	}
	mf := mfst.memoryFrameTOS
	T().Debugf("popping memory frame [%s]", mf.Name)
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
