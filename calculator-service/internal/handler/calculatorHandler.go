package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"calculator-service/internal/service"
)

type CalculatorHandler struct {
	calcService *service.CalculatorService
}

func NewCalculatorHandler(calc *service.CalculatorService) *CalculatorHandler {
	return &CalculatorHandler{calcService: calc}
}

func (h *CalculatorHandler) HandleCalculate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req struct {
		A  float64 `json:"a"`
		B  float64 `json:"b"`
		Op string  `json:"op"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	resultMsg, err := h.calcService.Process(r.Context(), req.A, req.B, req.Op)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, service.ErrDivisionByZero) || errors.Is(err, service.ErrUnknownOp) {
			status = http.StatusBadRequest
		}
		writeJSONError(w, status, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Calculation processed and queued",
		"details": resultMsg,
	})
}

func writeJSONError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
