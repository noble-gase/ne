package leveltree

// LevelNode 层级树泛型约束
type LevelNode[E comparable] interface {
	GetId() E
	GetPId() E
}

// LevelTree 菜单和组织单位等分类层级树
type LevelTree[T LevelNode[E], E comparable] struct {
	Data     T                  `json:"data"`
	Children []*LevelTree[T, E] `json:"children"`
}

func buildTree[T LevelNode[E], E comparable](classify map[E][]T, rootId E) []*LevelTree[T, E] {
	nodes := classify[rootId]
	count := len(nodes)
	root := make([]*LevelTree[T, E], 0, count)
	for i := 0; i < count; i++ {
		node := nodes[i]
		root = append(root, &LevelTree[T, E]{
			Data:     node,
			Children: buildTree(classify, node.GetId()),
		})
	}
	return root
}

// New 构建菜单和组织单位等分类层级树
//
//	data = 一组带有pid的数据
//	rootId = 树的起始ID
func New[T LevelNode[E], E comparable](data []T, rootId E) []*LevelTree[T, E] {
	classify := make(map[E][]T, 0)
	for _, v := range data {
		pid := v.GetPId()
		if _, ok := classify[pid]; !ok {
			classify[pid] = make([]T, 0)
		}
		classify[pid] = append(classify[pid], v)
	}
	return buildTree(classify, rootId)
}
