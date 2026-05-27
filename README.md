# 🦫 Golang Labs

Repositório dedicado ao estudo e desenvolvimento de aplicações utilizando a linguagem Go (Golang). O objetivo principal é consolidar os conceitos da linguagem, gestão de dependências com Go Modules e design de software focado em qualidade.

## 🚀 O que estou a estudar de momento

* **Metodologia:** Test-Driven Development (TDD) — Escrever testes antes da implementação do código de produção.
* **Biblioteca standard:** `net/http`, `encoding/json`, `net/http/httptest`, `testing`.
* **Bibliotecas externas:** `go-chi/chi/v5` para routing e middlewares HTTP.
* **Padrões:** Injeção de dependências, composição de middlewares, testes de integração de handlers.
* **Recursos:** *Learn Go With Tests* (por @quii).
* **Ferramentas:** Go CLI, `gopls` (Language Server).

## 📂 Estrutura do Repositório

* `/hello-world` — Fundamentos de sintaxe, variáveis curtas (`:=`), constantes, estruturas condicionais (`if`/`switch`), subtests e funções públicas/privadas.
* `/desafio1` — Servidor HTTP com `go-chi/chi/v5`: handlers, parâmetros de rota com regex, query parameters, encoding/decoding JSON, middlewares de logging e autenticação por header.

---

> Desenvolvido com foco em aprender a criar código limpo, rápido e testado em ambiente Go.