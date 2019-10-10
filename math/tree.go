package math

import "fmt"
type Object interface{}

type TreeNode struct {
	Data       Object
	LeftChild  *TreeNode
	RightChild *TreeNode
}

type Tree struct {
	RootNode *TreeNode
}

// Add
func (this *Tree) Add(object Object) {
	node := &TreeNode{Data: object}
	if this.RootNode == nil {
		this.RootNode = node
		return
	}
	queue := []*TreeNode{this.RootNode}
	for len(queue) != 0 {
		cur_node := queue[0]
		queue = queue[1:]

		if cur_node.LeftChild == nil {
			cur_node.LeftChild = node
			return
		} else {
			queue = append(queue, cur_node.LeftChild)
		}
		if cur_node.RightChild == nil {
			cur_node.RightChild = node
			return
		} else {
			queue = append(queue, cur_node.RightChild)
		}
	}
}

// BreadthTravel
func (this *Tree) BreadthTravel() {

	if this.RootNode == nil {
		return
	}
	queue := []*TreeNode{}
	queue = append(queue, this.RootNode)

	for len(queue) != 0 {
		cur_node := queue[0]
		queue = queue[1:]

		fmt.Printf("%v  ", cur_node.Data)

		if cur_node.LeftChild != nil {
			queue = append(queue, cur_node.LeftChild)
		}
		if cur_node.RightChild != nil {
			queue = append(queue, cur_node.RightChild)
		}
	}

}

// PreOrder
func (this *Tree) PreOrder(node *TreeNode) {
	if node == nil {
		return
	}
	this.PreOrder(node.LeftChild)
	this.PreOrder(node.RightChild)
}

// InOrder
func (this *Tree) InOrder(node *TreeNode) {
	if node == nil {
		return
	}
	this.InOrder(node.LeftChild)
	this.InOrder(node.RightChild)
}

// PostOrder
func (this *Tree) PostOrder(node *TreeNode)  {
	if node == nil {
		return
	}
	this.PostOrder(node.LeftChild)
	this.PostOrder(node.RightChild)
}
