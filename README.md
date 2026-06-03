# 🦫 Golang Labs

Repositório dedicado ao estudo e desenvolvimento de aplicações utilizando a linguagem Go (Golang). O objetivo principal é consolidar os conceitos da linguagem, gestão de dependências com Go Modules e design de software focado em qualidade.

## 🚀 O que estou a estudar de momento

* **Metodologia:** Test-Driven Development (TDD) — Escrever testes antes da implementação do código de produção. Testes unitários com *mocks* e testes de integração com bases de dados temporárias (`t.TempDir()`).
* **Biblioteca standard:** `net/http`, `encoding/json`, `net/http/httptest`, `testing`, `database/sql`.
* **Bibliotecas externas:** `go-chi/chi/v5` para routing, `modernc.org/sqlite` como driver puro em Go para SQLite.
* **Padrões:** Clean Architecture, Padrão Repository, Injeção de Dependências utilizando Interfaces (*Duck Typing*).
* **Recursos:** *Learn Go With Tests* (por @quii).
* **Ferramentas:** Go CLI, `gopls` (Language Server).
* **Arquitetura:** Separação estrita de responsabilidades: Handlers (HTTP) ➔ Services (Regras de Negócio) ➔ Repositories (Acesso a Dados). 
* **Segurança:** Hashing de passwords com `bcrypt`, geração de IDs e tokens aleatórios, e utilização de *Prepared Statements* contra vulnerabilidades de SQL Injection.

## 📂 Estrutura do Repositório

* `/hello-world` — Fundamentos de sintaxe, variáveis curtas (`:=`), constantes, estruturas condicionais (`if`/`switch`), subtests e funções públicas/privadas.
* `/desafio1` — Servidor HTTP com `go-chi/chi/v5`: handlers, parâmetros de rota com regex, query parameters, encoding/decoding JSON, middlewares de logging e autenticação por header.
* `/desafio2` — Serviço de autenticação em memória com separação de camadas (services e handlers), bcrypt para hashing de passwords, sessões com tokens hexadecimais, e endpoints REST para registo, login e perfil. Construído com TDD desde o serviço até aos handlers.
* `/desafio3` — Evolução da arquitetura introduzindo a camada de **Repository**. Persistência real de dados com **SQLite**. Uso intensivo de **Interfaces** para injetar implementações em memória durante os testes (rapidez) e o motor SQLite no ambiente real. Novos endpoints de Logout e Listagem.

---

> Desenvolvido com foco em aprender a criar código limpo, rápido e testado em ambiente Go.