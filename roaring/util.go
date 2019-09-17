package roaring

func highbits(x uint64) uint64 {
	return x >> 16
}

func lowbits(x uint64) uint16 {
	return uint16(x & 0xFFFF)
}
