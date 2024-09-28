# 高频 SQL 50 题（基础版）

## 查询

### [1757. 可回收且低脂的产品](https://leetcode.cn/problems/recyclable-and-low-fat-products/)

![image-20240928223135808](assets/image-20240928223135808.png)

![image-20240928223153762](assets/image-20240928223153762.png)

```sql
# Schema
Create table If Not Exists Products (product_id int, low_fats ENUM('Y', 'N'), recyclable ENUM('Y','N'))
Truncate table Products
insert into Products (product_id, low_fats, recyclable) values ('0', 'Y', 'N')
insert into Products (product_id, low_fats, recyclable) values ('1', 'Y', 'Y')
insert into Products (product_id, low_fats, recyclable) values ('2', 'N', 'Y')
insert into Products (product_id, low_fats, recyclable) values ('3', 'Y', 'Y')
insert into Products (product_id, low_fats, recyclable) values ('4', 'N', 'N')

# Result
SELECT product_id FROM Products WHERE low_fats = 'Y' AND recyclable = 'Y';
```

### [584. 寻找用户推荐人](https://leetcode.cn/problems/find-customer-referee/)
