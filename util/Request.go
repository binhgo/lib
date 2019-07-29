package util

import (
	"io/ioutil"
	"net/http"
)

func GetBodyBytes(r *http.Request) ([]byte, error) {
	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return nil, err
	}

	return b, nil
}
