package repository

import (
	"context"
	"time"

	"github.com/PickHD/LezPay/merchant/internal/config"
	"github.com/PickHD/LezPay/merchant/internal/model"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/sdk/trace"
)

type (
	// MerchantRepository is an interface that has all the function to be implemented inside merchant repository
	MerchantRepository interface {
		CreateMerchant(req *model.CreateMerchantRequest) (int64, bool, error)
		UpdateVerifiedMerchant(email string) (bool, error)
		GetMerchantDetailsByEmail(req *model.GetMerchantDetailsByEmailRequest) (*model.GetMerchantDetailsByEmailResponse, error)
		UpdateMerchantPasswordByEmail(req *model.UpdateMerchantPasswordByEmailRequest) (*model.UpdateMerchantPasswordByEmailResponse, error)
	}

	// MerchantRepositoryImpl is an app merchant struct that consists of all the dependencies needed for merchant repository
	MerchantRepositoryImpl struct {
		Context context.Context
		Config  *config.Configuration
		Logger  *logrus.Logger
		Tracer  *trace.TracerProvider
		DB      *pgxpool.Pool
		Redis   *redis.Client
	}
)

// NewMerchantRepository return new instances merchant repository
func NewMerchantRepository(ctx context.Context, config *config.Configuration, logger *logrus.Logger, tracer *trace.TracerProvider, db *pgxpool.Pool, rds *redis.Client) *MerchantRepositoryImpl {
	return &MerchantRepositoryImpl{
		Context: ctx,
		Config:  config,
		Logger:  logger,
		Tracer:  tracer,
		DB:      db,
		Redis:   rds,
	}
}

func (mr *MerchantRepositoryImpl) CreateMerchant(req *model.CreateMerchantRequest) (int64, bool, error) {
	tr := mr.Tracer.Tracer("Merchant-CreateMerchant Repository")
	_, span := tr.Start(mr.Context, "Start CreateMerchant")
	defer span.End()

	// begin tx
	tx, err := mr.DB.Begin(mr.Context)
	if err != nil {
		mr.Logger.Error("MerchantRepositoryImpl.CreateMerchant DB.Begin ERROR", err)

		return 0, false, err
	}

	var checkEmail string
	sqlCheck := `
		SELECT 
			email
		FROM
			merchant
		WHERE
			email = $1
	`
	row := tx.QueryRow(mr.Context, sqlCheck, req.Email)

	err = row.Scan(&checkEmail)
	if err != nil {
		// if email not found, create merchant
		if err.Error() == pgx.ErrNoRows.Error() {
			var id int64
			sql := `
					INSERT INTO 
						merchant (full_name,email,phone_number,password,is_verified) 
					VALUES 
						($1,$2,$3,$4,$5) RETURNING id
					`

			row := tx.QueryRow(mr.Context, sql, req.FullName, req.Email, req.PhoneNumber, req.Password, false)

			err = row.Scan(&id)
			if err != nil {
				// do rollback tx
				errRollback := tx.Rollback(mr.Context)
				if errRollback != nil {
					mr.Logger.Error("MerchantRepositoryImpl.CreateMerchant tx.Rollback ERROR ", errRollback)

					return 0, false, errRollback
				}

				mr.Logger.Error("MerchantRepositoryImpl.CreateMerchant row.Scan ERROR ", err)

				return 0, false, err
			}

			// do commit tx
			err = tx.Commit(mr.Context)
			if err != nil {
				mr.Logger.Error("MerchantRepositoryImpl.CreateMerchant tx.Commit ERROR", err)

				return 0, false, err
			}

			return id, true, nil
		}

		// do rollback tx
		errRollback := tx.Rollback(mr.Context)
		if errRollback != nil {
			mr.Logger.Error("MerchantRepositoryImpl.UpdateVerifiedMerchant tx.Rollback ERROR", errRollback)

			return 0, false, errRollback
		}

		mr.Logger.Error("MerchantRepositoryImpl.UpdateVerifiedMerchant row.Scan ERROR", err)

		return 0, false, err
	}

	// do rollback tx
	err = tx.Rollback(mr.Context)
	if err != nil {
		mr.Logger.Error("MerchantRepositoryImpl.UpdateVerifiedMerchant tx.Rollback ERROR", err)

		return 0, false, err
	}

	return 0, false, model.NewError(model.Validation, "Email already exists, please use another email instead")
}

func (mr *MerchantRepositoryImpl) UpdateVerifiedMerchant(email string) (bool, error) {
	tr := mr.Tracer.Tracer("Merchant-UpdateVerifiedMerchant Repository")
	_, span := tr.Start(mr.Context, "Start UpdateVerifiedMerchant")
	defer span.End()

	// begin tx
	tx, err := mr.DB.Begin(mr.Context)
	if err != nil {
		mr.Logger.Error("MerchantRepositoryImpl.UpdateVerifiedMerchant DB.Begin ERROR", err)

		return false, err
	}

	var checkEmail string
	sqlCheck := `
		SELECT 
			email
		FROM
			merchant
		WHERE
			email = $1
	`
	row := tx.QueryRow(mr.Context, sqlCheck, email)

	err = row.Scan(&checkEmail)
	if err != nil {
		// if data not found
		if err.Error() == pgx.ErrNoRows.Error() {
			// do rollback tx
			errRollback := tx.Rollback(mr.Context)
			if errRollback != nil {
				mr.Logger.Error("MerchantRepositoryImpl.UpdateVerifiedMerchant tx.Rollback ERROR", errRollback)

				return false, errRollback
			}

			return false, err
		}

		// do rollback tx
		errRollback := tx.Rollback(mr.Context)
		if errRollback != nil {
			mr.Logger.Error("MerchantRepositoryImpl.UpdateVerifiedMerchant tx.Rollback ERROR", errRollback)

			return false, errRollback
		}

		mr.Logger.Error("MerchantRepositoryImpl.UpdateVerifiedMerchant row.Scan ERROR", err)

		return false, err
	}

	sqlUpdate := `
		UPDATE 
			merchant
		SET
			is_verified = $1,
			updated_at = $2
		WHERE
			email = $3
	`

	_, err = tx.Exec(mr.Context, sqlUpdate, true, time.Now(), email)
	if err != nil {
		// do rollback tx
		errRollback := tx.Rollback(mr.Context)
		if errRollback != nil {
			mr.Logger.Error("MerchantRepositoryImpl.UpdateVerifiedMerchant tx.Rollback ERROR", errRollback)

			return false, errRollback
		}

		mr.Logger.Error("MerchantRepositoryImpl.UpdateVerifiedMerchant tx.Exec ERROR", err)

		return false, err
	}

	// do commit tx
	err = tx.Commit(mr.Context)
	if err != nil {
		mr.Logger.Error("MerchantRepositoryImpl.UpdateVerifiedMerchant tx.Commit ERROR", err)

		return false, err
	}

	return true, nil
}

func (mr *MerchantRepositoryImpl) GetMerchantDetailsByEmail(req *model.GetMerchantDetailsByEmailRequest) (*model.GetMerchantDetailsByEmailResponse, error) {
	tr := mr.Tracer.Tracer("Merchant-GetMerchantDetailsByEmail Repository")
	_, span := tr.Start(mr.Context, "Start GetMerchantDetailsByEmail")
	defer span.End()

	data := &model.GetMerchantDetailsByEmailResponse{}

	sql := `
		SELECT
			id,
			full_name,
			email,
			phone_number,
			password
		FROM 
			merchant
		WHERE
			email = $1
		AND
			is_verified = true
	`

	row := mr.DB.QueryRow(mr.Context, sql, req.Email)

	err := row.Scan(&data.ID, &data.FullName, &data.Email, &data.PhoneNumber, &data.Password)
	if err != nil {
		// if data not found
		if err.Error() == pgx.ErrNoRows.Error() {
			mr.Logger.Info("MerchantRepositoryImpl.GetMerchantDetailsByEmail email not found ", err)

			return nil, model.NewError(model.NotFound, "email not found")
		}

		mr.Logger.Error("MerchantRepositoryImpl.GetMerchantDetailsByEmail row.Scan ERROR ", err)

		return nil, err
	}

	return data, nil
}

func (mr *MerchantRepositoryImpl) UpdateMerchantPasswordByEmail(req *model.UpdateMerchantPasswordByEmailRequest) (*model.UpdateMerchantPasswordByEmailResponse, error) {
	tr := mr.Tracer.Tracer("Merchant-UpdateMerchantPasswordByEmail Repository")
	_, span := tr.Start(mr.Context, "Start UpdateMerchantPasswordByEmail")
	defer span.End()

	// begin tx
	tx, err := mr.DB.Begin(mr.Context)
	if err != nil {
		mr.Logger.Error("MerchantRepositoryImpl.UpdateMerchantPasswordByEmail DB.Begin ERROR", err)

		return nil, err
	}

	sqlUpdate := `
		UPDATE 
			merchant
		SET
			password = $1,
			updated_at = $2
		WHERE
			email = $3
	`

	_, err = tx.Exec(mr.Context, sqlUpdate, req.Password, time.Now(), req.Email)
	if err != nil {
		// do rollback tx
		errRollback := tx.Rollback(mr.Context)
		if errRollback != nil {
			mr.Logger.Error("MerchantRepositoryImpl.UpdateMerchantPasswordByEmail tx.Rollback ERROR", errRollback)

			return nil, errRollback
		}

		mr.Logger.Error("MerchantRepositoryImpl.UpdateMerchantPasswordByEmail tx.Exec ERROR", err)

		return nil, err
	}

	// do commit tx
	err = tx.Commit(mr.Context)
	if err != nil {
		mr.Logger.Error("MerchantRepositoryImpl.UpdateMerchantPasswordByEmail tx.Commit ERROR", err)

		return nil, err
	}

	return &model.UpdateMerchantPasswordByEmailResponse{
		Email: req.Email,
	}, nil
}
