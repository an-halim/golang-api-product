package utils

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CommonResponse struct {
	CorrelationID string      `json:"correlationid,omitempty"`
	Success       bool        `json:"success"`
	Error         string      `json:"error"`
	TIN           time.Time   `json:"tin,omitempty"`
	TOUT          time.Time   `json:"tout,omitempty"`
	Data          interface{} `json:"data"`
}

func Success(c *fiber.Ctx, code int, data interface{}, tin time.Time) error {
	return c.Status(code).JSON(CommonResponse{
		CorrelationID: uuid.New().String(),
		Success: true,
		Error: "",
		TIN: tin,
		TOUT: time.Now(),
		Data:data,
	})
}

func Failed(c *fiber.Ctx, code int, errorMessage string, tin time.Time) error {
	return c.Status(code).JSON(CommonResponse{
		CorrelationID: uuid.New().String(),
		Success: false,
		Error: errorMessage,
		TIN: tin,
		TOUT: time.Now(),
		Data: nil,
	})
}