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

type userRepo struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) service.UserRepository {
	return &userRepo{
		db: db,
	}
}

func (r *userRepo) Register(ctx context.Context, req models.RegisterUserRequest) (userId int64, errx serror.SError) {
	err := r.db.QueryRowxContext(ctx, queries.RegisterUser,
		req.Name,
		req.Email,
		req.Password,
		time.Now(),
		time.Now(),
	).Scan(&userId)
	if err != nil {
		errx = serror.NewFromError(err)
		errx.AddComments("[repository][Register] while ExecContext queries.RegisterUser")
		return
	}

	return
}

func (r *userRepo) GetUserByEmail(ctx context.Context, email string) (res models.User, errx serror.SError) {
	err := r.db.QueryRowxContext(ctx, queries.GetUserByEmail, email).StructScan(&res)
	if err != nil && err != sql.ErrNoRows {
		errx = serror.NewFromError(err)
		errx.AddCommentf("[repository][GetUserByEmail][Email: %d] while QueryRowxContext", email)
		return
	}

	return
}

func (r *userRepo) GetUserByID(ctx context.Context, userID int64) (res models.User, errx serror.SError) {
	err := r.db.QueryRowxContext(ctx, queries.GetUserByID, userID).StructScan(&res)
	if err != nil && err != sql.ErrNoRows {
		errx = serror.NewFromError(err)
		errx.AddCommentf("[repository][GetUserByEmail][UserID: %d] while QueryRowxContext", userID)
		return
	}

	return
}
