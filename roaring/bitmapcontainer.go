package roaring

type bitmapContainer struct {
	cardinality int
	bitmap      []uint64
}

func newBitmapContainer() *bitmapContainer {
	p := new(bitmapContainer)
	size := (1 << 16) / 64
	p.bitmap = make([]uint64, size, size)
	return p
}

func (bc *bitmapContainer) addReturnMinimized(i uint16) container {
	bc.add(i)
	// TODO
	// if bc.isFull() {
	// 	return newRunContainer16Range(0, MaxUint16)
	// }
	return bc
}

func (bc *bitmapContainer) add(i uint16) bool {
	x := int(i)
	previous := bc.bitmap[x/64]
	mask := uint64(1) << (uint(x) % 64)
	newb := previous | mask
	bc.bitmap[x/64] = newb
	bc.cardinality += int((previous ^ newb) >> (uint(x) % 64))
	return newb != previous
}

func (bc *bitmapContainer) contains(i uint16) bool {
	x := uint(i)
	w := bc.bitmap[x>>6]
	mask := uint64(1) << (x & 63)
	return (w & mask) != 0
}

func (bc *bitmapContainer) getCardinality() int {
	return bc.cardinality
}
