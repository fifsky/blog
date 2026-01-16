package openapi

import (
	"context"
	"fmt"
	"strconv"

	"app/config"
	"app/pkg/aesutil"
	"app/pkg/wechat"
	apiv1 "app/proto/gen/api/v1"
	"app/store"
	"github.com/samber/lo"
)

var _ apiv1.RemindServiceServer = (*Remind)(nil)

type Remind struct {
	apiv1.UnimplementedRemindServiceServer
	store *store.Store
	robot *wechat.Robot
	conf  *config.Config
}

func NewRemind(s *store.Store, robot *wechat.Robot, conf *config.Config) *Remind {
	return &Remind{store: s, robot: robot, conf: conf}
}

func (r *Remind) Change(ctx context.Context, req *apiv1.RemindActionRequest) (*apiv1.TextResponse, error) {
	id, err := aesutil.AesDecode(r.conf.Common.TokenSecret, req.Token)
	if err != nil {
		return nil, fmt.Errorf("token错误:%w", err)
	}

	remind, err := r.store.GetRemind(ctx, lo.Must(strconv.Atoi(id)))
	if err != nil {
		return nil, fmt.Errorf("记录未找到:%w", err)
	}

	if err := r.store.UpdateRemindStatus(ctx, remind.Id, 1); err != nil {
		return nil, err
	}
	_ = r.robot.Message("已确认收到提醒")
	return &apiv1.TextResponse{Text: "已确认收到提醒"}, nil
}

func (r *Remind) Delay(ctx context.Context, req *apiv1.RemindActionRequest) (*apiv1.TextResponse, error) {
	id, err := aesutil.AesDecode(r.conf.Common.TokenSecret, req.Token)
	if err != nil {
		return nil, fmt.Errorf("token错误:%w", err)
	}

	remind, err := r.store.GetRemind(ctx, lo.Must(strconv.Atoi(id)))
	if err != nil {
		return nil, fmt.Errorf("记录未找到:%w", err)
	}

	if err := r.store.UpdateRemindNextTime(ctx, remind.Id, remind.NextTime); err != nil {
		return nil, err
	}
	_ = r.robot.Message("将在10分钟后再次提醒")
	return &apiv1.TextResponse{Text: "将在10分钟后再次提醒"}, nil
}
