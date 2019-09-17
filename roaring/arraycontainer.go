package roaring

const (
	arrayDefaultMaxSize = 4096
)

type arrayContainer struct {
	content []uint16
}

func newArrayContainer() *arrayContainer {
	return &arrayContainer{
		content: make([]uint16, 0),
	}
}

func (ac *arrayContainer) addReturnMinimized(x uint16) container {
	// Special case adding to the end of the container.
	l := len(ac.content)
	if l > 0 && l < arrayDefaultMaxSize && ac.content[l-1] < x {
		ac.content = append(ac.content, x)
		return ac
	}

	loc := ac.binarySearch(x)
	if loc < 0 {
		if len(ac.content) >= arrayDefaultMaxSize {
			a := ac.toBitmapContainer()
			a.add(x)
			return a
		}
		s := ac.content
		i := -loc - 1
		s = append(s, 0)
		copy(s[i+1:], s[i:])
		s[i] = x
		ac.content = s
	}
	return ac
}

func (ac *arrayContainer) contains(x uint16) bool {
	return ac.binarySearch(x) >= 0
}

func (ac *arrayContainer) getCardinality() int {
	return len(ac.content)
}

func (ac *arrayContainer) toBitmapContainer() *bitmapContainer {
	bc := newBitmapContainer()
	bc.loadData(ac)
	return bc
}

func (bc *bitmapContainer) loadData(arrayContainer *arrayContainer) {
	bc.cardinality = arrayContainer.getCardinality()
	c := arrayContainer.getCardinality()
	for k := 0; k < c; k++ {
		x := arrayContainer.content[k]
		i := int(x) / 64
		bc.bitmap[i] |= (uint64(1) << uint(x%64))
	}
}

func (ac *arrayContainer) binarySearch(x uint16) int {
	low := 0
	high := len(ac.content) - 1
	for low+16 <= high {
		middleIndex := low + (high-low)/2
		middleValue := ac.content[middleIndex]
		if middleValue < x {
			low = middleIndex + 1
		} else if middleValue > x {
			high = middleIndex - 1
		} else {
			return middleIndex
		}
	}
	for ; low <= high; low++ {
		val := ac.content[low]
		if val >= x {
			if val == x {
				return low
			}
			break
		}
	}
	return -(low + 1)
}
