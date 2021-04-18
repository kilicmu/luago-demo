package binchunk

import (
	"encoding/binary"
	"math"
	"testing"
)

func TestReadByte(t *testing.T) {
	data := []byte{0x00, 0x01}
	reader := &reader{data}

	cases := []struct {
		Name   string
		Expect byte
	}{
		{"first binary", data[0]},
		{"second binary", data[1]},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			if ans := reader.readByte(); ans != c.Expect {
				t.Errorf("%s error: expect %d, but get %d", c.Name, c.Expect, ans)
			}
		})
	}
}

func TestReadBytes(t *testing.T) {
	data := []byte("hello world")
	reader := &reader{data}

	cases := []struct {
		Name   string
		Input  uint
		Expect string
	}{
		{"zero bytes", 0, ""},
		{"1 bytes", 1, string(data[0])},
		{"2 bytes", 2, string(data[1:3])},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			if ans := reader.readBytes(c.Input); string(ans) != c.Expect {
				t.Errorf("%s error: expect %s, but get %s", c.Name, c.Expect, string(ans))
			}
		})
	}
}

func TestReadUint32(t *testing.T) {
	data := []byte{
		0x00,
		0x00,
		0x12,
		0x34,
	}

	reader := &reader{
		data,
	}

	t.Run("test get uint32 by LE", func(t *testing.T) {
		exp := binary.LittleEndian.Uint32(data)
		ans := reader.readUint32()
		if ans != exp {
			t.Errorf("get uint32 error: expect %d, but get %d", exp, ans)
		}
		if len(reader.data) != 0 {
			t.Errorf("len not expect: get %d", len(data))
		}
	})
}

func TestReadUint64(t *testing.T) {
	data := []byte{
		0x00,
		0x00,
		0x00,
		0x00,
		0x00,
		0x00,
		0x12,
		0x34,
	}

	reader := &reader{data}

	t.Run("test get uint64 by LE", func(t *testing.T) {
		exp := binary.LittleEndian.Uint64(data)
		ans := reader.readUint64()
		if exp != ans {
			t.Errorf("get uint64 error: expect %d, but get %d", exp, ans)
		}
		if len(reader.data) != 0 {
			t.Errorf("len not expect: get %d", len(data))
		}
	})

}

func TestReadLuaInteger(t *testing.T) {
	data := []byte{
		0x00,
		0x00,
		0x00,
		0x00,
		0x00,
		0x00,
		0x12,
		0x34,
	}

	reader := &reader{data}

	t.Run("test get lua integer by LE", func(t *testing.T) {
		exp := int64(binary.LittleEndian.Uint64(data))
		ans := reader.readLuaInteger()
		if exp != ans {
			t.Errorf("get lua integer error: expect %d, but get %d", exp, ans)
		}
		if len(reader.data) != 0 {
			t.Errorf("len not expect: get %d", len(data))
		}
	})

}

func TestReadLuaNumber(t *testing.T) {
	data := []byte{
		0x00,
		0x00,
		0x00,
		0x00,
		0x00,
		0x00,
		0x12,
		0x34,
	}

	reader := &reader{data}

	t.Run("test get lua integer by LE", func(t *testing.T) {
		exp := math.Float64frombits(binary.LittleEndian.Uint64(data))
		ans := reader.readLuaNumber()
		if exp != ans {
			t.Errorf("get lua number error: expect %f, but get %f", exp, ans)
		}
		if len(reader.data) != 0 {
			t.Errorf("len not expect: get %d", len(data))
		}
	})

}

func TestReadString(t *testing.T) {

	t.Run("test get lua string", func(t *testing.T) {
		data := []byte{
			0x00,
		}

		reader := &reader{data}
		if ans := reader.readString(); ans != "" {
			t.Errorf("get nil string error: expect :\"\", but get %s", ans)
		}
		if len(reader.data) != 0 {
			t.Errorf("len not expect: get %d", len(data))
		}
	})

	t.Run("test get lua string len > 0xFF", func(t *testing.T) {
		data := []byte{0xFF}
		exp := string("hello world")
		data = append(data, []byte{0x0B, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}...)
		data = append(data, []byte(exp)...)

		reader := &reader{data}
		if ans := reader.readString(); ans != exp {
			t.Errorf("get nil string error: expect : %s, but get %s", exp, ans)
		}
		if len(reader.data) != 0 {
			t.Errorf("len not expect: get %d", len(data))
		}
	})

	t.Run("test get lua string len > 0xFF", func(t *testing.T) {
		data := []byte{0x0B}
		exp := string("hello world")
		data = append(data, []byte(exp)...)

		reader := &reader{data}
		if ans := reader.readString(); ans != exp {
			t.Errorf("get nil string error: expect : %s, but get %s", exp, ans)
		}
		if len(reader.data) != 0 {
			t.Errorf("len not expect: get %d", len(data))
		}
	})

}
