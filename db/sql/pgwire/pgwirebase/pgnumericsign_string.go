// Code generated by "stringer -type=PGNumericSign"; DO NOT EDIT.

package pgwirebase

import "strconv"

const (
    _PGNumericSign_name_0 = "PGNumericPos"
    _PGNumericSign_name_1 = "PGNumericNeg"
)

func (i PGNumericSign) String() string {
    switch {
    case i == 0:
        return _PGNumericSign_name_0
    case i == 16384:
        return _PGNumericSign_name_1
    default:
        return "PGNumericSign(" + strconv.FormatInt(int64(i), 10) + ")"
    }
}
