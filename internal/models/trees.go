package models

type Node struct {
	Name        string
	Description string
	Parent      *Node
	Children    []*Node
	Visible     bool
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

func (tn *TreeNode) Add(name, desc, parentId string, visible bool) {
	node := &Node{
		Name:        name,
		Description: desc,
		Children:    []*Node{},
		Parent:      &Node{},
		Visible:     visible}
	if parentId == "" {
		tn.root = node
	} else {
		parent, ok := tn.nodeTable[parentId]
		if !ok {
			return
		}
		node.Parent = parent
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

func (tr *TreeNode) GetNode(name string) *Node {
	node, ok := tr.nodeTable[name]
	if !ok {
		return nil
	}
	return node
}

func (tr *TreeNode) GetParentNode(name string) *Node {
	node, ok := tr.nodeTable[name]
	if !ok {
		return nil
	}
	return node.Parent
}
