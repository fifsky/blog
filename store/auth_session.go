package store

import (
	"context"
	"time"

	"app/store/model"
)

func (s *Store) CreateAuthSession(ctx context.Context, session *model.AuthSession) (int64, error) {
	query := `INSERT INTO auth_sessions (session_id, user_id, attempt_count, verified_score, status, expires_at) 
		VALUES (?, ?, ?, ?, ?, ?)`
	res, err := s.db.ExecContext(ctx, query, session.SessionId, session.UserId, session.AttemptCount,
		session.VerifiedScore, session.Status, session.ExpiresAt)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *Store) GetAuthSession(ctx context.Context, sessionID string) (*model.AuthSession, error) {
	query := `SELECT id, session_id, user_id, attempt_count, verified_score, status, expires_at, created_at 
		FROM auth_sessions WHERE session_id = ? AND status = 'active' AND expires_at > ?`
	row := s.db.QueryRowContext(ctx, query, sessionID, time.Now())

	var session model.AuthSession
	err := row.Scan(&session.Id, &session.SessionId, &session.UserId, &session.AttemptCount,
		&session.VerifiedScore, &session.Status, &session.ExpiresAt, &session.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (s *Store) UpdateAuthSession(ctx context.Context, session *model.AuthSession) error {
	query := `UPDATE auth_sessions SET attempt_count = ?, verified_score = ?, status = ? WHERE session_id = ?`
	_, err := s.db.ExecContext(ctx, query, session.AttemptCount, session.VerifiedScore, session.Status, session.SessionId)
	return err
}
