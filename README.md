# Simulados de concurso

Aplicação de simulados feita com Go, JavaScript e SQLite. O navegador baixa uma tentativa para responder offline. O Go calcula a nota após receber a entrega.

## Executar localmente

Requisitos: Go 1.27.4 e Git.

```bash
go run ./cmd/server
```

Abra `http://localhost:8080`. O servidor clona o repositório de provas em `./data/exams` e verifica se há atualizações ao iniciar e a cada minuto. Ele percorre a raiz e as subpastas em busca de JSONs. Depois de enviar uma prova, recarregue a página inicial após a próxima sincronização. As tentativas em andamento continuam usando a versão com que foram iniciadas.

Para trabalhar apenas com os arquivos do diretório atual, execute

```bash
EXAMS_DIR=. SYNC_ON_START=false SYNC_INTERVAL=0 go run ./cmd/server
```

## Responder sem conexão

No celular, acesse o mesmo domínio por **HTTPS**, abra uma prova com conexão e aguarde **Disponível offline**. O navegador guarda a interface e as 180 questões daquela tentativa. Em desenvolvimento, o navegador também permite esse recurso em `localhost`.

Você pode avançar, voltar, mudar uma resposta e fechar a aba. Para retomar, abra o site no mesmo navegador e escolha a tentativa em **Provas neste aparelho**. **Salvo neste aparelho** confirma a gravação local. **Sincronizado** confirma que o Go recebeu as escolhas. Limpar os dados do site antes do envio remove as respostas que só existem no celular.

Clique duas vezes em uma alternativa para deixá-la riscada como descartada; repita para restaurá-la. Um clique simples escolhe a alternativa e remove seu descarte. Se descartar a resposta marcada, a questão fica em branco. Com a alternativa em foco, a tecla **Delete** também alterna o descarte. Os descartes ficam salvos por tentativa neste navegador, inclusive offline.

Ao entregar sem conexão, a tentativa fica bloqueada para edição e aparece **Aguardando conexão para corrigir**. Quando a rede voltar, deixe o site aberto ou abra-o novamente. O servidor grava a entrega e calcula a nota usando a versão do gabarito salva no início da tentativa.

## Atualizar um serviço systemd existente

Consulte o caminho do binário em `ExecStart` com `systemctl cat caderno`. Compile a nova versão e instale-a nesse caminho. Por exemplo, para `ExecStart=/usr/local/bin/caderno`, execute

```bash
CGO_ENABLED=0 go build -trimpath -o ./caderno ./cmd/server
sudo install -m 0755 ./caderno /usr/local/bin/caderno
```

O usuário do serviço precisa ter acesso ao Git pelo `PATH`, acesso de leitura ao repositório remoto e permissão de escrita na cópia local. Para reservar um diretório próprio e manter o SQLite no caminho já configurado, abra `sudo systemctl edit caderno` e acrescente

```ini
[Service]
StateDirectory=caderno
Environment=EXAMS_DIR=/var/lib/caderno/exams
Environment=SYNC_ON_START=true
Environment=SYNC_INTERVAL=1m
```

O Git clona em um diretório vazio ou inexistente. Se já houver um checkout nesse caminho, a aplicação executa `git pull --ff-only` usando o upstream do branch atual. Para colocar a nova versão em execução em um serviço já ativo e acompanhar os avisos de sincronização, use

```bash
sudo systemctl restart caderno
journalctl -u caderno -f
```

## Configuração

| Variável | Padrão | Uso |
|---|---|---|
| `ADDR` | `:8080` | Endereço HTTP |
| `DATABASE_PATH` | `./data/provas.db` | Arquivo SQLite |
| `EXAMS_DIR` | `exams/` ao lado de `DATABASE_PATH` | Checkout dedicado às provas |
| `EXAMS_REPOSITORY_URL` | `https://github.com/Guiribei/tests.git` | Origem Git |
| `SYNC_ON_START` | `true` | Clona/atualiza a origem ao iniciar |
| `SYNC_INTERVAL` | `1m` | Intervalo entre consultas ao Git e recargas do catálogo. Use `0` para desativar a consulta periódica |
| `SECURE_COOKIES` | `false` | Deve ser `true` atrás de HTTPS |

## Formato das provas

`option-one.json`, `option-two.json` e `option-three.json` seguem o mesmo formato. Cada arquivo deve conter:

- `metadados`, inclusive as duas provas, pesos e durações;
- `textos_base`, que pode ser uma lista vazia;
- `conhecimentos_gerais`, com 80 questões;
- `conhecimentos_especificos`, com 100 questões;
- `gabarito`, com as chaves `"1"` a `"180"` e respostas de `A` a `E`.

O servidor registra o arquivo e o motivo de cada falha de validação. Um JSON inválido não impede que as demais provas sejam publicadas. Se uma prova já carregada receber uma edição inválida, a versão anterior continua disponível em memória. Uma remoção do repositório retira a prova da lista na próxima atualização.

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

Os testes validam os três arquivos de prova, a recarga do catálogo, o cálculo de padronização e o fluxo de persistência no SQLite.

## Docker

```bash
docker build -t caderno-sefaz .
docker run --rm -p 8080:8080 -v caderno-data:/data caderno-sefaz
```

A imagem inclui as provas da raiz para execução offline. Para sincronizar uma cópia independente do repositório ao iniciar e a cada minuto, execute

```bash
docker run --rm -p 8080:8080 -v caderno-data:/data \
  -e EXAMS_DIR=/data/exams -e SYNC_ON_START=true -e SYNC_INTERVAL=1m \
  -e SECURE_COOKIES=true caderno-sefaz
```

## Decisões do MVP

- Sem cadastro: a tentativa pertence a um cookie anônimo do navegador.
- Uma tentativa ativa por prova e navegador; ao voltar, ela é retomada.
- Questões em branco podem ser entregues e contam como erro.
- Cada tentativa guarda um snapshot do JSON para não mudar durante uma atualização do Git.
- O service worker guarda a interface e o IndexedDB guarda a prova baixada e as respostas pendentes. Os formulários HTML também funcionam online com JavaScript desativado.

