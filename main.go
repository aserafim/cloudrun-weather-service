package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

// Struct para resposta ViaCEP
type ViaCEPResponse struct {
	Localidade string `json:"localidade"`
	UF         string `json:"uf"`
	Erro       bool   `json:"erro,omitempty"`
}

// Struct para resposta WeatherAPI
type WeatherAPIResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

// Struct para resposta final JSON
type WeatherResponse struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

func main() {
	http.HandleFunc("/weather", weatherHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("Server listening on port", port)
	http.ListenAndServe(":"+port, nil)
}

func weatherHandler(w http.ResponseWriter, r *http.Request) {
	cep := r.URL.Query().Get("cep")

	// Valida CEP (exatamente 8 dígitos numéricos)
	matched, _ := regexp.MatchString(`^\d{8}$`, cep)
	if !matched {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	// Consulta ViaCEP
	viacepURL := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)
	resp, err := http.Get(viacepURL)
	if err != nil || resp.StatusCode != 200 {
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}
	defer resp.Body.Close()

	// Lê corpo para debug + decodifica
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "error reading zipcode response", http.StatusInternalServerError)
		return
	}
	// DEBUG: fmt.Println("ViaCEP response:", string(bodyBytes))

	// Recria o Body para decodificar
	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var viaCEP ViaCEPResponse
	if err := json.NewDecoder(resp.Body).Decode(&viaCEP); err != nil || viaCEP.Erro {
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}

	// Consulta WeatherAPI com cidade e país "Cidade,Brazil"
	apiKey := "2033a0227c094c97808185330252107"
	if apiKey == "" {
		http.Error(w, "missing WEATHERAPI_KEY", http.StatusInternalServerError)
		return
	}

	location := fmt.Sprintf("%s,Brazil", viaCEP.Localidade)
	weatherURL := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s", apiKey, location)
	weatherResp, err := http.Get(weatherURL)
	if err != nil || weatherResp.StatusCode != 200 {
		http.Error(w, "can not find weather data", http.StatusNotFound)
		return
	}
	defer weatherResp.Body.Close()

	var weather WeatherAPIResponse
	if err := json.NewDecoder(weatherResp.Body).Decode(&weather); err != nil {
		http.Error(w, "can not find weather data", http.StatusNotFound)
		return
	}

	// Converte temperaturas
	tempC := weather.Current.TempC
	tempF := tempC*1.8 + 32
	tempK := tempC + 273

	respJSON := WeatherResponse{
		TempC: round(tempC),
		TempF: round(tempF),
		TempK: round(tempK),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respJSON)
}

func round(f float64) float64 {
	res, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", f), 64)
	return res
}
