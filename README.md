# CPF/CNPJ Validator

Aplicação full-stack para cadastrar, validar e gerenciar CPFs e CNPJs, com filtros, ordenação e blocklist. Backend em Go (Gin + PostgreSQL), frontend em React (SPA) e tudo containerizado com Docker.

![Go](https://img.shields.io/badge/Go-1.26-00ADD8?style=flat&logo=go&logoColor=white)
![React](https://img.shields.io/badge/React-Vite-61DAFB?style=flat&logo=react&logoColor=black)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-4169E1?style=flat&logo=postgresql&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat&logo=docker&logoColor=white)
![Status](https://img.shields.io/badge/status-conclu%C3%ADdo-brightgreen)


## 🇧🇷 Português

### 📖 Sobre o projeto

Projeto baseado em um desafio técnico real de processo seletivo. O sistema valida o dígito verificador de CPFs e CNPJs, persiste os documentos no PostgreSQL e oferece uma interface web para gerenciá-los.

### ✅ Requisitos do desafio

| Requisito | Como foi atendido |
|---|---|
| Validação do dígito verificador na consulta e inclusão | Pacote `validator`, implementado sem bibliotecas externas e coberto por testes unitários |
| SPA com CRUD, filtros, ordenação e blocklist | Frontend em React (Vite) |
| Interface REST para integração com o front-end | API REST em Go com Gin |
| Rota de suporte com uptime e quantidade de consultas | `GET /status` |
| Dados em MongoDB ou PostgreSQL | PostgreSQL 16 |
| Conteinerização (desejável) | Docker e Docker Compose para API, frontend e banco |
| Testes (desejável) | Testes unitários do validador e testes de handlers com `httptest` e `sqlmock` |
| README com instruções de instalação e execução | Este documento |

### ✨ Funcionalidades

- ➕ Cadastro de CPF/CNPJ com validação do dígito verificador (aceita com ou sem pontuação)
- 🔎 Tipo (CPF ou CNPJ) detectado automaticamente
- 📋 Listagem com filtro por tipo e status, e ordenação por campo e direção
- 🚫 Marcação e desmarcação de documentos na blocklist
- ✏️ Atualização e exclusão de documentos
- 📊 Barra de status na interface com uptime e contagem de consultas

### 🚀 Como executar

#### Opção 1 — Com Docker (recomendado)

**Pré-requisitos:** [Docker](https://docs.docker.com/get-docker/) e Docker Compose.

```bash
git clone https://github.com/fshenrique6/cpf-cnpj-validator.git
cd cpf-cnpj-validator
docker compose up --build
```

| Serviço | Endereço |
|---|---|
| Frontend | http://localhost:3000 |
| API | http://localhost:8080 |
| PostgreSQL | `localhost:5433` (usuário `cpfcnpjapi`, banco `cpfcnpjdb`) |

A tabela é criada automaticamente a partir do `init.sql`. O script só roda na primeira criação do volume; para recriar o banco do zero:

```bash
docker compose down -v
docker compose up --build
```

> A API usa a porta `8080` e o frontend a `3000`. Se alguma estiver ocupada na sua máquina, pare o processo que a usa antes de subir.

#### Opção 2 — Ambiente local (sem Docker)

**Pré-requisitos:** Go 1.26+, Node.js 22+ e PostgreSQL 14+.

**1. Banco de dados**

```sql
CREATE USER cpfcnpjapi WITH PASSWORD 'cpfcnpjapi123';
CREATE DATABASE cpfcnpjdb OWNER cpfcnpjapi;
```

Depois, conectado ao `cpfcnpjdb`, execute o conteúdo de `init.sql`.

**2. Backend**

```bash
cd backend
cp .env.example .env
go run main.go
```

A API sobe em `http://localhost:8080`.

**3. Frontend**

```bash
cd frontend
npm install
npm run dev
```

A interface sobe em `http://localhost:5173`.

### ⚙️ Variáveis de ambiente

**Backend**

| Variável | Descrição | Padrão |
|---|---|---|
| `DB_HOST` | Host do PostgreSQL | — |
| `DB_PORT` | Porta do PostgreSQL | — |
| `DB_USER` | Usuário do banco | — |
| `DB_PASSWORD` | Senha do banco | — |
| `DB_NAME` | Nome do banco | — |
| `CORS_ALLOWED_ORIGINS` | Origens permitidas, separadas por vírgula | `http://localhost:5173` |

**Frontend**

| Variável | Descrição | Padrão |
|---|---|---|
| `VITE_API_URL` | URL base da API | `http://localhost:8080` |

No Docker, `VITE_API_URL` é gravada no momento do build (o Vite a embute no JavaScript gerado), então alterá-la exige `docker compose up --build`.

> As credenciais deste README são apenas para desenvolvimento local e não devem ser usadas em produção.

### 📚 Endpoints da API

| Método | Rota | Descrição |
|---|---|---|
| `POST` | `/documents` | Cadastra um documento |
| `GET` | `/documents` | Lista documentos (aceita filtros e ordenação) |
| `GET` | `/documents/:id` | Busca um documento por ID |
| `PUT` | `/documents/:id` | Atualiza o número do documento |
| `PATCH` | `/documents/:id/blocklist` | Marca ou desmarca na blocklist |
| `DELETE` | `/documents/:id` | Remove um documento |
| `GET` | `/status` | Uptime do servidor e quantidade de consultas |

#### Cadastrar

```bash
curl -X POST http://localhost:8080/documents \
  -H "Content-Type: application/json" \
  -d '{"number":"529.982.247-25"}'
```

```json
{
  "id": 1,
  "number": "52998224725",
  "type": "cpf",
  "blocklisted": false,
  "created_at": "2026-10-16T19:00:00Z",
  "updated_at": "2026-10-16T19:00:00Z"
}
```

O número é armazenado apenas com dígitos. Erros comuns:

| Status | Quando |
|---|---|
| `400` | Documento inválido (tamanho ou dígito verificador) ou ID não numérico |
| `404` | Documento não encontrado |
| `409` | Documento já cadastrado |

#### Listar com filtros

```bash
curl "http://localhost:8080/documents?type=cpf&blocklisted=false&sort_by=number&order=asc"
```

| Parâmetro | Valores | Padrão |
|---|---|---|
| `type` | `cpf`, `cnpj` | todos |
| `blocklisted` | `true`, `false` | todos |
| `sort_by` | `id`, `number`, `type`, `created_at`, `updated_at` | `created_at` |
| `order` | `asc`, `desc` | `desc` |

Campos de ordenação fora dessa lista retornam `400` (allowlist contra SQL injection). Com a tabela vazia, a resposta é `[]`.

#### Blocklist

```bash
curl -X PATCH http://localhost:8080/documents/1/blocklist \
  -H "Content-Type: application/json" \
  -d '{"blocklisted":true}'
```

#### Status

```bash
curl http://localhost:8080/status
```

```json
{
  "request_count": 12,
  "uptime": "1h45m42s"
}
```

O contador considera apenas as consultas às rotas de `/documents` desde o início do servidor e é seguro para acesso concorrente (`atomic.Int64`). Em modo de desenvolvimento, o React StrictMode executa efeitos duas vezes, então cada carregamento de página pode contar duas consultas. Em produção, conta uma.

### 🧪 Testes

```bash
cd backend
go test ./...
```

- **Validador:** tabelas de casos para CPF e CNPJ (válido, com e sem máscara, dígitos repetidos, tamanho incorreto, dígito verificador incorreto, vazio).
- **Handlers:** testes HTTP com `httptest`, usando `go-sqlmock` no lugar do banco. Cobrem sucesso, ID inválido, não encontrado, duplicidade (`23505` → `409`), falha genérica de banco (→ `500`), filtros, ordenação e a allowlist.

### 🏗️ Estrutura do projeto

```
cpf-cnpj-validator/
├── docker-compose.yml        # API + frontend + PostgreSQL
├── init.sql                  # Criação da tabela documents
├── backend/
│   ├── Dockerfile            # Build multi-stage
│   ├── main.go               # Rotas, CORS e middleware
│   ├── database/             # Conexão com o PostgreSQL
│   ├── models/               # Struct Document
│   ├── validator/            # Validação de CPF/CNPJ e testes
│   ├── handlers/             # Endpoints REST e testes
│   └── status/               # Middleware contador e rota /status
└── frontend/
    ├── Dockerfile            # Build do Vite servido por nginx
    ├── nginx.conf
    └── src/
        ├── api/              # Cliente HTTP
        ├── components/       # Formulário, tabela, filtros, barra de status
        └── utils/            # Formatação de CPF/CNPJ
```

### 🧠 Decisões técnicas

- **Validação sem biblioteca externa:** o algoritmo do dígito verificador (módulo 11 com pesos) e a rejeição de sequências repetidas foram implementados do zero, com testes.
- **Número normalizado:** o documento é guardado só com dígitos, com `UNIQUE` no banco, o que evita duplicidade por formatação (`529.982.247-25` e `52998224725` são o mesmo registro).
- **Tipo derivado:** o cliente não informa se é CPF ou CNPJ; a API deduz pelo tamanho e valida.
- **Queries parametrizadas e allowlist:** valores vão em placeholders (`$1`, `$2`) e campos de ordenação são validados contra uma lista fixa, pois nomes de coluna não podem ser parametrizados.
- **Erros do banco tratados por código:** apenas a violação de unicidade (`23505`) vira `409`; qualquer outra falha vira `500` com mensagem genérica, sem expor detalhes internos.
- **Healthcheck no Compose:** a API só sobe depois que o PostgreSQL aceita conexões, evitando falha de inicialização.
- **Rota de suporte com contador seguro:** `atomic.Int64` evita condição de corrida entre requisições concorrentes.

### 🔮 Próximos passos

- [ ] Paginação na listagem
- [ ] Testes do frontend (Vitest e Testing Library)
- [ ] Pipeline de CI (GitHub Actions) com `golangci-lint`, `go test` e `npm run lint`
- [ ] Máscara automática no campo de cadastro

### 👤 Autor

**Henrique Souza**
Projeto desenvolvido para fins de aprendizado e portfólio.

- GitHub: [@fshenrique6](https://github.com/fshenrique6)

---

## 🇺🇸 English

### 📖 About the project

Project based on a real technical challenge from a hiring process. The system validates the check digits of CPF and CNPJ numbers (Brazilian individual and company tax IDs), stores them in PostgreSQL, and provides a web interface to manage them.

### ✅ Challenge requirements

| Requirement | How it was met |
|---|---|
| Check-digit validation on lookup and insert | `validator` package, implemented without external libraries and covered by unit tests |
| SPA with CRUD, filters, sorting and blocklist | React (Vite) frontend |
| REST interface for front-end integration | REST API in Go with Gin |
| Support route with uptime and number of queries | `GET /status` |
| Data stored in MongoDB or PostgreSQL | PostgreSQL 16 |
| Containerization (nice to have) | Docker and Docker Compose for API, frontend and database |
| Tests (nice to have) | Validator unit tests and handler tests with `httptest` and `sqlmock` |
| README with install and run instructions | This document |

### ✨ Features

- ➕ CPF/CNPJ registration with check-digit validation (accepts formatted or plain digits)
- 🔎 Type (CPF or CNPJ) detected automatically
- 📋 Listing with filters by type and status, and sorting by field and direction
- 🚫 Add and remove documents from the blocklist
- ✏️ Update and delete documents
- 📊 Status bar in the UI with uptime and query count

### 🚀 Getting started

#### Option 1 — With Docker (recommended)

**Prerequisites:** [Docker](https://docs.docker.com/get-docker/) and Docker Compose.

```bash
git clone https://github.com/fshenrique6/cpf-cnpj-validator.git
cd cpf-cnpj-validator
docker compose up --build
```

| Service | Address |
|---|---|
| Frontend | http://localhost:3000 |
| API | http://localhost:8080 |
| PostgreSQL | `localhost:5433` (user `cpfcnpjapi`, database `cpfcnpjdb`) |

The table is created automatically from `init.sql`. The script only runs when the volume is first created; to recreate the database from scratch:

```bash
docker compose down -v
docker compose up --build
```

> The API uses port `8080` and the frontend uses `3000`. If either is taken on your machine, stop the process using it before starting.

#### Option 2 — Local environment (without Docker)

**Prerequisites:** Go 1.26+, Node.js 22+ and PostgreSQL 14+.

**1. Database**

```sql
CREATE USER cpfcnpjapi WITH PASSWORD 'cpfcnpjapi123';
CREATE DATABASE cpfcnpjdb OWNER cpfcnpjapi;
```

Then, connected to `cpfcnpjdb`, run the contents of `init.sql`.

**2. Backend**

```bash
cd backend
cp .env.example .env
go run main.go
```

The API starts at `http://localhost:8080`.

**3. Frontend**

```bash
cd frontend
npm install
npm run dev
```

The UI starts at `http://localhost:5173`.

### ⚙️ Environment variables

**Backend**

| Variable | Description | Default |
|---|---|---|
| `DB_HOST` | PostgreSQL host | — |
| `DB_PORT` | PostgreSQL port | — |
| `DB_USER` | Database user | — |
| `DB_PASSWORD` | Database password | — |
| `DB_NAME` | Database name | — |
| `CORS_ALLOWED_ORIGINS` | Comma-separated list of allowed origins | `http://localhost:5173` |

**Frontend**

| Variable | Description | Default |
|---|---|---|
| `VITE_API_URL` | API base URL | `http://localhost:8080` |

In Docker, `VITE_API_URL` is baked in at build time (Vite embeds it in the generated JavaScript), so changing it requires `docker compose up --build`.

> The credentials in this README are for local development only and must not be used in production.

### 📚 API endpoints

| Method | Route | Description |
|---|---|---|
| `POST` | `/documents` | Create a document |
| `GET` | `/documents` | List documents (supports filters and sorting) |
| `GET` | `/documents/:id` | Get a document by ID |
| `PUT` | `/documents/:id` | Update the document number |
| `PATCH` | `/documents/:id/blocklist` | Add to or remove from the blocklist |
| `DELETE` | `/documents/:id` | Delete a document |
| `GET` | `/status` | Server uptime and number of queries |

#### Create

```bash
curl -X POST http://localhost:8080/documents \
  -H "Content-Type: application/json" \
  -d '{"number":"529.982.247-25"}'
```

```json
{
  "id": 1,
  "number": "52998224725",
  "type": "cpf",
  "blocklisted": false,
  "created_at": "2026-10-16T19:00:00Z",
  "updated_at": "2026-10-16T19:00:00Z"
}
```

The number is stored as digits only. Common errors:

| Status | When |
|---|---|
| `400` | Invalid document (length or check digit) or non-numeric ID |
| `404` | Document not found |
| `409` | Document already registered |

#### List with filters

```bash
curl "http://localhost:8080/documents?type=cpf&blocklisted=false&sort_by=number&order=asc"
```

| Parameter | Values | Default |
|---|---|---|
| `type` | `cpf`, `cnpj` | all |
| `blocklisted` | `true`, `false` | all |
| `sort_by` | `id`, `number`, `type`, `created_at`, `updated_at` | `created_at` |
| `order` | `asc`, `desc` | `desc` |

Sort fields outside this list return `400` (allowlist against SQL injection). With an empty table, the response is `[]`.

#### Blocklist

```bash
curl -X PATCH http://localhost:8080/documents/1/blocklist \
  -H "Content-Type: application/json" \
  -d '{"blocklisted":true}'
```

#### Status

```bash
curl http://localhost:8080/status
```

```json
{
  "request_count": 12,
  "uptime": "1h45m42s"
}
```

The counter only counts queries to `/documents` routes since the server started, and is safe for concurrent access (`atomic.Int64`). In development mode, React StrictMode runs effects twice, so each page load may count two queries. In production, it counts one.

### 🧪 Tests

```bash
cd backend
go test ./...
```

- **Validator:** table-driven cases for CPF and CNPJ (valid, with and without mask, repeated digits, wrong length, wrong check digit, empty).
- **Handlers:** HTTP tests with `httptest`, using `go-sqlmock` in place of the database. They cover success, invalid ID, not found, duplicates (`23505` → `409`), generic database failure (→ `500`), filters, sorting and the allowlist.

### 🏗️ Project structure

```
cpf-cnpj-validator/
├── docker-compose.yml        # API + frontend + PostgreSQL
├── init.sql                  # documents table creation
├── backend/
│   ├── Dockerfile            # Multi-stage build
│   ├── main.go               # Routes, CORS and middleware
│   ├── database/             # PostgreSQL connection
│   ├── models/               # Document struct
│   ├── validator/            # CPF/CNPJ validation and tests
│   ├── handlers/             # REST endpoints and tests
│   └── status/               # Counter middleware and /status route
└── frontend/
    ├── Dockerfile            # Vite build served by nginx
    ├── nginx.conf
    └── src/
        ├── api/              # HTTP client
        ├── components/       # Form, table, filters, status bar
        └── utils/            # CPF/CNPJ formatting
```

### 🧠 Technical decisions

- **Validation without external libraries:** the check-digit algorithm (modulo 11 with weights) and the rejection of repeated-digit sequences were implemented from scratch, with tests.
- **Normalized number:** the document is stored as digits only, with a `UNIQUE` constraint, which prevents duplicates caused by formatting (`529.982.247-25` and `52998224725` are the same record).
- **Derived type:** the client does not say whether it is a CPF or CNPJ; the API infers it from the length and validates it.
- **Parameterized queries and allowlist:** values go in placeholders (`$1`, `$2`), and sort fields are checked against a fixed list, since column names cannot be parameterized.
- **Database errors handled by code:** only a unique violation (`23505`) becomes `409`; any other failure becomes `500` with a generic message, without exposing internal details.
- **Compose healthcheck:** the API only starts once PostgreSQL accepts connections, avoiding startup failures.
- **Safe support-route counter:** `atomic.Int64` avoids race conditions between concurrent requests.

### 🔮 Roadmap

- [ ] Pagination on the list
- [ ] Frontend tests (Vitest and Testing Library)
- [ ] CI pipeline (GitHub Actions) with `golangci-lint`, `go test` and `npm run lint`
- [ ] Automatic input mask on the registration field

### 👤 Author

**Henrique Souza**
Project built for learning and portfolio purposes.

- GitHub: [@fshenrique6](https://github.com/fshenrique6)
