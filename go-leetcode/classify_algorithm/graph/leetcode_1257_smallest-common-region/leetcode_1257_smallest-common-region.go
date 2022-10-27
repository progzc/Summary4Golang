package leetcode_1257_smallest_common_region

// 1257. 最小公共区域
// https://leetcode.cn/problems/smallest-common-region/

// findSmallestRegion
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func findSmallestRegion(regions [][]string, region1 string, region2 string) string {
	p := make(map[string]string)
	for _, r := range regions {
		f := r[0]
		for _, s := range r[1:] {
			p[s] = f
		}
	}

	visited := make(map[string]bool)
	visited[region1] = true
	visited[region2] = true
	for region1 != "" {
		region1 = p[region1]
		if visited[region1] {
			return region1
		}
		visited[region1] = true
	}
	for region2 != "" {
		region2 = p[region2]
		if visited[region2] {
			return region2
		}
		visited[region2] = true
	}
	return ""
}
