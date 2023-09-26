package repository

import (
	"context"
	"time"

	"github.com/PickHD/LezPay/customer/internal/v1/config"
	"github.com/PickHD/LezPay/customer/internal/v1/helper"
	"github.com/PickHD/LezPay/customer/internal/v1/model"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/protobuf/proto"

	walletpb "github.com/PickHD/LezPay/customer/pkg/proto/v1/wallet"
)

type (
	// CustomerRepository is an interface that has all the function to be implemented inside customer repository
	CustomerRepository interface {
		CreateCustomer(req *model.CreateCustomerRequest) (int64, bool, error)
		UpdateVerifiedCustomer(email string) (bool, error)
		GetCustomerIDByEmail(req *model.GetCustomerIDByEmailRequest) (int64, error)
		GetCustomerDetailsByEmail(req *model.GetCustomerDetailsByEmailRequest) (*model.GetCustomerDetailsByEmailResponse, error)
		UpdateCustomerPasswordByEmail(req *model.UpdateCustomerPasswordByEmailRequest) (*model.UpdateCustomerPasswordByEmailResponse, error)
		ProduceCustomerTopupTransaction(msg *walletpb.TopupTransactionMessage) error
	}

	// CustomerRepositoryImpl is an app customer struct that consists of all the dependencies needed for customer repository
	CustomerRepositoryImpl struct {
		Context       context.Context
		Config        *config.Configuration
		Logger        *logrus.Logger
		Tracer        *trace.TracerProvider
		DB            *pgxpool.Pool
		Redis         *redis.Client
		KafkaProducer *kafka.Writer
	}
)

// NewCustomerRepository return new instances customer repository
func NewCustomerRepository(ctx context.Context, config *config.Configuration, logger *logrus.Logger, tracer *trace.TracerProvider, db *pgxpool.Pool, rds *redis.Client, kafkaProducer *kafka.Writer) *CustomerRepositoryImpl {
	return &CustomerRepositoryImpl{
		Context:       ctx,
		Config:        config,
		Logger:        logger,
		Tracer:        tracer,
		DB:            db,
		Redis:         rds,
		KafkaProducer: kafkaProducer,
	}
}

func (cr *CustomerRepositoryImpl) CreateCustomer(req *model.CreateCustomerRequest) (int64, bool, error) {
	tr := cr.Tracer.Tracer("Customer-CreateCustomer Repository")
	_, span := tr.Start(cr.Context, "Start CreateCustomer")
	defer span.End()

	// begin tx
	tx, err := cr.DB.Begin(cr.Context)
	if err != nil {
		cr.Logger.Error("CustomerRepositoryImpl.CreateCustomer DB.Begin ERROR ", err)

		return 0, false, err
	}

	var checkEmail string
	sqlCheck := `
		SELECT 
			email
		FROM
			customer
		WHERE
			email = $1
	`
	row := tx.QueryRow(cr.Context, sqlCheck, req.Email)

	err = row.Scan(&checkEmail)
	if err != nil {
		// if email not found, create customer
		if err.Error() == pgx.ErrNoRows.Error() {
			sql := `
					INSERT INTO 
						customer (id,full_name,email,phone_number,password,pin,is_verified) 
					VALUES 
						($1,$2,$3,$4,$5,$6,$7)
					`

			id, err := helper.GenerateSnowflakeID()
			if err != nil {
				// do rollback tx
				errRollback := tx.Rollback(cr.Context)
				if errRollback != nil {
					cr.Logger.Error("CustomerRepositoryImpl.CreateCustomer tx.Rollback ERROR ", errRollback)

					return 0, false, errRollback
				}

				cr.Logger.Error("CustomerRepositoryImpl.CreateCustomer GenerateSnowflakeID ERROR ", err)

				return 0, false, err
			}

			_, err = tx.Exec(cr.Context, sql, id, req.FullName, req.Email, req.PhoneNumber, req.Password, req.Pin, false)
			if err != nil {
				// do rollback tx
				errRollback := tx.Rollback(cr.Context)
				if errRollback != nil {
					cr.Logger.Error("CustomerRepositoryImpl.CreateCustomer tx.Rollback ERROR ", errRollback)

					return 0, false, errRollback
				}

				cr.Logger.Error("CustomerRepositoryImpl.CreateCustomer tx.Exec ERROR ", err)

				return 0, false, err
			}

			// do commit tx
			err = tx.Commit(cr.Context)
			if err != nil {
				cr.Logger.Error("CustomerRepositoryImpl.CreateCustomer tx.Commit ERROR", err)

				return 0, false, err
			}

			return id, true, nil
		}

		// do rollback tx
		errRollback := tx.Rollback(cr.Context)
		if errRollback != nil {
			cr.Logger.Error("CustomerRepositoryImpl.CreateCustomer tx.Rollback ERROR", errRollback)

			return 0, false, errRollback
		}

		cr.Logger.Error("CustomerRepositoryImpl.CreateCustomer row.Scan ERROR", err)

		return 0, false, err
	}

	// do rollback tx
	err = tx.Rollback(cr.Context)
	if err != nil {
		cr.Logger.Error("CustomerRepositoryImpl.CreateCustomer tx.Rollback ERROR", err)

		return 0, false, err
	}

	return 0, false, model.NewError(model.Validation, "Email already exists, please use another email instead")
}

func (cr *CustomerRepositoryImpl) UpdateVerifiedCustomer(email string) (bool, error) {
	tr := cr.Tracer.Tracer("Customer-UpdateVerifiedCustomer Repository")
	_, span := tr.Start(cr.Context, "Start UpdateVerifiedCustomer")
	defer span.End()

	// begin tx
	tx, err := cr.DB.Begin(cr.Context)
	if err != nil {
		cr.Logger.Error("CustomerRepositoryImpl.UpdateVerifiedCustomer DB.Begin ERROR", err)

		return false, err
	}

	var checkEmail string
	sqlCheck := `
		SELECT 
			email
		FROM
			customer
		WHERE
			email = $1
	`
	row := tx.QueryRow(cr.Context, sqlCheck, email)

	err = row.Scan(&checkEmail)
	if err != nil {
		// if data not found
		if err.Error() == pgx.ErrNoRows.Error() {
			// do rollback tx
			errRollback := tx.Rollback(cr.Context)
			if errRollback != nil {
				cr.Logger.Error("CustomerRepositoryImpl.UpdateVerifiedCustomer tx.Rollback ERROR", errRollback)

				return false, errRollback
			}

			return false, err
		}

		// do rollback tx
		errRollback := tx.Rollback(cr.Context)
		if errRollback != nil {
			cr.Logger.Error("CustomerRepositoryImpl.UpdateVerifiedCustomer tx.Rollback ERROR", errRollback)

			return false, errRollback
		}

		cr.Logger.Error("CustomerRepositoryImpl.UpdateVerifiedCustomer row.Scan ERROR", err)

		return false, err
	}

	sqlUpdate := `
		UPDATE 
			customer
		SET
			is_verified = $1,
			updated_at = $2
		WHERE
			email = $3
	`

	_, err = tx.Exec(cr.Context, sqlUpdate, true, time.Now(), email)
	if err != nil {
		// do rollback tx
		errRollback := tx.Rollback(cr.Context)
		if errRollback != nil {
			cr.Logger.Error("CustomerRepositoryImpl.UpdateVerifiedCustomer tx.Rollback ERROR", errRollback)

			return false, errRollback
		}

		cr.Logger.Error("CustomerRepositoryImpl.UpdateVerifiedCustomer tx.Exec ERROR", err)

		return false, err
	}

	// do commit tx
	err = tx.Commit(cr.Context)
	if err != nil {
		cr.Logger.Error("CustomerRepositoryImpl.UpdateVerifiedCustomer tx.Commit ERROR", err)

		return false, err
	}

	return true, nil
}

func (cr *CustomerRepositoryImpl) GetCustomerIDByEmail(req *model.GetCustomerIDByEmailRequest) (int64, error) {
	tr := cr.Tracer.Tracer("Customer-GetCustomerIDByEmail Repository")
	_, span := tr.Start(cr.Context, "Start GetCustomerIDByEmail")
	defer span.End()

	var getCustomerID int64
	sqlCheck := `
		SELECT 
			id
		FROM
			customer
		WHERE
			email = $1
	`
	row := cr.DB.QueryRow(cr.Context, sqlCheck, req.Email)

	err := row.Scan(&getCustomerID)
	if err != nil {
		cr.Logger.Error("CustomerRepositoryImpl.GetCustomerIDByEmail row.Scan ERROR", err)

		return 0, err
	}

	return getCustomerID, nil
}

func (cr *CustomerRepositoryImpl) GetCustomerDetailsByEmail(req *model.GetCustomerDetailsByEmailRequest) (*model.GetCustomerDetailsByEmailResponse, error) {
	tr := cr.Tracer.Tracer("Customer-GetCustomerDetailsByEmail Repository")
	_, span := tr.Start(cr.Context, "Start GetCustomerDetailsByEmail")
	defer span.End()

	data := &model.GetCustomerDetailsByEmailResponse{}

	sql := `
		SELECT
			id,
			full_name,
			email,
			phone_number,
			password,
			pin
		FROM 
			customer
		WHERE
			email = $1
		AND
			is_verified = true
	`

	row := cr.DB.QueryRow(cr.Context, sql, req.Email)

	err := row.Scan(&data.ID, &data.FullName, &data.Email, &data.PhoneNumber, &data.Password, &data.Pin)
	if err != nil {
		// if data not found
		if err.Error() == pgx.ErrNoRows.Error() {
			cr.Logger.Info("CustomerRepositoryImpl.GetCustomerDetailsByEmail email not found ", err)

			return nil, model.NewError(model.NotFound, "email not found")
		}

		cr.Logger.Error("CustomerRepositoryImpl.GetCustomerDetailsByEmail row.Scan ERROR ", err)

		return nil, err
	}

	return data, nil
}

func (cr *CustomerRepositoryImpl) UpdateCustomerPasswordByEmail(req *model.UpdateCustomerPasswordByEmailRequest) (*model.UpdateCustomerPasswordByEmailResponse, error) {
	tr := cr.Tracer.Tracer("Customer-UpdateCustomerPasswordByEmail Repository")
	_, span := tr.Start(cr.Context, "Start UpdateCustomerPasswordByEmail")
	defer span.End()

	// begin tx
	tx, err := cr.DB.Begin(cr.Context)
	if err != nil {
		cr.Logger.Error("CustomerRepositoryImpl.UpdateCustomerPasswordByEmail DB.Begin ERROR", err)

		return nil, err
	}

	sqlUpdate := `
		UPDATE 
			customer
		SET
			password = $1,
			updated_at = $2
		WHERE
			email = $3
	`

	_, err = tx.Exec(cr.Context, sqlUpdate, req.Password, time.Now(), req.Email)
	if err != nil {
		// do rollback tx
		errRollback := tx.Rollback(cr.Context)
		if errRollback != nil {
			cr.Logger.Error("CustomerRepositoryImpl.UpdateCustomerPasswordByEmail tx.Rollback ERROR", errRollback)

			return nil, errRollback
		}

		cr.Logger.Error("CustomerRepositoryImpl.UpdateCustomerPasswordByEmail tx.Exec ERROR", err)

		return nil, err
	}

	// do commit tx
	err = tx.Commit(cr.Context)
	if err != nil {
		cr.Logger.Error("CustomerRepositoryImpl.UpdateCustomerPasswordByEmail tx.Commit ERROR", err)

		return nil, err
	}

	return &model.UpdateCustomerPasswordByEmailResponse{
		Email: req.Email,
	}, nil
}

func (cr *CustomerRepositoryImpl) ProduceCustomerTopupTransaction(msg *walletpb.TopupTransactionMessage) error {
	tr := cr.Tracer.Tracer("Customer-ProduceCustomerTopupTransaction Repository")
	_, span := tr.Start(cr.Context, "Start ProduceCustomerTopupTransaction")
	defer span.End()

	b, err := proto.Marshal(msg)
	if err != nil {
		cr.Logger.Error("CustomerRepositoryImpl.ProduceCustomerTopupTransaction Marshal proto TopupTransactionMessage ERROR, ", err)
		return err
	}

	key, err := helper.GenerateByteSnowflakeID()
	if err != nil {
		cr.Logger.Error("CustomerRepositoryImpl.ProduceCustomerTopupTransaction GenerateByteSnowflakeID ERROR, ", err)
		return err
	}

	err = cr.KafkaProducer.WriteMessages(cr.Context,
		kafka.Message{
			Key:   key,
			Value: b,
			Topic: cr.Config.Kafka.TopicTopupTransaction,
		},
	)
	if err != nil {
		cr.Logger.Error("CustomerRepositoryImpl.ProduceCustomerTopupTransaction KafkaProducer.WriteMessages ERROR, ", err)
		return err
	}

	defer cr.KafkaProducer.Close()

	return nil
}
