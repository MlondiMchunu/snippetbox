package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"snippetbox.mlodev.net/internal/assert"
)

func TestPing(t *testing.T) {
	rr := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, "/", nil)

	if err != nil {
		t.Fatal(err)
	}
	ping(rr, req)
	rs := rr.Result()

	assert.Equal(t, rs.StatusCode, http.StatusOK)

}
