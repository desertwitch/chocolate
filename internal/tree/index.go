package tree

type index[K comparable, T any] map[K]Node[K, T]

func (idx *index[K, T]) find(id K) Node[K, T] {
	if idx == nil {
		return nil
	}

	m := *idx
	if v, ok := m[id]; ok {
		return v
	}

	return nil
}

func (idx *index[K, T]) insert(id K, v Node[K, T]) bool {
	if idx == nil {
		return false
	}

	m := *idx
	m[id] = v
	return true
}
