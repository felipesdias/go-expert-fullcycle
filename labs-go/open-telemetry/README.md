# Sistema de Consulta de Clima por CEP

Este projeto consiste em um sistema distribuído com dois serviços em Go para consultar a temperatura a partir de um CEP. Ele utiliza OpenTelemetry e Zipkin para tracing distribuído.

## Pré-requisitos

* Docker
* Docker Compose
* Uma chave de API da [WeatherAPI](https://www.weatherapi.com/)

## Como Executar

1.  **Configure a Chave de API:**
    Crie um arquivo `.env` com o conteudo `WEATHER_API_KEY={SUA_CHAVE_DA_WEATHER_API}`.

2.  **Suba os contêineres:**
    A partir da raiz do projeto, execute o comando:
    ```bash
    docker-compose up --build -d
    ```
    Isso irá construir as imagens dos serviços Go e iniciar todos os contêineres (`service-a`, `service-b`, `otel-collector`, `zipkin`).

## Como Testar

Após os contêineres estarem rodando, você pode testar a aplicação fazendo uma requisição `POST` para o `Serviço A`.

**Exemplo com cURL:**

```bash
# CEP Válido (Exemplo: Porto Alegre)
curl --location 'http://localhost:8080/' \
--header 'Content-Type: application/json' \
--data '{
    "cep": "90010040"
}'

# CEP Inválido (formato)
curl --location 'http://localhost:8080/' \
--header 'Content-Type: application/json' \
--data '{
    "cep": "123"
}'

# CEP Inexistente
curl --location 'http://localhost:8080/' \
--header 'Content-Type: application/json' \
--data '{
    "cep": "11111111"
}'
