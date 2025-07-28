# Go Rate Limiter

Este projeto implementa um Rate Limiter em Go, projetado para atuar como um middleware em um servidor web. Ele controla o número de requisições por segundo com base no endereço IP do cliente ou em um token de API fornecido.

A persistência do estado do limiter é realizada através do Redis, com uma arquitetura baseada em "strategy" que permite a substituição do Redis por outro mecanismo de armazenamento.

## Funcionalidades

-   Limitação de requisições por endereço IP.
-   Limitação de requisições por token de API (`API_KEY` no header).
-   As configurações de token se sobrepõem às configurações de IP.
-   Período de bloqueio configurável para IPs ou tokens que excedem o limite.
-   Configuração via variáveis de ambiente ou arquivo `.env`.
-   Armazenamento de dados no Redis.
-   Interface de armazenamento (`LimiterStorage`) para fácil extensibilidade.

## Estrutura

-   `cmd/server/main.go`: Ponto de entrada da aplicação. Inicializa o servidor web e o middleware.
-   `internal/config`: Carregamento e parsing das configurações.
-   `internal/storage`: Abstração (`storage.go`) e implementação (`redis.go`) da camada de persistência.
-   `internal/limiter`: Lógica principal do rate limiter.
-   `internal/middleware`: Middleware HTTP que integra o limiter ao servidor.
-   `Dockerfile`: Define a imagem Docker para a aplicação.
-   `docker-compose.yml`: Orquestra os contêineres da aplicação e do Redis.

## Configuração

Crie um arquivo `.env` na raiz do projeto ou defina as seguintes variáveis de ambiente:

| Variável                     | Descrição                                                                            | Padrão no `.env` de exemplo |
| ---------------------------- | -------------------------------------------------------------------------------------- | --------------------------- |
| `REDIS_ADDR`                 | Endereço do servidor Redis.                                                            | `redis:6379`                |
| `REDIS_PASSWORD`             | Senha do Redis.                                                                        | Vazio                       |
| `REDIS_DB`                   | Número do banco de dados Redis.                                                        | `0`                         |
| `IP_REQUESTS_PER_SECOND`     | Número máximo de requisições por segundo para um único IP.                             | `5`                         |
| `IP_BLOCK_DURATION_SECONDS`  | Duração do bloqueio (em segundos) para um IP que excedeu o limite.                     | `300`                       |
| `TOKEN_LIMITS`               | Limites para tokens específicos. Formato: `token1:limite1,token2:limite2`.             | `token1:20`            |
| `TOKEN_BLOCK_DURATION_SECONDS` | Duração do bloqueio (em segundos) para um token que excedeu o limite.                | `600`                       |

## Como Executar

**Pré-requisitos:**
* Docker
* Docker Compose

1.  Clone o repositório.
2.  Certifique-se de que o arquivo `.env` está configurado corretamente na raiz do projeto.
3.  Execute o seguinte comando na raiz do projeto:

    ```bash
    docker-compose up --build
    ```

A aplicação estará disponível em `http://localhost:8080`.

## Como Testar

### Testes Automatizados

Para executar os testes unitários, utilize o comando:

```bash
go test ./...
```

### Testes Manuais

Você pode usar ferramentas como `curl` para testar o rate limiter.

**Teste de Limite por IP (padrão: 5 req/s):**

```bash
# Execute 6 vezes rapidamente. A 6ª requisição deve falhar.
for i in {1..6}; do curl -i http://localhost:8080/; done

# Resposta da 6ª requisição (ou posterior):
# HTTP/1.1 429 Too Many Requests
# Content-Type: application/json
# {"message": "you have reached the maximum number of requests or actions allowed within a certain time frame"}
```

**Teste de Limite por Token (padrão `token1`: 4 req/s):**

```bash
# Execute 5 vezes rapidamente. A 5ª requisição deve falhar.
for i in {1..5}; do curl -i -H "API_KEY: token1" http://localhost:8080/; done

# Resposta da 5ª requisição (ou posterior):
# HTTP/1.1 429 Too Many Requests
# ...
```