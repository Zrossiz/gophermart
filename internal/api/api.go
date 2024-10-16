package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Zrossiz/gophermart/internal/config"
)

type ApiService struct {
	cfg *config.Config
}

type OrderResponse struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual,omitempty"`
}

func New(cfg *config.Config) *ApiService {
	return &ApiService{cfg: cfg}
}

func (a *ApiService) UpdateOrder(orderID int) (string, float64, error) {
	url := fmt.Sprintf("%s/api/orders/%d", a.cfg.AcccrualSystemAddress, orderID)

	resp, err := http.Get(url)
	if err != nil {
		return "", 0.00, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", 0.00, err
	}

	var orderResp OrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&orderResp); err != nil {
		return "", 0.00, err
	}

	return orderResp.Status, orderResp.Accrual, nil
}
