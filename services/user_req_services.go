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
type UserReqServices interface {
	CreateNew(context.Context, dt.UserRequest, lib.DbConnection) dt.UserResponse
}

// GetEmployeeNIKService is a concrete implementation of QueueServices
type UserReqServices struct{}

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

func validateData(req dt.UserRequest) bool, string {
	err := validation.Errors{
		"name": validation.Validate(req.Name, validation.Required, validation.Length(1, 50), is.Letter),
		"email": validation.Validate(req.Email, validation.Required, is.Email),
		"phone": validation.Validate(req.Phone, validation.Required, is.Digit),
	}.Filter()

	ResponseCode := false
	ResponseDesc := err
	if err == nil || err == "" {
		
		ResponseCode = true
		ResponseDesc = "Success"
	}

	return ResponseCode, ResponseDesc
}

func structToArr(structName interface{}) map[string]string {
	// refl := reflect.TypeOf(dt.SodogiColumns).Name()
	var reflectValue = reflect.ValueOf(structName)

	if reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}

	var reflectType = reflectValue.Type()
	var ArrStruct = make(map[string]string)
	for i := 0; i < reflectValue.NumField(); i++ {
		ArrStruct[reflectType.Field(i).Name] = reflectValue.Field(i).Interface().(string)
	}
	return ArrStruct
}

func (UserReqServices) GetData(ctx context.Context, req dt.UserRequest, dbConn lib.DbConnection) (resp dt.UserResponse) {
	code, desc := validateData(req)

	if code {
		willBeInsert := structToArr(req)
		dbConn.InsertData("users", willBeInsert)
	}
	return
}