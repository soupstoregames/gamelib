package data

type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type Bitfield1[T Unsigned] struct {
	Raw T
}

func NewBitfield1[T Unsigned](value T) Bitfield1[T] {
	return Bitfield1[T]{
		Raw: value,
	}
}

func (b *Bitfield1[T]) Set(i int) {
	b.Raw |= 1 << i
}

func (b *Bitfield1[T]) Clear(i int) {
	b.Raw &^= 1 << i
}

func (b *Bitfield1[T]) Toggle(i int) {
	b.Raw ^= 1 << i
}

func (b *Bitfield1[T]) Has(i int) bool {
	return b.Raw&(1<<i) != 0
}
