package server

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"remote-server-control/internal/handlers"
	"strings"
	"testing"
)

type ErrorResponseMessage struct {
	Err string `json:"message"`
}

type ResponseMessage struct {
	Msg CommandOutPut `json:"message"`
}

type CommandOutPut struct {
	Stder  string `json:"stderr"`
	Stdout string `json:"stdout"`
}

func TestErrorCaseExecuteRemoteCommand(t *testing.T) {
	TestCasesWithErrors := []struct {
		body string
	}{
		{
			//interrupted: execution time limit exceeded
			body: `{"cmd": "ping ya.ru", "os": "linux","stdin": ""}`,
		},
		{
			//not vallid wrong os type.
			body: `{"cmd": "echo test", "os": "windows","stdin": ""}`,
		},
		{
			// not vallid failed execute command
			body: `{"cmd": "t", "os": "linux","stdin": ""}`,
		},
		{
			// not vallid empty command  name
			body: `{"cmd": "", "os": "linux","stdin": ""}`,
		},
		{
			// wrong type for field
			body: `{"cmd": 23, "os": "linux","stdin": ""}`,
		},
		{
			// not vallid wrong field
			body: `{"name": "", "ostype": "linux","pipe": ""}`,
		},
	}

	expectedErrorResult := []struct {
		statusCode int
		err        error
	}{
		{
			statusCode: http.StatusBadRequest,
			err:        errors.New("interrupted: execution time limit exceeded"),
		},
		{
			statusCode: http.StatusBadRequest,
			err:        errors.New("wrong os type. This machine works on Os: linux"),
		},
		{
			statusCode: http.StatusBadRequest,
			err:        errors.New("failed execute command: exec: \"t\": executable file not found in $PATH"),
		},
		{
			statusCode: http.StatusBadRequest,
			err:        errors.New("empty command name"),
		},
		{
			statusCode: http.StatusBadRequest,
			err:        errors.New("wrong type provided for field cmd, expected string"),
		},
		{
			statusCode: http.StatusBadRequest,
			err:        errors.New("json: unknown field \"name\""),
		},
	}

	TestCases := []struct {
		body string
	}{
		{
			body: `{"cmd": "tr a-z A-Z", "os": "linux","stdin": "test"}`,
		},
		{
			body: `{"cmd": "echo test", "os": "linux","stdin": ""}`,
		},
		{
			body: `{"cmd": "ping -c 3 ya.ru", "os": "linux","stdin": ""}`,
		},
	}

	expectedResult := []struct {
		statusCode int
		out        CommandOutPut
	}{
		{
			statusCode: http.StatusOK,
			out: CommandOutPut{
				Stder:  "",
				Stdout: "TEST",
			},
		},
		{
			statusCode: http.StatusOK,
			out: CommandOutPut{
				Stder:  "",
				Stdout: "test\n",
			},
		},
		{
			statusCode: http.StatusOK,
		},
	}

	mux := http.NewServeMux()

	mux.HandleFunc(_apiPrifix+"/remote-execution", handlers.ExecuteRemoteCommand)

	for i, tc := range TestCasesWithErrors {
		reader := strings.NewReader(tc.body)

		r, _ := http.NewRequest(http.MethodPost, "/api/v1/remote-execution", reader)

		w := httptest.NewRecorder()

		mux.ServeHTTP(w, r)

		resp := w.Result()

		if resp.StatusCode != expectedErrorResult[i].statusCode {
			t.Errorf("TestCase: %v, expected: %v, received: %v \n",
				i,
				expectedErrorResult[i].statusCode,
				resp.StatusCode,
			)
		}

		var e ErrorResponseMessage

		if err := json.Unmarshal(w.Body.Bytes(), &e); err != nil {
			log.Fatalln(err)
		}

		if e.Err != expectedErrorResult[i].err.Error() {
			t.Errorf("TestCase: %v, expected: %v, received: %v \n",
				i,
				e.Err,
				expectedErrorResult[i].err.Error(),
			)
		}
	}

	for i, tc := range TestCases {
		reader := strings.NewReader(tc.body)

		r, _ := http.NewRequest(http.MethodPost, "/api/v1/remote-execution", reader)

		w := httptest.NewRecorder()

		mux.ServeHTTP(w, r)

		resp := w.Result()

		if resp.StatusCode != expectedResult[i].statusCode {
			t.Errorf("TestCase: %v, expected: %v, received: %v \n",
				i,
				expectedResult[i].statusCode,
				resp.StatusCode,
			)
		}

		if i == 2 { //request with ping
			return
		}

		var res ResponseMessage

		if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
			log.Fatalln(err)
		}

		if res.Msg.Stdout != expectedResult[i].out.Stdout {
			t.Errorf("TestCase: %v, expected: %v, received: %v \n",
				i,
				expectedResult[i].out.Stdout,
				res.Msg.Stdout,
			)
		}
	}
}
