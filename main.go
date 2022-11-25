package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"

	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// to run a function to use this mock runtime API:
// AWS_LAMBDA_FUNCTION_NAME=z_test AWS_LAMBDA_RUNTIME_API=localhost:8999/z_test ./build/z_test
func main() {
	// serviceAPIRouter := mux.NewRouter()

	// serviceAPIRouter.Handle("/{SERVICE_NAME}").Methods(http.MethodPost)

	// services := flag.String("services", "./services.json", "path to services descriptor JSON file")
	// flag.Parse()

	// servicesFile, err := os.Open(*services)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// servicesData, err := io.ReadAll(servicesFile)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// services := make([]Service)
	// json.Unmarshal(servicesData)

	runtimeAPIRouter := mux.NewRouter()

	runtimeAPIRouter.Handle("/{FUNCTION_NAME}/2018-06-01/runtime/invocation/next", NextHandler{}).Methods(http.MethodGet)
	runtimeAPIRouter.Handle("/{FUNCTION_NAME}/2018-06-01/runtime/invocation/{INVOCATION_ID}/response",
		ResponseHandler{}).Methods(http.MethodPost)
	runtimeAPIRouter.Handle("/{FUNCTION_NAME}/2018-06-01/runtime/invocation/{INVOCATION_ID}/error",
		ErrorHandler{}).Methods(http.MethodPost)
	runtimeAPIRouter.Handle("/{FUNCTION_NAME}/2018-06-01/runtime/init/error", InitErrorHandler{}).Methods(http.MethodPost)

	runtimeAPISrv := &http.Server{
		Handler:      runtimeAPIRouter,
		Addr:         "localhost:8999",
		WriteTimeout: 600 * time.Second,
		ReadTimeout:  600 * time.Second,
	}

	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGCHLD,
		syscall.SIGABRT,
	)
	go func() {
		log.Info("started basic server on localhost:8999")
		log.Fatal(runtimeAPISrv.ListenAndServe())
	}()
	s := <-sigC
	fmt.Println(s.String())
}

type NextHandler struct{}

func (h NextHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	log.Info("running lambda")

	writer.Header().Set("Lambda-Runtime-Aws-Request-Id", uuid.NewString())
	writer.Header().Set("Lambda-Runtime-Deadline-Ms", "5000")

	functionName := mux.Vars(request)["FUNCTION_NAME"]

	log.Info(functionName)

	// pop event from queue

	evt := events.APIGatewayProxyRequest{
		Headers:         map[string]string{"Authorization": "Bearer 12345"},
		Body:            "lalalalala",
		IsBase64Encoded: false,
		RequestContext:  events.APIGatewayProxyRequestContext{},
	}

	evtBody, err := json.Marshal(evt)
	if err != nil {
		log.Errorf("%s", err.Error())
		return
	}

	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write(evtBody)
	if err != nil {
		log.Errorf("%s", err.Error())
		return
	}
}

type ErrorHandler struct{}

func (h ErrorHandler) ServeHTTP(writer http.ResponseWriter, _ *http.Request) {
	log.Errorf("encountered error")
	errorBody, err := json.Marshal(map[string]interface{}{"StatusResponse": "202", "ErrorResponse": "hectic"})
	if err != nil {
		log.Errorf("%s", err.Error())
		return
	}
	http.Error(writer, string(errorBody), http.StatusAccepted)
}

type InitErrorHandler struct{}

func (h InitErrorHandler) ServeHTTP(writer http.ResponseWriter, _ *http.Request) {
	log.Errorf("encountered error")
	errorBody, err := json.Marshal(map[string]interface{}{"StatusResponse": "202", "ErrorResponse": "hectic"})
	if err != nil {
		log.Errorf("%s", err.Error())
		return
	}
	http.Error(writer, string(errorBody), http.StatusAccepted)
}

type ResponseHandler struct{}

func (h ResponseHandler) ServeHTTP(writer http.ResponseWriter, _ *http.Request) {
	log.Errorf("response")
	writer.WriteHeader(http.StatusAccepted)
}

