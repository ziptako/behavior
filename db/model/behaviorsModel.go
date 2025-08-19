package model

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ BehaviorsModel = (*customBehaviorsModel)(nil)

type (
	// BehaviorsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customBehaviorsModel.
	BehaviorsModel interface {
		behaviorsModel
		// FindList 分页查询行为数据列表
		FindList(ctx context.Context, key string, userId int64, startTime, endTime *time.Time, page, pageSize int32) ([]*Behaviors, error)
		// CountList 统计符合条件的行为数据总数
		CountList(ctx context.Context, key string, userId int64, startTime, endTime *time.Time) (int64, error)
	}

	customBehaviorsModel struct {
		*defaultBehaviorsModel
	}
)

// NewBehaviorsModel returns a model for the database table.
func NewBehaviorsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) BehaviorsModel {
	return &customBehaviorsModel{
		defaultBehaviorsModel: newBehaviorsModel(conn, c, opts...),
	}
}

// FindList 分页查询行为数据列表
func (m *customBehaviorsModel) FindList(ctx context.Context, key string, userId int64, startTime, endTime *time.Time, page, pageSize int32) ([]*Behaviors, error) {
	// 构建查询条件
	var conditions []string
	var args []interface{}
	argIndex := 1

	// 基础条件：未删除
	conditions = append(conditions, "deleted_at IS NULL")

	// 按key过滤
	if key != "" {
		conditions = append(conditions, fmt.Sprintf("key = $%d", argIndex))
		args = append(args, key)
		argIndex++
	}

	// 按用户ID过滤
	if userId > 0 {
		conditions = append(conditions, fmt.Sprintf("user_id = $%d", argIndex))
		args = append(args, userId)
		argIndex++
	}

	// 按时间范围过滤
	if startTime != nil {
		conditions = append(conditions, fmt.Sprintf("created_at >= $%d", argIndex))
		args = append(args, *startTime)
		argIndex++
	}
	if endTime != nil {
		conditions = append(conditions, fmt.Sprintf("created_at <= $%d", argIndex))
		args = append(args, *endTime)
		argIndex++
	}

	// 分页参数
	offset := (page - 1) * pageSize
	args = append(args, pageSize, offset)

	// 构建完整查询
	whereClause := strings.Join(conditions, " AND ")
	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d",
		behaviorsRows, m.table, whereClause, argIndex, argIndex+1)

	var resp []*Behaviors
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, args...)
	return resp, err
}

// CountList 统计符合条件的行为数据总数
func (m *customBehaviorsModel) CountList(ctx context.Context, key string, userId int64, startTime, endTime *time.Time) (int64, error) {
	// 构建查询条件
	var conditions []string
	var args []interface{}
	argIndex := 1

	// 基础条件：未删除
	conditions = append(conditions, "deleted_at IS NULL")

	// 按key过滤
	if key != "" {
		conditions = append(conditions, fmt.Sprintf("key = $%d", argIndex))
		args = append(args, key)
		argIndex++
	}

	// 按用户ID过滤
	if userId > 0 {
		conditions = append(conditions, fmt.Sprintf("user_id = $%d", argIndex))
		args = append(args, userId)
		argIndex++
	}

	// 按时间范围过滤
	if startTime != nil {
		conditions = append(conditions, fmt.Sprintf("created_at >= $%d", argIndex))
		args = append(args, *startTime)
		argIndex++
	}
	if endTime != nil {
		conditions = append(conditions, fmt.Sprintf("created_at <= $%d", argIndex))
		args = append(args, *endTime)
		argIndex++
	}

	// 构建完整查询
	whereClause := strings.Join(conditions, " AND ")
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s", m.table, whereClause)

	var count int64
	err := m.QueryRowNoCacheCtx(ctx, &count, query, args...)
	return count, err
}
