package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ElfAstAhe/goph-keeper/api/rest/dto"
	"github.com/ElfAstAhe/goph-keeper/pkg/utils"
)

type GophKeeperSimpleClient struct {
	baseURL    string
	httpClient *http.Client
	token      string
}

func NewGophKeeperSimpleClient(baseURL string) *GophKeeperSimpleClient {
	return &GophKeeperSimpleClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 600 * time.Second,
		},
	}
}

// Login - Метод для логина (чтобы получить токен)
func (c *GophKeeperSimpleClient) Login(ctx context.Context, username, encryptedPassword string) error {
	payload := dto.NewLoginDto(username, encryptedPassword)

	var res = new(dto.LoginResultDto)
	err := c.doRequest(ctx, http.MethodPost, "/api/auth/login", payload, res, false)
	if err != nil {
		return NewClientError("Login", -1, "", err)
	}

	c.token = res.Token // Запоминаем токен для следующих вызовов

	return nil
}

// Register - Метод регистрации нового пользователя
func (c *GophKeeperSimpleClient) Register(ctx context.Context, username, password, person, eMail string) (*dto.RegisterResultDto, error) {
	payload := dto.NewRegisterDto(username, password, person, eMail)

	var res = new(dto.RegisterResultDto)
	err := c.doRequest(ctx, http.MethodPost, "/api/auth/register", payload, res, false)
	if err != nil {
		return nil, NewClientError("Register", -1, "", err)
	}

	return res, nil
}

// GetProfile - Метод для получения профиля
func (c *GophKeeperSimpleClient) GetProfile(ctx context.Context) (*dto.UserDto, error) {
	var res = new(dto.UserDto)
	err := c.doRequest(ctx, http.MethodGet, "/api/users/profile", nil, res, true)
	if err != nil {
		return nil, NewClientError("GetProfile", -1, "", err)
	}

	return res, nil
}

func (c *GophKeeperSimpleClient) ChangeKeys(ctx context.Context) (*dto.ChangeKeysResultDto, error) {
	var res = new(dto.ChangeKeysResultDto)
	err := c.doRequest(ctx, http.MethodPut, "/api/users/keys", nil, res, true)
	if err != nil {
		return nil, NewClientError("ChangeKeys", -1, "", err)
	}

	return res, nil
}

func (c *GophKeeperSimpleClient) UpdatePassword(ctx context.Context, oldPassword, newPassword string) error {
	payload := dto.NewUpdatePasswordDto(oldPassword, newPassword)
	err := c.doRequest(ctx, http.MethodPut, "/api/users/password", payload, nil, true)
	if err != nil {
		return NewClientError("UpdatePassword", -1, "", err)
	}

	return nil
}

func (c *GophKeeperSimpleClient) Get(ctx context.Context, id string) (*dto.UserDataDto, error) {
	var res = new(dto.UserDataDto)
	err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/api/users/data/%s", id), nil, res, true)
	if err != nil {
		return nil, NewClientError("Get", -1, "", err)
	}

	return res, nil
}

func (c *GophKeeperSimpleClient) GetByKey(ctx context.Context, dataKind string, name string) (*dto.UserDataDto, error) {
	var res = new(dto.UserDataDto)
	err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/api/users/data/%s/%s", dataKind, name), nil, res, true)
	if err != nil {
		return nil, NewClientError("GetByKey", -1, "", err)
	}

	return res, nil
}

func (c *GophKeeperSimpleClient) ListAll(ctx context.Context) ([]*dto.UserDataDto, error) {
	var res []*dto.UserDataDto
	err := c.doRequest(ctx, http.MethodGet, "/api/users/data/list", nil, &res, true)
	if err != nil {
		return nil, NewClientError("ListAll", -1, "", err)
	}

	return res, nil
}

func (c *GophKeeperSimpleClient) Save(ctx context.Context, data *dto.UserDataDto) (*dto.UserDataDto, error) {
	if data == nil {
		return nil, NewClientError("Save", -1, "data is nil", nil)
	}
	if data.ID == "" {
		return c.create(ctx, data)
	}

	return c.change(ctx, data)
}

func (c *GophKeeperSimpleClient) create(ctx context.Context, data *dto.UserDataDto) (*dto.UserDataDto, error) {
	var res = new(dto.UserDataDto)
	err := c.doRequest(ctx, http.MethodPost, "/api/users/data", data, res, true)
	if err != nil {
		return nil, NewClientError("create", -1, "", err)
	}

	return res, nil
}

func (c *GophKeeperSimpleClient) change(ctx context.Context, data *dto.UserDataDto) (*dto.UserDataDto, error) {
	var res = new(dto.UserDataDto)
	err := c.doRequest(ctx, http.MethodPut, fmt.Sprintf("/api/users/data/%s", data.ID), data, res, true)
	if err != nil {
		return nil, NewClientError("change", -1, "", err)
	}

	return res, nil
}

func (c *GophKeeperSimpleClient) Delete(ctx context.Context, id string) error {
	err := c.doRequest(ctx, http.MethodDelete, fmt.Sprintf("/api/users/data/%s", id), nil, nil, true)
	if err != nil {
		return NewClientError("Delete", -1, "", err)
	}

	return nil
}

// doRequest -  Универсальный внутренний метод для запросов
func (c *GophKeeperSimpleClient) doRequest(ctx context.Context, method, path string, body any, target any, useAuth bool) error {
	var bodyReader *bytes.Buffer
	if body != nil {
		b, _ := json.Marshal(body)
		bodyReader = bytes.NewBuffer(b)
	} else {
		bodyReader = bytes.NewBuffer([]byte{})
	}

	req, _ := http.NewRequestWithContext(ctx, method, c.baseURL+path, bodyReader)
	req.Header.Set("Content-Type", "application/json")

	if useAuth && c.token != "" {
		req.AddCookie(&http.Cookie{
			Name:     utils.DefaultCookieName,
			Value:    c.token,
			SameSite: http.SameSiteLaxMode,
		})
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return NewClientError("doRequest", -1, "request", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return NewClientError("doRequest", resp.StatusCode, "response", nil)
	}

	if target != nil {
		return json.NewDecoder(resp.Body).Decode(target)
	}

	return nil
}
