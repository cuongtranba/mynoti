package usecase

import (
	"io"
	"net/http"

	"github.com/cuongtranba/mynoti/internal/domain"
	"github.com/cuongtranba/mynoti/pkg/app_context"
	"github.com/pkg/errors"
)

type htmlFetcher struct {
	httpClient *http.Client
}

func NewHtmlFetcher(httpClient *http.Client) domain.HtmlFetcher {
	return &htmlFetcher{
		httpClient: httpClient,
	}
}

func (h *htmlFetcher) Fetch(ctx *app_context.AppContext, url string) (string, error) {
	resp, err := h.httpClient.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	isErrorStatus := resp.StatusCode >= 400
	if isErrorStatus {
		return "", errors.Errorf("fetch html failed with status code:%d ", resp.StatusCode)
	}
	return string(body), nil
}
