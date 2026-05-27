# 🛡️ Desafio 1

## Visão Geral

Este documento descreve os requisitos técnicos e a arquitectura de implementação para o **Desafio 1**, um exercício de engenharia orientado para a construção de um servidor HTTP de alto desempenho em Go, utilizando a biblioteca `go-chi/chi/v5`.

O objectivo central é demonstrar precisão na definição de rotas, validação de dados, gestão de middlewares e correcta utilização dos métodos HTTP — competências fundamentais para qualquer sistema de entrada (*entry point*) em produção.

---

## 🧰 Stack Tecnológica

| Componente | Versão |
|---|---|
| Linguagem | Go 1.22+ |
| Router | `github.com/go-chi/chi/v5` |
| Porta | `:8080` |

---

## 🏗️ Arquitectura e Estrutura de Rotas

## 1. Navegação Básica — *Warm-up*

### `GET /`

Retorna uma resposta de verificação de disponibilidade do sistema.

- **Código de estado:** `200 OK`
- **Content-Type:** `text/plain`

**Corpo da Resposta:**
  `Olá Mundo!`

---

### `GET /health`

Retorna o estado operacional do servidor em formato JSON, incluindo o tempo de actividade (*uptime*) desde o arranque.

- **Código de estado:** `200 OK`
- **Content-Type:** `application/json`

**Exemplo de resposta:**
```json
{
  "status": "up",
  "uptime": "3h14m52s"
}
```

---

## 2. Correspondência Dinâmica de Padrões

### `GET /hello/{name}`

Retorna uma saudação personalizada com base no parâmetro de rota fornecido.

**Exemplo:**
- Pedido: `GET /hello/Pedro`
- Resposta: `Olá, Pedro!`

---

### `GET /user/id`

Esta rota utiliza uma **expressão regular** para garantir que apenas identificadores numéricos são aceites. Pedidos com valores não numéricos não devem ser correspondidos por esta rota.

| Pedido | Comportamento |
|---|---|
| `GET /user/123` | `200 OK` → `User Profile: 123` |
| `GET /user/admin` | `404 Not Found` (sem correspondência de rota) |

---

### `GET /search`

Processa parâmetros de consulta (*query parameters*) a partir do URL.

**Parâmetros suportados:**

| Parâmetro | Obrigatório | Valor por Omissão |
|---|---|---|
| `q` | Sim | — |
| `page` | Não | `1` |

**Exemplos:**

| URL | Resposta |
|---|---|
| `/search?q=golang&page=2` | `Searching for 'golang' on page 2` |
| `/search?q=golang` | `Searching for 'golang' on page 1` |

---

## 3. Troca de Dados — JSON e Métodos HTTP

### `POST /echo`

Recebe um corpo JSON, valida que não está vazio e retorna o mesmo *payload* enriquecido com um campo de timestamp de processamento.

**Corpo do pedido:**
```json
{
  "payload": "dados_de_exemplo"
}
```

**Corpo da resposta:**
```json
{
  "payload": "dados_de_exemplo",
  "processed_at": "2025-01-15T10:30:00Z"
}
```

**Regras de validação:**
- Caso o corpo esteja ausente ou vazio → `400 Bad Request`
- Caso o JSON seja inválido → `400 Bad Request`

> ⚠️ Esta rota está **protegida pelo middleware de autenticação** (ver Secção 4).

---

## 4. Middlewares

### 4.1 Middleware de Registo (*Logging*)

**Âmbito de aplicação:** Global — todos os pedidos ao servidor.

Cada pedido deve gerar uma entrada no terminal no seguinte formato:

```
[METHOD] - /path - Latency
```

**Exemplo:**
```
[GET] - /health - 1.204ms
[POST] - /echo - 3.812ms
```

> 💡 A latência deve ser calculada utilizando `time.Since()`.

---

### 4.2 Middleware de Autenticação por Cabeçalho (*Token Auth*)

**Âmbito de aplicação:** Restrito — aplicado apenas à rotas `/echo`.

Este middleware inspecciona o cabeçalho HTTP `X-App-Token` em cada pedido recebido.

| Condição | Comportamento |
|---|---|
| `X-App-Token: secret123` | Pedido autorizado — prossegue normalmente |
| Cabeçalho ausente | `401 Unauthorized` |
| Valor incorrecto | `401 Unauthorized` |

**Exemplo de resposta de erro:**
```json
{
  "error": "Unauthorized: token inválido ou ausente"
}
```

> 💡 **Boas práticas:** Em ambiente de produção, o token deve ser carregado a partir de variáveis de ambiente (e.g., `os.Getenv("APP_TOKEN")`), e nunca codificado directamente no código-fonte (*hardcoded*).

---

## 🚀 Como Executar

```bash
# Instalar dependências
go mod tidy

# Iniciar o servidor
go run main.go

# O servidor ficará disponível em:
# http://localhost:8080
```

---

## 🧪 Exemplos de Teste com cURL

```bash
# Verificação do sistema
curl http://localhost:8080/

# Estado de saúde
curl http://localhost:8080/health

# Saudação personalizada
curl http://localhost:8080/hello/Pedro

# Perfil de utilizador (numérico)
curl http://localhost:8080/user/42

# Pesquisa com parâmetros
curl "http://localhost:8080/search?q=golang&page=3"

# Echo com autenticação (pedido válido)
curl -X POST http://localhost:8080/echo \
  -H "Content-Type: application/json" \
  -H "X-App-Token: secret123" \
  -d '{"payload": "teste"}'

# Echo sem token (deve retornar 401)
curl -X POST http://localhost:8080/echo \
  -H "Content-Type: application/json" \
  -d '{"payload": "teste"}'

```

---

## 📋 Sumário de Rotas

| Método | Rota | Middleware Auth | Descrição |
|---|---|---|---|
| `GET` | `/` | ✗ | Verificação de disponibilidade |
| `GET` | `/health` | ✗ | Estado e uptime do servidor |
| `GET` | `/hello/{name}` | ✗ | Saudação personalizada |
| `GET` | `/user/{id:[0-9]+}` | ✗ | Perfil de utilizador (ID numérico) |
| `GET` | `/search` | ✗ | Pesquisa com query parameters |
| `POST` | `/echo` | ✅ | Eco de payload JSON com timestamp |
