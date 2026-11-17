# Caderno — simulados SEFAZ/SC 2026

Aplicação de simulados feita com Go, `net/http`, HTMX e SQLite. As provas são arquivos JSON versionados no repositório; o progresso de cada navegador é salvo automaticamente e a correção só aparece depois da finalização.

## Executar localmente

Requisitos: Go 1.27.4 e Git.

```bash
go run ./cmd/server
```

Abra `http://localhost:8080`. Por padrão, a aplicação procura os arquivos `*.json` no diretório atual e cria o banco em `./data/provas.db`.

Para usar um checkout separado e atualizá-lo ao iniciar:

```bash
EXAMS_DIR=./data/exams \
EXAMS_REPOSITORY_URL=https://github.com/Guiribei/tests.git \
SYNC_ON_START=true \
go run ./cmd/server
```

Se `EXAMS_DIR` ainda não existir, a aplicação clona o repositório. Se já for um checkout Git, executa `git pull --ff-only`. Uma falha de sincronização não apaga os arquivos locais existentes.

## Configuração

| Variável | Padrão | Uso |
|---|---|---|
| `ADDR` | `:8080` | Endereço HTTP |
| `DATABASE_PATH` | `./data/provas.db` | Arquivo SQLite |
| `EXAMS_DIR` | `.` | Diretório dos JSONs |
| `EXAMS_REPOSITORY_URL` | `https://github.com/Guiribei/tests.git` | Origem Git |
| `SYNC_ON_START` | `false` | Clona/atualiza a origem ao iniciar |
| `SECURE_COOKIES` | `false` | Deve ser `true` atrás de HTTPS |

## Formato das provas

`option-one.json` é o contrato canônico. O antigo `option-two.json` foi convertido para o mesmo formato, preservando suas questões e seu gabarito. Cada arquivo deve conter:

- `metadados`, inclusive as duas provas, pesos e durações;
- `textos_base`, que pode ser uma lista vazia;
- `conhecimentos_gerais`, com 80 questões;
- `conhecimentos_especificos`, com 100 questões;
- `gabarito`, com as chaves `"1"` a `"180"` e respostas de `A` a `E`.

O carregamento falha de forma explícita se houver numeração duplicada, questão sem cinco alternativas, gabarito incompleto ou total diferente de 180.

## Nota

A aplicação sempre mostra os acertos e os pontos brutos ponderados, no máximo 280 (`gerais + 2 × específicas`).

O edital usa nota padronizada por prova:

```text
NP = ((A - média do grupo) / desvio-padrão do grupo) × 10 + 50
nota final = NP gerais + 2 × NP específicas
```

O mínimo de habilitação é 150. Como a nota depende de todos os candidatos, a aplicação só apresenta uma **estimativa** quando a instalação possui ao menos duas tentativas concluídas e desvio-padrão não nulo nas duas provas. A interface nunca chama essa estimativa de resultado oficial.

## Testes

```bash
go test ./...
```

Os testes validam os dois arquivos de prova, o cálculo de padronização e o fluxo de persistência no SQLite.

## Docker

```bash
docker build -t caderno-sefaz .
docker run --rm -p 8080:8080 -v caderno-data:/data caderno-sefaz
```

Para sincronizar uma cópia independente do repositório a cada inicialização:

```bash
docker run --rm -p 8080:8080 -v caderno-data:/data \
  -e EXAMS_DIR=/data/exams -e SYNC_ON_START=true \
  -e SECURE_COOKIES=true caderno-sefaz
```

## Decisões do MVP

- Sem cadastro: a tentativa pertence a um cookie anônimo do navegador.
- Uma tentativa ativa por prova e navegador; ao voltar, ela é retomada.
- Questões em branco podem ser finalizadas e contam como erro.
- Cada tentativa guarda um snapshot do JSON para não mudar durante uma atualização do Git.
- HTMX é usado para salvar respostas e navegar sem recarregar a página inteira; sem JavaScript, os links e formulários essenciais continuam funcionando.

