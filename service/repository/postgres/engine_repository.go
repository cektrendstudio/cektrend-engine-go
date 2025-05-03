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

type engineRepo struct {
	db *sqlx.DB
}

func NewEngineRepository(db *sqlx.DB) service.EngineRepository {
	return &engineRepo{
		db: db,
	}
}

func (r *engineRepo) InsertPhisingWebReport(ctx context.Context, req models.CreatePhishingWebReportRequest) (id int64, errx serror.SError) {
	err := r.db.QueryRowxContext(ctx, queries.InsertPhisingWebReport,
		req.SiteURL,
		req.ImageURL,
		req.Version,
		time.Now(),
		time.Now(),
	).Scan(&id)
	if err != nil {
		errx = serror.NewFromError(err)
		errx.AddComments("[repository][InsertPhisingWebReport] while ExecContext queries.InsertPhisingWebReport")
		return
	}

	return
}

func (r *engineRepo) GetPhishingReportByURL(ctx context.Context, url string) (res models.GetPhisingWebReportResponse, errx serror.SError) {
	err := r.db.QueryRowxContext(ctx, queries.GetPhishingReportByURL, url).StructScan(&res)
	if err != nil && err != sql.ErrNoRows {
		errx = serror.NewFromError(err)
		errx.AddCommentf("[repository][GetPhishingReportByURL][UserID: %d] while QueryRowxContext", url)
		return
	}

	return
}

func (r *engineRepo) GetPhishingReportByURLs(ctx context.Context, urls []string, version int8) (res []models.GetPhisingWebReportResponse, errx serror.SError) {
	query, args, err := sqlx.In(queries.GetPhishingReportByURLs, version, urls)
	if err != nil {
		errx = serror.NewFromError(err)
		errx.AddComments("[repository][GetPhishingReportByURLs] while preparing IN query")
		return
	}

	query = r.db.Rebind(query)
	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		errx = serror.NewFromError(err)
		errx.AddComments("[repository][GetPhishingReportByURLs] while executing query")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var report models.GetPhisingWebReportResponse
		if err := rows.StructScan(&report); err != nil {
			errx = serror.NewFromError(err)
			errx.AddComments("[repository][GetPhishingReportByURLs] while scanning row")
			return
		}
		res = append(res, report)
	}

	return
}
