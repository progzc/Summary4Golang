package bit

// 常用的位操作: &、|、~、^

// 1.异或的特性（力扣示例：136、268、389、421）
// x ^ 0 = x
// x ^ 11111……1111 = ~x
// x ^ (~x) = 11111……1111
// x ^ x = 0
// a ^ b = c  => a ^ c = b  => b ^ c = a (交换律)
// a ^ b ^ c = a ^ (b ^ c) = (a ^ b）^ c (结合律)

// 2.构造特殊 Mask，将特殊位置放 0 或 1
// 将 x 最右边的 n 位清零， x & ( ~0 << n )
// 获取 x 的第 n 位值(0 或者 1)，(x >> n) & 1
// 获取 x 的第 n 位的幂值，x & (1 << (n - 1))
// 仅将第 n 位置为 1，x | (1 << n)
// 仅将第 n 位置为 0，x & (~(1 << n))
// 将 x 最高位至第 n 位(含)清零，x & ((1 << n) - 1)
// 将第 n 位至第 0 位(含)清零，x & (~((1 << (n + 1)) - 1)）

// 3.有特殊意义的 & 位操作运算（力扣示例：201、260、318、371、397、461、693）
// X & 1 == 1 判断是否是奇数(偶数)
// X & = (X - 1) 将最低位(LSB)的 1 清零
// X & -X 得到最低位(LSB)的 1
// X & ~X = 0
