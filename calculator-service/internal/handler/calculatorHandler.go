package Internal

import (
	"calculator-service/internal/entity"
	"encoding/json"
	"io"
	"net/http"
)

type Service interface {
	CalculatorService(input entity.Input) (output entity.Output)
}

type Handler struct {
	service Service
}

func New(s Service) *Handler {
	return &Handler{s}
}

func (h *Handler) CalculatorHandler(w http.ResponseWriter, req *http.Request) {

	if req.Method != "POST" {
		http.Error(w, "Not POST method", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != io.EOF && err != nil {
		panic(err)
	}
	defer req.Body.Close()

	var input entity.Input
	err = json.Unmarshal(body, &input)
	if err != nil {
		panic(err)
	}

	result := h.service.CalculatorService(input)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		panic(err)
	}

}
