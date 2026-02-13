# Guia de Implementação - Como Usar Este Projeto

## 🎯 Objetivo

Este documento explica como usar o **DEVELOPMENT_PLAN.md** para implementar o projeto passo a passo, utilizando IA (ChatGPT, Claude, Copilot) ou desenvolvendo manualmente.

---

## 📋 Preparação Inicial

### 1. Configurar Ambiente Local

```bash
# Criar diretório do projeto
mkdir microservices-observability-portfolio
cd microservices-observability-portfolio

# Inicializar Git
git init
git branch -M main

# Copiar o DEVELOPMENT_PLAN.md para referência
# (você já tem ele aqui)
```

### 2. Criar Estrutura Base

```bash
# Criar diretórios principais
mkdir -p services/{product-api,product-consumer}
mkdir -p infrastructure/{postgres,kafka,observability}
mkdir -p scripts
mkdir -p sre/{slos,runbooks,postmortems}
```

---

## 🤖 Prompt Template para Implementação com IA

Use este template para cada fase do projeto. Copie e cole para seu assistente de IA favorito (ChatGPT, Claude, etc):

### Template Base de Prompt

```
Estou implementando um projeto de portfólio de observabilidade para demonstrar 
competências em OpenTelemetry, ClickHouse, Prometheus, Grafana e SRE practices.

Tenho um plano de desenvolvimento completo em DEVELOPMENT_PLAN.md.

CONTEXTO DO PROJETO:
- Arquitetura de microserviços (Python + Golang)
- Event streaming com Kafka
- Stack de observabilidade completa
- Foco em demonstrar habilidades para vaga de Observability Engineer

FASE ATUAL: [NÚMERO E NOME DA FASE]

PROMPT DA FASE:
[COPIAR O PROMPT DA SEÇÃO "Prompt para Implementação" DA FASE]

INSTRUÇÕES ADICIONAIS:
1. Gerar código completo e funcional
2. Seguir best practices da linguagem
3. Incluir comentários explicativos
4. Criar arquivos na estrutura correta
5. Fornecer comandos de teste
6. Sugerir próximos passos
7. Evitar comentarios desnecessario
8. Evitar emoji em logs
9. Manter limpo

Por favor, implemente esta fase completamente.
```

---

## 📝 Exemplo Prático - Fase 1

### Prompt Completo para Fase 1

```
Estou implementando um projeto de portfólio de observabilidade para demonstrar 
competências em OpenTelemetry, ClickHouse, Prometheus, Grafana e SRE practices.

CONTEXTO DO PROJETO:
- Arquitetura de microserviços (Python FastAPI + Golang)
- Event streaming com Kafka
- PostgreSQL para persistência
- Stack de observabilidade completa (OpenTelemetry, ClickHouse, Prometheus, Grafana)
- Foco em demonstrar habilidades para vaga de Observability Engineer na Wolt

FASE ATUAL: Fase 1 - Setup Inicial da Infraestrutura

OBJETIVO:
Configurar Docker Compose com todos os serviços de infraestrutura necessários 
(Kafka, Zookeeper, PostgreSQL).

TAREFAS:
1. Criar docker-compose.yml com:
   - Zookeeper (Confluent)
   - Kafka (Confluent) com configurações:
     * KAFKA_ADVERTISED_LISTENERS para localhost e rede interna
     * Auto-criação de tópicos habilitada
     * Replication factor adequado para desenvolvimento
   - PostgreSQL 15:
     * Database: products_db
     * User/Password configuráveis via .env
     * Volume persistente
     * Health check
   - Redes:
     * kafka-network (Kafka e serviços)
     * db-network (PostgreSQL e consumer)

2. Criar infrastructure/postgres/init.sql:
   - CREATE TABLE products com campos:
     * id UUID PRIMARY KEY
     * name VARCHAR(255) NOT NULL
     * price DECIMAL(10,2) NOT NULL
     * description TEXT
     * created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
   - Índices apropriados

3. Criar .env.example:
   - Todas as variáveis de ambiente necessárias
   - Documentação inline

4. Criar .gitignore:
   - Python (__pycache__, .pytest_cache, venv)
   - Go (vendor/, *.exe)
   - Docker (.env, volumes/)
   - IDE (.vscode/, .idea/)

REQUISITOS:
- Usar versões estáveis e recentes
- Health checks em todos os serviços
- Restart policies adequados
- Volumes nomeados para persistência
- Configurações otimizadas para desenvolvimento local

CRITÉRIOS DE ACEITE:
- [ ] Docker Compose sobe todos os serviços sem erros
- [ ] Kafka está acessível na porta 9092
- [ ] PostgreSQL está acessível na porta 5432
- [ ] Health checks estão funcionando
- [ ] Tabela products é criada automaticamente
- [ ] .env.example documenta todas as variáveis

Por favor, gere todos os arquivos necessários com código completo e funcional.
```

---

## 🔄 Workflow Recomendado

### Para Cada Fase:

#### 1. **Preparação** (5-10 min)
```bash
# Criar branch para a fase
git checkout -b phase-[NÚMERO]-[nome-curto]

# Exemplo:
git checkout -b phase-1-infrastructure
```

#### 2. **Leitura** (10-15 min)
- Ler a seção da fase no DEVELOPMENT_PLAN.md
- Entender objetivos e critérios de aceite
- Revisar comandos de teste

#### 3. **Implementação** (tempo varia por fase)

**Opção A: Com IA**
```
1. Copiar o template de prompt acima
2. Preencher com informações da fase
3. Enviar para ChatGPT/Claude/Copilot
4. Revisar código gerado
5. Criar arquivos localmente
6. Ajustar conforme necessário
```

**Opção B: Manual**
```
1. Criar arquivos conforme especificação
2. Implementar código seguindo o prompt da fase
3. Consultar documentação oficial quando necessário
```

#### 4. **Teste** (15-30 min)
```bash
# Executar comandos de teste da fase
# (copiados da seção "Comandos de Teste")

# Exemplo Fase 1:
docker-compose up -d
docker-compose ps
docker-compose logs kafka
```

#### 5. **Validação** (10 min)
- Verificar todos os critérios de aceite ✅
- Garantir que tudo funciona conforme esperado
- Documentar problemas encontrados e soluções

#### 6. **Commit** (5 min)
```bash
# Adicionar arquivos
git add .

# Commit descritivo
git commit -m "feat(phase-1): implement infrastructure setup

- Add docker-compose.yml with Kafka, Zookeeper, PostgreSQL
- Create database init script
- Add .env.example and .gitignore
- Configure health checks and volumes

✅ All acceptance criteria met"

# Merge para main
git checkout main
git merge phase-1-infrastructure
```

#### 7. **Documentação** (10 min)
```bash
# Atualizar README.md com:
# - O que foi implementado
# - Como rodar
# - Screenshots (se aplicável)
```

---

## 📊 Tracking de Progresso

### Checklist Geral

Copie isso para um arquivo `PROGRESS.md` e vá marcando:

```markdown
# Progresso do Projeto

## Fases Básicas (1-8)
- [ ] Fase 1: Infraestrutura (Docker Compose) - ⏱️ 2-3h
- [ ] Fase 2: API Python (FastAPI + OTel) - ⏱️ 4-5h
- [ ] Fase 3: Kafka Producer (Python) - ⏱️ 2-3h
- [ ] Fase 4: Estrutura Go (Base + OTel) - ⏱️ 3-4h
- [ ] Fase 5: Repository PostgreSQL (Go) - ⏱️ 3-4h
- [ ] Fase 6: Kafka Consumer (Go completo) - ⏱️ 4-5h
- [ ] Fase 7: Chaos Engineering - ⏱️ 2-3h
- [ ] Fase 8: Testes e Documentação - ⏱️ 4-5h

## Fases de Observabilidade (9-12)
- [ ] Fase 9: OpenTelemetry + Collector - ⏱️ 4-6h
- [ ] Fase 10: ClickHouse + Traces/Logs - ⏱️ 5-7h
- [ ] Fase 11: Prometheus + Grafana - ⏱️ 4-5h
- [ ] Fase 12: SLIs/SLOs + Alerting - ⏱️ 3-4h

## Finalização
- [ ] README.md completo
- [ ] Screenshots de dashboards
- [ ] Vídeo demo (opcional)
- [ ] Repositório público no GitHub
- [ ] LinkedIn post sobre o projeto

**Total de Horas Investidas**: _____ / 54h
**Data de Início**: _____
**Data de Conclusão**: _____
```

---

## 🎓 Dicas de Aprendizado

### Aprenda Enquanto Implementa

Para cada fase, dedique tempo para:

1. **Entender os Conceitos** (20% do tempo)
   - Por que essa tecnologia?
   - Como ela funciona?
   - Quais são as alternativas?

2. **Implementar** (50% do tempo)
   - Código funcional
   - Testes passando
   - Documentado

3. **Experimentar** (30% do tempo)
   - Mudar configurações
   - Quebrar e consertar
   - Otimizar
   - Fazer perguntas: "E se...?"

### Recursos por Fase

**Fases 1-3 (Python/Kafka)**:
- [FastAPI Docs](https://fastapi.tiangolo.com/)
- [Kafka Python](https://kafka-python.readthedocs.io/)
- [Pydantic](https://docs.pydantic.dev/)

**Fases 4-6 (Go/Consumer)**:
- [Go by Example](https://gobyexample.com/)
- [Effective Go](https://go.dev/doc/effective_go)
- [kafka-go](https://github.com/segmentio/kafka-go)

**Fases 9-10 (OpenTelemetry/ClickHouse)**:
- [OpenTelemetry Docs](https://opentelemetry.io/docs/)
- [ClickHouse Docs](https://clickhouse.com/docs)
- [OTLP Specification](https://opentelemetry.io/docs/specs/otlp/)

**Fases 11-12 (Prometheus/SRE)**:
- [Prometheus Docs](https://prometheus.io/docs/)
- [Grafana Docs](https://grafana.com/docs/)
- [Google SRE Book](https://sre.google/sre-book/table-of-contents/)

---

## 🚀 Aceleradores

### Scripts Auxiliares

**scripts/new-phase.sh**:
```bash
#!/bin/bash
# Criar estrutura para nova fase

PHASE_NUM=$1
PHASE_NAME=$2

echo "🚀 Iniciando Fase $PHASE_NUM: $PHASE_NAME"

# Criar branch
git checkout -b "phase-$PHASE_NUM-$PHASE_NAME"

# Criar diretório de trabalho temporário
mkdir -p .phase-$PHASE_NUM

echo "✅ Branch criada: phase-$PHASE_NUM-$PHASE_NAME"
echo "📁 Diretório: .phase-$PHASE_NUM"
echo ""
echo "Próximo passo: Ler DEVELOPMENT_PLAN.md - Fase $PHASE_NUM"
```

**scripts/complete-phase.sh**:
```bash
#!/bin/bash
# Finalizar fase

PHASE_NUM=$1

echo "✅ Finalizando Fase $PHASE_NUM"

# Rodar testes
echo "🧪 Rodando testes..."
# (adicionar comandos de teste aqui)

# Commit
echo "💾 Fazendo commit..."
git add .
git commit -m "feat(phase-$PHASE_NUM): complete implementation"

# Merge
echo "🔀 Merging to main..."
git checkout main
git merge "phase-$PHASE_NUM-*"

echo "🎉 Fase $PHASE_NUM concluída!"
echo ""
echo "Atualizar PROGRESS.md ✓"
```

---

## 💡 Prompts Rápidos por Tipo de Tarefa

### Quando Precisar de Código

```
Contexto: Projeto de observabilidade com microserviços (Python + Go)

Preciso implementar [DESCRIÇÃO].

Requisitos:
- [REQ 1]
- [REQ 2]
- [REQ 3]

Por favor, gere o código completo com:
1. Imports necessários
2. Error handling
3. Comentários explicativos
4. Type hints/annotations
5. Testes básicos (opcional)
```

### Quando Precisar de Configuração

```
Estou configurando [TECNOLOGIA] para [OBJETIVO].

Contexto do ambiente:
- Docker Compose
- Development environment
- [OUTRAS TECNOLOGIAS RELACIONADAS]

Preciso de:
1. Arquivo de configuração completo
2. Explicação de cada seção
3. Best practices aplicadas
4. Variáveis de ambiente documentadas
```

### Quando Tiver Erro

```
Estou na Fase [N] do projeto e encontrei o seguinte erro:

```
[COPIAR ERRO COMPLETO]
```

Contexto:
- O que estava tentando fazer: [DESCRIÇÃO]
- Arquivos envolvidos: [LISTA]
- Configuração relevante: [DESCRIÇÃO]

Logs adicionais:
```
[COPIAR LOGS]
```

Como posso resolver?
```

### Quando Quiser Otimizar

```
Tenho este código funcionando:

```[LINGUAGEM]
[CÓDIGO]
```

Como posso otimizá-lo considerando:
1. Performance
2. Legibilidade
3. Manutenibilidade
4. Best practices de [LINGUAGEM]

Mantenha a funcionalidade atual.
```

---

## 📸 Documentação Visual

### O Que Capturar em Cada Fase

**Fase 1-3**: 
- Screenshot de `docker-compose ps` com todos serviços UP
- Postman/curl testando API

**Fase 6**:
- Fluxo end-to-end funcionando
- Logs de ambos serviços

**Fase 9-10**:
- Traces no Jaeger
- Query results no ClickHouse

**Fase 11**:
- Dashboards Grafana (múltiplos)
- Métricas em tempo real

**Fase 12**:
- Alerts firing
- SLO dashboard

### Criar Diretório de Assets

```bash
mkdir -p docs/screenshots/{phase-1,phase-2,...,phase-12}
mkdir -p docs/diagrams
mkdir -p docs/videos
```

---

## 🎯 Meta Final: Projeto Completo

### Quando Terminar Todas as Fases

#### 1. Revisar README.md Principal

Deve incluir:
- [ ] Badges (build status, coverage, etc)
- [ ] Arquitetura clara com diagrama
- [ ] Quick start em < 5 minutos
- [ ] Screenshots dos dashboards
- [ ] Link para demo video (se tiver)
- [ ] Seção "What I Learned"
- [ ] Tecnologias com links
- [ ] Contato/Social links

#### 2. Preparar para Compartilhar

```bash
# Criar tag de release
git tag -a v1.0.0 -m "Complete observability platform portfolio"
git push origin v1.0.0

# Criar Release no GitHub com:
# - Descrição completa
# - Screenshots
# - Link para demo
# - Highlights do projeto
```

#### 3. Escrever Post LinkedIn

Template:
```
🚀 Projeto Completo: Plataforma de Observabilidade Escalável

Acabei de finalizar um projeto de portfólio demonstrando práticas modernas 
de observabilidade e SRE.

🔧 Stack Técnica:
• Go + Python para microserviços
• OpenTelemetry para instrumentação
• ClickHouse processando bilhões de eventos
• Kafka para event streaming
• Prometheus + Grafana para métricas
• SLIs/SLOs com error budgets

💡 Destaques:
• Distributed tracing end-to-end
• Chaos engineering para testes de resiliência
• Dashboards customizados com Grafana
• Alerting baseado em SLOs

📊 Resultados:
• [X] horas de desenvolvimento
• [X] linhas de código
• [X] dashboards
• [X] queries de análise

🔗 GitHub: [link]
📹 Demo: [link]

#observability #sre #golang #python #opentelemetry #clickhouse #devops

[Tags de pessoas/empresas relevantes]
```

#### 4. Preparar para Entrevistas

Estude para falar sobre:
- Decisões arquiteturais
- Desafios enfrentados
- Trade-offs considerados
- O que faria diferente
- Como escalaria para produção
- Lições aprendidas

---

## 🎁 Bônus: Melhorias Futuras

Depois de completar as 12 fases, considere adicionar:

### Kubernetes Deployment
- Helm charts
- HPA (Horizontal Pod Autoscaler)
- Network policies
- Service mesh (Istio/Linkerd)

### CI/CD Pipeline
- GitHub Actions
- Automated testing
- Container scanning
- Deployment automation

### Advanced Features
- Multi-tenancy
- Authentication/Authorization
- Rate limiting
- API Gateway (Kong/Traefik)

### Análise Avançada
- Anomaly detection com ML
- Trace analysis com grafos
- Cost optimization dashboard
- Capacity planning tools

---

## 📞 Suporte

### Comunidades para Ajuda

- **OpenTelemetry**: [CNCF Slack](https://slack.cncf.io) #otel
- **ClickHouse**: [ClickHouse Slack](https://clickhouse.com/slack)
- **Go**: [Gophers Slack](https://gophers.slack.com)
- **SRE**: [SRE Weekly](https://sreweekly.com)

### Quando Estiver Travado

1. **Consulte documentação oficial**
2. **Busque issues similares no GitHub**
3. **Use IA para debugging** (ChatGPT/Claude)
4. **Simplifique o problema**
5. **Pare e volte depois** (às vezes ajuda!)

---

## ✅ Checklist de Lançamento

Antes de considerar o projeto "completo":

- [ ] Todas as 12 fases implementadas
- [ ] Testes passando
- [ ] Documentação completa
- [ ] README com screenshots
- [ ] .env.example atualizado
- [ ] Sem credenciais hardcoded
- [ ] Código formatado e linted
- [ ] Git history limpo
- [ ] GitHub repository público
- [ ] License file (MIT recomendado)
- [ ] CONTRIBUTING.md (se open source)
- [ ] Demo gravado (opcional)
- [ ] LinkedIn post publicado
- [ ] Adicionado ao currículo/portfólio

---

**Boa sorte com seu projeto de portfólio! 🚀**

Lembre-se: O objetivo não é apenas completar, mas **aprender** e **demonstrar** suas habilidades. Cada fase é uma oportunidade de aprofundar conhecimento.

**"The best way to learn is by doing."** - Go build something amazing! 💪
