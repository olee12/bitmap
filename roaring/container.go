package roaring

type container interface {
	addReturnMinimized(x uint16) container
	contains(x uint16) bool
}
