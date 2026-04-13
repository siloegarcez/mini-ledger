# Mini-Ledger

API HTTP para gerenciar contas e transacoes financeiras (desafio Pismo), com persistencia em PostgreSQL.


## Stack

- **Linguagem:** Go
- **HTTP:** Chi (com Huma para contrato da API)
- **Banco de dados:** PostgreSQL
- **Arquitetura:** Arquitetura em camadas (Layered Architecture)
- **Testes:** `go test` + Testify + Testcontainers
- **Qualidade:** `golangci-lint` + SQLFluff (SQL linter)
- **Containerização:** Docker + Docker Compose

## Ferramentas e Frameworks

-  **[Huma](https://huma.rocks/)**: Escolhi utilizar Huma para definir a API de forma declarativa, garantindo a geração automática da especificação OpenAPI 3.1. Huma facilita a manutenção da documentação e a evolução da API sem desincronizar a documentação da implementação.

- **[Testify](https://github.com/stretchr/testify)** Uma ótima biblioteca para facilitar a escrita de testes unitários e de integração, com recursos como asserções e mocks.

- **[Zerolog](https://github.com/rs/zerolog)**: Uma biblioteca de logging leve e eficiente com zero alocação, que suporta logs estruturados e é fácil de configurar para diferentes níveis de log além de existir suporte para integração com os principais sistemas de logs (Axiom, Datadog, etc).

- **[Migrate](https://github.com/golang-migrate/migrate)** Utilizado para gerencia migrações do banco de dados. 

- **[go-jet](https://github.com/go-jet/jet)** Uma solução completa para acesso eficiente e de alto desempenho a banco de dados, combinando um construtor de SQL com type safety, geração de código e mapeamento automático dos resultados das consultas.

## Decisões técnicas

### 💰 Representação de Valores Monetários

Esta codebase **não utiliza `float32` ou `float64`** para representar valores monetários. Tipos de ponto flutuante podem introduzir erros de precisão devido à sua representação binária, o que é inadequado para operações financeiras.

Para garantir precisão e consistência:

- Valores monetários são manipulados utilizando **`int64`**, representando a menor unidade da moeda (centavos).
- A API faz o parse de valores monetários como **string**, evitando qualquer conversão implícita para ponto flutuante durante o parsing.
- O parsing é realizado utilizando `json.Number`, garantindo que os valores não sejam convertidos para `float` em nenhum momento do fluxo.

Por simplicidade, o sistema assume a moeda **BRL**, que possui **2 casas decimais fixas**. Esse padrão é validado e aplicado consistentemente em toda a aplicação.

**Exemplo:**

```json
{
  "amount": 123.45
}
```

É internamente representado como:
```
amount = 12345
scale = 2
currency = "BRL"
```

## Pre-requisitos

- Go 1.26+
- Docker e Docker Compose
- Make

## Como rodar com Docker Compose

1. Suba API e banco:

```bash
make compose-up
```

ou

```bash
docker-compose up --build
```

2. API disponivel em:

- `http://localhost:8080`

3. Documentação disponível em:

- `http://localhost:8080/docs`

## Como rodar localmente

### 1) Subir apenas o banco

```bash
make db-dev
```

### 2) Rodar a API

A API é uma CLI que aceita flags para configurar porta e outras opções.
Para ver os comandos e flags disponiveis:

```bash
go run ./cmd/api/main.go --help
```

Sem hot reload:

```bash
go run ./cmd/api/main.go --port 8080
```

Com hot reload:

```bash
make run-watch
```

## Testes e cobertura

Rodar todos os testes do modulo `internal` (unitários e de integração):

```bash
make test
```
Rodar todos os testes do modulo `internal` (unitários apenas):
```bash
make test-unit
```
Gerar cobertura de testes (arquivo `coverage.out` + resumo por funcao):

```bash
make test-coverage
```

## Estrutura do projeto

- `cmd/api/main.go`: ponto de entrada da aplicacao (startup, DB, migracoes, servidor)
- `config/`: configuracao da aplicacao e ambientes (`dev`, `staging`, `prod`)
- `internal/server/`: inicializacao HTTP, middlewares, Huma/Chi e registro de rotas
- `internal/handler/`: camada HTTP (DTOs, handlers e mapeamentos de request/response)
- `internal/service/`: regras de negocio e validacoes de caso de uso
- `internal/domain/`: entidades, value objects e contratos (interfaces)
- `internal/repository/postgres/`: implementacao de repositorios com PostgreSQL/Jet
- `internal/database/`: conexao e utilitarios de banco
- `internal/shared/dberrs/`: traducao de erros SQL para erros de dominio/aplicacao
- `migrations/`: scripts SQL de schema e dados iniciais
- `gen/`: codigo gerado pelo Jet (models/tabelas)
- `scripts/jet/jet.go`: gerador das estruturas Jet


## Rotas da API (com exemplos)

### 1) Criar conta

`POST /accounts`

Request:

```json
{
  "document_number": "12345678900"
}
```

Response `201 Created`:

```json
{
  "account_id": 1,
  "document_number": "12345678900"
}
```

Erros comuns:

- `409 Conflict`: documento ja cadastrado
- `422 Unprocessable Entity`: documento invalido

### 2) Buscar conta por ID

`GET /accounts/{account_id}`

Exemplo:

```bash
curl http://localhost:8080/accounts/1
```

Response `200 OK`:

```json
{
  "account_id": 1,
  "document_number": "12345678900"
}
```

Erros comuns:

- `404 Not Found`: conta nao encontrada

### 3) Criar transacao

`POST /transactions`

Request:

```json
{
  "account_id": 1,
  "operation_type_id": 4,
  "amount": 123.45
}
```

Response `201 Created`:

```json
{
  "transaction_id": 1,
  "account_id": 1,
  "operation_type_id": 4,
  "amount": 123.45
}
```

Erros comuns:

- `400 Bad Request`: `account_id` invalido ou `operation_type_id` inexistente
- `422 Unprocessable Entity`: valor invalido (ex.: formato/valor nao aceito)

## Tipos de operacao disponiveis

As migracoes inserem os seguintes `operation_type_id`:

- `1` = `PURCHASE` (debito)
- `2` = `INSTALLMENT PURCHASE` (debito)
- `3` = `WITHDRAWAL` (debito)
- `4` = `PAYMENT` (credito)
