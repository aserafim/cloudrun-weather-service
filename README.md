# 🌤️ cloudrun-weather-service

Este serviço recebe um **CEP brasileiro válido** e retorna a temperatura atual da cidade correspondente, em graus Celsius, Fahrenheit e Kelvin.

Ele utiliza as APIs:
- [ViaCEP](https://viacep.com.br/) para obter o nome da cidade a partir do CEP.
- [WeatherAPI](https://www.weatherapi.com/) para buscar a temperatura atual da cidade.

---

## 🚀 Como usar

### Requisitos

- Go 1.18+
- Docker (opcional, para build containerizado)
- Conta gratuita na [WeatherAPI](https://www.weatherapi.com/) (você precisará de uma API Key)
- (opcional) GCP com Cloud Run habilitado

---

## 💻 Executando localmente

```bash
go run main.go
````

Por padrão o serviço roda em `http://localhost:8080`.

### Exemplo de chamada

```bash
curl "http://localhost:8080/weather?cep=06020194"
```

#### Resposta:

```json
{
  "temp_C": 24.0,
  "temp_F": 75.2,
  "temp_K": 297.0
}
```

---

## 🐳 Executando com Docker

### Build da imagem

```bash
docker build -t cloudrun-weather-service .
```

### Rodando o container

```bash
docker run -p 8080:8080 cloudrun-weather-service
```

---

## ☁️ Serviço no Google Cloud Run

1. Exemplo de chamada

```bash
curl https://cloudrun-weather-service-10747099608.us-central1.run.app/weather?cep=01311000
```

2. Retorno:

```bash
{"temp_C":17.2,"temp_F":62.96,"temp_K":290.2}
```
---

## 📂 Estrutura do Projeto

```
.
├── Dockerfile
├── go.mod
├── main.go
├── main_test.go
└── README.md
```
---

## 🧪 Testes

Execute os testes unitários com:

```bash
go test -v
```

---

## 📄 Licença

[MIT](LICENSE)

---

## ✨ Autor

Alefe Serafim – [LinkedIn](https://www.linkedin.com/in/alefeserafim)
