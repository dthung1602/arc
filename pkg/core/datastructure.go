package core

type KeySpace map[string]InternalType

var rootKeySpace KeySpace

func GetRootKeySpace() *KeySpace {
	if rootKeySpace == nil {
		rootKeySpace = make(KeySpace)
	}
	return &rootKeySpace
}

func (space *KeySpace) Get(key []byte) InternalType {
	val := (*space)[string(key)]
	if val == nil {
		return nil
	}
	return val.Clone()
}

func (space *KeySpace) Set(key []byte, val InternalType) InternalType {
	ret := space.Get(key)
	(*space)[string(key)] = val.Clone()
	return ret
}

//

type InternalType interface {
	Clone() InternalType
}

type InternalBytes []byte

func (b InternalBytes) Clone() InternalType {
	clone := make(InternalBytes, len(b))
	copy(clone, b)
	return clone
}

type InternalList []InternalBytes

func (l InternalList) Clone() InternalType {
	clone := make(InternalList, len(l))
	for i, b := range l {
		clone[i] = b.Clone().(InternalBytes)
	}
	return clone
}
