package req

import (
	"encoding/json"
	"net/http"
)

func Decode[T any](r *http.Request) (T, error) {
	var payload T
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return payload, err
	}
	return payload, nil
}