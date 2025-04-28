package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/cektrendstudio/cektrend-engine-go/models"
	"github.com/cektrendstudio/cektrend-engine-go/pkg/serror"
	"github.com/cektrendstudio/cektrend-engine-go/service"
	"github.com/cektrendstudio/cektrend-engine-go/service/repository/queries"

	"github.com/jmoiron/sqlx"
)

type authRepo struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) service.AuthRepository {
	return &authRepo{
		db: db,
	}
}

func (r *authRepo) CreateToken(ctx context.Context, req models.CreateTokenRequest) (tokenID int64, errx serror.SError) {
	err := r.db.QueryRowxContext(ctx, queries.CreateToken,
		req.UserID,
		req.AccessToken,
		req.RefreshToken,
		req.AccessTokenExpiresAt,
		req.RefreshTokenExpiresAt,
		time.Now(),
	).Scan(&tokenID)
	if err != nil {
		errx = serror.NewFromError(err)
		errx.AddComments("[repository][CreateToken] while ExecContext queries.CreateToken")
		return
	}

	return
}

func (r *authRepo) RevokedInactiveToken(ctx context.Context, req models.RevokedInactiveTokenRequest) (errx serror.SError) {
	_, err := r.db.ExecContext(ctx, queries.RevokedInactiveToken,
		req.UserID,
		req.TokenID,
		time.Now(),
	)
	if err != nil {
		errx = serror.NewFromError(err)
		errx.AddComments("[repository][RevokedInactiveToken] while ExecContext queries.RevokedInactiveToken")
		return
	}

	return
}

func (r *authRepo) ValidateAccessToken(ctx context.Context, req models.ValidateAccessTokenRequest) (res models.ValidateAccessTokenResponse, errx serror.SError) {
	err := r.db.QueryRowxContext(ctx, queries.ValidateAccessToken,
		req.AccessToken,
		req.UserID,
	).StructScan(&res)
	if err != nil && err != sql.ErrNoRows {
		errx = serror.NewFromError(err)
		errx.AddCommentf("[repository][ValidateAccessToken][UserID: %d] while QueryRowxContext", req.UserID)
		return
	}

	return
}

func (r *authRepo) ValidateRefreshToken(ctx context.Context, req models.ValidateRefreshTokenRequest) (res models.ValidateRefreshTokenResponse, errx serror.SError) {
	err := r.db.QueryRowxContext(ctx, queries.ValidateRefreshToken,
		req.RefreshToken,
		req.UserID,
	).StructScan(&res)
	if err != nil && err != sql.ErrNoRows {
		errx = serror.NewFromError(err)
		errx.AddCommentf("[repository][ValidateRefreshToken][UserID: %d] while QueryRowxContext", req.UserID)
		return
	}

	return
}

func (r *authRepo) UpdateToken(ctx context.Context, req models.UpdateTokenRequest) (errx serror.SError) {
	_, err := r.db.ExecContext(ctx, queries.UpdateTokenByTokenID,
		req.TokenID,
		req.AccessToken,
		req.AccessTokenExpiresAt,
		req.RefreshToken,
		req.RefreshTokenExpiresAt,
		time.Now(),
	)
	if err != nil {
		errx = serror.NewFromError(err)
		errx.AddComments("[repository][UpdateAccessToken] while ExecContext queries.UpdateTokenByTokenID")
		return
	}

	return
}
