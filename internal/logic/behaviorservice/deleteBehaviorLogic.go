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

type DeleteBehaviorLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteBehaviorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteBehaviorLogic {
	return &DeleteBehaviorLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// DeleteBehavior 删除行为数据
func (l *DeleteBehaviorLogic) DeleteBehavior(in *behavior.DeleteBehaviorRequest) (*behavior.DeleteBehaviorResponse, error) {
	// 参数验证
	if in.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "[DB001] 行为数据ID必须大于0")
	}

	// 先检查数据是否存在
	behaviorData, err := l.svcCtx.BehaviorsModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "[DB002] 行为数据不存在")
		}
		eInfo := "[DB003] 查询行为数据失败"
		l.Logger.Errorf("%v: %v", eInfo, err)
		return nil, status.Error(codes.Internal, eInfo)
	}

	// 检查是否已经被删除
	if behaviorData.DeletedAt.Valid {
		return nil, status.Error(codes.NotFound, "[DB004] 行为数据已被删除")
	}

	// 执行删除操作
	err = l.svcCtx.BehaviorsModel.Delete(l.ctx, in.Id)
	if err != nil {
		eInfo := "[DB005] 删除行为数据失败"
		l.Logger.Errorf("%v: %v", eInfo, err)
		return nil, status.Error(codes.Internal, eInfo)
	}

	return &behavior.DeleteBehaviorResponse{
		Success: true,
	}, nil
}
