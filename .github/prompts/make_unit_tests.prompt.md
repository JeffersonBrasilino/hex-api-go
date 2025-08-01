---
mode: agent
model: GPT-4o
description: Generate unit tests for the provided code files.
---

Analise o arquivo, e gere um arquivo de teste unitário:

* O teste deve cobrir 100% das linhas de código.
* Cada método e função deve ser checado seus parametros de entrada(se houver) e seus parametros de saída(se houver)
* Gere um arquivo de teste por arquivo de código.
* não use o pacote reflect, use apenas o pacote testing
* Use somente as structs e funções que estão no arquivo analisado, não importe outros pacotes.
* Use o nome do arquivo analisado como base para o nome do arquivo de teste, adicionando o sufixo `_test.go`.
* não reescreva as structs e funções, apenas as chame.
* Não gere testes para funções que não estão no arquivo analisado.
* gere testes de sucesso e falha, se houver possibilidade de falha.