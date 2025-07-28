# Go Expert Full Cycle

Este repositório contém os projetos e desafios desenvolvidos durante o curso Go Expert da Full Cycle.

## Projetos

Aqui está uma visão geral dos projetos incluídos neste repositório:

### 1. Banco de Dados - Desafio Client-Server API

* **Localização:** `banco-dados/desafio-client-server-api/`
* **Descrição:** Uma aplicação cliente-servidor para consulta de cotação do dólar. O servidor busca a cotação de uma API externa, persiste em um banco de dados SQLite e a disponibiliza via um endpoint. O cliente consome essa API e salva a cotação em um arquivo `cotacao.txt`.
* **Tecnologias:** Go, SQLite.

### 2. Desafios Técnicos

#### Rate Limiter

* **Localização:** `desafios-tecnicos/rate-limiter/`
* **Descrição:** Um middleware de Rate Limiter em Go que controla o número de requisições por IP ou token de API. O estado do limiter é persistido no Redis. As configurações podem ser feitas via variáveis de ambiente ou arquivo `.env`.
* **Tecnologias:** Go, Redis, Docker.

#### Stress Test

* **Localização:** `desafios-tecnicos/stress-test/`
* **Descrição:** Uma ferramenta de linha de comando para realizar testes de carga em serviços web. É possível configurar a URL do serviço, o número total de requests e a quantidade de chamadas simultâneas. Ao final, um relatório é exibido com o tempo total, quantidade de requests com status 200 e a distribuição de outros status.
* **Tecnologias:** Go, Docker.

### 3. Labs Go

#### Concorrência em um Leilão

* **Localização:** `labs-go/concorrencia-leilao/`
* **Descrição:** Uma API de leilão que utiliza concorrência para processamento de lances em lote e para o fechamento automático de leilões após um tempo determinado. Os dados são armazenados no MongoDB.
* **Tecnologias:** Go, MongoDB, Docker, Gin.

#### Deploy no Cloud Run

* **Localização:** `labs-go/deploy-cloud-run/`
* **Descrição:** Uma aplicação de consulta de temperatura baseada em CEP, criada para ser implantada no Google Cloud Run. O serviço utiliza a API do ViaCEP para obter a localização e a WeatherAPI para obter a temperatura.
* **Tecnologias:** Go, Docker, Google Cloud Run.

#### OpenTelemetry

* **Localização:** `labs-go/open-telemetry/`
* **Descrição:** Um sistema distribuído para consulta de clima por CEP, composto por dois serviços. O projeto utiliza OpenTelemetry com Zipkin para tracing distribuído, permitindo observar a comunicação entre os serviços.
* **Tecnologias:** Go, Docker, OpenTelemetry, Zipkin.