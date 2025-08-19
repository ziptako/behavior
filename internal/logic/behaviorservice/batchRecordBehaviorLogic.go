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

type BatchRecordBehaviorLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchRecordBehaviorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchRecordBehaviorLogic {
	return &BatchRecordBehaviorLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// BatchRecordBehavior 批量记录行为数据
func (l *BatchRecordBehaviorLogic) BatchRecordBehavior(in *behavior.BatchRecordBehaviorRequest) (*behavior.BatchRecordBehaviorResponse, error) {
	// 参数验证
	if len(in.Behaviors) == 0 {
		return nil, status.Error(codes.InvalidArgument, "[BRB001] 批量记录数据不能为空")
	}
	if len(in.Behaviors) > 100 {
		return nil, status.Error(codes.InvalidArgument, "[BRB002] 批量记录数量不能超过100条")
	}

	var successList []*behavior.RecordBehaviorResponse
	var failList []*behavior.RecordBehaviorRequest
	successCount := int64(0)
	failCount := int64(0)

	// 逐个处理每条记录
	for _, behaviorReq := range in.Behaviors {
		// 验证单条记录
		if behaviorReq.Key == "" {
			l.Logger.Infof("批量记录中发现无效数据：行为标识键为空，用户ID: %d", behaviorReq.UserId)
			failList = append(failList, behaviorReq)
			failCount++
			continue
		}
		if behaviorReq.UserId <= 0 {
			l.Logger.Infof("批量记录中发现无效数据：用户ID无效，Key: %s", behaviorReq.Key)
			failList = append(failList, behaviorReq)
			failCount++
			continue
		}
		if behaviorReq.Data == "" {
			l.Logger.Infof("批量记录中发现无效数据：行为数据为空，Key: %s, 用户ID: %d", behaviorReq.Key, behaviorReq.UserId)
			failList = append(failList, behaviorReq)
			failCount++
			continue
		}

		// 构建数据库模型
		behaviorData := &model.Behaviors{
			Key:       behaviorReq.Key,
			UserId:    behaviorReq.UserId,
			Data:      behaviorReq.Data,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: sql.NullTime{Valid: false}, // 未删除
		}

		// 插入数据库
		result, err := l.svcCtx.BehaviorsModel.Insert(l.ctx, behaviorData)
		if err != nil {
			l.Logger.Errorf("批量记录失败，Key: %s, 用户ID: %d, 错误: %v", behaviorReq.Key, behaviorReq.UserId, err)
			failList = append(failList, behaviorReq)
			failCount++
			continue
		}

		// 获取插入的ID
		id, err := result.LastInsertId()
		if err != nil {
			// 如果无法获取ID，仍然认为插入成功，但ID为0
			l.Logger.Infof("无法获取插入的行为数据ID，Key: %s, 用户ID: %d, 错误: %v", behaviorReq.Key, behaviorReq.UserId, err)
			id = 0
		}

		successList = append(successList, &behavior.RecordBehaviorResponse{
			Success: true,
			Id:      id,
		})
		successCount++
	}

	return &behavior.BatchRecordBehaviorResponse{
		SuccessCount: successCount,
		FailCount:    failCount,
		SuccessList:  successList,
		FailList:     failList,
	}, nil
}
