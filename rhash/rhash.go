package rhash

import (
	"sync"
)

const (
	sizeofUInt64  = 64
	divideByShift = 4
	moduloByAnd   = 15
)
// Bitmap64 base class for using bitmap as item filter
type Bitmap64 struct {
	sync.RWMutex
	bitCount int64
	rhmap map[uint64]map[uint16]uint16
}

// NewBitmap64 create a Bitmap64
func NewBitmap64() *Bitmap64 {
	return &Bitmap64{
		rhmap: make(map[uint64]map[uint16]uint16),
		bitCount: 0,
	}
}

func highbits(x uint64) uint64 {
	return x>>16
}

func lowbits(x uint64) uint16 {
	return uint16(x & 0xFFFF)
}

// IsSet checks if position bit is set
func (b *Bitmap64) IsSet(x uint64) bool {
	hb := highbits(x)
	b.RLock()
	defer b.RUnlock()
	im,ok := b.rhmap[hb]
	if !ok {
		return false
	}
	lb := lowbits(x)
	id, pos := lb >> divideByShift , lb & moduloByAnd
	return (im[id] >> pos) & (0x1) != 0
}

// SetBit SetBit
func (b *Bitmap64) SetBit(x uint64, val uint8) bool {
	if val == 0x01 {
		return b.Set(x)
	}
	if val == 0x00 {
		return b.Clear(x)
	}
	return false
}

// Set set indexth bit to 1
func (b *Bitmap64) Set(x uint64) bool {
	hb := highbits(x)
	b.Lock()
	defer b.Unlock()
	im, ok := b.rhmap[hb]
	if !ok {
		b.rhmap[hb] = make(map[uint16]uint16)
		im = b.rhmap[hb]
	}
	lb := lowbits(x)
	id,pos := lb >> divideByShift, lb & moduloByAnd
	if (im[id] & (0x01 << pos)) == 0 {
		b.bitCount++
	}
	im[id] |= (0x1 << pos)
	return true
}

func (b *Bitmap64) Size() int64 {
	b.RLock()
	defer b.RUnlock()
	return b.bitCount
}

// Clear set indexth bit to 0
func (b *Bitmap64) Clear(x uint64) bool {
	hb := highbits(x)
	b.Lock()
	defer b.Unlock()
	im, ok := b.rhmap[hb]
	if !ok {
		return false
	}

	lb := lowbits(x)
	id,pos := lb >> divideByShift, lb & moduloByAnd
	if (im[id] & (0x01 << pos)) != 0 {
		b.bitCount--
	}
	im[id] &^= (0x01 << pos)
	return true
}





