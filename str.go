package apiver

import (
	"errors"
	"strconv"
)

type IntS int64

func (i *IntS) UnmarshalJSON(d []byte) error {
	if len(d) > 2 && d[0] == '"' {
		d = d[1 : len(d)-1]
	}
	err := error(nil)
	di, err := strconv.ParseInt(string(d), 10, 64)
	if err != nil {
		return errors.New("not an integer")
	}
	*i = IntS(di)
	return nil
}

type UIntS uint64

func (i *UIntS) UnmarshalJSON(d []byte) error {
	if len(d) > 2 && d[0] == '"' {
		d = d[1 : len(d)-1]
	}
	err := error(nil)
	di, err := strconv.ParseUint(string(d), 10, 64)
	if err != nil {
		return errors.New("not an integer")
	}
	*i = UIntS(di)
	return nil
}

type FloatS float64

func (i *FloatS) UnmarshalJSON(d []byte) error {
	if len(d) > 2 && d[0] == '"' {
		d = d[1 : len(d)-1]
	}
	err := error(nil)
	di, err := strconv.ParseFloat(string(d), 64)
	if err != nil {
		return errors.New("not a number")
	}
	*i = FloatS(di)
	return nil
}

// BoolS unmarshals a JSON bool or string as a boolean.
// Note strings are unmarshaled like Perl: 0, "0" and "" are false; all other values are true.
type BoolS bool

func (i *BoolS) UnmarshalJSON(d []byte) error {
	// TODO determine if empty JSON arrays are false in Perl and should be false
	sd := string(d)
	if sd == `""` {
		*i = BoolS(false)
		return nil
	}
	if sd == `"0"` {
		*i = BoolS(false)
		return nil
	}
	if di, err := strconv.ParseFloat(sd, 64); err == nil && di == 0.0 {
		*i = BoolS(false)
		return nil
	}
	*i = BoolS(true)
	return nil
}
