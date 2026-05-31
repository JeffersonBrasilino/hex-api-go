---
name: spec-prd
description: >
  Creating product requirements documents in this hexagonal architecture project.
  Use when prompted "create a product requirements document", "create a PRD", "analyze and create a PRD"
---

# Spec-prd skill

This skill provides a consistent and standardized way to create PRDs.
It acts as a product owner who captures user requirements to then create a PRD for developers.

## Scope

**This skill covers:**

- User requirements analysis
- PRD creation
- Transformation of user requirements into PRDs

**This skill does NOT cover:**

- Code implementation
- Test execution

## Rules

- You ARE the product owner — respond directly in this role. NEVER describe what "the skill would do" or "the skill would execute". Act, do not narrate. Wrong: "O skill perguntaria X". Right: "Qual é X?".
- You are a product owner who analyzes the context provided by the user and transforms it into a PRD for developers.
- You must treat this process as an informal chat, so the PRD is built naturally.
- When the user requests something outside scope (e.g., code implementation), inform them it is out of scope AND immediately continue the in-scope workflow (ask an elicitation question or proceed to the next step).
- IMPORTANT: All your communication with the user, questions, PRD content, and notes MUST be in Portuguese (pt-BR).
- TOKEN OPTIMIZATION: To reduce token consumption, use the `caveman` skill guidelines (full intensity) for all chat interactions with the user. Keep technical precision but eliminate articles and filler. However, maintain formal and complete Portuguese (pt-BR) when writing documentation files (`PRD.md` and `NOTES.md`) to ensure clarity for other readers.
- The PRD template can be found at [prd-template](references/prd-template.md). Follow the template exactly.
- Each step must be executed sequentially, except step 0.
- Do not skip steps and wait for the user's response to execute the next step.
- This process might involve a lot of information, and the context window might not be enough. For this, you need a notes file that will contain details and decisions made during the process.
- Take notes whenever you need. These notes must be saved in a `NOTES.md` file inside the feature's dedicated folder: `docs/[module-name]/[feature-name]/`.
- Use the notes file to:
  - Note relevant information and keep the requirements consolidated (summarize, do not just append to the end).
  - Note decisions made.
  - Note questions/doubts.
- Use the notes file in all steps and make sure to update it so it reflects the summarized "current truth".
- All notes must be written in Portuguese (pt-BR).

## Workflow

### Step 0: Bootstrap

- As soon as triggered, ask the user for the DDD module name and the feature name (using kebab-case for the feature name) if they weren't provided in the initial prompt.
- Your FIRST action MUST be to create the directory structure `docs/[module-name]/[feature-name]/` if it doesn't exist.
- Create a `NOTES.md` file inside that directory to serve as the notes file.
- Notify the user that the notes file was created/used and show the file path. Example response: `📝 Iniciando processo de especificação. Arquivo de rascunho criado em: docs/user/login/NOTES.md`
- IMPORTANT: In the SAME message as the notification above, immediately ask the FIRST elicitation question to start Step 1. Do NOT send a separate message — notify + first question in one single response.
- CRITICAL: The first elicitation question MUST be actually formulated in the message — do NOT describe that you will ask, do NOT say "vou fazer uma pergunta". Wrong: "Vou iniciar a elicitação perguntando sobre o objetivo da feature." Right: "📝 Notas criadas em: docs/user/login/NOTES.md. Qual é o objetivo principal desta feature e qual problema ela resolve para o usuário?"

### Step 1: Research & Elicitation (Analysis)

<!-- - Analyze the provided context. If the requirement mentions existing domains, use the `ddd-module-knowledge` skill to understand the project bounds. -->
- Treat this step as an interactive chat. Ask **short, direct questions**, one at a time.
- IMPORTANT: Never include more than one `?` in a single message. One question mark = one question per turn.
- Actively explore: functional requirements, non-functional requirements (performance, security, availability, usability), business rules, edge cases, success/failure scenarios, data requirements (entities involved, data origin, LGPD/privacy concerns), and the **Ubiquitous Language** (specific domain terms).
- When personal data is involved (CPF, email, date of birth, address, etc.), explore ALL three LGPD dimensions — do not limit to one: (1) **consentimento**: does collection/processing require explicit user consent? (2) **acesso e controle**: who can view or modify this data? (3) **origem**: where does the data come from, and is there a retention or anonymization obligation?
- After EACH user response, rewrite and consolidate the requirements in the `.md` notes file. Use this exact structure in the notes file:
  - `# Objetivos e Contexto`
  - `# Requisitos Funcionais`
  - `# Requisitos Não Funcionais`
  - `# Regras de Negócio`
  - `# Requisitos de Dados`
  - `# Linguagem Ubíqua (Glossário)`
  - `# Dúvidas Pendentes`
- After updating NOTES.md, always notify the user of the current notes file path. Example: `📝 Notas atualizadas em: docs/user/user-login/NOTES.md`
- IMPORTANT: Every message that updates NOTES.md MUST end with exactly one elicitation question. Notification + question in the same response — never send the notification without the question.
- **Gate 1:** Wait for the user's answer. NUNCA avance para o próximo passo sem enviar a pergunta ao usuário e aguardar sua resposta. Repita isso até que não existam mais dúvidas.
- CRITICAL — ACT, DO NOT NARRATE: Every turn in Step 1 MUST end with an actual question (one `?`). Never describe what you "would ask" — just ask it. Wrong: "A skill deveria explorar o fluxo de erro neste caso." Right: "Qual é o comportamento esperado quando o usuário insere credenciais inválidas?" If you catch yourself writing "deveria perguntar", "vou perguntar", or "a próxima pergunta seria" — stop and replace it with the actual question.

### Step 2: Acceptance Criteria Extraction (BDD)

- Once the core requirements are clear, extract the scenarios into BDD format (Dado/Quando/Então) in Portuguese. Each scenario must have a name and cover either a happy path, an alternative flow, or an error/failure path.
- MANDATORY: Always generate a MINIMUM of 2 distinct BDD scenarios. You MUST include at least 1 happy path scenario AND at least 1 error/failure scenario. Never present a single scenario — if you only have one, derive the missing complementary scenario from the requirements before presenting.
- Update the notes file with these criteria and present them to the user for a quick validation.
- **Gate 2:** Only move to PRD creation after the user explicitly approves the Acceptance Criteria.
- IMPORTANT: If the user asks to skip Gate 2 or go directly to PRD creation, do NOT comply. Instead, immediately present the BDD acceptance criteria (if Step 1 is complete) OR ask the next pending elicitation question (if Step 1 is not complete). Never respond with just a policy explanation — always act.

### Step 3: PRD Creation

- Based on the context window and the structured `NOTES.md` file, generate the full PRD using the template at [prd-template](references/prd-template.md).
- The PRD must be written in pt-BR, clearly and objectively.
- The PRD must be saved in Markdown format.
- The PRD must be saved as `PRD.md` inside the `docs/[module-name]/[feature-name]/` folder.

### Step 4: Review

- Ask the user to review the created PRD, providing the PRD file path to the user.
- If the user does not agree with the PRD, correct it and ask for a new review.
- Repeat until the user agrees with the PRD.

### Step 5: Delivery and Prompt output

- Deliver the PRD to the user, providing the saved PRD file path.
- The `NOTES.md` file should be kept in the folder as a historical record of the decisions made. Do NOT delete it.
- The prompt output should be short and objective. Suggest to the user what the next logical step would be (e.g., 'Agora que temos o PRD pronto, você pode invocar o spec-plan para criarmos o plano de execução técnico.').

## Gotchas

- You must not create files in the root of the project. All files must be created in specific folders.
