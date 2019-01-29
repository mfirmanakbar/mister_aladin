package services

import (
	"context"
	"fmt"
	"reflect"

	dt "test/datastruct"
	lib "test/lib"

	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

// DoReqServices provides operations for endpoint
type FetchReqServices interface {
	CreateNew(context.Context, dt.UserRequest, lib.DbConnection) dt.UserResponse
}

// GetEmployeeNIKService is a concrete implementation of QueueServices
type FetchReqServices struct{}

// GetEmployeeNIK service for counting queue
func (UserReqServices) CreateNew(ctx context.Context, req dt.UserRequest, dbConn lib.DbConnection) (resp dt.UserResponse) {
	err := validation.Errors{
		"name": validation.Validate(req.Name, validation.Required, validation.Length(1, 50), is.Letter),
		"email": validation.Validate(req.Email, validation.Required, is.Email),
		"phone": validation.Validate(req.Phone, validation.Required, is.Digit),
	}.Filter()
	resp.ResponseCode = "-1"
	resp.ResponseDesc = err
	if err == nil || err == "" {

		resp.ResponseCode = "1"
		resp.ResponseDesc = "Success"
	}
	
	return
}

func (UserReqServices) GoFetch(ctx context.Context, req dt.UserRequest, dbConn lib.DbConnection) (resp dt.UserResponse) {
	
}