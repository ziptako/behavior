package behaviorservicelogic

import (
	"context"
	"errors"

	"github.com/ziptako/behavior/behavior"
	"github.com/ziptako/behavior/db/model"
	"github.com/ziptako/behavior/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BatchDeleteBehaviorLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchDeleteBehaviorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteBehaviorLogic {
	return &BatchDeleteBehaviorLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// BatchDeleteBehavior 批量删除行为数据
func (l *BatchDeleteBehaviorLogic) BatchDeleteBehavior(in *behavior.BatchDeleteBehaviorRequest) (*behavior.BatchDeleteBehaviorResponse, error) {
	// 参数验证
	if len(in.Ids) == 0 {
		return nil, status.Error(codes.InvalidArgument, "[BDB001] 批量删除ID列表不能为空")
	}
	if len(in.Ids) > 100 {
		return nil, status.Error(codes.InvalidArgument, "[BDB002] 批量删除数量不能超过100条")
	}

	var failedIds []int64
	deletedCount := int64(0)

	// 逐个处理每个ID
	for _, id := range in.Ids {
		// 验证ID
		if id <= 0 {
			l.Logger.Infof("批量删除中发现无效ID: %d", id)
			failedIds = append(failedIds, id)
			continue
		}

		// 先检查数据是否存在
		behaviorData, err := l.svcCtx.BehaviorsModel.FindOne(l.ctx, id)
		if err != nil {
			if errors.Is(err, model.ErrNotFound) {
				l.Logger.Infof("批量删除中发现不存在的数据，ID: %d", id)
				failedIds = append(failedIds, id)
				continue
			}
			l.Logger.Errorf("批量删除时查询数据失败，ID: %d, 错误: %v", id, err)
			failedIds = append(failedIds, id)
			continue
		}

		// 检查是否已经被删除
		if behaviorData.DeletedAt.Valid {
			l.Logger.Infof("批量删除中发现已删除的数据，ID: %d", id)
			failedIds = append(failedIds, id)
			continue
		}

		// 执行删除操作
		err = l.svcCtx.BehaviorsModel.Delete(l.ctx, id)
		if err != nil {
			l.Logger.Errorf("批量删除失败，ID: %d, 错误: %v", id, err)
			failedIds = append(failedIds, id)
			continue
		}

		deletedCount++
	}

	return &behavior.BatchDeleteBehaviorResponse{
		DeletedCount: deletedCount,
		FailedIds:    failedIds,
	}, nil
}
