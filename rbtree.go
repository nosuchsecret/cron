package cron
import (
	//"fmt"
	//"sync"
//	"math/rand"
)

type RbtreeKey uint
type RbtreeKeyInt int

type RbtreeNode struct {
	key    RbtreeKey
	left   *RbtreeNode
	right  *RbtreeNode
	parent *RbtreeNode
	color  int
	data   interface{}
}

//type  RbtreeInsertPt func(root, node, sentinel RbtreeNode)

type Rbtree struct {
	num      int
	root     *RbtreeNode
	sentinel *RbtreeNode
	insert   func(root, node, sentinel *RbtreeNode)
}

func RbtreeInit(i func(root, node, sentinel *RbtreeNode)) (*Rbtree) {
	tree := &Rbtree{}
	s := &RbtreeNode{}
	rbtBlack(s)
	tree.root = s
	tree.sentinel = s
	tree.insert = i
	tree.num = 0
	return tree
}

func rbtRed(node *RbtreeNode) {
	node.color = 1
}

func rbtBlack(node *RbtreeNode) {
	node.color = 0
}

func rbtIsRed(node *RbtreeNode) bool {
	if node.color > 0 {
		return true
	}
	return false
}

func rbtIsBlack(node *RbtreeNode) bool {
	return !rbtIsRed(node)
}

func rbtCopyColor(n1, n2 *RbtreeNode) {
	n1.color = n2.color
}

func (tree *Rbtree) RbtreeMin(node, sentinel *RbtreeNode) *RbtreeNode {
	for {
		if node.left == sentinel {
			break
		}
		node = node.left
	}
	return node
}

func (tree *Rbtree)RbtreeInsert(node *RbtreeNode) {
	var root **RbtreeNode
	var temp, sentinel *RbtreeNode

	root = &tree.root
	sentinel = tree.sentinel
	if *root == sentinel {
		node.parent = nil
		node.left = sentinel
		node.right = sentinel
		rbtBlack(node)
		*root = node
		tree.num++

		return
	}

	tree.insert(*root, node, sentinel)

	for {
		if node != *root && rbtIsRed(node.parent) {
			if node.parent == node.parent.parent.left {
				temp = node.parent.parent.right

				if rbtIsRed(temp) {
					rbtBlack(node.parent)
					rbtBlack(temp)
					rbtRed(node.parent.parent)
					node = node.parent.parent
				} else {
					if node == node.parent.right {
						node = node.parent
						RbtreeLeftRotate(root, sentinel, node)
					}

					rbtBlack(node.parent)
					rbtRed(node.parent.parent)
					RbtreeRightRotate(root, sentinel, node.parent.parent)
				}
			} else {
				temp = node.parent.parent.left
				if rbtIsRed(temp) {
					rbtBlack(node.parent)
					rbtBlack(temp)
					rbtRed(node.parent.parent)
					node = node.parent.parent
				} else {
					if node == node.parent.left {
						node = node.parent
						RbtreeRightRotate(root, sentinel, node)

					}

					rbtBlack(node.parent)
					rbtRed(node.parent.parent)
					RbtreeLeftRotate(root, sentinel, node.parent.parent)
				}

			}
		} else {
			break
		}
	}

	rbtBlack(*root)
	tree.num++
}

func RbtreeInsertValue(temp, node, sentinel *RbtreeNode) {
	var p **RbtreeNode

	for {
		if node.key < temp.key {
			p = &temp.left
		} else {
			p = &temp.right
		}

		if *p == sentinel {
			break
		}

		temp = *p
	}

	*p = node
	node.parent = temp
	node.left = sentinel
	node.right = sentinel
	rbtRed(node)
}

func (tree *Rbtree) RbtreeDelete(node *RbtreeNode) {
	var red bool
	var root **RbtreeNode
	var sentinel, subst, temp, w *RbtreeNode

	root = &tree.root

	sentinel = tree.sentinel

	if node.left == sentinel {
		temp = node.right
		subst = node
	} else if node.right == sentinel {
		temp = node.left
		subst = node
	} else {
		subst = tree.RbtreeMin(node.right, sentinel)

		if subst.left != sentinel {
			temp = subst.left
		} else {
			temp = subst.right
		}
	}

	if subst == *root {
		*root = temp
		rbtBlack(temp)

		node.left = nil
		node.right = nil
		node.parent = nil
		node.key = 0

		tree.num--

		return
	}

	red = rbtIsRed(subst)

	if subst == subst.parent.left {
		subst.parent.left = temp

	} else {
		subst.parent.right = temp
	}

	if subst == node {
		temp.parent = subst.parent
	} else {
		if subst.parent == node {
			temp.parent = subst
		} else {
			temp.parent = subst.parent
		}

		subst.left = node.left
		subst.right = node.right
		subst.parent = node.parent
		rbtCopyColor(subst, node)

		if node == *root {
			*root = subst
		} else {
			if node == node.parent.left {
				node.parent.left = subst
			} else {
				node.parent.right = subst
			}
		}

		if subst.left != sentinel {
			subst.left.parent = subst
		}

		if subst.right != sentinel {
			subst.right.parent = subst
		}
	}

	node.left = nil
	node.right = nil
	node.parent = nil
	node.key = 0

	if red {
		tree.num--
		return
	}

	for {
		if temp != *root && rbtIsBlack(temp) {
			if temp == temp.parent.left {
				w = temp.parent.right
				if rbtIsRed(w) {
					rbtBlack(w)
					rbtRed(temp.parent)
					RbtreeLeftRotate(root, sentinel, temp.parent)
					w = temp.parent.right
				}

				if rbtIsBlack(w.left) && rbtIsBlack(w.right) {
					rbtRed(w)
					temp = temp.parent
				} else {
					if rbtIsBlack(w.right) {
						rbtBlack(w.left)
						rbtRed(w)
						RbtreeRightRotate(root, sentinel, w)
						w = temp.parent.right
					}

					rbtCopyColor(w, temp.parent)
					rbtBlack(temp.parent)
					rbtBlack(w.right)
					RbtreeLeftRotate(root, sentinel, temp.parent)
					temp = *root
				}
			} else {
				w = temp.parent.left
				if rbtIsRed(w) {
					rbtBlack(w)
					rbtRed(temp.parent)
					RbtreeRightRotate(root, sentinel, temp.parent)
					w = temp.parent.left
				}

				if rbtIsBlack(w.left) && rbtIsBlack(w.right) {
					rbtRed(w)
					temp = temp.parent
				} else {
					if rbtIsBlack(w.left) {
						rbtBlack(w.right)
						rbtRed(w)
						RbtreeLeftRotate(root, sentinel, w)
						w = temp.parent.left
					}

					rbtCopyColor(w, temp.parent)
					rbtBlack(temp.parent)
					rbtBlack(w.left)
					RbtreeRightRotate(root, sentinel, temp.parent)
					temp = *root


				}
			}
		} else {
			break
		}
	}

	rbtBlack(temp)
	tree.num--
}

func RbtreeLeftRotate(root **RbtreeNode, sentinel, node *RbtreeNode) {
	var temp *RbtreeNode
	temp = node.right
	node.right = temp.left
	if temp.left != sentinel {
		temp.left.parent = node
	}

	temp.parent = node.parent

	if node == *root {
		*root = temp
	} else if node == node.parent.left {
		node.parent.left = temp
	} else {
		node.parent.right = temp
	}

	temp.left = node
	node.parent = temp
}

func RbtreeRightRotate(root **RbtreeNode, sentinel, node *RbtreeNode) {
	var temp *RbtreeNode

	temp = node.left
	node.left = temp.right

	if temp.right != sentinel {
		temp.right.parent = node
	}

	temp.parent = node.parent

	if node == *root {
		*root = temp
	} else if node == node.parent.right {
		node.parent.right = temp
	} else {
		node.parent.left = temp
	}

	temp.right = node
	node.parent = temp
}

func testInsert(temp, node, sentinel *RbtreeNode) {
	var p **RbtreeNode
	var nd, td int
	for {
		if node.key < temp.key {
			//fmt.Println("set p to left")
			p = &temp.left
		} else if node.key > temp.key {
			//fmt.Println("set p to right")
			p = &temp.right
		} else {
			// ==
			//do something
			if i, ok := node.data.(int); ok {
				nd = i
			} else {
				nd = -1
			}
			if i, ok := temp.data.(int); ok {
				td = i
			} else {
				td = -1
			}

			if nd < td {
				p = &temp.left
			} else {
				p = &temp.right
			}
		}

		if *p == sentinel {
			break
		}

		temp = *p
	}

	*p = node
	node.parent = temp
	node.left = sentinel
	node.right = sentinel
	rbtRed(node)
}

func (tree *Rbtree) NodeNum() int {
	return tree.num
}

func (tree *Rbtree) FindMin() *RbtreeNode{
	if tree.root == tree.sentinel {
		return nil
	}
	return tree.RbtreeMin(tree.root, tree.sentinel)
}

//func Rbtree_main(n int) {
//	fmt.Println("n is", n)
//	tree := &Rbtree{}
//	sen := &RbtreeNode{}
//	RbtreeInit(tree, sen, testInsert)
//	c := 0
//	r := rand.New(rand.NewSource(99))
//	for {
//		rbn := &RbtreeNode{
//			data: r.Int31(),
//		}
//		RbtreeInsert(tree, rbn)
//		c++
//		if c > n {
//			break
//		}
//	}
//	c = 0
//	for {
//		node := RbtreeMin(tree.root, tree.sentinel)
//		RbtreeDelete(tree, node)
//		c++
//		if c > n {
//			break
//		}
//	}
//
//}


//func main() {
//	arr := make([]RbtreeNode, 10)
//	arr[0].data = 0
//	arr[1].data = 1
//	arr[2].data = 2
//	arr[3].data = 3
//	arr[4].data = 4
//	arr[5].data = 5
//
//	tree := RbtreeInit(testInsert)
//
//	tree.RbtreeInsert(&arr[2])
//	tree.RbtreeInsert(&arr[1])
//	tree.RbtreeInsert(&arr[0])
//	tree.RbtreeInsert(&arr[3])
//	tree.RbtreeInsert(&arr[5])
//	tree.RbtreeInsert(&arr[4])
//
//
//	fmt.Println(tree.FindMin())
//	node := tree.RbtreeMin(tree.root, tree.sentinel)
//	tree.RbtreeDelete(node)
//	fmt.Println(tree.FindMin())
//	node = tree.RbtreeMin(tree.root, tree.sentinel)
//	tree.RbtreeDelete(node)
//	fmt.Println(tree.FindMin())
//	node = tree.RbtreeMin(tree.root, tree.sentinel)
//	tree.RbtreeDelete(node)
//
//}
