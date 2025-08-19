package behaviorservicelogic

import (
	"context"
	"time"

	"github.com/ziptako/behavior/behavior"
	"github.com/ziptako/behavior/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ListBehaviorsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListBehaviorsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListBehaviorsLogic {
	return &ListBehaviorsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ListBehaviors 分页查询行为数据列表
func (l *ListBehaviorsLogic) ListBehaviors(in *behavior.ListBehaviorsRequest) (*behavior.ListBehaviorsResponse, error) {
	// 参数验证
	if in.Page <= 0 {
		in.Page = 1
	}
	if in.PageSize <= 0 {
		in.PageSize = 10
	}
	if in.PageSize > 100 {
		return nil, status.Error(codes.InvalidArgument, "[LB001] 每页数量不能超过100")
	}

	// 处理时间参数
	var startTime, endTime *time.Time
	if in.StartTime > 0 {
		t := time.UnixMilli(in.StartTime)
		startTime = &t
	}
	if in.EndTime > 0 {
		t := time.UnixMilli(in.EndTime)
		endTime = &t
	}

	// 查询数据列表
	behaviorList, err := l.svcCtx.BehaviorsModel.FindList(l.ctx, in.Key, in.UserId, startTime, endTime, in.Page, in.PageSize)
	if err != nil {
		eInfo := "[LB002] 查询行为数据列表失败"
		l.Logger.Errorf("%v: %v", eInfo, err)
		return nil, status.Error(codes.Internal, eInfo)
	}

	// 查询总数
	total, err := l.svcCtx.BehaviorsModel.CountList(l.ctx, in.Key, in.UserId, startTime, endTime)
	if err != nil {
		eInfo := "[LB003] 查询行为数据总数失败"
		l.Logger.Errorf("%v: %v", eInfo, err)
		return nil, status.Error(codes.Internal, eInfo)
	}

	// 转换为proto格式
	var list []*behavior.Behavior
	for _, item := range behaviorList {
		list = append(list, &behavior.Behavior{
			Id:        item.Id,
			Key:       item.Key,
			UserId:    item.UserId,
			Data:      item.Data,
			CreatedAt: item.CreatedAt.UnixMilli(),
			UpdatedAt: item.UpdatedAt.UnixMilli(),
		})
	}

	return &behavior.ListBehaviorsResponse{
		List:  list,
		Total: total,
		Page:  in.Page,
		Size:  in.PageSize,
	}, nil
}
