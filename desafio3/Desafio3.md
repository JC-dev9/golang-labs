# 🗄️ Desafio 3

## Visão Geral

Este documento descreve os requisitos técnicos e a arquitectura de implementação para o **Desafio 3**

O sistema de autenticação construído no Desafio 2 funciona — mas tem uma limitação crítica: reiniciar o servidor apaga tudo. Utilizadores, sessões, estado. A aplicação não tem memória. Este desafio resolve esse problema, e fá-lo de uma forma que introduz uma nova camada arquitectural com responsabilidades próprias.

---

## 🧰 Stack Tecnológica

| Componente | Versão |
|---|---|
| Linguagem | Go 1.22+ |
| Router | `github.com/go-chi/chi/v5` |
| Base de Dados | SQLite3 |
| Hashing | `golang.org/x/crypto/bcrypt` |
| Porta | `:8080` |

---

## 🏗️ Uma Nova Camada — O Repositório

No Desafio 2, o `AuthService` acumulava duas responsabilidades distintas: aplicava as regras de negócio *e* geria directamente o armazenamento dos dados (o mapa em memória). Para um exercício simples, era aceitável. À medida que a complexidade cresce, essa acumulação torna-se um problema.

Introduzimos agora uma terceira camada: o **Repositório**.

A arquitectura passa a ter três camadas com contratos bem definidos:

```
Handler  →  Service  →  Repository  →  Base de Dados
```

**O Repositório tem uma única responsabilidade:** abstrair o acesso aos dados. Não conhece regras de negócio. Não sabe o que significa "autenticar um utilizador". Sabe apenas como persistir e recuperar entidades da base de dados — nada mais.

**O Service mantém a sua responsabilidade inalterada:** aplicar as regras de negócio. A diferença é que agora, em vez de manipular um mapa directamente, delega todas as operações de leitura e escrita no repositório. O serviço não sabe — nem deve saber — se os dados estão guardados em SQLite, PostgreSQL, ou em memória.

É precisamente aqui que reside o valor da **interface do repositório**. O `AuthService` não depende de uma implementação concreta — depende de um contrato. Qualquer estrutura que satisfaça esse contrato pode ser injectada no serviço, seja uma implementação SQLite para produção, seja uma implementação em memória para testes. O serviço é completamente indiferente.

Esta separação é o que torna o sistema **testável, substituível e escalável** de forma independente em cada camada.

> O *handler* não decide. O serviço não persiste. O repositório não tem opinião sobre regras de negócio.

---

## 1. Os Modelos

Neste desafio, os modelos deixam de ser estruturas informais definidas onde for conveniente — passam a ser **entidades de domínio** com uma definição canónica própria.

### `User`

Representa um utilizador registado no sistema.

| Campo | Tipo | Notas |
|---|---|---|
| `ID` | `string` | Identificador único |
| `Username` | `string` | Nome de utilizador — único na base de dados |
| `Password` | `string` | Nunca deve ser exposto em respostas JSON |
| `CreatedAt` | `time.Time` | Timestamp de criação |

### `Session`

Representa uma sessão activa associada a um utilizador autenticado.

| Campo | Tipo | Notas |
|---|---|---|
| `Token` | `string` | Identificador único da sessão |
| `UserID` | `string` | Referência ao utilizador proprietário |
| `CreatedAt` | `time.Time` | Timestamp de criação |

> No Desafio 2, as sessões eram geridas num mapa interno ao `AuthService`. São agora entidades de primeira classe, persistidas na base de dados como qualquer outra.

---

## 2. A Interface do Repositório

O `AuthService` deve depender de uma **interface**, não de uma implementação concreta. Essa interface define o contrato que qualquer camada de persistência deve satisfazer para ser utilizável pelo serviço.

### `UserRepository`

| Método | Descrição |
|---|---|
| `Create(user User) error` | Persiste um novo utilizador |
| `FindByUsername(username string) (User, error)` | Procura um utilizador pelo nome |
| `FindByID(id string) (User, error)` | Procura um utilizador pelo identificador |

### `SessionRepository`

| Método | Descrição |
|---|---|
| `Create(session Session) error` | Persiste uma nova sessão |
| `FindByToken(token string) (Session, error)` | Procura uma sessão pelo token |
| `DeleteByToken(token string) error` | Remove uma sessão (logout) |

---

## 3. O `AuthService` — Adaptado

O `AuthService` mantém a mesma lógica de negócio do Desafio 2. A única diferença é que recebe as implementações dos repositórios por injecção de dependência — e passa a delegar neles todas as operações sobre dados.

A assinatura dos seus métodos públicos não precisa de mudar. O que muda é o que acontece internamente.

---

## 4. Os Endpoints

Os endpoints do Desafio 2 mantêm-se inalterados. São adicionados dois novos:

### `POST /api/user/register`

Comportamento idêntico ao Desafio 2. Os dados passam agora a ser persistidos em SQLite.

---

### `POST /api/user/login`

Comportamento idêntico ao Desafio 2. A sessão criada passa a ser persistida na base de dados.

---

### `GET /api/user/profile`

Comportamento idêntico ao Desafio 2. O token é agora validado contra a base de dados.

---

### `POST /api/user/logout` *(novo)*

Encerra a sessão activa do utilizador autenticado.

**Cabeçalho obrigatório:**
```
X-Session-Token: f7a3c9e1b2d4...
```

**Resposta de sucesso (`200 OK`):**
```json
{
  "message": "sessão terminada com sucesso"
}
```

**Respostas de erro:**

| Cenário | Código |
|---|---|
| Token ausente | `401 Unauthorized` |
| Token inválido ou inexistente | `401 Unauthorized` |

---

### `GET /api/users` *(novo)*

Retorna a lista de todos os utilizadores registados no sistema.

**Cabeçalho obrigatório:**
```
X-Session-Token: f7a3c9e1b2d4...
```

**Resposta de sucesso (`200 OK`):**
```json
[
  {
    "id": "a1b2c3d4",
    "username": "admin",
    "created_at": "2025-01-15T09:00:00Z"
  },
  {
    "id": "e5f6g7h8",
    "username": "intern1",
    "created_at": "2025-01-15T10:30:00Z"
  }
]
```

**Respostas de erro:**

| Cenário | Código |
|---|---|
| Token ausente ou inválido | `401 Unauthorized` |

> Este endpoint requer autenticação. Qualquer utilizador com sessão activa pode consultá-lo.

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

# Autenticar
curl -X POST http://localhost:8080/api/user/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "password123"}'

# Obter perfil (substituir TOKEN pelo valor recebido no login)
curl http://localhost:8080/api/user/profile \
  -H "X-Session-Token: TOKEN_AQUI"

# Listar todos os utilizadores
curl http://localhost:8080/api/users \
  -H "X-Session-Token: TOKEN_AQUI"

# Terminar sessão
curl -X POST http://localhost:8080/api/user/logout \
  -H "X-Session-Token: TOKEN_AQUI"

# Verificar que o token já não é válido após logout (deve retornar 401)
curl http://localhost:8080/api/user/profile \
  -H "X-Session-Token: TOKEN_AQUI"
```

---

## 📋 Sumário de Endpoints

| Método | Rota | Autenticação | Descrição |
|---|---|---|---|
| `POST` | `/api/user/register` | ✗ | Registo de novo utilizador |
| `POST` | `/api/user/login` | ✗ | Autenticação e emissão de token |
| `GET` | `/api/user/me` | ✅ | Perfil do utilizador autenticado |
| `POST` | `/api/user/logout` | ✅ | Encerramento de sessão |
| `GET` | `/api/users` | ✅ | Listagem de todos os utilizadores |
