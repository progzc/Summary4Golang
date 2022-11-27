package leetcode_0235_lowest_common_ancestor_of_a_binary_search_tree

// 235. 二叉搜索树的最近公共祖先
// https://leetcode.cn/problems/lowest-common-ancestor-of-a-binary-search-tree/

// 题目同:
// 236. 二叉树的最近公共祖先
// https://leetcode.cn/problems/lowest-common-ancestor-of-a-binary-tree/

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// lowestCommonAncestor 递归
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 思路：
//	定义：f(x)表示x节点的子树中是否包含p节点或q节点,包含则为true，否则为false
//	那么：最近公共祖先一定满足如下条件
//		(f(x_L) && f(x_R)) || ((x==p||x==q) && (f(x_L)||f(x_R)))
//	解释：一共只有如下4种情况
//	f(x_L) && f(x_R)：
//		a.x的左子树包含p,且x的右子树包含q，则x即为最近公共祖先
//		b.x的左子树包含q,且x的右子树包含p，则x即为最近公共祖先
//	(x==p||x==q) && (f(x_L)||f(x_R)：
//		c.x为p节点，且x的左或右子树包含q，则p即为最近公共祖先
//		d.x为q节点，且x的左或右子树包含p，则q即为最近公共祖先
// 注意: 此种解法与 【236. 二叉树的最近公共祖先】相同, 没有充分利用到 【二叉搜索树】这个条件。
//		虽然题目可以通过，但是效率还可以提高。
func lowestCommonAncestor(root, p, q *TreeNode) *TreeNode {
	if root == nil {
		return nil
	}
	if root.Val == p.Val || root.Val == q.Val {
		return root
	}

	left := lowestCommonAncestor(root.Left, p, q)
	right := lowestCommonAncestor(root.Right, p, q)
	if left != nil && right != nil {
		return root
	}
	if left == nil {
		return right
	}
	return left
}

// lowestCommonAncestor_4
// 注意：对比与 lowestCommonAncestor 写法的异同
func lowestCommonAncestor_4(root, p, q *TreeNode) *TreeNode {
	if root == nil {
		return nil
	}
	if root == p || root == q {
		return root
	}
	x := lowestCommonAncestor(root.Left, p, q)
	y := lowestCommonAncestor(root.Right, p, q)
	if x != nil && y != nil {
		return root
	}
	if x == nil {
		return y
	}

	if y == nil {
		return x
	}
	return nil
}

// lowestCommonAncestor_2 二次遍历
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 思路：
//	a.分别找到从根节点到p,q节点的路径
//	b.找到两个路径中的最后一个相同的节点
func lowestCommonAncestor_2(root, p, q *TreeNode) *TreeNode {
	var getPath func(root, target *TreeNode) []*TreeNode
	getPath = func(root, target *TreeNode) []*TreeNode {
		var path []*TreeNode
		node := root
		for node != target {
			path = append(path, node)
			if target.Val < node.Val {
				node = node.Left
			} else {
				node = node.Right
			}
		}
		path = append(path, node)
		return path
	}

	var ans *TreeNode
	pPath := getPath(root, p)
	qPath := getPath(root, q)
	for i := 0; i < len(pPath) && i < len(qPath) && pPath[i] == qPath[i]; i++ {
		ans = pPath[i]
	}
	return ans
}

// lowestCommonAncestor_3 一次遍历
// 时间复杂度: O(n)
// 空间复杂度: O(1)
// 思路：在二次遍历的基础上可以同时进行搜索
func lowestCommonAncestor_3(root, p, q *TreeNode) *TreeNode {
	ans := root
	for {
		if p.Val < ans.Val && q.Val < ans.Val {
			ans = ans.Left
		} else if p.Val > ans.Val && q.Val > ans.Val {
			ans = ans.Right
		} else {
			return ans
		}
	}
}
