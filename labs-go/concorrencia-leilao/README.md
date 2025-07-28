# Leilão Full Cycle - Fechamento Automático

Este documento descreve a implementação da funcionalidade de fechamento automático de leilões e fornece um guia sobre como configurar e executar o projeto em um ambiente de desenvolvimento.

## Ideia da Solução

O objetivo principal era adicionar uma rotina que fechasse automaticamente os leilões após um período pré-definido, sem alterar o fluxo existente de criação de leilões e lances.

Para alcançar isso, a solução foi implementada utilizando os seguintes conceitos:

1.  **Go Routine e Concorrência**: Ao iniciar o repositório da entidade `Auction`, uma nova *go routine* é lançada para trabalhar em segundo plano, de forma concorrente com a aplicação principal. Isso garante que a verificação dos leilões não impacte a performance das requisições HTTP.

2.  **Cálculo de Tempo via Variável de Ambiente**: A duração de um leilão é definida pela variável de ambiente `AUCTION_INTERVAL`. Uma função auxiliar (`getAuctionInterval`) lê essa variável e a converte para um `time.Duration`. Caso a variável não esteja definida, um valor padrão é utilizado, tornando a configuração flexível.

3.  **Verificação Periódica com `Ticker`**: Dentro da go routine, um `time.Ticker` é utilizado para executar a lógica de verificação em intervalos regulares (a cada 5 segundos). A cada "tick", o sistema executa os seguintes passos:

      * Busca por todos os leilões que ainda estão com o status `Active` no banco de dados.
      * Para cada leilão ativo, calcula-se o tempo final somando o timestamp de sua criação com o `AUCTION_INTERVAL`.
      * Se o tempo atual for posterior ao tempo final do leilão, significa que ele expirou.
      * O status do leilão expirado é então atualizado para `Completed` no banco de dados.

Essa abordagem garante que o sistema seja autônomo e resiliente, fechando os leilões de forma assíncrona e automática.

## Ambiente de Desenvolvimento

Siga os passos abaixo para executar a aplicação e os testes localmente.

### Pré-requisitos

Para executar o projeto, você precisará ter instalado em sua máquina:

  * Go (versão 1.20 ou superior)
  * Docker
  * Docker Compose

### Executando a Aplicação

1.  Clone o repositório para a sua máquina local.

2.  Na raiz do projeto, suba os contêineres da aplicação e do banco de dados utilizando o Docker Compose:

    ```bash
    docker-compose up --build
    ```

    Este comando irá construir a imagem da aplicação Go, iniciar o servidor e o banco de dados MongoDB. A API estará disponível no endereço `http://localhost:8080`.

### Executando os Testes

Para validar a nova funcionalidade de fechamento automático e garantir a integridade das demais partes do sistema, siga os passos:

1.  **Inicie o serviço do MongoDB**: É necessário que o banco de dados esteja em execução para os testes de integração. Inicie-o em modo *detached* (`-d`):

    ```bash
    docker-compose up -d mongodb
    ```

2.  **Rode os testes**: Com o banco de dados no ar, execute o comando de teste do Go a partir da raiz do projeto:

    ```bash
    go test -v ./...
    ```

    Este comando percorrerá todos os pacotes e executará seus respectivos testes. Você verá a saída do teste `TestAuctionRepository_CloseExpiredAuctions`, confirmando que a automação de fechamento de leilão funciona como esperado.

3.  **(Opcional) Desligue os contêineres**: Após a conclusão dos testes, você pode parar e remover os contêineres:

    ```bash
    docker-compose down
    ```