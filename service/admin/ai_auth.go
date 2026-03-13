package admin

import (
	"context"
	"time"

	"app/config"
	"app/pkg/errors"
	apiv1 "app/proto/gen/api/v1"
	"app/proto/gen/types"
	"app/store"
	"app/store/model"
)

var _ apiv1.AIAuthServiceServer = (*AIAuth)(nil)

type AIAuth struct {
	apiv1.UnimplementedAIAuthServiceServer
	store *store.Store
	conf  *config.Config
}

func NewAIAuth(s *store.Store, conf *config.Config) *AIAuth {
	return &AIAuth{store: s, conf: conf}
}

func (a *AIAuth) SetupProfile(ctx context.Context, req *apiv1.ProfileSetupRequest) (*types.IDResponse, error) {
	user := GetLoginUser(ctx)
	if user == nil {
		return nil, errors.ErrUnauthorized
	}

	threshold := req.VerificationThreshold
	if threshold <= 0 {
		threshold = 0.80
	}

	maxAttempts := req.MaxAttempts
	if maxAttempts <= 0 {
		maxAttempts = 3
	}

	profile := &model.AuthProfile{
		UserId:                user.Id,
		IdentityDescription:   req.IdentityDescription,
		VerificationThreshold: float64(threshold),
		MaxAttempts:           int(maxAttempts),
	}

	existing, err := a.store.GetAuthProfile(ctx, user.Id)
	if err == nil {
		profile.Id = existing.Id
		if err := a.store.UpdateAuthProfile(ctx, profile); err != nil {
			return nil, errors.InternalServer("UPDATE_ERROR", "更新失败")
		}
		return &types.IDResponse{Id: int32(existing.Id)}, nil
	}

	id, err := a.store.CreateAuthProfile(ctx, profile)
	if err != nil {
		return nil, errors.InternalServer("CREATE_ERROR", "创建失败")
	}

	return &types.IDResponse{Id: int32(id)}, nil
}

func (a *AIAuth) GetProfile(ctx context.Context, req *apiv1.GetProfileRequest) (*apiv1.ProfileResponse, error) {
	user := GetLoginUser(ctx)
	if user == nil {
		return nil, errors.ErrUnauthorized
	}

	profile, err := a.store.GetAuthProfile(ctx, user.Id)
	if err != nil {
		return nil, errors.BadRequest("PROFILE_NOT_FOUND", "未设置身份特征")
	}

	return &apiv1.ProfileResponse{
		Id:                    int32(profile.Id),
		IdentityDescription:   profile.IdentityDescription,
		VerificationThreshold: float64(profile.VerificationThreshold),
		MaxAttempts:           int32(profile.MaxAttempts),
		CreatedAt:             profile.CreatedAt.Format(time.DateTime),
		UpdatedAt:             profile.UpdatedAt.Format(time.DateTime),
	}, nil
}
