package thrift

import (
	"bytes"
	"testing"
)

func TestCompactProtocol(t *testing.T) {
	testProtocol(t, NewCompactProtocol())
}

func TestCompactList(t *testing.T) {
	p := NewCompactProtocol()

	tests := []struct {
		values []byte
		bytes  []byte
	}{
		{[]byte{}, []byte{3}},
		{[]byte{64}, []byte{19, 64}},
		{[]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
			[]byte{243, 17, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}},
	}

	for _, exp := range tests {
		expValue := exp.values
		expBytes := exp.bytes
		b := &bytes.Buffer{}
		if err := p.WriteListBegin(b, TypeByte, len(expValue)); err != nil {
			t.Fatalf("WriteListBegin returned an error: %+v", err)
		}
		for _, v := range expValue {
			if err := p.WriteByte(b, v); err != nil {
				t.Fatalf("WriteByte returned an error: %+v", err)
			}
		}
		if err := p.WriteListEnd(b); err != nil {
			t.Fatalf("WriteListEnd returned an error: %+v", err)
		}
		out := b.Bytes()
		if bytes.Compare(out, expBytes) != 0 {
			t.Fatalf("WriteListBegin wrote %+v which did not match expected %+v", out, expBytes)
		}

		b = bytes.NewBuffer(expBytes)
		etype, size, err := p.ReadListBegin(b)
		if err != nil {
			t.Fatalf("ReadListBegin returned an error: %+v", err)
		} else if etype != TypeByte {
			t.Fatalf("ReadListBegin returned wrong type %d instead of %d", etype, TypeByte)
		} else if size != len(expValue) {
			t.Fatalf("ReadListBegin returned wrong size %d insted of %d", size, len(expValue))
		}
		for i := 0; i < size; i++ {
			if v, err := p.ReadByte(b); err != nil {
				t.Fatalf("ReadByte returned an error: %+v", err)
			} else if v != expValue[i] {
				t.Fatalf("ReadByte returned wrong value %d insted of %d", v, expBytes[i])
			}
		}
		if err := p.ReadListEnd(b); err != nil {
			t.Fatalf("ReadListEnd returned an error: %+v", err)
		}
	}
}

func TestCompactString(t *testing.T) {
	p := NewCompactProtocol()

	expStrings := map[string][]byte{
		"":    []byte{0},
		"foo": []byte{3, 102, 111, 111},
	}

	for expValue, expBytes := range expStrings {
		b := &bytes.Buffer{}
		err := p.WriteString(b, expValue)
		if err != nil {
			t.Fatalf("WriteString returned an error: %+v", err)
		}
		out := b.Bytes()
		if bytes.Compare(out, expBytes) != 0 {
			t.Fatalf("WriteString wrote %+v which did not match expected %+v", out, expBytes)
		}

		b = bytes.NewBuffer(expBytes)
		v, err := p.ReadString(b)
		if err != nil {
			t.Fatalf("ReadString returned an error: %+v", err)
		}
		if v != expValue {
			t.Fatalf("ReadString returned the wrong value %d instead of %d", v, expValue)
		}
	}
}

func TestCompactI16(t *testing.T) {
	p := NewCompactProtocol()

	exp := map[int16][]byte{
		0:     []byte{0},
		-1:    []byte{1},
		1:     []byte{2},
		12345: []byte{242, 192, 1},
	}

	for expValue, expBytes := range exp {
		b := &bytes.Buffer{}
		err := p.WriteI16(b, expValue)
		if err != nil {
			t.Fatalf("WriteI16 returned an error: %+v", err)
		}
		out := b.Bytes()
		if bytes.Compare(out, expBytes) != 0 {
			t.Fatalf("WriteI16 wrote %+v which did not match expected %+v", out, expBytes)
		}

		b = bytes.NewBuffer(expBytes)
		v, err := p.ReadI16(b)
		if err != nil {
			t.Fatalf("ReadI16 returned an error: %+v", err)
		}
		if v != expValue {
			t.Fatalf("ReadI16 returned the wrong value %d instead of %d", v, expValue)
		}
	}
}

func TestCompactI32(t *testing.T) {
	p := NewCompactProtocol()

	exp := map[int32][]byte{
		0:          []byte{0},
		-1:         []byte{1},
		1:          []byte{2},
		1234567890: []byte{164, 139, 176, 153, 9},
	}

	for expValue, expBytes := range exp {
		b := &bytes.Buffer{}
		err := p.WriteI32(b, expValue)
		if err != nil {
			t.Fatalf("WriteI32 returned an error: %+v", err)
		}
		out := b.Bytes()
		if bytes.Compare(out, expBytes) != 0 {
			t.Fatalf("WriteI32 wrote %+v which did not match expected %+v", out, expBytes)
		}

		b = bytes.NewBuffer(expBytes)
		v, err := p.ReadI32(b)
		if err != nil {
			t.Fatalf("Read32 returned an error: %+v", err)
		}
		if v != expValue {
			t.Fatalf("Read32 returned the wrong value %d instead of %d", v, expValue)
		}
	}
}
