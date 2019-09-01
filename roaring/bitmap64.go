package roaring

import (
	"sync"
)

// Bitmap64 base class for using bitmap as item filter
type Bitmap64 struct {
	sync.RWMutex
	keys       []uint64
	containers []container
}

// NewBitmap64 create a Bitmap64
func NewBitmap64() *Bitmap64 {
	return &Bitmap64{
		keys:       make([]uint64, 0),
		containers: make([]container, 0),
	}
}

// IsSet checks if position bit is set
func (b *Bitmap64) IsSet(x uint64) bool {
	hb := highbits(x)
	b.RLock()
	defer b.RUnlock()
	c := b.getContainer(hb)
	return c != nil && c.contains(lowbits(x))
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
	i := b.getContainerIndex(hb)
	if i >= 0 {
		var c container = b.getContainerAtIndex(i).addReturnMinimized(lowbits(x))
		b.setContainerAtIndex(i, c)
	} else {
		ac := newArrayContainer()
		b.insertContainerAt(-i-1, hb, ac.addReturnMinimized(lowbits(x)))
	}
	return true
}

// Clear set indexth bit to 0
func (b *Bitmap64) Clear(x uint64) bool {
	// TODO
	return true
}

func (b *Bitmap64) getContainer(hb uint64) container {
	i := b.binarySearch(hb)
	if i < 0 {
		return nil
	}
	return b.containers[i]
}

func (b *Bitmap64) getContainerIndex(hb uint64) int {
	// before the binary search, we optimize for frequent cases
	size := len(b.keys)
	if (size == 0) || (b.keys[size-1] == hb) {
		return size - 1
	}
	return b.binarySearch(hb)
}

func (b *Bitmap64) binarySearch(hb uint64) int {
	low := 0
	high := len(b.keys) - 1
	for low+16 <= high {
		middleIndex := low + (high-low)/2 // avoid overflow
		middleValue := b.keys[middleIndex]
		if middleValue < hb {
			low = middleIndex + 1
		} else if middleValue > hb {
			high = middleIndex - 1
		} else {
			return middleIndex
		}
	}
	for ; low <= high; low++ {
		val := b.keys[low]
		if val >= hb {
			if val == hb {
				return low
			}
			break
		}
	}
	return -(low + 1)
}

func (b *Bitmap64) setContainerAtIndex(i int, c container) {
	b.containers[i] = c
}

func (b *Bitmap64) getContainerAtIndex(i int) container {
	return b.containers[i]
}

func (b *Bitmap64) insertContainerAt(i int, hb uint64, value container) {
	b.keys = append(b.keys, 0)
	b.containers = append(b.containers, nil)

	copy(b.keys[i+1:], b.keys[i:])
	copy(b.containers[i+1:], b.containers[i:])

	b.keys[i] = hb
	b.containers[i] = value
}
