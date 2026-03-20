# School Management System

Um sistema de gerenciamento escolar desenvolvido em **Go** (Golang), construído para facilitar a interação entre Alunos (`USER`), Professores (`TEACHER`) e Profissionais da saúde/pedagogia (`PROFESSIONAL`).

O sistema oferece recursos robustos para criação de atividades/provas, submissão de respostas, geração de métricas de desempenho (dashboards) e agendamento de consultas com profissionais.

## 🚀 Tecnologias Utilizadas

*   **Linguagem:** [Go](https://golang.org/)
*   **Framework Web:** [Gin](https://github.com/gin-gonic/gin) para roteamento HTTP rápido e flexível.
*   **ORM:** [GORM](https://gorm.io/) para mapeamento objeto-relacional com o banco de dados.
*   **Autenticação:** JWT (JSON Web Tokens) e roteamento protegido por **Cargos (Roles)**.
*   **Banco de Dados:** PostgreSQL (Típico com GORM, mas adaptável via Dialectors).

## 🛠️ Arquitetura do Projeto

O projeto segue princípios de **Clean Architecture** e **Domain-Driven Design (DDD)** simplificado. A base de código está dividida em módulos independentes (`internal/`), garantindo baixo acoplamento:

*   **`model/`**: Definição das estruturas (structs) do domínio e esquemas do banco de dados (GORM tags).
*   **`dto/`** (Data Transfer Objects): Contratos de entrada (Request) e saída (Response) das APIs, validando dados antes de chegarem à regra de negócio.
*   **`repository/`**: Camada de acesso aos dados. Isola as consultas ao banco de dados (GORM) do resto do sistema.
*   **`service/`**: Camada que contém as **regras de negócio**. Orquestra validações e chamadas ao repositório.
*   **`handler/`** (Controllers): Recebe as requisições HTTP do Gin, extrai os parâmetros/body, chama o `service` e retorna as respostas HTTP padronizadas.
*   **`router/`**: Configuração e mapeamento das rotas da aplicação, junto com a injeção de middlewares de autenticação.

## 📦 Principais Módulos e Funcionalidades

### 1. 👥 Autenticação e Usuários (`/internal/auth` e `/internal/user`)
*   **Login e Registro:** Geração de Tokens JWT.
*   **Controle de Acesso (RBAC):** Proteção de rotas baseadas na role do usuário (`USER`, `TEACHER`, `PROFESSIONAL`).
*   **Gestão de Usuários:** CRUD de alunos, professores e listagem de profissionais.

### 2. 📝 Atividades e Provas (`/internal/activity`)
**Para Professores (`TEACHER`):**
*   Criação e edição completa de Provas (Atividades), contendo Exercícios e Alternativas.
*   Ativação e Desativação de testes.
*   **Dashboard da Atividade:** Visualização de notas, alunos que entregaram, média da sala e as questões com maior índice de erro.
*   **Métricas da Sala / Ranking:** Acompanhamento do desempenho geral da turma, top 3 matérias com mais erros e ranking de alunos.

**Para Alunos (`USER`):**
*   Listagem de Provas Ativas e envio de Submissões (respostas).
*   **Dashboard do Aluno:** Acompanhamento pessoal das provas concluídas, nota média, porcentagem de acerto por **matéria da prova** e por **matéria isolada da questão**.

### 3. 📅 Agendamentos (`/internal/appointment` ou relacionado)
*   **Marcação de Consultas:** Alunos/Usuários podem agendar horários com Profissionais disponíveis.
*   Visualização de agendas e cancelamentos.

## ⚙️ Como Executar o Projeto Localmente

1. **Clone o repositório:**
   ```bash
   git clone https://github.com/komatsuhenry1/school-management-system.git
   cd school-management-system
   ```

2. **Configure as Variáveis de Ambiente:**
   Crie um arquivo `.env` na raiz do projeto contendo as credenciais do seu banco de dados e chaves secretas. Exemplo:
   ```env
   DB_HOST=localhost
   DB_USER=postgres
   DB_PASSWORD=sua_senha
   DB_NAME=school_db
   DB_PORT=5432
   JWT_SECRET=super_secret_key
   PORT=8080
   ```

3. **Baixe as dependências:**
   ```bash
   go mod tidy
   ```

4. **Execute o servidor:**
   ```bash
   go run main.go
   ```
   *O servidor iniciará na porta configurada (ex: `http://localhost:8080`).*

## 🛡️ Middlewares
*   **`AuthRoles(roles...)`**: Verifica se a requisição possui um token válido no header `Authorization` e garante que a `Role` (cargo) decodificada seja compatível com a rota acessada.
*   **Rate Limiter / Cors**: (Se configurados no `main.go`).
