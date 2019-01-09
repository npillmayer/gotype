// Code generated from PMMPVar.g4 by ANTLR 4.7.2. DO NOT EDIT.

package grammar // PMMPVar
import "github.com/antlr/antlr4/runtime/Go/antlr"

// PMMPVarListener is a complete listener for a parse tree produced by PMMPVarParser.
type PMMPVarListener interface {
	antlr.ParseTreeListener

	// EnterVariable is called when entering the variable production.
	EnterVariable(c *VariableContext)

	// EnterPathtag is called when entering the pathtag production.
	EnterPathtag(c *PathtagContext)

	// EnterSimpletag is called when entering the simpletag production.
	EnterSimpletag(c *SimpletagContext)

	// EnterSuffix is called when entering the suffix production.
	EnterSuffix(c *SuffixContext)

	// EnterSubscript is called when entering the subscript production.
	EnterSubscript(c *SubscriptContext)

	// ExitVariable is called when exiting the variable production.
	ExitVariable(c *VariableContext)

	// ExitPathtag is called when exiting the pathtag production.
	ExitPathtag(c *PathtagContext)

	// ExitSimpletag is called when exiting the simpletag production.
	ExitSimpletag(c *SimpletagContext)

	// ExitSuffix is called when exiting the suffix production.
	ExitSuffix(c *SuffixContext)

	// ExitSubscript is called when exiting the subscript production.
	ExitSubscript(c *SubscriptContext)
}
