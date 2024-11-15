package delivery

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cuongtranba/mynoti/internal/domain/mocks"
	"github.com/cuongtranba/mynoti/internal/usecase"
	"github.com/stretchr/testify/mock"
)

func TestNewServer(t *testing.T) {
	repoMock := new(mocks.ComicRepository)
	repoMock.On("Save", mock.Anything, mock.Anything).Return(nil)

	payload := Comic{
		Url:         "http://example.com",
		Name:        "Test Comic",
		Description: "Test Description",
		Html:        "<div>Test</div>",
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/subscribe", bytes.NewBuffer(payloadBytes))
	rec := httptest.NewRecorder()

	server := NewServer(":8080", usecase.NewComicUseCase(repoMock))
	server.server.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, rec.Code)
	}
}
