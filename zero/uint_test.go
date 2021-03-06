package zero

import (
	"encoding/json"
	"testing"
)

var (
	uintJSON     = []byte(`12345`)
	nullUintJSON = []byte(`{"Uint":12345,"Valid":true}`)
	zeroUJSON    = []byte(`0`)
)

func TestUintFrom(t *testing.T) {
	i := UintFrom(12345)
	assertUint(t, i, "UintFrom()")

	zero := UintFrom(0)
	if zero.Valid {
		t.Error("UintFrom(0)", "is valid, but should be invalid")
	}
}

func TestUintFromPtr(t *testing.T) {
	n := uint(12345)
	iptr := &n
	i := UintFromPtr(iptr)
	assertUint(t, i, "UintFromPtr()")

	null := UintFromPtr(nil)
	assertNullUint(t, null, "UintFromPtr(nil)")
}

func TestUnmarshalUint(t *testing.T) {
	var i Uint
	err := json.Unmarshal(uintJSON, &i)
	maybePanic(err)
	assertUint(t, i, "uint json")

	var ni Uint
	err = json.Unmarshal(nullUintJSON, &ni)
	maybePanic(err)
	assertUint(t, ni, "sql.NullUint json")

	var zero Uint
	err = json.Unmarshal(zeroUJSON, &zero)
	maybePanic(err)
	assertNullUint(t, zero, "zero json")

	var null Uint
	err = json.Unmarshal(nullJSON, &null)
	maybePanic(err)
	assertNullUint(t, null, "null json")

	var badType Uint
	err = json.Unmarshal(boolJSON, &badType)
	if err == nil {
		panic("err should not be nil")
	}
	assertNullUint(t, badType, "wrong type json")

	var invalid Uint
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("expected json.SyntaxError, not %T", err)
	}
	assertNullUint(t, invalid, "invalid json")
}

func TestUnmarshalNonUintegerNumber(t *testing.T) {
	var i Uint
	err := json.Unmarshal(floatJSON, &i)
	if err == nil {
		panic("err should be present; non-uinteger number coerced to uint")
	}
}

func TestTextUnmarshalUint(t *testing.T) {
	var i Uint
	err := i.UnmarshalText([]byte("12345"))
	maybePanic(err)
	assertUint(t, i, "UnmarshalText() uint")

	var zero Uint
	err = zero.UnmarshalText([]byte("0"))
	maybePanic(err)
	assertNullUint(t, zero, "UnmarshalText() zero uint")

	var blank Uint
	err = blank.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullUint(t, blank, "UnmarshalText() empty uint")

	var null Uint
	err = null.UnmarshalText([]byte("null"))
	maybePanic(err)
	assertNullUint(t, null, `UnmarshalText() "null"`)
}

func TestMarshalUint(t *testing.T) {
	i := UintFrom(12345)
	data, err := json.Marshal(i)
	maybePanic(err)
	assertJSONEquals(t, data, "12345", "non-empty json marshal")

	// invalid values should be encoded as 0
	null := NewUint(0, false)
	data, err = json.Marshal(null)
	maybePanic(err)
	assertJSONEquals(t, data, "0", "null json marshal")
}

func TestMarshalUintText(t *testing.T) {
	i := UintFrom(12345)
	data, err := i.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "12345", "non-empty text marshal")

	// invalid values should be encoded as zero
	null := NewUint(0, false)
	data, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "0", "null text marshal")
}

func TestUintPouinter(t *testing.T) {
	i := UintFrom(12345)
	ptr := i.Ptr()
	if *ptr != 12345 {
		t.Errorf("bad %s uint: %#v ≠ %d\n", "pouinter", ptr, 12345)
	}

	null := NewUint(0, false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s uint: %#v ≠ %s\n", "nil pouinter", ptr, "nil")
	}
}

func TestUintIsZero(t *testing.T) {
	i := UintFrom(12345)
	if i.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := NewUint(0, false)
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}

	zero := NewUint(0, true)
	if !zero.IsZero() {
		t.Errorf("IsZero() should be true")
	}
}

func TestUintScan(t *testing.T) {
	var i Uint
	err := i.Scan(12345)
	maybePanic(err)
	assertUint(t, i, "scanned uint")

	var null Uint
	err = null.Scan(nil)
	maybePanic(err)
	assertNullUint(t, null, "scanned null")
}

func TestUintSetValid(t *testing.T) {
	change := NewUint(0, false)
	assertNullUint(t, change, "SetValid()")
	change.SetValid(12345)
	assertUint(t, change, "SetValid()")
}

func assertUint(t *testing.T, i Uint, from string) {
	if i.Uint != 12345 {
		t.Errorf("bad %s uint: %d ≠ %d\n", from, i.Uint, 12345)
	}
	if !i.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullUint(t *testing.T, i Uint, from string) {
	if i.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}
