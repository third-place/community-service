package util

import (
	"encoding/json"
	"net/http"
)

func WriteResponse(w http.ResponseWriter, object interface{}, err error) {
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	if object != nil {
		data, _ := json.Marshal(object)
		_, _ = w.Write(data)
	}
}
