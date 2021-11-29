package tests

import (
	"SB/service/repository/db"
	"SB/service/repository/messenger"
	"SB/service/repository/persistence"
	"SB/service/repository/token"
	"SB/service/repository/training"
	"SB/service/service"
	"SB/service/service/tests/mocks"
	"b.yadro.com/storlib/logger"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type (
	ResponseTest struct{ *httptest.ResponseRecorder }
	JSON         map[string]interface{}
	JSONArray    []interface{}
	IJSON        interface {
		AsReader() io.Reader
	}
)

var (
	get  = handler(http.MethodGet)
	post = handler(http.MethodPost)
	put  = handler(http.MethodPut)
	del  = handler(http.MethodDelete)
)

//func configureEnvironment(t *testing.T) (*echo.Echo, func()) {
//	controller := gomock.NewController(t)
//
//	srv := NewServer("test", nil, &config)
//	srv.Start()
//
//	return srv.ServerApi(), func() {
//		controller.Finish()
//		rcTeardown()
//		fmt.Println("teardown")
//	}
//}

func startEchoServer(e *echo.Echo, port string) {
	go func(e *echo.Echo, port string) {
		err := e.Start(":" + port)
		if err != http.ErrServerClosed {
			logger.Fatal()
		}
	}(e, port)
}

func configureEnvironment(t *testing.T) (*echo.Echo, func()) {
	psqlDb, mock, err := sqlmock.New()
	assert.NoError(t, err)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: psqlDb,
	}), &gorm.Config{})

	assert.NoError(t, err)

	address := flag.String("address", "", "address to listen")
	port := flag.String("port", "3000", "port to listen")

	log.Info(fmt.Sprintf(
		"address: %s, port: %s", *address, *port,
	))

	persistent := persistence.NewPersistent(gormDB)
	usrMgr := db.NewDbManager(persistent)
	tknMgr := token.NewTokenManager(persistent)
	trainingMgr := training.NewTrainingManager(persistent)
	messenger := messenger.NewMessenger(persistent)

	server := service.NewServer(*address, *port, usrMgr, tknMgr, trainingMgr, messenger)
	log.Info("starting a server")

	mocks.ExpectCheckPassword(mock)
	mocks.ExpectMockGetUserAuthParams(mock)

	server.Start()

	return server.ServerApi(), func() {
		server.Stop()
		fmt.Println("teardown")
	}
}

func handler(method string) func(h http.Handler, path string, body ...IJSON) *ResponseTest {
	return func(h http.Handler, path string, body ...IJSON) *ResponseTest { return do(h, method, path, body...) }
}

func do(h http.Handler, method, path string, body ...IJSON) *ResponseTest {
	if body == nil {
		body = append(body, JSON{})
	}

	var (
		resp = ResponseTest{httptest.NewRecorder()}
		req  = httptest.NewRequest(method, path, body[0].AsReader())
	)

	h.ServeHTTP(resp, req)
	return &resp
}

func (j JSON) AsReader() io.Reader {
	if j == nil {
		return nil
	}

	var out bytes.Buffer
	if err := json.NewEncoder(&out).Encode(j); err != nil {
		panic(err.Error())
	}
	return &out
}
