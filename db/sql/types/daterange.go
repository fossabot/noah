package types

import (
	"database/sql/driver"
	"github.com/Ready-Stock/Noah/db/sql/pgio"
	"github.com/pkg/errors"
)

type Daterange struct {
	Lower     Date
	Upper     Date
	LowerType BoundType
	UpperType BoundType
	Status    Status
}

func (dst *Daterange) Set(src interface{}) error {
	return errors.Errorf("cannot convert %v to Daterange", src)
}

func (dst *Daterange) Get() interface{} {
	switch dst.Status {
	case Present:
		return dst
	case Null:
		return nil
	default:
		return dst.Status
	}
}

func (src *Daterange) AssignTo(dst interface{}) error {
	return errors.Errorf("cannot assign %v to %T", src, dst)
}

func (dst *Daterange) DecodeText(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Daterange{Status: Null}
		return nil
	}

	utr, err := ParseUntypedTextRange(string(src))
	if err != nil {
		return err
	}

	*dst = Daterange{Status: Present}

	dst.LowerType = utr.LowerType
	dst.UpperType = utr.UpperType

	if dst.LowerType == Empty {
		return nil
	}

	if dst.LowerType == Inclusive || dst.LowerType == Exclusive {
		if err := dst.Lower.DecodeText(ci, []byte(utr.Lower)); err != nil {
			return err
		}
	}

	if dst.UpperType == Inclusive || dst.UpperType == Exclusive {
		if err := dst.Upper.DecodeText(ci, []byte(utr.Upper)); err != nil {
			return err
		}
	}

	return nil
}

func (dst *Daterange) DecodeBinary(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Daterange{Status: Null}
		return nil
	}

	ubr, err := ParseUntypedBinaryRange(src)
	if err != nil {
		return err
	}

	*dst = Daterange{Status: Present}

	dst.LowerType = ubr.LowerType
	dst.UpperType = ubr.UpperType

	if dst.LowerType == Empty {
		return nil
	}

	if dst.LowerType == Inclusive || dst.LowerType == Exclusive {
		if err := dst.Lower.DecodeBinary(ci, ubr.Lower); err != nil {
			return err
		}
	}

	if dst.UpperType == Inclusive || dst.UpperType == Exclusive {
		if err := dst.Upper.DecodeBinary(ci, ubr.Upper); err != nil {
			return err
		}
	}

	return nil
}

func (src Daterange) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	switch src.LowerType {
	case Exclusive, Unbounded:
		buf = append(buf, '(')
	case Inclusive:
		buf = append(buf, '[')
	case Empty:
		return append(buf, "empty"...), nil
	default:
		return nil, errors.Errorf("unknown lower bound type %v", src.LowerType)
	}

	var err error

	if src.LowerType != Unbounded {
		buf, err = src.Lower.EncodeText(ci, buf)
		if err != nil {
			return nil, err
		} else if buf == nil {
			return nil, errors.Errorf("Lower cannot be null unless LowerType is Unbounded")
		}
	}

	buf = append(buf, ',')

	if src.UpperType != Unbounded {
		buf, err = src.Upper.EncodeText(ci, buf)
		if err != nil {
			return nil, err
		} else if buf == nil {
			return nil, errors.Errorf("Upper cannot be null unless UpperType is Unbounded")
		}
	}

	switch src.UpperType {
	case Exclusive, Unbounded:
		buf = append(buf, ')')
	case Inclusive:
		buf = append(buf, ']')
	default:
		return nil, errors.Errorf("unknown upper bound type %v", src.UpperType)
	}

	return buf, nil
}

func (src Daterange) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	var rangeType byte
	switch src.LowerType {
	case Inclusive:
		rangeType |= lowerInclusiveMask
	case Unbounded:
		rangeType |= lowerUnboundedMask
	case Exclusive:
	case Empty:
		return append(buf, emptyMask), nil
	default:
		return nil, errors.Errorf("unknown LowerType: %v", src.LowerType)
	}

	switch src.UpperType {
	case Inclusive:
		rangeType |= upperInclusiveMask
	case Unbounded:
		rangeType |= upperUnboundedMask
	case Exclusive:
	default:
		return nil, errors.Errorf("unknown UpperType: %v", src.UpperType)
	}

	buf = append(buf, rangeType)

	var err error

	if src.LowerType != Unbounded {
		sp := len(buf)
		buf = pgio.AppendInt32(buf, -1)

		buf, err = src.Lower.EncodeBinary(ci, buf)
		if err != nil {
			return nil, err
		}
		if buf == nil {
			return nil, errors.Errorf("Lower cannot be null unless LowerType is Unbounded")
		}

		pgio.SetInt32(buf[sp:], int32(len(buf[sp:])-4))
	}

	if src.UpperType != Unbounded {
		sp := len(buf)
		buf = pgio.AppendInt32(buf, -1)

		buf, err = src.Upper.EncodeBinary(ci, buf)
		if err != nil {
			return nil, err
		}
		if buf == nil {
			return nil, errors.Errorf("Upper cannot be null unless UpperType is Unbounded")
		}

		pgio.SetInt32(buf[sp:], int32(len(buf[sp:])-4))
	}

	return buf, nil
}

// Scan implements the database/sql Scanner interface.
func (dst *Daterange) Scan(src interface{}) error {
	if src == nil {
		*dst = Daterange{Status: Null}
		return nil
	}

	switch src := src.(type) {
	case string:
		return dst.DecodeText(nil, []byte(src))
	case []byte:
		srcCopy := make([]byte, len(src))
		copy(srcCopy, src)
		return dst.DecodeText(nil, srcCopy)
	}

	return errors.Errorf("cannot scan %T", src)
}

// Value implements the database/sql/driver Valuer interface.
func (src Daterange) Value() (driver.Value, error) {
	return EncodeValueText(src)
}
