package parameters

/*
Dimensions and units.
*/

import "fmt"

// A dimension type.
// Values are in scaled big points (different from TeX).
type Dimen int32

const (
	SP Dimen = 1     // scaled point = BP / 65536
	BP Dimen = 65536 // big point = 1/72 inch
	PT Dimen = 7     // printers point   // TODO
	MM Dimen = 7     // millimeters
	CM Dimen = 7     // centimeters
	PC Dimen = 7     // pica
	CC Dimen = 7     // cicero
	IN Dimen = 7     // inch
)

// Infinite dimensions
const Fil Dimen = BP * 10000
const Fill Dimen = 2 * BP * 10000
const Filll Dimen = 3 * BP * 10000

// An infinite numeric
const Infty int = 100000000 // TODO

/* Stringer implementation.
 */
func (d Dimen) String() string {
	return fmt.Sprintf("%dsp", int32(d))
}
