package leetcode_1192_critical_connections_in_a_network

// 1192.查找集群内的「关键连接」
// https://leetcode-cn.com/problems/critical-connections-in-a-network/

// criticalConnections tarjan强连通分量算法
// 思路：所有在环上的边都不是关键连接，所有不在环上的边都是关键连接
//	问题转化为：如何寻找环?如何将整个环看作一个点?
func criticalConnections(n int, connections [][]int) [][]int {
	// TODO
	return nil
}
