#### Decision Notes (NOTES.md) — Template

Durable decision record for a PRD in progress. It survives context compaction and stays as history
after the PRD is written. Write in pt-BR.

**How to use it (progressive disclosure):**

- Created lazily — only on **Trigger A** (context about to be compacted / long conversation) or
  **Trigger B** (after PRD approval). Not at bootstrap.
- **Never overwrite.** If `NOTES.md` already exists, merge into it: update each section to reflect the
  current consolidated truth and append new decisions to the log. Do not discard prior content.
- Consolidate, don't just pile up: rewrite each section so it reads as the *current* state. The only
  append-only section is "Registro de Decisões".
- Path: `docs/[module-name]/[feature-name]/NOTES.md`.

```markdown
# NOTES — [Funcionalidade]

> Rascunho de decisões. Fonte de verdade enquanto o PRD não está aprovado.
> Última atualização: [YYYY-MM-DD]

## Objetivos e Contexto

[Problema, valor e objetivo consolidados.]

## Requisitos Funcionais

- [RF consolidado]

## Requisitos Não Funcionais

- [RNF consolidado — performance, segurança, disponibilidade, usabilidade]

## Regras de Negócio

- [RN consolidada]

## Requisitos de Dados

- **Entidades:** [...]
- **Origem:** [...]
- **Privacidade / LGPD:** [consentimento, acesso/controle, origem/retenção — ou "sem dados pessoais"]

## Linguagem Ubíqua (Glossário)

- **[Termo]:** [definição acordada]

## Critérios de Aceitação (rascunho BDD)

- **Cenário — [Nome]:** Dado [...], quando [...], então [...].

## Dúvidas Pendentes

- [ ] [Dúvida ainda aberta]

## Contexto Técnico (informado pelo usuário, fora do PRD)

[Tudo que o usuário mencionou e que não pertence ao PRD — tecnologias, padrões, protocolos,
estruturas de cache, formatos de token, rotas, TTLs, nomes de campos técnicos etc. Este é o
insumo bruto que a skill `sdd-plan` consome para decidir arquitetura. Nunca descarte: se o
usuário citou, registre aqui mesmo que pareça óbvio ou redundante com o PRD.]

- [Detalhe técnico 1]

## Registro de Decisões (append-only)

- [YYYY-MM-DD] [Decisão tomada e por quê.]
```
