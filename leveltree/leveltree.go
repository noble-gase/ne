package leveltree

// LevelNode 层级树泛型约束
type LevelNode[E comparable] interface {
	GetID() E
	GetPid() E
}

// LevelTree 菜单或分类层级树
type LevelTree[T LevelNode[E], E comparable] struct {
	Data     T
	Children []*LevelTree[T, E]
}

// New 构建菜单或分类层级树（data=按pid归类后的数据, pid=树的起始ID）
func New[T LevelNode[E], E comparable](data map[E][]T, pid E) []*LevelTree[T, E] {
	nodes := data[pid]
	count := len(nodes)
	root := make([]*LevelTree[T, E], 0, count)
	for i := 0; i < count; i++ {
		node := nodes[i]
		root = append(root, &LevelTree[T, E]{
			Data:     node,
			Children: New(data, node.GetID()),
		})
	}
	return root
}
