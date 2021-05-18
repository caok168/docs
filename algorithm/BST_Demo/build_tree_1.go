package main

func buildTree(inorder []int, postorder []int) *TreeNode {
	idxMap := map[int]int{}
	for i, k := range inorder{
		idxMap[k] = i
	}

	var build func(int, int) *TreeNode
	build = func(inorderLeft, inorderRight int) *TreeNode {
		if inorderLeft > inorderRight {
			return nil
		}

		val := postorder[len(postorder) - 1]
		postorder = postorder[:len(postorder) - 1]
		treeIndex := idxMap[val]
		root := &TreeNode{
			Val: val,
		}

		root.Right = build(treeIndex + 1, inorderRight)
		root.Left = build(inorderLeft, treeIndex - 1)

		return root
	}

	return build(0, len(inorder) - 1)
}
