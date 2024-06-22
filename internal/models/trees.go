package models

type Node struct {
	Name        string
	Description string
	Children    []*Node
}

type TreeNode struct {
	nodeTable map[string]*Node
	root      *Node
}

func InitTree() *TreeNode {
	items := make(map[string]*Node)
	node := &Node{
		Name:     "0",
		Children: []*Node{}}

	items[node.Name] = node

	treeNode := TreeNode{
		nodeTable: items,
		root:      node,
	}
	return &treeNode
}

func (tn *TreeNode) Add(name, desc, parentId string) {
	node := &Node{
		Name:        name,
		Description: desc,
		Children:    []*Node{}}
	if parentId == "" {
		tn.root = node
	} else {
		parent, ok := tn.nodeTable[parentId]
		if !ok {
			return
		}
		parent.Children = append(parent.Children, node)
	}
	tn.nodeTable[name] = node
}

func (tr *TreeNode) GetNodeChild(name string) []*Node {
	nodes := make([]*Node, 0, 10)

	parent, ok := tr.nodeTable[name]
	if !ok || parent.Children == nil {
		return nil
	}
	nodes = append(nodes, parent.Children...)

	return nodes
}

// func ShowNode(node *Node, prefix string) {
// 	if prefix == "" {
// 		fmt.Printf("%v\n\n", node.name)
// 	} else {
// 		fmt.Printf("%v %v\n\n", prefix, node.name)
// 	}
// 	for _, n := range node.children {
// 		ShowNode(n, prefix+"--")
// 	}
// }
