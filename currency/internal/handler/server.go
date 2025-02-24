package handler

import (
	"github.com/vctrl/currency-service/currency/internal/service"
	"github.com/vctrl/currency-service/pkg/currency"

	"go.uber.org/zap"
)

type CurrencyServer struct {
	currency.UnimplementedCurrencyServiceServer
	service service.Currency
	logger  *zap.Logger
}

func NewCurrencyServer(svc service.Currency, logger *zap.Logger) CurrencyServer {
	return CurrencyServer{
		service: svc,
		logger:  logger,
	}
}
