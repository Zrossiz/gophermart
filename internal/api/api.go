package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/Zrossiz/gophermart/internal/apperrors"
	"github.com/Zrossiz/gophermart/internal/config"
	"github.com/Zrossiz/gophermart/internal/dto"
)

type APIService struct {
	cfg *config.Config
}

func New(cfg *config.Config) *APIService {
	return &APIService{cfg: cfg}
}

func (a *APIService) UpdateOrder(orderID int) (dto.ExternalOrderResponse, error) {
	var orderResp dto.ExternalOrderResponse
	url := fmt.Sprintf("%s/api/orders/%d", a.cfg.AcccrualSystemAddress, orderID)

	resp, err := http.Get(url)
	if err != nil {
		return orderResp, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return orderResp, err
	}

	if err := json.NewDecoder(resp.Body).Decode(&orderResp); err != nil {
		return orderResp, err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return orderResp, err
		}

		if err := json.Unmarshal(data, &orderResp); err != nil {
			return orderResp, err
		}

		return orderResp, nil
	case http.StatusNoContent:
		return orderResp, nil
	case http.StatusTooManyRequests:
		retryHeader := resp.Header.Get("Retry-After")
		retryAfter, err := strconv.Atoi(retryHeader)
		if err != nil {
			return orderResp, apperrors.ErrTooManyRequests
		}

		go func(wait time.Duration) {
			time.Sleep(wait)
		}(time.Duration(retryAfter) * time.Second)

		return orderResp, apperrors.ErrTooManyRequests
	}

	return orderResp, nil
}
