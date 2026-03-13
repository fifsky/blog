package store

import (
	"context"

	"app/store/model"
)

func (s *Store) GetAuthProfile(ctx context.Context, userID int) (*model.AuthProfile, error) {
	query := `SELECT id, user_id, identity_description, verification_threshold, max_attempts, created_at, updated_at 
		FROM user_auth_profiles WHERE user_id = ?`
	row := s.db.QueryRowContext(ctx, query, userID)

	var p model.AuthProfile
	err := row.Scan(&p.Id, &p.UserId, &p.IdentityDescription, &p.VerificationThreshold,
		&p.MaxAttempts, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *Store) CreateAuthProfile(ctx context.Context, p *model.AuthProfile) (int64, error) {
	query := `INSERT INTO user_auth_profiles (user_id, identity_description, verification_threshold, max_attempts) 
		VALUES (?, ?, ?, ?)`
	res, err := s.db.ExecContext(ctx, query, p.UserId, p.IdentityDescription, p.VerificationThreshold, p.MaxAttempts)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *Store) UpdateAuthProfile(ctx context.Context, p *model.AuthProfile) error {
	query := `UPDATE user_auth_profiles SET identity_description = ?, verification_threshold = ?, max_attempts = ? 
		WHERE user_id = ?`
	_, err := s.db.ExecContext(ctx, query, p.IdentityDescription, p.VerificationThreshold, p.MaxAttempts, p.UserId)
	return err
}
