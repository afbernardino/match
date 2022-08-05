package partners_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"match/cmd/pkg/controller/partners"
	"match/cmd/pkg/controller/partners/mock"
	"match/cmd/pkg/models"
	"match/cmd/pkg/repository"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
)

const testMatchRequestBody string = `
{
    "materials": [1, 2],
    "address": {
		"lat": 1.1,
		"long": 1.2
	},
	"square_meters": 5,
	"phone_number": "+351912345678"
}
`

func TestGetMatches_InvalidBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	db := mock.NewMockDatabase(ctrl)

	handler := partners.NewHandler(db)
	rr := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodPost, "/partners/match", strings.NewReader(""))

	handler.GetMatches(rr, req)

	expectedCode := http.StatusBadRequest
	if rr.Code != expectedCode {
		t.Errorf("status code mismatch: want %v got %v", expectedCode, rr.Code)
	}

	expectedBody := `{"error":"bad_request"}`
	if rr.Body.String() != expectedBody {
		t.Errorf("body mismatch: want %v got %v", expectedBody, rr.Body.String())
	}
}

func TestGetMatches_NoAddress(t *testing.T) {
	ctrl := gomock.NewController(t)
	db := mock.NewMockDatabase(ctrl)

	handler := partners.NewHandler(db)
	rr := httptest.NewRecorder()

	reqBody := `
	{
		"materials": [1, 2],
		"square_meters": 5,
		"phone_number": "+351912345678"
	}
	`
	req := httptest.NewRequest(http.MethodPost, "/partners/match", strings.NewReader(reqBody))

	handler.GetMatches(rr, req)

	expectedCode := http.StatusBadRequest
	if rr.Code != expectedCode {
		t.Errorf("status code mismatch: want %v got %v", expectedCode, rr.Code)
	}

	expectedBody := `{"error":"bad_request"}`
	if rr.Body.String() != expectedBody {
		t.Errorf("body mismatch: want %v got %v", expectedBody, rr.Body.String())
	}
}

func TestGetMatches_NoMaterials(t *testing.T) {
	ctrl := gomock.NewController(t)
	db := mock.NewMockDatabase(ctrl)

	handler := partners.NewHandler(db)
	rr := httptest.NewRecorder()

	reqBody := `
	{
		"address": {
			"lat": 1.1,
			"long": 1.2
		},
		"square_meters": 5,
		"phone_number": "+351912345678"
	}
	`
	req := httptest.NewRequest(http.MethodPost, "/partners/match", strings.NewReader(reqBody))

	handler.GetMatches(rr, req)

	expectedCode := http.StatusBadRequest
	if rr.Code != expectedCode {
		t.Errorf("status code mismatch: want %v got %v", expectedCode, rr.Code)
	}

	expectedBody := `{"error":"bad_request"}`
	if rr.Body.String() != expectedBody {
		t.Errorf("body mismatch: want %v got %v", expectedBody, rr.Body.String())
	}
}

func TestGetMatches_DatabaseFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	db := mock.NewMockDatabase(ctrl)

	db.EXPECT().
		GetMatches(gomock.Any(), []uint{1, 2}, float32(1.1), float32(1.2)).
		Return(nil, errors.New("some error"))

	handler := partners.NewHandler(db)
	rr := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodPost, "/partners/match", strings.NewReader(testMatchRequestBody))

	handler.GetMatches(rr, req)

	expectedCode := http.StatusInternalServerError
	if rr.Code != expectedCode {
		t.Errorf("status code mismatch: want %v got %v", expectedCode, rr.Code)
	}

	expectedBody := `{"error":"internal_server_error"}`
	if rr.Body.String() != expectedBody {
		t.Errorf("body mismatch: want %v got %v", expectedBody, rr.Body.String())
	}
}

func TestGetMatches_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	db := mock.NewMockDatabase(ctrl)

	p := models.Partner{
		ID: 3,
		Categories: []models.Category{
			{
				ID:          4,
				PartnerID:   3,
				Description: "category 4",
			},
		},
		Materials: []models.Material{
			{
				ID:          1,
				PartnerID:   3,
				Description: "material 1",
			},
			{
				ID:          2,
				PartnerID:   3,
				Description: "material 2",
			},
		},
		Address: models.Address{
			Lat:  1.1,
			Long: 1.2,
		},
		Radius: 100,
		Rating: 5,
	}

	db.EXPECT().
		GetMatches(gomock.Any(), []uint{1, 2}, float32(1.1), float32(1.2)).
		Return([]models.Partner{p}, nil)

	handler := partners.NewHandler(db)
	rr := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodPost, "/partners/match", strings.NewReader(testMatchRequestBody))

	handler.GetMatches(rr, req)

	expectedCode := http.StatusOK
	if rr.Code != expectedCode {
		t.Errorf("status code mismatch: want %v got %v", expectedCode, rr.Code)
	}

	expectedBodyJson := `
	[
		{
			"id": 3,
			"categories": [
				{
					"id": 4,
					"description": "category 4"
				}
			],
			"materials": [
				{
					"id": 1,
					"description": "material 1"
				},
				{
					"id": 2,
					"description": "material 2"
				}
			],
			"address": {
				"lat": 1.1,
				"long": 1.2
			},
			"radius": 100,
			"rating": 5
		}
	]
	`
	buffer := new(bytes.Buffer)
	_ = json.Compact(buffer, []byte(expectedBodyJson))
	expectedBody := buffer.String()

	if rr.Body.String() != expectedBody {
		t.Errorf("body mismatch: want %v got %v", expectedBody, rr.Body.String())
	}
}

func TestGetPartnerById_InvalidId(t *testing.T) {
	ctrl := gomock.NewController(t)
	db := mock.NewMockDatabase(ctrl)

	handler := partners.NewHandler(db)
	rr := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodGet, "/partners/a", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "a"})

	handler.GetPartnerById(rr, req)

	expectedCode := http.StatusBadRequest
	if rr.Code != expectedCode {
		t.Errorf("status code mismatch: want %v got %v", expectedCode, rr.Code)
	}

	expectedBody := `{"error":"bad_request"}`
	if rr.Body.String() != expectedBody {
		t.Errorf("body mismatch: want %v got %v", expectedBody, rr.Body.String())
	}
}

func TestGetPartnerById_PartnerNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	db := mock.NewMockDatabase(ctrl)

	db.EXPECT().
		GetPartnerById(gomock.Any(), uint(1)).
		Return(models.Partner{}, repository.ErrNotFound)

	handler := partners.NewHandler(db)
	rr := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodGet, "/partners/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	handler.GetPartnerById(rr, req)

	expectedCode := http.StatusNotFound
	if rr.Code != expectedCode {
		t.Errorf("status code mismatch: want %v got %v", expectedCode, rr.Code)
	}

	expectedBody := `{"error":"not_found"}`
	if rr.Body.String() != expectedBody {
		t.Errorf("body mismatch: want %v got %v", expectedBody, rr.Body.String())
	}
}

func TestGetPartnerById_DatabaseFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	db := mock.NewMockDatabase(ctrl)

	db.EXPECT().
		GetPartnerById(gomock.Any(), uint(1)).
		Return(models.Partner{}, errors.New("some error"))

	handler := partners.NewHandler(db)
	rr := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodGet, "/partners/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	handler.GetPartnerById(rr, req)

	expectedCode := http.StatusInternalServerError
	if rr.Code != expectedCode {
		t.Errorf("status code mismatch: want %v got %v", expectedCode, rr.Code)
	}

	expectedBody := `{"error":"internal_server_error"}`
	if rr.Body.String() != expectedBody {
		t.Errorf("body mismatch: want %v got %v", expectedBody, rr.Body.String())
	}
}

func TestGetPartnerById_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	db := mock.NewMockDatabase(ctrl)

	p := models.Partner{
		ID: 3,
		Categories: []models.Category{
			{
				ID:          4,
				PartnerID:   3,
				Description: "category 4",
			},
		},
		Materials: []models.Material{
			{
				ID:          1,
				PartnerID:   3,
				Description: "material 1",
			},
			{
				ID:          2,
				PartnerID:   3,
				Description: "material 2",
			},
		},
		Address: models.Address{
			Lat:  1.1,
			Long: 1.2,
		},
		Radius: 100,
		Rating: 5,
	}

	db.EXPECT().
		GetPartnerById(gomock.Any(), uint(1)).
		Return(p, nil)

	handler := partners.NewHandler(db)
	rr := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodGet, "/partners/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	handler.GetPartnerById(rr, req)

	expectedCode := http.StatusOK
	if rr.Code != expectedCode {
		t.Errorf("status code mismatch: want %v got %v", expectedCode, rr.Code)
	}

	expectedBodyJson := `
	{
		"id": 3,
		"categories": [
			{
				"id": 4,
				"description": "category 4"
			}
		],
		"materials": [
			{
				"id": 1,
				"description": "material 1"
			},
			{
				"id": 2,
				"description": "material 2"
			}
		],
		"address": {
			"lat": 1.1,
			"long": 1.2
		},
		"radius": 100,
		"rating": 5
	}
	`
	buffer := new(bytes.Buffer)
	_ = json.Compact(buffer, []byte(expectedBodyJson))
	expectedBody := buffer.String()

	if rr.Body.String() != expectedBody {
		t.Errorf("body mismatch: want %v got %v", expectedBody, rr.Body.String())
	}
}
