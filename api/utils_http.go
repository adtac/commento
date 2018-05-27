package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
)

type response map[string]interface{}

// TODO: Add tests in utils_http_test.go

func unmarshalBody(r *http.Request, x interface{}) error {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Errorf("cannot read POST body: %v\n", err)
		return errorInternal
	}

	if err = json.Unmarshal(b, x); err != nil {
		return errorInvalidJSONBody
	}

	xv := reflect.Indirect(reflect.ValueOf(x))
	for i := 0; i < xv.NumField(); i++ {
		if xv.Field(i).IsNil() {
			return errorMissingField
		}
	}

	return nil
}

func writeBody(w http.ResponseWriter, x map[string]interface{}) error {
	resp, err := json.Marshal(x)
	if err != nil {
		w.Write([]byte(`{"success":false,"message":"Some internal error occurred"}`))
		logger.Errorf("cannot marshal response: %v\n")
		return errorInternal
	}

	w.Write(resp)
	return nil
}
