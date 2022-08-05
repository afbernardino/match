package response_test

import (
	"net/http/httptest"
	"testing"

	"match/cmd/pkg/controller/response"
)

func TestWrite(t *testing.T) {
	rr := httptest.NewRecorder()

	b := `{"some_key","some_value"}`
	response.Write(rr, []byte(b))

	if rr.Body.String() != b {
		t.Errorf("returned unexpected body: want %v got %v", b, rr.Body.String())
	}
}

func TestWriteInternalServerError(t *testing.T) {
	rr := httptest.NewRecorder()

	response.WriteInternalServerError(rr)

	b := `{"error":"internal_server_error"}`
	if rr.Body.String() != b {
		t.Errorf("returned unexpected body: want %v got %v", b, rr.Body.String())
	}
}
