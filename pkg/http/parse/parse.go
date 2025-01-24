package parse

import (
	"errors"
	"net/http"
	"strconv"
)

var InvalidParam = errors.New("invalid query param")

func Uint64FromQueryParam(r *http.Request, param string) (uint64, error) {
	idStr := r.URL.Query().Get(param)
	var val uint64

	if len(idStr) == 0 {
		return val, InvalidParam
	}

	val, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return val, InvalidParam
	}

	return val, nil
}
