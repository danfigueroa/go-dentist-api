# Dental SaaS Platform

Uma plataforma SaaS modular para gestão completa de clínicas odontológicas, construída com **Go** e **AWS DynamoDB**.

## 🏗️ Arquitetura Modular

O projeto foi estruturado com uma arquitetura modular que permite fácil expansão e manutenção:

```
dental-saas/
├── cmd/                    # Ponto de entrada da aplicação
│   └── main.go
├── shared/                 # Recursos compartilhados entre módulos
│   ├── config/            # Configurações (DynamoDB, etc.)
│   └── router/            # Router principal
├── modules/               # Módulos de funcionalidades
│   ├── dental/           # Módulo odontológico
│   │   ├── models/       # Modelos de dados
│   │   ├── handlers/     # Controladores HTTP
│   │   └── router/       # Rotas específicas do módulo
│   └── financial/        # Módulo financeiro (em desenvolvimento)
│       └── models/       # Modelos financeiros
├── docs/                 # Documentação Swagger
└── docker-compose.yml    # Configuração Docker
```

## 📦 Módulos Disponíveis

### 1. Módulo Dental (Implementado)
Gerenciamento completo das operações odontológicas:
- **Dentistas**: Cadastro, consulta, atualização e remoção
- **Pacientes**: Gestão de informações dos pacientes
- **Procedimentos**: Catálogo de procedimentos odontológicos
- **Agendamentos**: Sistema de agendamento de consultas

### 2. Módulo Financeiro (Estrutura Criada)
Gestão financeira da clínica:
- **Receitas**: Controle de entradas financeiras
- **Despesas**: Gestão de gastos (materiais, aluguel, funcionários, etc.)
- **Notas Fiscais**: Emissão e controle de notas fiscais
- **Relatórios**: Análises financeiras (planejado)

## 🚀 Como Executar

### Pré-requisitos
- Docker e Docker Compose
- Go 1.22+ (para desenvolvimento)

### Executando com Docker Compose

1. Clone o repositório:
```bash
git clone <repository-url>
cd dental-saas
```

2. Execute com Docker Compose:
```bash
docker-compose up -d
```

3. Acesse a aplicação:
- **API**: http://localhost:8080
- **Documentação Swagger**: http://localhost:8080/swagger/
- **DynamoDB Local**: http://localhost:8000

### Executando em Desenvolvimento

1. Inicie o DynamoDB Local:
```bash
docker-compose up dynamodb-local -d
```

2. Execute a aplicação:
```bash
go run cmd/main.go
```

## 📚 API Endpoints

### Informações Gerais
- `GET /health` - Status da aplicação
- `GET /api/v1` - Informações da API e módulos disponíveis

### Módulo Dental (`/api/v1/dental`)

#### Dentistas
- `POST /api/v1/dental/dentist` - Criar dentista
- `GET /api/v1/dental/dentist` - Listar todos os dentistas
- `GET /api/v1/dental/dentist/{id}` - Buscar dentista por ID
- `GET /api/v1/dental/dentist/name/{name}` - Buscar dentista por nome
- `GET /api/v1/dental/dentist/cro/{cro}` - Buscar dentista por CRO
- `PUT /api/v1/dental/dentist/{id}` - Atualizar dentista
- `DELETE /api/v1/dental/dentist/{id}` - Remover dentista

#### Pacientes, Procedimentos e Agendamentos
*Rotas similares serão migradas para a nova estrutura modular*

### Módulo Financeiro (`/api/v1/financial`)
*Em desenvolvimento - estrutura de modelos criada*

## 🛠️ Tecnologias Utilizadas

- **Go 1.22**: Linguagem de programação
- **Gorilla Mux**: Router HTTP
- **AWS SDK Go v2**: Cliente DynamoDB
- **DynamoDB Local**: Banco de dados NoSQL
- **Swagger**: Documentação da API
- **Docker**: Containerização
- **UUID**: Geração de identificadores únicos

## 🔧 Configuração

### Variáveis de Ambiente
- `DYNAMODB_ENDPOINT`: Endpoint do DynamoDB (padrão: http://localhost:8000)

### Tabelas DynamoDB
As seguintes tabelas são criadas automaticamente:

**Módulo Dental:**
- `Dentists`
- `Patients`
- `Procedures`
- `Appointments`

**Módulo Financeiro:**
- `Expenses`
- `Revenues`
- `Invoices`

## 🚧 Roadmap

### Próximas Implementações
1. **Migração Completa do Módulo Dental**
   - Migrar handlers de pacientes, procedimentos e agendamentos
   - Atualizar rotas para nova estrutura

2. **Implementação do Módulo Financeiro**
   - Handlers para receitas, despesas e notas fiscais
   - Relatórios financeiros
   - Integração com APIs de pagamento

3. **Novos Módulos Futuros**
   - Módulo de Estoque
   - Módulo de Relatórios e Analytics
   - Módulo de Comunicação (SMS/Email)
   - Módulo de Agendamento Online

### Melhorias Técnicas
- Implementação de middleware de autenticação
- Testes unitários e de integração
- Logging estruturado
- Métricas e monitoramento
- Cache Redis
- API Rate Limiting

## 📝 Contribuição

Este é um projeto pessoal em desenvolvimento. Sugestões e melhorias são bem-vindas!

## 📄 Licença

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para mais detalhes.

---

**Desenvolvido por Daniel Figueroa**  
📧 danielmfigueroa@gmail.com
