package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateHandler_ServeHTTP(t *testing.T) {
	testCases := []struct {
		method       string
		expectedCode int
		expectedBody string
		path         string
	}{
		{method: http.MethodGet, expectedCode: http.StatusMethodNotAllowed, path: "/"},
		{method: http.MethodPut, expectedCode: http.StatusMethodNotAllowed, path: "/"},
		{method: http.MethodDelete, expectedCode: http.StatusMethodNotAllowed, path: "/"},
		{method: http.MethodPost, expectedCode: http.StatusOK, path: "/counter/Alloc/20"},
	}

	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
			r := httptest.NewRequest(tc.method, tc.path, nil)
			w := httptest.NewRecorder()
			h := http.StripPrefix("/", NewUpdateHandler())
			h.ServeHTTP(w, r)

			assert.Equal(t, tc.expectedCode, w.Code, "Код ответа не совпадает с ожидаемым")
		})
	}
}
