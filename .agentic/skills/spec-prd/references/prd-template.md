#### Product Requirements Document (PRD) Template

This is a template for a Product Requirements Document (PRD). 
This template must be used to document all new features that are to be built.
This template must be writed in pt-BR.

```markdown

# [Título da Funcionalidade/Produto] - Documento de Requisitos do Produto (PRD)

| Campo | Valor |
|-------|-------|
| **Autor** | [Seu nome] |
| **Última Atualização** | [YYYY-MM-DD] |
| **Status** | Draft/Em Revisão/Aprovado/Em Desenvolvimento |
| **Versão** | [Versão] |

## 1. Visão Geral (Overview)

[Forneça um resumo de alto nível da funcionalidade. Qual problema está sendo resolvido? Por que estamos construindo isso? Que valor isso entrega aos usuários? Mantenha-o conciso - 2-3 parágrafos no máximo.]

## 2. Objetivos

[Quais são os objetivos específicos e mensuráveis desta funcionalidade? O que define o sucesso? Use o formato de meta SMART onde aplicável.]
- [Objetivo 1]
- [Objetivo 2]
- [Objetivo 3]

## 3. Histórias de Usuário

[Capture os requisitos do usuário no formato padrão: "Como um [usuário], eu quero [ação], para que [benefício]."]

**Épicos:**
- **Épico 1:** [Título do épico]
    - Como um [usuário], eu quero [ação], para que [benefício].
    - Como um [usuário], eu quero [ação], para que [benefício].
- **Épico 2:** [Título do épico]
    - Como um [usuário], eu quero [ação], para que [benefício].

**Requisitos Não-Funcionais (NFRs):**
- [Requisitos de desempenho]
- [Requisitos de segurança]
- [Requisitos de escalabilidade]
- [Requisitos de acessibilidade]

## 4. Regras de Negócio

[Documente quaisquer regras de negócio, restrições ou políticas que se aplicam a esta funcionalidade. Isso é especialmente importante para indústrias regulamentadas.]
- [Regra 1]
- [Regra 2]
- [Regra 3]

## 5. Restrições Arquiteturais

[Decisões e restrições arquiteturais que afetam a implementação desta funcionalidade. Estas restrições orientam o agente de planejamento técnico e reduzem ambiguidades durante a criação do plano de especificação.]

- **Localização de responsabilidades:** [Indicar se lógicas específicas devem residir no domínio, aplicação ou infraestrutura. Ex: "hashing de senha deve ficar na infraestrutura, não no domínio".]
- **Contratos existentes:** [Indicar se contratos/interfaces existentes podem ser alterados ou se devem permanecer inalterados. Ex: "não alterar UserRepository existente".]
- **Organização de arquivos:** [Preferências de granularidade de arquivos. Ex: "múltiplas interfaces coesas podem coexistir no mesmo arquivo" ou "uma interface por arquivo".]
- **Escopo de alteração:** [Definir se a feature pode gerar ripple effects em código existente ou deve ser isolada. Ex: "não alterar handlers existentes" ou "permitido adicionar dependências aos handlers existentes".]
- **Bibliotecas e dependências:** [Restrições sobre uso de bibliotecas externas nas camadas. Ex: "nenhuma lib de terceiros no domínio além do ddgo".]

## 6. Requisitos de Dados

[Especifique os requisitos relacionados a dados, incluindo fontes de dados, transformações de dados, políticas de retenção de dados e considerações de privacidade de dados.]
- **Fontes de dados:** [Liste as fontes de dados]
- **Transformações de dados:** [Descreva quaisquer transformações de dados necessárias]
- **Retenção de dados:** [Por quanto tempo os dados devem ser retidos?]
- **Privacidade:** [Considerações sobre PII, conformidade com GDPR/LGPD]

## 7. Critérios de Aceitação

[Condições claras e testáveis que devem ser atendidas para que a funcionalidade seja considerada completa. Formato: "Dado [contexto], quando [ação], então [resultado]."]

**Critérios Funcionais:**
- **Cenário 1:** Dado [contexto], quando [ação], então [resultado].
- **Cenário 2:** Dado [contexto], quando [ação], então [resultado].
- **Cenário 3:** Dado [contexto], quando [ação], então [resultado].

**Critérios Não-Funcionais:**
- [Critérios de aceitação de desempenho]
- [Critérios de aceitação de segurança]
- [Critérios de aceitação de usabilidade]

## 8. Design & Experiência do Usuário

### Interface do Usuário (UI)
- [Componentes de UI principais e seu propósito]
- [Fluxos do usuário e padrões de interação]
- [Localização de wireframes/mockups]

### Experiência do Usuário (UX)
- [Experiência do usuário desejada e resposta emocional]
- [Abordagem de tratamento de erros]
- [Mecanismos de feedback]

### Arquitetura
- [Decisões arquitetônicas de alto nível]
- [Pontos de integração com sistemas existentes]
- [Considerações sobre a pilha de tecnologia]

### Diagramas
- [Diagramas em mermaid de arquitetura, fluxo, sequência, etc.]

## 9. Critérios de Lançamento

[Quais condições devem ser atendidas antes que esta funcionalidade possa ser lançada em produção? Isso deve se alinhar com o processo de lançamento da sua organização.]
- Todos os critérios de aceitação foram atendidos
- Revisão de código concluída e aprovada
- Testes de QA concluídos com aprovação
- Testes de desempenho concluídos
- Testes de segurança concluídos
- Documentação concluída
- Materiais de treinamento disponíveis (se aplicável)
- Revisão legal/de conformidade (se aplicável)

## 10. Considerações Futuras

[Documente quaisquer melhorias futuras conhecidas, funcionalidades relacionadas ou "nice-to-haves" que estão fora do escopo deste lançamento, mas devem ser rastreadas para consideração futura.]
- [Consideração futura 1]
- [Consideração futura 2]
- [Consideração futura 3]

## 11. Perguntas Abertas

[Quaisquer perguntas que surgiram durante a coleta de requisitos que precisam de respostas antes que o desenvolvimento possa prosseguir.]
- [Pergunta 1]
- [Pergunta 2]
- [Pergunta 3]

```