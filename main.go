package main

import (
	"encoding/json"
	_ "encoding/base64"
	"io/ioutil"
	"net/http"
	
	"os"
	"strconv"

	"github.com/gorilla/mux"

	conf "test/config"
	dt "test/datastruct"
	lib "test/lib"
	log "test/logging"
	services "test/services"
	"test/transport"

	httptransport "github.com/go-kit/kit/transport/http"

)

func initHandlers(dbConn lib.DbConnection) {
	
	var userServices services.UserReqServices
	userServices = services.UserReqServices{}

	NewUserHandler := httptransport.NewServer(
		transport.CreateUserEndpoint(userServices, dbConn),
		transport.CreateUserDecodeRequest,
		transport.CreateUserEncodeResponse,
	)

	FetchHandler := httptransport.NewServer(
		transport.CreateUserEndpoint(userServices, dbConn),
		transport.CreateUserDecodeRequest,
		transport.CreateUserEncodeResponse,
	)

	http.Handle("/api/insert/", NewUserHandler)
	
	/*
		Router with dynamic param
		by : Raditya Pratama
		5 Oktober 2018
	*/
	rtr := mux.NewRouter()
	rtr.HandleFunc("/api/DownloadFile/{type:[a-zA-z0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		param := mux.Vars(r)
		typeFile := param["type"]
		resCode, resDesc := services.GetFileFTP(typeFile, dbConn)

		result := map[string]string{
			"respCode": strconv.Itoa(resCode),
			"respDesc": resDesc,
		}
		services.SendResponses(w, result)
	})
	
	http.Handle("/", rtr)
}

func decodeBody(r *http.Request, request interface{}) {
	
	var body []byte

		//decode request body
	body, err := ioutil.ReadAll(r.Body)
	// log.Logf("Get Query Request Body : %#v %s", body, string(body[:]))
	if err != nil {
		log.Errorf("Error when get data %s", err.Error())
	}

	if err = json.Unmarshal(body, &request); err != nil {
		log.Errorf("Error when Unmarshall data %s", err.Error())
	}
	// log.Logf("After Unmarshall : %#v %s", request, string(body[:]))
}

func main() {
	lib.LoadConfiguration()

	// initiate Service Database connection
	dbConn := lib.InitDb()

	// Register and Initiate Listener
	initHandlers(dbConn)
	
	var err error

	err = http.ListenAndServe(conf.Param.ListenPort, nil)

	if err != nil {
		log.Errorf("Unable to start the server %v", err)
		os.Exit(1)
	}
}
