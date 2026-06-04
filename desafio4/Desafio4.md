# 🖥️ Desafio 4

## Visão Geral

Este documento descreve os requisitos técnicos e a arquitectura de implementação para o **Desafio 4: A Interface**, um exercício orientado para a transição de uma API JSON para uma aplicação web com **renderização de HTML no servidor**.

A lógica de negócio construída nos desafios anteriores mantém-se intacta. O que muda é a camada de apresentação: em vez de serializar dados em JSON e devolvê-los a um cliente externo, o servidor passa a ser responsável por produzir HTML completo, pronto a ser interpretado directamente pelo browser.

---

## 🧰 Stack Tecnológica

| Componente | Versão |
|---|---|
| Linguagem | Go 1.22+ |
| Router | `github.com/go-chi/chi/v5` |
| Templating | `html/template` *(stdlib)* |
| Base de Dados | SQLite3 |
| Hashing | `golang.org/x/crypto/bcrypt` |
| Porta | `:8080` |

---

## 🏗️ Uma Nova Responsabilidade — O Handler como Apresentador

Nos desafios anteriores, os *handlers* terminavam o seu trabalho ao serializar uma struct em JSON e escrever o resultado no corpo da resposta. Era simples e uniforme — todos os *handlers* falavam a mesma língua.

Neste desafio, os *handlers* assumem uma responsabilidade adicional: **a apresentação**. Em vez de produzir dados, produzem documentos HTML completos, construídos a partir de *templates* que combinam estrutura estática com dados dinâmicos injectados em tempo de execução.

A separação de camadas mantém-se rigorosamente igual:

```
Handler  →  Service  →  Repository  →  Base de Dados
   ↓
Template
```

O *handler* continua a não ter opinião sobre regras de negócio. Continua a invocar o serviço, a receber dados ou erros, e a produzir uma resposta. A diferença está no formato dessa resposta: em vez de `application/json`, o `Content-Type` é agora `text/html`, e o corpo é gerado pela execução de um *template*.

Os *templates* não contêm lógica de negócio. Não tomam decisões. Recebem um conjunto de dados — uma struct, um mapa, um valor simples — e renderizam HTML com base nesses dados. Qualquer lógica condicional presente num *template* deve ser estritamente de apresentação: mostrar ou esconder um elemento, iterar uma lista, formatar um valor.

> O *template* não decide o que mostrar — decide como mostrar o que lhe é dado.

---

## 1. O Sistema de Templates

A biblioteca padrão do Go inclui `html/template`, um engine de templating. Todo o conteúdo dinâmico injectado num template é automaticamente sanitizado, prevenindo ataques de *Cross-Site Scripting* (XSS) sem qualquer esforço adicional.

### Organização

Os ficheiros de *template* devem ser tratados como recursos estáticos da aplicação — não como código Go. Devem estar organizados de forma coerente e ser carregados pelo servidor no arranque.

---

## 2. Gestão de Sessão

Neste desafio não há cliente JavaScript a gerir cabeçalhos programaticamente. O browser envia apenas o que o HTML lhe instrui a enviar — e formulários HTML enviam pedidos `POST` com corpo `application/x-www-form-urlencoded`, não JSON.

A sessão é mantida através do cabeçalho `X-Session-Token`, que deve ser propagado manualmente entre páginas. Após um login bem-sucedido, o token deve estar acessível ao utilizador para que possa ser incluído em pedidos subsequentes.

> Reflectir sobre como um browser comunica estado entre pedidos, e que mecanismos o HTML disponibiliza para isso, é parte integrante deste desafio.

---

## 3. Formulários e Métodos HTTP

Os formulários HTML apenas suportam `GET` e `POST`. Os *handlers* que recebem submissões de formulários devem processar corpos `application/x-www-form-urlencoded` em vez de JSON — o que implica uma diferença na forma como os dados são extraídos do pedido.

Em caso de erro de validação ou autenticação, o servidor deve **re-renderizar o formulário** com uma mensagem de erro contextual visível na página — não redirecionar para uma página de erro genérica.

---

## 4. As Páginas

### `GET /user/register` — Formulário de Registo

Renderiza o formulário de criação de conta.

**Campos do formulário:**
- `username`
- `password`

**Comportamento após submissão (`POST /user/register`):**

| Resultado | Comportamento |
|---|---|
| Registo bem-sucedido | Redireccionamento para `/user/login` |
| Username demasiado curto | Re-renderiza o formulário com mensagem de erro |
| Username já existente | Re-renderiza o formulário com mensagem de erro |

---

### `GET /user/login` — Formulário de Autenticação

Renderiza o formulário de autenticação.

**Campos do formulário:**
- `username`
- `password`

**Comportamento após submissão (`POST /user/login`):**

| Resultado | Comportamento |
|---|---|
| Autenticação bem-sucedida | Redireccionamento para `/user/profile` com o token disponível |
| Credenciais inválidas | Re-renderiza o formulário com mensagem de erro |

---

### `GET /user/profile` — Perfil do Utilizador

Renderiza a página de perfil do utilizador autenticado.

**Cabeçalho obrigatório:**
```
X-Session-Token: f7a3c9e1b2d4...
```

**Conteúdo da página:**
- Nome de utilizador
- Data de criação da conta
- Opção para terminar sessão

**Comportamento:**

| Cenário | Comportamento |
|---|---|
| Token válido | Renderiza a página de perfil |
| Token ausente ou inválido | Redireccionamento para `/user/login` |

---

## 🚀 Como Executar

```bash
# Instalar dependências
go mod tidy

# Iniciar o servidor
go run main.go

# Aceder à aplicação em:
# http://localhost:8080
```

---

## 📋 Sumário de Rotas

| Método | Rota | Autenticação | Descrição |
|---|---|---|---|
| `GET` | `/user/register` | ✗ | Formulário de registo |
| `POST` | `/user/register` | ✗ | Submissão do formulário de registo |
| `GET` | `/user/login` | ✗ | Formulário de autenticação |
| `POST` | `/user/login` | ✗ | Submissão do formulário de autenticação |
| `GET` | `/user/profile` | ✅ `X-Session-Token` | Página de perfil do utilizador |
