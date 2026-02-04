package treekit

// Node 节点泛型约束
type Node[E comparable] interface {
	ID() E
	BelongTo() E
}

// TreeNode 树定义
type TreeNode[T Node[E], E comparable] struct {
	ID       E                 `json:"id"`
	BelongTo E                 `json:"belong_to,omitempty"`
	Data     T                 `json:"data"`
	Kids     []*TreeNode[T, E] `json:"kids,omitempty"`
}

func buildTree[T Node[E], E comparable](group map[E][]T, rootId E) []*TreeNode[T, E] {
	nodes := group[rootId]
	count := len(nodes)
	root := make([]*TreeNode[T, E], 0, count)
	for _, v := range nodes {
		root = append(root, &TreeNode[T, E]{
			ID:       v.ID(),
			BelongTo: v.BelongTo(),
			Data:     v,
			Kids:     buildTree(group, v.ID()),
		})
	}
	return root
}

// NewTree 构建一颗树
//
//	[data] 一组带有`pid`的数据
//	[rootId] 树的起始ID
func NewTree[T Node[E], E comparable](data []T, rootId E) []*TreeNode[T, E] {
	group := make(map[E][]T, 0)
	for _, v := range data {
		pid := v.BelongTo()
		if _, ok := group[pid]; !ok {
			group[pid] = make([]T, 0)
		}
		group[pid] = append(group[pid], v)
	}
	return buildTree(group, rootId)
}
