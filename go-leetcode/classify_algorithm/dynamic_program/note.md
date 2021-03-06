# 动态规划
- [力扣上的DP问题分类汇总](https://leetcode-cn.com/circle/article/NfHhXD/)
- [leetcode动态规划题目总结](https://leetcode-cn.com/circle/article/2Xxlw3/)
- [背包问题九讲](https://github.com/tianyicui/pack/blob/master/V2.pdf)
- [背包问题解法总结](https://oi-wiki.org/dp/dynamic/)

## 1. 背包DP
[背包问题题目的归纳解法](https://leetcode-cn.com/problems/last-stone-weight-ii/solution/yi-pian-wen-zhang-chi-tou-bei-bao-wen-ti-5lfv/)
### 1.1 分类解题模板
  - 背包问题大体的解题模板是两层循环，分别遍历物品nums和背包容量target，然后写转移方程
  - 根据背包的分类我们确定物品和容量遍历的先后顺序，根据问题的分类我们确定状态转移方程的写法
### 1.2 首先是背包分类的模板：
  - 0/1背包：外循环nums,内循环target,target倒序且target>=nums\[i\];
  - 完全背包：外循环nums,内循环target,target正序且target>=nums\[i\];
  - 组合背包：外循环target,内循环nums,target正序且target>=nums\[i\];
  - 分组背包：这个比较特殊，需要三重循环：外循环背包bags,内部两层循环根据题目的要求转化为1,2,3三种背包类型的模板
### 1.3 然后是问题分类的模板：
  - 最值问题: dp\[i\] = max/min(dp\[i\], dp\[i-nums\]+1)或dp\[i\] = max/min(dp\[i\], dp\[i-num\]+nums);
  - 存在问题(bool)：dp\[i\]=dp\[i\]||dp\[i-num\];
  - 组合问题：dp\[i\]+=dp\[i-num\];

### 1.4 经典例题
  - 0/1背包
    - [416.分割等和子集](https://leetcode-cn.com/problems/partition-equal-subset-sum/)
      - 0/1背包存在性问题
    - [494.目标和](https://leetcode-cn.com/problems/target-sum/)
      - 0/1背包不考虑元素顺序的组合问题
    - [1049.最后一块石头的重量II](https://leetcode-cn.com/problems/last-stone-weight-ii/)
      - 0/1背包最值问题

  - 完全背包
    - [279.完全平方数](https://leetcode-cn.com/problems/perfect-squares/)
      - 完全背包最值问题
    - [322.零钱兑换](https://leetcode-cn.com/problems/coin-change/)
      - 完全背包最值问题
    - [518.零钱兑换II](https://leetcode-cn.com/problems/coin-change-2/)
      - 完全背包不考虑顺序的组合问题

  - 组合背包
    - [377.组合总和Ⅳ](https://leetcode-cn.com/problems/combination-sum-iv/)
      - 考虑顺序的组合问题
  
  - 分组背包
    - [1155.掷骰子的N种方法](https://leetcode-cn.com/problems/number-of-dice-rolls-with-target-sum/)
      - 分组0/1背包的组合问题
    
 - 多维背包
   - [879.盈利计划](https://leetcode-cn.com/problems/profitable-schemes/)
      - 多维0/1背包
## 2. 线性DP
  - [53.最大子数组和](https://leetcode-cn.com/problems/maximum-subarray/)
  - [300.最长递增子序列](https://leetcode-cn.com/problems/longest-increasing-subsequence/)
  - [1143.最长公共子序列](https://leetcode-cn.com/problems/longest-common-subsequence/)
    - 二维动态规划
## 3. 区间DP
  - [5.最长回文子串](https://leetcode-cn.com/problems/longest-palindromic-substring/)
  - [516.最长回文子序列](https://leetcode-cn.com/problems/longest-palindromic-subsequence/)
  - [647.回文子串](https://leetcode-cn.com/problems/palindromic-substrings/)
  - [730.统计不同回文子序列](https://leetcode-cn.com/problems/count-different-palindromic-subsequences/)
