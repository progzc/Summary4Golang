package leetcode_0536_construct_binary_tree_from_string

import (
	"strconv"
	"strings"
)

// 0536. 从字符串生成二叉树
// https://leetcode.cn/problems/construct-binary-tree-from-string/

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// str2tree 递归
// 时间复杂度: O(nlog(n))
// 空间复杂度: O(log(n))
func str2tree(s string) *TreeNode {
	// 当字符串长度为0时，构造空结点
	if len(s) == 0 {
		return nil
	}

	// 找到'('第一次出现的位置
	pos := strings.Index(s, "(")
	// 如果没找到，说明字符串里不再包含子树信息，为叶子节点，用它来直接构并返回
	if pos == -1 {
		val, _ := strconv.Atoi(s)
		return &TreeNode{Val: val}
	}

	// 构建根节点
	rootVal, _ := strconv.Atoi(s[0:pos])
	root := &TreeNode{Val: rootVal}
	// 记录起始位置，从pos开始；需要右括号的数量
	start, count := pos, 0
	for i := pos; i < len(s); i++ {
		if s[i] == '(' {
			count++
		} else if s[i] == ')' {
			count--
		}

		if count == 0 && start == pos {
			// 当count 为0，且起始位置是从第一个'(' 开始的，那么就去构建左子树
			root.Left = str2tree(s[start+1 : i])
			// 构建完之后更新起始位置
			start = i + 1
		} else if count == 0 && start != pos {
			// 如果起始位置不是第一次出现'(',就去构建右子树
			root.Right = str2tree(s[start+1 : i])
		}
	}

	return root
}
