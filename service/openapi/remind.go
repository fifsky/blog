package openapi

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"app/config"
	"app/pkg/aesutil"
	apiv1 "app/proto/gen/api/v1"
	remindpkg "app/service/remind"
	"app/store"

	"github.com/samber/lo"
)

var _ apiv1.RemindServiceHTTPServer = (*Remind)(nil)

type Remind struct {
	store *store.Store
	conf  *config.Config
}

func NewRemind(s *store.Store, conf *config.Config) *Remind {
	return &Remind{store: s, conf: conf}
}

func (r *Remind) Change(ctx context.Context, req *apiv1.RemindActionRequest) (*apiv1.TextResponse, error) {
	id, err := aesutil.AesDecode(r.conf.Common.TokenSecret, req.GetToken())
	if err != nil {
		return nil, fmt.Errorf("token错误:%w", err)
	}

	remind, err := r.store.GetRemind(ctx, lo.Must(strconv.Atoi(id)))
	if err != nil {
		return nil, fmt.Errorf("记录未找到:%w", err)
	}

	nextTime := remindpkg.NextTimeFromRule(time.Now(), remind)

	// 判断是否是固定时间任务
	isFixedDate := len(remind.Cron) >= 10 && remind.Cron[4] == '-'

	if isFixedDate {
		// 固定时间任务确认收到后，置为状态 3 (已完成)
		if err := r.store.UpdateRemindStatus(ctx, remind.Id, 3); err != nil {
			return nil, err
		}
	} else {
		// 周期性任务，恢复状态 1，并更新下一次时间
		if err := r.store.UpdateRemindStatus(ctx, remind.Id, 1); err != nil {
			return nil, err
		}
		if err := r.store.UpdateRemindNextTime(ctx, remind.Id, nextTime); err != nil {
			return nil, err
		}
	}

	return apiv1.TextResponse_builder{Text: "已确认收到提醒"}.Build(), nil
}

func (r *Remind) Delay(ctx context.Context, req *apiv1.RemindActionRequest) (*apiv1.TextResponse, error) {
	id, err := aesutil.AesDecode(r.conf.Common.TokenSecret, req.GetToken())
	if err != nil {
		return nil, fmt.Errorf("token错误:%w", err)
	}

	remind, err := r.store.GetRemind(ctx, lo.Must(strconv.Atoi(id)))
	if err != nil {
		return nil, fmt.Errorf("记录未找到:%w", err)
	}

	nextTime := time.Now().Add(10 * time.Minute)
	if err := r.store.UpdateRemindNextTime(ctx, remind.Id, nextTime); err != nil {
		return nil, err
	}
	return apiv1.TextResponse_builder{Text: "将在10分钟后再次提醒"}.Build(), nil
}
