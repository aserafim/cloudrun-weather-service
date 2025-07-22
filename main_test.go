package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Testa o endpoint /weather com um CEP válido
func TestWeatherHandlerSuccess(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/weather?cep=06020194", nil)
	rec := httptest.NewRecorder()

	weatherHandler(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("esperado status 200, obtido %d", res.StatusCode)
	}

	var wr WeatherResponse
	if err := json.NewDecoder(res.Body).Decode(&wr); err != nil {
		t.Fatalf("erro ao decodificar resposta JSON: %v", err)
	}

	if wr.TempC == 0 {
		t.Errorf("TempC deve ser diferente de 0")
	}
}

// Testa o endpoint com um CEP inválido
func TestWeatherHandlerInvalidCEP(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/weather?cep=123", nil)
	rec := httptest.NewRecorder()

	weatherHandler(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Errorf("esperado 422 para CEP inválido, obtido %d", rec.Code)
	}
}

// Testa o endpoint sem CEP
func TestWeatherHandlerMissingCEP(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/weather", nil)
	rec := httptest.NewRecorder()

	weatherHandler(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Errorf("esperado 422 para CEP ausente, obtido %d", rec.Code)
	}
}
