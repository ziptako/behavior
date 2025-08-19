# Behavior Service Database
```shell
 goctl model pg datasource -url "postgres://admin:123123@localhost:5432/edu?sslmode=disable" -c -d model --idea -s behavior --strict --style goZero -t *
```
## 概述

Behavior Service 数据库用于存储和管理用户行为数据，支持行为记录、查询、统计和分析功能。

## 数据库结构

### 1. 核心表

#### behaviors 表
- **用途**: 存储用户行为记录
- **主要字段**:
  - `id`: 主键ID
  - `key`: 行为标识键（如：login, logout, page_view, click等）
  - `user_id`: 用户ID
  - `data`: 行为数据（JSONB格式，支持复杂查询）
  - `created_at`: 创建时间
  - `updated_at`: 更新时间
  - `deleted_at`: 软删除时间

#### behavior_stats 表
- **用途**: 缓存行为统计数据，提高查询性能
- **主要字段**:
  - `id`: 主键ID
  - `key`: 行为标识键
  - `user_id`: 用户ID（NULL表示全局统计）
  - `date`: 统计日期
  - `count`: 统计数量
  - `created_at`: 创建时间
  - `updated_at`: 更新时间

### 2. 视图

#### active_behaviors 视图
- **用途**: 查询活跃的行为数据（排除软删除）
- **字段**: id, key, user_id, data, created_at, updated_at

#### behavior_stats_summary 视图
- **用途**: 行为统计汇总
- **字段**: key, user_id, total_count, first_date, last_date, active_days

### 3. 函数

#### get_user_behavior_stats 函数
```sql
SELECT * FROM behavior.get_user_behavior_stats(
    p_user_id => 123,
    p_key => 'login',
    p_start_date => '2024-01-01',
    p_end_date => '2024-12-31'
);
```
- **用途**: 获取用户行为统计
- **参数**:
  - `p_user_id`: 用户ID
  - `p_key`: 行为键（可选）
  - `p_start_date`: 开始日期（可选）
  - `p_end_date`: 结束日期（可选）

#### get_behavior_trend 函数
```sql
SELECT * FROM behavior.get_behavior_trend(
    p_key => 'login',
    p_user_id => 123,
    p_start_date => '2024-01-01',
    p_end_date => '2024-12-31'
);
```
- **用途**: 获取行为趋势数据
- **参数**:
  - `p_key`: 行为键
  - `p_user_id`: 用户ID（可选）
  - `p_start_date`: 开始日期（可选）
  - `p_end_date`: 结束日期（可选）

## 特性

### 1. 自动统计
- 通过触发器自动维护 `behavior_stats` 表
- 插入行为数据时自动更新统计
- 删除行为数据时自动减少统计

### 2. 软删除
- 支持软删除机制，删除的数据不会物理删除
- 通过 `deleted_at` 字段标记删除状态

### 3. 性能优化
- 针对常用查询场景创建了多个索引
- 使用 GIN 索引优化 JSONB 字段查询
- 统计表缓存提高查询性能

### 4. 数据完整性
- 时间戳逻辑性约束
- JSON 数据有效性约束
- 统计数量非负约束
- 唯一性约束

## 使用示例

### 1. 插入行为数据
```sql
INSERT INTO behavior.behaviors (key, user_id, data) VALUES
('login', 123, '{"ip": "192.168.1.1", "device": "mobile"}'),
('page_view', 123, '{"page": "/dashboard", "duration": 30}'),
('click', 123, '{"element": "button", "position": {"x": 100, "y": 200}}');
```

### 2. 查询用户行为
```sql
-- 查询用户所有行为
SELECT * FROM behavior.active_behaviors WHERE user_id = 123;

-- 查询特定类型行为
SELECT * FROM behavior.active_behaviors 
WHERE user_id = 123 AND key = 'login'
ORDER BY created_at DESC;

-- 查询时间范围内的行为
SELECT * FROM behavior.active_behaviors 
WHERE user_id = 123 
  AND created_at >= '2024-01-01'
  AND created_at < '2024-02-01';
```

### 3. 查询统计数据
```sql
-- 查询用户行为统计
SELECT * FROM behavior.behavior_stats 
WHERE user_id = 123 
ORDER BY date DESC;

-- 查询全局行为统计
SELECT key, SUM(count) as total_count
FROM behavior.behavior_stats 
GROUP BY key 
ORDER BY total_count DESC;
```

### 4. JSON 数据查询
```sql
-- 查询包含特定字段的行为
SELECT * FROM behavior.active_behaviors 
WHERE data ? 'ip';

-- 查询特定字段值的行为
SELECT * FROM behavior.active_behaviors 
WHERE data->>'device' = 'mobile';

-- 复杂 JSON 查询
SELECT * FROM behavior.active_behaviors 
WHERE data->'position'->>'x' = '100';
```

## 部署说明

1. 确保 PostgreSQL 版本支持 JSONB 类型（9.4+）
2. 执行 `init.sql` 文件创建数据库结构
3. 根据业务需求调整初始化数据
4. 配置适当的数据库连接参数

## 维护建议

1. 定期清理过期的软删除数据
2. 监控统计表的数据量，必要时进行归档
3. 根据查询模式调整索引策略
4. 定期分析查询性能并优化

## 扩展性

该数据库设计支持以下扩展：
- 添加新的行为类型
- 扩展 JSON 数据结构
- 增加更多统计维度
- 支持分区表（按时间分区）
- 支持数据归档策略