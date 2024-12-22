package routes

import (
	"Yandex_Calc/internal/eval"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// Request структура для входящих запросов
type Request struct {
	Expression string `json:"expression"`
}

// Response структура для ответа
type Response struct {
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

// SetupRoutes настраивает маршруты для сервиса
func SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/calculate", calculateHandler)
	return mux
}

// calculateHandler обрабатывает запросы к /api/v1/calculate
func calculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	expr := strings.TrimSpace(req.Expression)
	if expr == "" {
		respondWithError(w, http.StatusUnprocessableEntity, "Expression is not valid")
		return
	}

	result, err := eval.Eval(expr)
	if err != nil {
		if errors.Is(err, eval.ErrInvalidExpression) {
			respondWithError(w, http.StatusUnprocessableEntity, "Expression is not valid")
		} else {
			respondWithError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	respondWithResult(w, eval.BigratToFloat(result))
}

// respondWithError отвечает ошибкой
func respondWithError(w http.ResponseWriter, statusCode int, errorMessage string) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(Response{Error: errorMessage})
}

// respondWithResult отвечает результатом
func respondWithResult(w http.ResponseWriter, result float64) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{Result: formatResult(result)})
}

// formatResult форматирует результат
func formatResult(result float64) string {
	return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.2f", result), "0"), ".")
}
