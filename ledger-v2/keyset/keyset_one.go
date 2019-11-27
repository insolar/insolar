package keyset

// creates an immutable set with a single key
func SoloKeySet(k Key) KeySet {
	return oneKeySet{k}
}

var _ KeyList = oneKeySet{}

type oneKeySet struct {
	key Key
}

func (v oneKeySet) EnumKeys(fn func(k Key) bool) bool {
	return fn(v.key)
}

func (v oneKeySet) Count() int {
	return 1
}

func (v oneKeySet) EnumRawKeys(fn func(k Key, exclusive bool) bool) bool {
	return fn(v.key, false)
}

func (v oneKeySet) RawKeyCount() int {
	return 1
}

func (v oneKeySet) IsNothing() bool {
	return false
}

func (v oneKeySet) IsEverything() bool {
	return false
}

func (v oneKeySet) IsOpenSet() bool {
	return false
}

func (v oneKeySet) Contains(k Key) bool {
	return v.key == k
}

func (v oneKeySet) ContainsAny(ks KeySet) bool {
	return ks.Contains(v.key)
}

func (v oneKeySet) SupersetOf(ks KeySet) bool {
	if ks.IsOpenSet() {
		return false
	}

	switch ks.RawKeyCount() {
	case 0:
		return true
	case 1:
		return ks.Contains(v.key)
	default:
		return false
	}
}

func (v oneKeySet) SubsetOf(ks KeySet) bool {
	return ks.Contains(v.key)
}

func (v oneKeySet) Equal(ks KeySet) bool {
	if ks.IsOpenSet() || v.RawKeyCount() != ks.RawKeyCount() {
		return false
	}
	return ks.Contains(v.key)
}

func (v oneKeySet) EqualInverse(ks KeySet) bool {
	if !ks.IsOpenSet() || v.RawKeyCount() != ks.RawKeyCount() {
		return false
	}
	return !ks.Contains(v.key)
}

func (v oneKeySet) Inverse() KeySet {
	return exclusiveKeySet{basicKeySet{v.key: {}}}
}

func (v oneKeySet) Union(ks KeySet) KeySet {
	switch {
	case ks.Contains(v.key):
		return ks
	case ks.IsOpenSet():
		switch ks.RawKeyCount() {
		case 0:
			panic("illegal state")
		case 1:
			return Everything()
		}
		return ks.Union(v)

	case ks.RawKeyCount() == 0:
		return v
	}
	return inclusiveKeySet{keyUnion(basicKeySet{v.key: {}}, ks)}
}

func (v oneKeySet) Intersect(ks KeySet) KeySet {
	switch {
	case ks.Contains(v.key):
		return v
	case ks.RawKeyCount() == 0:
		return ks
	default:
		return Nothing()
	}
}

func (v oneKeySet) Subtract(ks KeySet) KeySet {
	if ks.Contains(v.key) {
		return Nothing()
	}
	return v
}
