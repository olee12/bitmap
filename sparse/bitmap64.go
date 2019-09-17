package sparse

import (
	"sync"
)

const (
	sizeofUInt64  = 64
	divideByShift = 6
	moduloByAnd   = 63
)

// Bitmap64 base class for using bitmap as item filter
type Bitmap64 struct {
	sync.RWMutex
	data     map[uint64]uint64
	bitCount int
}

// NewBitmap64 create a Bitmap64
func NewBitmap64() *Bitmap64 {
	bf := &Bitmap64{
		data: make(map[uint64]uint64),
	}
	return bf
}

// NewBitmap64Size create a Bitmap64
func NewBitmap64Size(size int) *Bitmap64 {
	bf := &Bitmap64{
		data: make(map[uint64]uint64),
	}
	return bf
}

// IsSet checks if position bit is set
func (bf *Bitmap64) IsSet(index uint64) bool {
	id, pos := index>>divideByShift, index&moduloByAnd
	bf.RLock()
	defer bf.RUnlock()
	val, ok := bf.data[id]
	if !ok {
		return false
	}
	return ((val >> pos) & 0x01) != 0
}

//SetBit SetBit
func (bf *Bitmap64) SetBit(index uint64, val uint8) bool {
	if val == 0x01 {
		return bf.Set(index)
	} else if val == 0x00 {
		return bf.Clear(index)
	}
	return false
}

// Set set indexth bit to 1
func (bf *Bitmap64) Set(index uint64) bool {
	id, pos := index>>divideByShift, index&moduloByAnd
	bf.Lock()
	defer bf.Unlock()
	if (bf.data[id] & (0x01 << pos)) == 0 {
		bf.bitCount++
	}
	bf.data[id] |= (0x01 << pos)
	return true
}

// Clear set indexth bit to 0
func (bf *Bitmap64) Clear(index uint64) bool {
	id, pos := index>>divideByShift, index&moduloByAnd
	bf.Lock()
	defer bf.Unlock()
	_, ok := bf.data[id]
	if !ok {
		return false
	}
	if (bf.data[id] & (0x01 << pos)) != 0 {
		bf.bitCount--
	}
	bf.data[id] &^= (0x01 << pos)
	return true
}

//Size Size (should be correct)
func (bf *Bitmap64) Size() int {
	return bf.bitCount
}
