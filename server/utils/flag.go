package utils

type flagSize interface {
	uint8 | uint16 | uint32 | uint64
}

type Flag[T flagSize] struct {
	Flag T
}

func (f *Flag[T]) SetFlag(bit T) {
	f.Flag |= bit
}

func (f *Flag[T]) ClearFlag(bit T) {
	f.Flag &= ^bit
}

func (f *Flag[T]) HasFlag(bit T) bool {
	return f.Flag&bit != 0
}
