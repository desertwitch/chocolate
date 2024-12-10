package tree

type AttributeSelector interface {
	Has(interface{}) bool
}

type Node[K comparable, T AttributeSelector] interface {
	GetID() K
	GetParentID() K
	GetChildren() []Node[K, T]
	GetParent() Node[K, T]
	AddChildren(...Node[K, T])
	ReplaceChildren(...Node[K, T])
	setParent(v Node[K, T])
	GetData() T
	SetData(T)
}

type node[K comparable, T AttributeSelector] struct {
	id       K
	pid      K
	parent   Node[K, T]
	data     T
	children []Node[K, T]
}

func (n *node[K, T]) GetID() K {
	return n.id
}

func (n *node[K, T]) GetParentID() K {
	return n.pid
}

func (n *node[K, T]) GetChildren() []Node[K, T] {
	return n.children
}

func (n *node[K, T]) GetParent() Node[K, T] {
	return n.parent
}

func (n *node[K, T]) AddChildren(children ...Node[K, T]) {
	if n.children == nil {
		n.children = []Node[K, T]{}
	}
	n.children = append(n.children, children...)
}

func (n *node[K, T]) ReplaceChildren(children ...Node[K, T]) {
	n.children = []Node[K, T]{}
	n.AddChildren(children...)
}

func (n *node[K, T]) setParent(parent Node[K, T]) {
	if parent == nil || parent.GetID() == n.GetID() {
		return
	}
	n.parent = parent
	n.pid = parent.GetID()
}

func (n *node[K, T]) GetData() T {
	return n.data
}

func (n *node[K, T]) SetData(data T) {
	n.data = data
}
