package binchunk

import (
	"encoding/binary"
	"math"
)

type reader struct {
	data []byte // 待解析二进制数据
}

func (self *reader) readByte() byte {
	b := self.data[0]
	self.data = self.data[1:]
	return b
}

func (self *reader) readBytes(n uint) []byte {
	bs := self.data[0:n]
	self.data = self.data[n:]
	return bs
}

func (self *reader) readUint32() uint32 {
	bs := binary.LittleEndian.Uint32(self.data)
	self.data = self.data[4:]
	return bs
}

func (self *reader) readUint64() uint64 {
	bs := binary.LittleEndian.Uint64(self.data)
	self.data = self.data[8:]
	return bs
}

func (self *reader) readLuaInteger() int64 {
	return int64(self.readUint64())
}

func (self *reader) readLuaNumber() float64 {
	return math.Float64frombits(self.readUint64())
}

func (self *reader) readString() string {
	size := uint(self.readByte())
	if size == 0x00 {
		return ""
	}

	if size == 0xFF {
		size = uint(self.readUint64())
	}

	bs := self.readBytes(size)
	return string(bs)
}
