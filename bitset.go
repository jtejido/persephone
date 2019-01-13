package persephone

const (
	maxUint      = ^uint(0)
	wordSize     = uint(64)
	log2WordSize = uint(6)
)

// Internal bitset type for states and inputs for efficiency.
type bitSet struct {
	length uint
	set    []uint64
}

func (b *bitSet) contains(i uint) bool {
	if i >= b.length {
		return false
	}
	return b.set[i>>log2WordSize]&(1<<(i&(wordSize-1))) != 0
}

func (b *bitSet) add(i uint) {
	b.extendSetMaybe(i)
	b.set[i>>log2WordSize] |= 1 << (i & (wordSize - 1))
}

func (b *bitSet) len() uint {
	return b.length
}

func (b *bitSet) extendSetMaybe(i uint) {
	if i >= b.length {
		nsize := wordsNeeded(i + 1)
		if b.set == nil {
			b.set = make([]uint64, nsize)
		} else if cap(b.set) >= nsize {
			b.set = b.set[:nsize] // resize
		} else if len(b.set) < nsize {
			newset := make([]uint64, nsize, 2*nsize) // capacity x 2
			copy(newset, b.set)
			b.set = newset
		}
		b.length = i + 1
	}
}

func wordsNeeded(i uint) int {
	if i > (maxUint - wordSize + 1) {
		return int(maxUint >> log2WordSize)
	}
	return int((i + (wordSize - 1)) >> log2WordSize)

}
