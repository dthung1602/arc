package core

type HashMap map[string][]byte

var hashMapInstance = make(HashMap)

func (m *HashMap) Get(key []byte) []byte {
	val := (*m)[string(key)]
	if val == nil {
		return nil
	}
	data := make([]byte, len(val))
	copy(data, val) // TODO should copy?
	return data
}

func (m *HashMap) Set(key []byte, val []byte) {
	data := make([]byte, len(val))
	copy(data, val) // TODO should copy?
	(*m)[string(key)] = data
}
