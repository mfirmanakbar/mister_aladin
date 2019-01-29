package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	dt "test/datastruct"
	ex "test/error"
	lib "test/lib"
	logger "test/logging"
	"test/services"

	"github.com/go-kit/kit/endpoint"
)

// GetDoDecodeRequest : request param for queue list using JSON format place in body
func CreateUserDecodeRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var request dt.UserRequest

	var body []byte

	//decode request body
	body, err := ioutil.ReadAll(r.Body)
	logger.Logf("GetDoDecodeRequest : %s", string(body[:]))
	if err != nil {
		return ex.Errorc(dt.ErrInvalidFormat).Rem("Unable to read request body"), nil
	}

	if err = json.Unmarshal(body, &request); err != nil {
		return ex.Error(err, dt.ErrInvalidFormat).Rem("Failed decoding json message"), nil
	}

	return request, nil
}

// GetDoEncodeResponse : response using JSON format
func CreateUserEncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	var body []byte
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	body, err := json.Marshal(&response)
	logger.Logf("GetDoEncodeResponse : %s", string(body[:]))

	if err != nil {
		return err
	}

	//w.Header().Set("X-Checksum", cm.Cksum(body))

	var e = response.(dt.UserResponse).ResponseCode

	if e <= 500 {
		w.WriteHeader(http.StatusOK)
	} else if e <= 900 {
		w.WriteHeader(http.StatusBadRequest)
	} else if e <= 998 {
		w.WriteHeader(http.StatusUnauthorized)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}

	_, err = w.Write(body)

	return err
}

// GetDoFileEndpoint call Queue List
func CreateUserEndpoint(svc services.DoReqServices, dbConn lib.DbConnection) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(dt.DoRequest)
		if ok {

			return svc.GetDoData(ctx, req, dbConn), nil

		}
		logger.Error("Unhandled error occured: request is in unknown format")
		fmt.Printf("%#v\n", req)
		fmt.Printf("%#v\n", ok)
		var resp = dt.DoResponse{}
		resp.ResponseCode = dt.ErrOthers
		return resp, nil
	}
}
