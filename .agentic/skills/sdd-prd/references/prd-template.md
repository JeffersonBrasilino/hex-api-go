#### Product Requirements Document (PRD) — Template

Universal schema for every PRD produced by `spec-prd-v2`. Fill every section; if a section does not
apply, write "Não se aplica" and a one-line reason — do not delete it. Write in pt-BR. Keep it
product-focused: describe **behavior and intent**, never implementation (no architecture, files,
libraries, layers, or contracts).

```markdown
# [Título da Funcionalidade] — Documento de Requisitos do Produto (PRD)

| Campo                  | Valor                          |
|------------------------|--------------------------------|
| **Autor**              | [Nome]                         |
| **Última Atualização** | [YYYY-MM-DD]                   |
| **Status**             | Draft / Em Revisão / Aprovado  |
| **Versão**             | [Versão]                       |
| **Tipo**               | feat / fix / refactor / perf / chore / docs / test / build / ci |

## 1. Visão Geral

[Qual problema esta funcionalidade resolve? Por que estamos construindo isso? Que valor entrega aos
usuários? Máximo 2-3 parágrafos.]

## 2. Objetivos

[Objetivos específicos e mensuráveis. O que define sucesso?]

- [Objetivo 1]
- [Objetivo 2]

## 3. Histórias de Usuário

[Formato: "Como um [usuário], eu quero [ação], para que [benefício]."]

- **Épico 1:** [Título]
    - Como um [usuário], eu quero [ação], para que [benefício].
- **Épico 2:** [Título]
    - Como um [usuário], eu quero [ação], para que [benefício].

## 4. Requisitos Funcionais

[O que o sistema deve fazer — comportamentos observáveis pelo usuário ou por sistemas integrados.
Sem detalhe técnico.]

- **RF-01:** [Comportamento esperado]
- **RF-02:** [Comportamento esperado]

## 5. Requisitos Não Funcionais

[Qualidade e restrições em linguagem de negócio. Sem decisões de implementação.]

- **RNF-01 (Performance):** [Ex: responder em até X s para Y usuários simultâneos]
- **RNF-02 (Segurança):** [Ex: dados sensíveis não expostos em logs ou erros]
- **RNF-03 (Disponibilidade):** [Ex: disponível 24/7]
- **RNF-04 (Usabilidade):** [Ex: mensagens de erro claras e orientadoras]

## 6. Regras de Negócio

[Restrições, políticas e invariantes de negócio. Sem detalhes de implementação.]

- **RN-01:** [Regra]
- **RN-02:** [Regra]

## 7. Requisitos de Dados

[Dados manipulados pela funcionalidade. Foco em origem, responsabilidade e privacidade.]

- **Entidades envolvidas:** [Quais entidades de negócio são criadas, alteradas ou consultadas]
- **Origem dos dados:** [Usuário final, sistema externo, evento de domínio, etc.]
- **Privacidade / LGPD:** [Há dados pessoais? Requer consentimento? Quem acessa/edita? Há obrigação
  de anonimização ou retenção limitada? Se não houver dados pessoais, declare explicitamente.]

## 8. Critérios de Aceitação

[Condições testáveis no formato BDD (Dado/Quando/Então). Mínimo de 2 cenários: ao menos 1 caminho
feliz e ao menos 1 de erro/borda. Cada cenário tem nome.]

- **Cenário 1 — [Nome]:** Dado [contexto], quando [ação], então [resultado esperado].
- **Cenário 2 — [Nome — erro/borda]:** Dado [contexto inválido ou falha], quando [ação], então
  [comportamento de erro esperado].

## 9. Fora de Escopo

[O que explicitamente NÃO entra nesta entrega — evita ambiguidade e scope creep.]

- [Item fora de escopo 1]
- [Item fora de escopo 2]

## 10. Perguntas Abertas

[Dúvidas que precisam de resposta antes ou durante o desenvolvimento.]

- [Pergunta 1]

## 11. Considerações Futuras

[Melhorias conhecidas e funcionalidades relacionadas fora deste lançamento, mas que valem rastrear.]

- [Consideração 1]
```
