CREATE SCHEMA IF NOT EXISTS behavior;

-- =========================================================
-- 1. Behavior Service 表 - 行为数据管理
-- =========================================================

-- 行为数据表
CREATE TABLE behavior.behaviors
(
    id         BIGSERIAL PRIMARY KEY,
    key        VARCHAR(100) NOT NULL CHECK (LENGTH(TRIM(key)) > 0),
    user_id    BIGINT       NOT NULL,
    data       JSONB        NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    
    -- 确保时间戳的逻辑性
    CONSTRAINT chk_behaviors_timestamps CHECK (
        created_at <= updated_at AND
        (deleted_at IS NULL OR deleted_at >= created_at)
    ),
    
    -- 确保data字段是有效的JSON
    CONSTRAINT chk_behaviors_data_valid_json CHECK (
        data IS NOT NULL
    )
);

-- =========================================================
-- 2. 触发器函数
-- =========================================================

-- 创建触发器函数自动更新updated_at字段
CREATE OR REPLACE FUNCTION behavior.update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- =========================================================
-- 3. 触发器
-- =========================================================

-- 行为数据表触发器 - 自动更新updated_at
CREATE TRIGGER trigger_update_behaviors_updated_at
    BEFORE UPDATE ON behavior.behaviors
    FOR EACH ROW
    EXECUTE FUNCTION behavior.update_updated_at_column();

-- =========================================================
-- 4. 索引优化
-- =========================================================

-- 行为数据表基础索引
CREATE INDEX idx_behaviors_key ON behavior.behaviors (key) WHERE deleted_at IS NULL;
CREATE INDEX idx_behaviors_user_id ON behavior.behaviors (user_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_behaviors_key_user_id ON behavior.behaviors (key, user_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_behaviors_created_at ON behavior.behaviors (created_at) WHERE deleted_at IS NULL;
CREATE INDEX idx_behaviors_updated_at ON behavior.behaviors (updated_at);
CREATE INDEX idx_behaviors_deleted_at ON behavior.behaviors (deleted_at) WHERE deleted_at IS NOT NULL;
CREATE INDEX idx_behaviors_active ON behavior.behaviors (id) WHERE deleted_at IS NULL;

-- 时间范围查询优化索引
CREATE INDEX idx_behaviors_key_time_range ON behavior.behaviors (key, created_at) WHERE deleted_at IS NULL;
CREATE INDEX idx_behaviors_user_time_range ON behavior.behaviors (user_id, created_at) WHERE deleted_at IS NULL;
CREATE INDEX idx_behaviors_key_user_time_range ON behavior.behaviors (key, user_id, created_at) WHERE deleted_at IS NULL;

-- JSON数据查询优化索引（GIN索引用于JSONB字段）
CREATE INDEX idx_behaviors_data_gin ON behavior.behaviors USING GIN (data) WHERE deleted_at IS NULL;

-- =========================================================
-- 5. 表和字段注释
-- =========================================================

-- 行为数据表注释
COMMENT ON TABLE behavior.behaviors IS '行为数据表，存储用户行为记录';
COMMENT ON COLUMN behavior.behaviors.id IS '主键ID';
COMMENT ON COLUMN behavior.behaviors.key IS '行为标识键，用于区分不同类型的行为';
COMMENT ON COLUMN behavior.behaviors.user_id IS '用户ID，关联用户表';
COMMENT ON COLUMN behavior.behaviors.data IS '行为数据，以JSONB格式存储，支持复杂查询';
COMMENT ON COLUMN behavior.behaviors.created_at IS '创建时间';
COMMENT ON COLUMN behavior.behaviors.updated_at IS '更新时间，通过触发器自动维护';
COMMENT ON COLUMN behavior.behaviors.deleted_at IS '软删除时间，NULL表示未删除';

-- =========================================================
-- 7. 备注说明
-- =========================================================

/*
数据库设计说明：

1. 核心设计理念：
   - 专注于行为数据的存储和基础查询
   - 移除复杂的统计缓存功能，保持数据库设计简洁
   - 统计需求可通过Redis等缓存系统实现

2. 主要特性：
   - 软删除支持（deleted_at字段）
   - JSONB数据存储，支持灵活的行为数据结构
   - 自动时间戳维护
   - 完善的索引优化，支持高效查询

3. 扩展性：
   - 如需统计功能，建议使用Redis进行实时统计
   - 可通过应用层实现复杂的数据分析
   - 数据库专注于数据存储和基础查询性能

4. 性能优化：
   - 针对常用查询场景创建了复合索引
   - GIN索引支持JSONB字段的高效查询
   - 软删除索引优化，避免查询已删除数据

5. 使用建议：
   - 行为数据直接存储到behaviors表
   - 统计需求通过Redis等缓存系统实现
   - 复杂分析可通过数据仓库或分析系统处理
*/