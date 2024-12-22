package tree

type NodeSelector[T any] func(data T) bool

type Tree[K comparable, T any] struct {
	root Node[K, T]
	keys *index[K, T]
}

func NewTree[K comparable, T any]() *Tree[K, T] {
	return &Tree[K, T]{
		keys: &index[K, T]{},
	}
}

func (t *Tree[K, T]) Root() Node[K, T] {
	return t.root
}

func (t *Tree[K, T]) Add(id, pid K, data T) (added bool, exists bool) {
	child := &node[K, T]{
		id:   id,
		pid:  pid,
		data: data,
	}

	if t.keys.find(id) != nil {
		exists = true
		return
	}

	if t.root == nil {
		t.root = child
	} else {
		parent := t.keys.find(pid)
		if parent == nil {
			if t.root.GetParentID() == id {
				t.reroot(child)
			} else {
				return
			}
		} else {
			if t.root.GetParentID() == id {
				return
			}
			child.setParent(parent)
			parent.AddChildren(child)
		}
	}

	t.keys.insert(id, child)

	added = true
	return
}

func (t *Tree[K, T]) AddChild(child Node[K, T]) (added bool, exists bool) {
	if t.keys.find(child.GetID()) != nil {
		exists = true
		return
	}

	if t.root == nil {
		t.root = child
	} else {
		parent := t.keys.find(child.GetParentID())
		if parent == nil {
			if t.root.GetParentID() == child.GetID() {
				t.reroot(child)
			} else {
				return
			}
		} else {
			if t.root.GetParentID() == child.GetID() {
				return
			}
			child.setParent(parent)
			parent.AddChildren(child)
		}
	}

	t.keys.insert(child.GetID(), child)

	added = true
	return
}

func (t *Tree[K, T]) reroot(root Node[K, T]) {
	t.root.setParent(root)
	root.AddChildren(t.root)
	t.root = root
}

func (t *Tree[K, T]) Find(id K) (Node[K, T], bool) {
	if f := t.keys.find(id); f != nil {
		return f, true
	}
	return nil, false
}

func findAllBy[K comparable, T any](n Node[K, T], selector NodeSelector[T], c chan<- Node[K, T], topdown bool) {
	if topdown {
		if selector(n.GetData()) {
			c <- n
		}
	}
	for _, child := range n.GetChildren() {
		findAllBy(child, selector, c, topdown)
	}
	if !topdown {
		if selector(n.GetData()) {
			c <- n
		}
	}
}

func (t *Tree[K, T]) FindAllBy(selector NodeSelector[T], topdown bool) <-chan Node[K, T] {
	ret := make(chan Node[K, T])
	go func() {
		findAllBy(t.root, selector, ret, topdown)
		close(ret)
	}()

	return ret
}
