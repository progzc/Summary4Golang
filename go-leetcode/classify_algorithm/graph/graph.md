#图的解法
##1.并查集
- 相关文档：https://leetcode.cn/circle/discuss/qmjuMW/
- 需要掌握的点：
    - 基本解法
    - 进一步优化方法
        - 两种合并优化（union方法）
            - 按树大小合并
            - 按高度合并（推荐）：不带路径压缩
        - 一种查询优化（find方法）
            - 路径压缩
    - 关于 *按秩合并* 和 *按高度合并* 的区别
        - 按秩合并：按高度合并+带路径压缩
    - 关于面试：在机试中如果写出*按秩合并的并查集*就是最完美解决方案了
    
- 代码模板
```java
class UnionFind{
    private int[] parent, rank, size; // 实际代码中，按秩求并和按大小求并选择其一
    public UnionFind(int[] parent) {
        this.parent = parent;
        this.rank = new int[parent.length];
        this.size = new int[parent.length];
        Arrays.fill(rank, 1); // 实际代码中，按秩求并和按大小求并选择其一
        Arrays.fill(size, 1); // 实际代码中，按秩求并和按大小求并选择其一
    }
    // 直接查找
    public int findDirect(int x) {
        if(parent[x] == x) return x;
        return findDirect(parent[x]);
    }
    // 带路径压缩的查找
    public int find(int x) {
        if(parent[x] == x) return x;
        return parent[x] = find(parent[x]);
    }
    // 直接求并
    public void unionDirect(int x, int y) {
        int xRoot = find(x), yRoot = find(y);
        if(xRoot != yRoot){
            parent[yRoot] = xRoot;
        }
    }
    // 按大小求并
    public void unionBySize(int x, int y){
        int xRoot = find(x), yRoot = find(y);
        if(xRoot != yRoot) { // 根节点不同才求并
            if(size[yRoot] <= size[xRoot]){
                parent[yRoot] = xRoot;
                size[xRoot] += size[yRoot];
            } else {
                parent[xRoot] = yRoot;
                size[yRoot] += size[xRoot];
            }
        }
    }
    // 按秩求并
    public void union(int x, int y){
        int xRoot = find(x), yRoot = find(y);
        if( xRoot != yRoot){
            if(rank[yRoot] <= rank[xRoot]) {
                parent[yRoot] = xRoot;
            } else {
                parent[xRoot] = yRoot;
            }
            if(rank[xRoot] == rank[yRoot]) {
                rank[xRoot]++;
            }
        }
    }
}
```
- 相关题目
    - [128.最长连续数列](https://leetcode-cn.com/problems/longest-consecutive-sequence/)
    - [200.岛屿数量](https://leetcode-cn.com/problems/number-of-islands/)
    - [323. 无向图中连通分量的数目](https://leetcode.cn/problems/number-of-connected-components-in-an-undirected-graph/)
    - [399. 除法求值](https://leetcode.cn/problems/evaluate-division/)
    - [547.省份数量](https://leetcode-cn.com/problems/number-of-provinces/)
    - [684. 冗余连接](https://leetcode.cn/problems/redundant-connection/)
    - [695. 岛屿的最大面积](https://leetcode.cn/problems/max-area-of-island/)
    - [785. 判断二分图](https://leetcode.cn/problems/is-graph-bipartite/)
    - [839. 相似字符串组](https://leetcode.cn/problems/similar-string-groups/)
    - [1631. 最小体力消耗路径](https://leetcode.cn/problems/path-with-minimum-effort/)
##2.tarjan强连通分量
##3.Dijkstra最短路径
- [宫水三叶】涵盖所有的「存图方式」与「最短路算法（详尽注释）」](https://leetcode.cn/problems/network-delay-time/solution/gong-shui-san-xie-yi-ti-wu-jie-wu-chong-oghpz/)
- 适用于路径权重都是正数
- 本质是广度优先搜索+贪心算法
- 单源
- 代码模板
```c#
public int Dijkstra(int k){ //k表示初始点， this.n, n = 终点
        Array.Fill(distance, int.MaxValue/2); //搜索最短路径前以inf来代替未记录的两节点间距
        distance[k] = 0; //k代表，出发点，从k到它自己的距离为0

        for(int i = 1; i <= n; i++){
            int t = -1; //两者不相邻，先初始-1
            for(int j = 1; j <= n; j++){
                //t == -1为true，第一次可执行，后面被j赋值
                if(!visited[j] && (t == -1 || distance[t] > distance[j])) //找还未访问的节点
                    t = j;
            }
            visited[t] = true; //表示为已访问元素

            for(int j = 1; j <= n; j++) //对比路径大小
                distance[j] = Math.Min(distance[j], distance[t] + nodes[t][j]); //比较最短距离
        }
        return distance[n];
    }
```
- 相关题目
    - [743. 网络延迟时间](https://leetcode.cn/problems/network-delay-time/)
##4.Bellman-Ford算法
- 适用于路径权重包含负数
- 单源
- 相关题目
    - [787. K 站中转内最便宜的航班](https://leetcode.cn/problems/cheapest-flights-within-k-stops/)
##5.Floyd-Warshall算法
    - 多源
    
