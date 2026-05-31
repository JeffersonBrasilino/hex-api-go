#### Product Requirements Document (PRD) Template

This is a template for a Product Requirements Document (PRD).
This template must be used to document all new features that are to be built.
This template must be written in pt-BR.

```markdown

# [Título da Funcionalidade] - Documento de Requisitos do Produto (PRD)

| Campo | Valor |
|-------|-------|
| **Autor** | [Nome] |
| **Última Atualização** | [YYYY-MM-DD] |
| **Status** | Draft / Em Revisão / Aprovado |
| **Versão** | [Versão] |

## 1. Visão Geral

[Qual problema esta funcionalidade resolve? Por que estamos construindo isso? Que valor entrega aos usuários? 2-3 parágrafos no máximo.]

## 2. Objetivos

[Quais são os objetivos específicos e mensuráveis desta funcionalidade? O que define o sucesso?]

- [Objetivo 1]
- [Objetivo 2]
- [Objetivo 3]

## 3. Histórias de Usuário

[Capture os requisitos do usuário no formato padrão: "Como um [usuário], eu quero [ação], para que [benefício]."]

- **Épico 1:** [Título do épico]
    - Como um [usuário], eu quero [ação], para que [benefício].
    - Como um [usuário], eu quero [ação], para que [benefício].
- **Épico 2:** [Título do épico]
    - Como um [usuário], eu quero [ação], para que [benefício].

## 4. Requisitos Funcionais

[O que o sistema deve fazer. Descreva comportamentos observáveis pelo usuário ou por sistemas integrados, sem detalhe de implementação técnica.]

- **RF-01:** [Descrição do comportamento esperado]
- **RF-02:** [Descrição do comportamento esperado]
- **RF-03:** [Descrição do comportamento esperado]

## 5. Requisitos Não Funcionais

[Qualidade e restrições do sistema: performance, segurança, disponibilidade, usabilidade. Descreva em linguagem de negócio, sem decisões de implementação.]

- **RNF-01 (Performance):** [Ex: a operação deve responder em até X segundos para Y usuários simultâneos]
- **RNF-02 (Segurança):** [Ex: dados sensíveis não devem ser expostos em logs ou respostas de erro]
- **RNF-03 (Disponibilidade):** [Ex: a funcionalidade deve estar disponível 24/7]
- **RNF-04 (Usabilidade):** [Ex: mensagens de erro devem ser claras e orientar o usuário sobre como proceder]

## 6. Regras de Negócio

[Documente restrições, políticas e invariantes de negócio que se aplicam a esta funcionalidade. Sem detalhes de implementação.]

- **RN-01:** [Regra]
- **RN-02:** [Regra]
- **RN-03:** [Regra]

## 7. Requisitos de Dados

[Dados manipulados pela funcionalidade. Foco em origem, responsabilidade e privacidade — sem descrever transformações técnicas.]

- **Entidades envolvidas:** [Quais entidades de domínio são criadas, alteradas ou consultadas]
- **Origem dos dados:** [Quem fornece os dados: usuário final, sistema externo, evento de domínio, etc.]
- **Privacidade / LGPD:** [Há dados pessoais envolvidos? Requer consentimento? Há obrigação de anonimização ou retenção limitada?]

## 8. Critérios de Aceitação

[Condições claras e testáveis que definem quando a funcionalidade está completa. Formato BDD: Dado/Quando/Então.]

- **Cenário 1 — [Nome]:** Dado [contexto], quando [ação], então [resultado esperado].
- **Cenário 2 — [Nome]:** Dado [contexto], quando [ação], então [resultado esperado].
- **Cenário 3 — [Nome — caminho de erro]:** Dado [contexto com dado inválido ou condição de falha], quando [ação], então [comportamento de erro esperado].

## 9. Perguntas Abertas

[Dúvidas que surgiram durante a coleta de requisitos e que precisam ser respondidas antes do desenvolvimento prosseguir.]

- [Pergunta 1]
- [Pergunta 2]

## 10. Considerações Futuras

[Melhorias conhecidas e funcionalidades relacionadas que estão fora do escopo deste lançamento, mas devem ser rastreadas para o futuro.]

- [Consideração 1]
- [Consideração 2]

```
