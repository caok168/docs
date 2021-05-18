package main
/***
中序遍历
*/

func inorderTraversal(root *TreeNode) []int {
	res := []int{}
	var inorder func(*TreeNode)
	inorder = func(node *TreeNode) {
		if node == nil {
			return
		}

		inorder(node.Left)
		res = append(res, node.Val)
		inorder(node.Right)
	}

	inorder(root)
	return res
}

