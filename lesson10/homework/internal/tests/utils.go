package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"homework10/internal/adapters/adrepo"
	"homework10/internal/adapters/userrepo"
	"homework10/internal/app"
	"homework10/internal/ports/httpgin"
	"io"
	"net/http"
	"net/http/httptest"
)

type adData struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Text      string `json:"text"`
	AuthorID  int64  `json:"author_id"`
	Published bool   `json:"published"`
}

type userData struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	UserID   int64  `json:"user_id"`
}

type adResponse struct {
	Data adData `json:"data"`
}

type adDeleteResponse struct {
	ID      int64  `json:"id"`
	Success string `json:"success"`
}

type userDeleteResponse struct {
	ID      int64  `json:"id"`
	Success string `json:"success"`
}

type adsResponse struct {
	Data []adData `json:"data"`
}

type userResponse struct {
	Data userData `json:"data"`
}

var (
	ErrBadRequest = fmt.Errorf("bad request")
	ErrForbidden  = fmt.Errorf("forbidden")
)

type testClient struct {
	client  *http.Client
	baseURL string
}

func getTestClient() *testClient {
	server := httpgin.NewHTTPServer(":18080", app.NewApp(adrepo.New(), userrepo.New()))
	testServer := httptest.NewServer(server.Handler)

	return &testClient{
		client:  testServer.Client(),
		baseURL: testServer.URL,
	}
}

func (tc *testClient) getResponse(req *http.Request, out any) error {
	resp, err := tc.client.Do(req)
	if err != nil {
		return fmt.Errorf("unexpected error: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusBadRequest {
			return ErrBadRequest
		}
		if resp.StatusCode == http.StatusForbidden {
			return ErrForbidden
		}
		return fmt.Errorf("unexpected status code: %s", resp.Status)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("unable to read response: %w", err)
	}

	err = json.Unmarshal(respBody, out)
	if err != nil {
		return fmt.Errorf("unable to unmarshal: %w", err)
	}

	return nil
}

func (tc *testClient) createAd(userID int64, title string, text string) (adResponse, error) {
	body := map[string]any{
		"user_id": userID,
		"title":   title,
		"text":    text,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return adResponse{}, fmt.Errorf("unable to marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, tc.baseURL+"/api/v1/ads", bytes.NewReader(data))
	if err != nil {
		return adResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	var response adResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return adResponse{}, err
	}

	return response, nil
}

func (tc *testClient) deleteAd(adID any, userID int64) (adDeleteResponse, error) {
	body := map[string]any{
		"ad_id":   adID,
		"user_id": userID,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return adDeleteResponse{}, fmt.Errorf("unable to marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf(tc.baseURL+"/api/v1/ads/%s", adID), bytes.NewReader(data))
	if err != nil {
		return adDeleteResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	var response adDeleteResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return adDeleteResponse{}, err
	}

	return response, nil
}

func (tc *testClient) deleteUser(userID any) (userDeleteResponse, error) {
	body := map[string]any{
		"ad_id": userID,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return userDeleteResponse{}, fmt.Errorf("unable to marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf(tc.baseURL+"/api/v1/user/%s", userID), bytes.NewReader(data))
	if err != nil {
		return userDeleteResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	var response userDeleteResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return userDeleteResponse{}, err
	}

	return response, nil
}

func (tc *testClient) changeAdStatus(userID int64, adID any, published bool) (adResponse, error) {
	body := map[string]any{
		"user_id":   userID,
		"published": published,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return adResponse{}, fmt.Errorf("unable to marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf(tc.baseURL+"/api/v1/ads/%s/status", adID), bytes.NewReader(data))
	if err != nil {
		return adResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	var response adResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return adResponse{}, err
	}

	return response, nil
}

func (tc *testClient) updateAd(userID int64, adID any, title string, text string) (adResponse, error) {
	body := map[string]any{
		"user_id": userID,
		"title":   title,
		"text":    text,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return adResponse{}, fmt.Errorf("unable to marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf(tc.baseURL+"/api/v1/ads/%s", adID), bytes.NewReader(data))
	if err != nil {
		return adResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	var response adResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return adResponse{}, err
	}

	return response, nil
}

func (tc *testClient) createUser(nick string, email string) (userResponse, error) {
	body := map[string]any{
		"nickname": nick,
		"email":    email,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return userResponse{}, fmt.Errorf("unable to marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, tc.baseURL+"/api/v1/user", bytes.NewReader(data))
	if err != nil {
		return userResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	var response userResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return userResponse{}, err
	}

	return response, nil
}

func (tc *testClient) listAds() (adsResponse, error) {
	req, err := http.NewRequest(http.MethodGet, tc.baseURL+"/api/v1/ads", nil)
	if err != nil {
		return adsResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	var response adsResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return adsResponse{}, err
	}

	return response, nil
}

func (tc *testClient) updateUser(userID any, nickname string, email string, authorID int64) (userResponse, error) {
	body := map[string]any{
		"nickname": nickname,
		"email":    email,
		"user_id":  authorID,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return userResponse{}, fmt.Errorf("unable to marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf(tc.baseURL+"/api/v1/user/%s", userID), bytes.NewReader(data))
	if err != nil {
		return userResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	var response userResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return userResponse{}, err
	}

	return response, nil
}

func (tc *testClient) getUser(userID any) (userResponse, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(tc.baseURL+"/api/v1/user/%s", userID), nil)
	if err != nil {
		return userResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	var response userResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return userResponse{}, err
	}

	return response, nil
}

func (tc *testClient) getAd(adID any) (adResponse, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(tc.baseURL+"/api/v1/ads/%s", adID), nil)
	if err != nil {
		return adResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	var response adResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return adResponse{}, err
	}

	return response, nil
}

func (tc *testClient) listFilterAds(pub string, date string, author string) (adsResponse, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(tc.baseURL+"/api/v1/ads/?pub=%s&date=%s&author=%s", pub, date, author), nil)
	if err != nil {
		return adsResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	var response adsResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return adsResponse{}, err
	}

	return response, nil
}
