package math

import (
	"fmt"
)

func main() {
	tree := Tree{}
	tree.Add(0)
	tree.Add(1)
	tree.Add(2)
	tree.Add(3)
	tree.Add(4)
	tree.Add(5)
	tree.Add(6)
	tree.Add(7)
	tree.Add(8)
	tree.Add(9)

	tree.PreOrder(tree.RootNode)
	fmt.Println("")

	tree.InOrder(tree.RootNode)
	fmt.Println("")

	tree.PostOrder(tree.RootNode)
	fmt.Println("")
}
