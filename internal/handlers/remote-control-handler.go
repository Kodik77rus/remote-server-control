package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	runner "remote-server-control/internal/command-runner"
	"runtime"
	"strings"
	// "log"
)

const _os = runtime.GOOS

var errEmptyBody = errors.New("empty body")

type errExecCommandValidator struct {
	osType string
}

type unmarshalTypeError struct {
	msg          string
	unmarshalErr *json.UnmarshalTypeError
}

type linuxCommnd struct {
	name  string
	flags []string
}

type requestBody struct {
	Cmd   string `json:"cmd"`
	Os    string `json:"os"`
	Stdin string `json:"stdin"`
}

//excute os command
func ExecuteRemoteCommand(w http.ResponseWriter, r *http.Request) {
	// log.Printf("%+v", r) // logger create middleware

	if r.Method != http.MethodPost {
		errorResponse(w, http.ErrBodyNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	var parsedbody requestBody
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&parsedbody); err != nil {
		if errors.As(err, &unmarshalErr) {
			errorResponse(
				w,
				unmarshalTypeError{
					msg:          "wrong type provided for field",
					unmarshalErr: unmarshalErr,
				},
				http.StatusBadRequest,
			)
			return
		}

		errorResponse(w, err, http.StatusBadRequest)
		return
	}

	if err := bodyValidator(&parsedbody); err != nil {
		errorResponse(w, err, http.StatusBadRequest)
		return
	}

	command := commandParser(&parsedbody)
	runner := runner.New(command, _os)

	outPut := runner.Run()

	sendRespone(w, outPut, http.StatusOK)
}

func (e unmarshalTypeError) Error() string {
	return fmt.Sprintf("%v %v, expected %v", e.msg, e.unmarshalErr.Field, e.unmarshalErr.Type)
}

func (e errExecCommandValidator) Error() string {
	return fmt.Sprintf("wrong os type. This machine works on Os: %s", e.osType)
}

func (x *requestBody) IsStructureEmpty() bool {
	return reflect.DeepEqual(x, requestBody{}) //
}

func bodyValidator(rb *requestBody) error {
	if rb.IsStructureEmpty() {
		return errEmptyBody
	}

	if ok := rb.Os == _os; !ok {
		return errExecCommandValidator{osType: _os}
	}

	return nil
}

func commandParser(b *requestBody) *runner.Command {
	switch _os {
	case "windows":
		return &runner.Command{
			Name: windowsCommandParser(b),
		}
	case "linux":
		command := linuxCommandParser(b)
		return &runner.Command{
			Name:  command.name,
			Flags: command.flags,
		}
	default:
		return nil
	}
}

func linuxCommandParser(b *requestBody) *linuxCommnd {
	flags := strings.Split(b.Cmd, " ")
	return &linuxCommnd{
		name:  flags[0],
		flags: flags[1:],
	}
}

func windowsCommandParser(b *requestBody) string {
	return fmt.Sprintf("%s %s", b.Cmd, b.Stdin)
}

func sendRespone(w http.ResponseWriter, r *runner.CommandRunner, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)

	resp := make(map[string]map[string]string)
	resp["message"] = map[string]string{
		"stdout": (r.CommandOutPut.StdOut).String(),
		"stderr": (r.CommandOutPut.StdErr).String(),
	}

	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}

func errorResponse(w http.ResponseWriter, err error, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)

	resp := make(map[string]string)
	resp["message"] = err.Error()

	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}
