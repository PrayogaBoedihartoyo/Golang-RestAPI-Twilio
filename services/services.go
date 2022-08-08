package services

import (
	"context"
	"database/sql"
	"github.com/go-playground/validator/v10"
	"log"
	"main/controllers"
	"main/entity"
	"main/helper"
	"main/repository"
	"time"
)

type OTPservices interface {
	Create(ctx context.Context, request entity.Request) (entity.Response, error)
	Verification(ctx context.Context, request entity.Verification) (entity.Response, error)
}

type OTPserviceImplementation struct {
	OTPrepository repository.OTPrepository
	db            *sql.DB
	validate      *validator.Validate
}

func NewOTPserviceImplementation(OTPrepository repository.OTPrepository, db *sql.DB, validate *validator.Validate) *OTPserviceImplementation {
	return &OTPserviceImplementation{OTPrepository: OTPrepository, db: db, validate: validate}
}

func (service *OTPserviceImplementation) Create(ctx context.Context, request entity.Request) (entity.Response, error) {
	err := service.validate.Struct(request)
	helper.HandlePanic(err)

	requestClient := entity.Request{
		Phone: request.Phone,
	}

	requestFinal := helper.RequestToUser(requestClient)
	requestFinal, err = service.OTPrepository.Create(ctx, service.db, requestFinal)
	helper.HandlePanic(err)

	err = controllers.SendOTP(requestClient.Phone)
	helper.HandlePanic(err)

	return helper.UserToResponse(requestFinal), nil
}

func (service *OTPserviceImplementation) Verification(ctx context.Context, request entity.Verification) (entity.Response, error) {
	err := service.validate.Struct(request)
	if err != nil {
		log.Println("SERVICE", err)
	}

	requestClient := entity.Verification{
		Id:         request.Id,
		Code:       request.Code,
		Phone:      request.Phone,
		Receiver:   request.Receiver,
		Payload:    request.Payload,
		VerifiedAt: time.Now(),
		ExpiredAt:  time.Now(),
	}
	_, err = service.OTPrepository.Verification(ctx, service.db, requestClient)
	if err != nil {
		log.Println("SERVICE 1", err)
	}

	controllers.CheckOTP(requestClient)

	return helper.RequestVerificationToResponse(requestClient), nil
}
