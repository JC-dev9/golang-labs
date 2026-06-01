# 🧠 Desafio 2

## Visão Geral

Este documento descreve os requisitos técnicos e a arquitectura de implementação para o **Desafio 2**, uma continução do exercício anterior com adição de uma nova layer **lógica de negócio**.

O objectivo é construir um **Serviço de Login em Memória** — mas, mais importante do que o que se constrói, é **como se constrói**.

---

## 🧰 Stack Tecnológica

| Componente | Versão |
|---|---|
| Linguagem | Go 1.22+ |
| Router | `github.com/go-chi/chi/v5` |
| Hashing | `golang.org/x/crypto/bcrypt` *(opcional, mas recomendado)* |
| Porta | `:8080` |

---

## 🏗️ Uma Nova Forma de Pensar — Separação de Responsabilidades

No Desafio 1, toda a lógica residia directamente nos *handlers*. Era uma abordagem adequada para o contexto — roteamento, parsing de parâmetros, respostas simples. Não havia lógica de negócio real para separar.

Este desafio é diferente. A partir daqui, existem **decisões de negócio** a tomar: um *username* é válido? Uma *password* está correcta? Uma sessão existe? Essas decisões **não pertencem ao handler**.

A arquitectura que se pede assenta numa separação clara e deliberada entre duas camadas com responsabilidades distintas e intransferíveis:

**Os *handlers* são responsáveis exclusivamente pela camada HTTP.** O seu trabalho começa e termina no protocolo: receber o pedido, fazer o *parse* do corpo, invocar o serviço, e traduzir o resultado numa resposta HTTP. Um *handler* não sabe o que significa "um utilizador já existe" — sabe apenas que, quando o serviço comunica esse facto, deve responder com um determinado código de estado. Não contém condições de negócio. Não valida regras de domínio. É deliberadamente burro.

**O `AuthService` é o único detentor da lógica de negócio.** É ele que decide se um *username* cumpre os requisitos mínimos, se as credenciais são válidas, se um token de sessão deve ser emitido. Crucialmente, o serviço é **completamente agnóstico à camada de transporte** — não conhece `http.Request`, não escreve `http.ResponseWriter`, não sabe sequer que existe um servidor HTTP. Recebe dados primitivos, aplica regras, e retorna resultados ou erros de domínio.

Esta divisão não é burocracia arquitectural. É o que torna cada camada **testável de forma independente**, e o que garante que, no futuro, a lógica de negócio pode ser reutilizada noutros contextos sem qualquer alteração — seja uma CLI, um *worker* assíncrono, ou um endpoint gRPC.

> O *handler* não decide — executa. O serviço não comunica — decide.

---

## 1. O `AuthService` — A Lógica de Negócio

O `AuthService` é a **fonte de verdade** (*source of truth*) da aplicação. Mantém em memória os utilizadores registados e as sessões activas.

### Método `Register(username, password string) (User, error)`

**Regras de negócio:**

| Condição | Comportamento |
|---|---|
| `username` com menos de 4 caracteres | Retorna erro de validação |
| `username` já existente | Retorna `ErrUserExists` |
| Registo bem-sucedido | Retorna o utilizador criado |

---

### Método `Login(username, password string) (string, error)`

**Regras de negócio:**

| Condição | Comportamento |
|---|---|
| Utilizador não existe | Retorna `ErrUnauthorized` |
| Password incorrecta | Retorna `ErrUnauthorized` |
| Autenticação bem-sucedida | Retorna um token de sessão |


---

### Erros de Domínio

O serviço deve expor erros semânticos próprios. São estes erros que permitem aos *handlers* mapear cada situação de falha para o código de estado HTTP correcto — sem que o *handler* precise de conter qualquer conhecimento sobre as regras que os originaram.

---

## 2. Configuração Inicial — *Bootstrap*

No arranque da aplicação, deve ser criado um utilizador administrador para fins de desenvolvimento e teste:

| Campo | Valor |
|---|---|
| Username | `admin` |
| Password | `password123` |


---

## 3. Os Endpoints

### `POST /api/user/register`

Regista um novo utilizador no sistema.

**Corpo do pedido:**
```json
{
  "username": "intern1",
  "password": "safe-password"
}
```

**Resposta de sucesso (`201 Created`):**
```json
{
  "id": "a1b2c3d4",
  "username": "intern1",
  "created_at": "2025-01-15T10:30:00Z"
}
```

**Respostas de erro:**

| Cenário | Código |
|---|---|
| `username` com menos de 4 caracteres | `400 Bad Request` |
| `username` já registado | `409 Conflict` |
| Corpo JSON inválido ou ausente | `400 Bad Request` |

---

### `POST /api/user/login`

Autentica um utilizador existente e emite um token de sessão.

**Corpo do pedido:**
```json
{
  "username": "admin",
  "password": "password123"
}
```

**Resposta de sucesso (`200 OK`):**
```json
{
  "token": "f7a3c9e1b2d4..."
}
```

**Respostas de erro:**

| Cenário | Código |
|---|---|
| Credenciais inválidas | `401 Unauthorized` |
| Corpo JSON inválido ou ausente | `400 Bad Request` |

---

### `GET /api/user/profile`

Retorna o perfil do utilizador autenticado com base no token de sessão activo.

**Cabeçalho obrigatório:**
```
X-Session-Token: f7a3c9e1b2d4...
```

**Resposta de sucesso (`200 OK`):**
```json
{
  "id": "a1b2c3d4",
  "username": "admin",
  "created_at": "2025-01-15T09:00:00Z"
}
```

**Respostas de erro:**

| Cenário | Código |
|---|---|
| Cabeçalho `X-Session-Token` ausente | `401 Unauthorized` |
| Token inválido ou inexistente | `401 Unauthorized` |

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
# Registar um novo utilizador
curl -X POST http://localhost:8080/api/user/register \
  -H "Content-Type: application/json" \
  -d '{"username": "intern1", "password": "safe-password"}'

# Tentar registar username duplicado (deve retornar 409)
curl -X POST http://localhost:8080/api/user/register \
  -H "Content-Type: application/json" \
  -d '{"username": "intern1", "password": "outra-password"}'

# Tentar registar username demasiado curto (deve retornar 400)
curl -X POST http://localhost:8080/api/user/register \
  -H "Content-Type: application/json" \
  -d '{"username": "ab", "password": "safe-password"}'

# Autenticar com o utilizador admin
curl -X POST http://localhost:8080/api/user/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "password123"}'

# Autenticar com credenciais inválidas (deve retornar 401)
curl -X POST http://localhost:8080/api/user/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "errada"}'

# Obter perfil do utilizador autenticado (substituir TOKEN pelo valor recebido no login)
curl http://localhost:8080/api/user/profile \
  -H "X-Session-Token: TOKEN_AQUI"

# Tentar aceder a /me sem token (deve retornar 401)
curl http://localhost:8080/api/user/profile
```

---

## 📋 Sumário de Endpoints

| Método | Rota | Autenticação | Descrição |
|---|---|---|---|
| `POST` | `/api/user/register` | ✗ | Registo de novo utilizador |
| `POST` | `/api/user/login` | ✗ | Autenticação e emissão de token |
| `GET` | `/api/user/profile` | ✅ `X-Session-Token` | Perfil do utilizador autenticado |

