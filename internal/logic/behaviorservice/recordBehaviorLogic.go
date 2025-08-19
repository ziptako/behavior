package behaviorservicelogic

import (
	"context"
	"database/sql"
	"time"

	"github.com/ziptako/behavior/behavior"
	"github.com/ziptako/behavior/db/model"
	"github.com/ziptako/behavior/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RecordBehaviorLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRecordBehaviorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecordBehaviorLogic {
	return &RecordBehaviorLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// RecordBehavior 记录行为数据
func (l *RecordBehaviorLogic) RecordBehavior(in *behavior.RecordBehaviorRequest) (*behavior.RecordBehaviorResponse, error) {
	// 参数验证
	if in.Key == "" {
		return nil, status.Error(codes.InvalidArgument, "[RB001] 行为标识键不能为空")
	}
	if in.UserId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "[RB002] 用户ID必须大于0")
	}
	if in.Data == "" {
		return nil, status.Error(codes.InvalidArgument, "[RB003] 行为数据不能为空")
	}

	// 构建数据库模型
	behaviorData := &model.Behaviors{
		Key:       in.Key,
		UserId:    in.UserId,
		Data:      in.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: sql.NullTime{Valid: false}, // 未删除
	}

	// 插入数据库
	result, err := l.svcCtx.BehaviorsModel.Insert(l.ctx, behaviorData)
	if err != nil {
		eInfo := "[RB004] 记录行为数据失败"
		l.Logger.Errorf("%v: %v", eInfo, err)
		return nil, status.Error(codes.Internal, eInfo)
	}

	// 获取插入的ID
	id, err := result.LastInsertId()
	if err != nil {
		// 如果无法获取ID，仍然返回成功，但ID为0
		l.Logger.Infof("无法获取插入的行为数据ID: %v", err)
		id = 0
	}

	return &behavior.RecordBehaviorResponse{
		Success: true,
		Id:      id,
	}, nil
}
