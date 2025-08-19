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

type GetBehaviorLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetBehaviorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBehaviorLogic {
	return &GetBehaviorLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetBehavior 根据ID获取行为数据详情
func (l *GetBehaviorLogic) GetBehavior(in *behavior.GetBehaviorRequest) (*behavior.Behavior, error) {
	// 参数验证
	if in.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "[GB001] 行为数据ID必须大于0")
	}

	// 从数据库查询行为数据
	behaviorData, err := l.svcCtx.BehaviorsModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "[GB002] 行为数据不存在")
		}
		eInfo := "[GB003] 查询行为数据失败"
		l.Logger.Errorf("%v: %v", eInfo, err)
		return nil, status.Error(codes.Internal, eInfo)
	}

	// 检查是否已软删除
	if behaviorData.DeletedAt.Valid {
		return nil, status.Error(codes.NotFound, "[GB004] 行为数据已被删除")
	}

	// 转换为proto格式
	return &behavior.Behavior{
		Id:        behaviorData.Id,
		Key:       behaviorData.Key,
		UserId:    behaviorData.UserId,
		Data:      behaviorData.Data,
		CreatedAt: behaviorData.CreatedAt.UnixMilli(),
		UpdatedAt: behaviorData.UpdatedAt.UnixMilli(),
	}, nil
}
