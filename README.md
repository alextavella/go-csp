# Bank API

Esta é uma API simples de banco construída com **Fiber (Go)**, utilizando o modelo **CSP (Communicating Sequential Processes)** para lidar com operações de depósito e saque de maneira concorrente e segura. A API também inclui um endpoint para consultar o saldo atual.

## Tecnologias Utilizadas

- Go
- Fiber
- Goroutines e Channels (CSP)
- Mutex para controle de concorrência
- Hey para testes de estresse

## Instalação

1. Clone o repositório:

   ```bash
   git clone https://github.com/alextavella/go-csp.git bank-api
   cd bank-api
   ```

2. Instale as dependências:
   ```bash
   go mod tidy
   ```

## Rodando com Docker Compose

1. Suba a aplicação com o Docker Compose:

   ```bash
   docker-compose up --build
   ```

2. Verifique se a API está rodando:

   ```bash
   curl http://localhost:3000/balance
   ```

3. Instale as dependências:
   ```bash
   go mod tidy
   ```

## Rodando a API

Execute a aplicação com o comando:

```bash
go run cmd/api/main.go
```

A API estará disponível em: `http://localhost:3000`

## Endpoints

### Depósito

**POST /deposit**

Realiza um depósito na conta.

**Payload:**

```json
{
  "amount": 100
}
```

**Resposta:**

```json
{
  "status": "success",
  "message": "Deposit successful"
}
```

### Saque

**POST /withdraw**

Realiza um saque da conta.

**Payload:**

```json
{
  "amount": 50
}
```

**Resposta:**

```json
{
  "status": "success",
  "message": "Withdraw attempted"
}
```

### Consultar Saldo

**GET /balance**

Retorna o saldo atual da conta.

**Resposta:**

```json
{
  "user_id": "alextavella",
  "balance": 1000
}
```

## Testes de Estresse

A ferramenta **hey** foi utilizada para realizar testes de estresse na API.

### Instalação do hey

```bash
go install github.com/rakyll/hey@latest
```

### Teste de Depósito

```bash
hey -n 10000 -c 100 -m POST -H "Content-Type: application/json" -d '{"amount": 10}' http://localhost:3000/deposit
```

### Teste de Saque

```bash
hey -n 10000 -c 100 -m POST -H "Content-Type: application/json" -d '{"amount": 5}' http://localhost:3000/withdraw
```

## Licença

Este projeto está licenciado sob a MIT License.
