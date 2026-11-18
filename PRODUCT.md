# Product

<!-- impeccable:product-schema 1 -->

## Platform

web

## Users

O proprietário usa o Caderno para estudar para vários concursos. Faz simulados completos em sessões que pode interromper e retomar, inclusive no celular sem conexão.

## Product Purpose

Caderno é o aplicativo de simulados pessoais da marca Computadora. O usuário escolhe uma prova, responde ou pula questões na ordem que preferir e pode voltar às respostas antes de entregar. Depois da correção, confere o gabarito, a pontuação e o desempenho por disciplina.

O sucesso do produto significa conseguir fazer e revisar as provas com continuidade, sem pistas prematuras e sem depender de uma plataforma comercial de concursos.

## Positioning

O produto transforma arquivos JSON versionados em Git em provas completas e retomáveis. Cada tentativa mantém um snapshot do caderno utilizado, oculta o gabarito até a entrega e calcula os resultados conforme a estrutura e os pesos do edital.

## Operating Context

- O uso é individual e pode acontecer em sessões longas ou interrompidas.
- As três provas disponíveis hoje seguem o edital da SEFAZ/SC 2026. Cada uma tem 80 questões de conhecimentos gerais com peso 1 e 100 de conhecimentos específicos com peso 2.
- O usuário pode responder em sequência ou navegar diretamente para qualquer questão.
- Para usar offline, o usuário abre a tentativa com conexão e aguarda “Disponível offline”. Pode fechar a aba e retomar em “Provas neste aparelho”, no mesmo navegador. A sessão é anônima.
- “Salvo neste aparelho” confirma a gravação local. “Sincronizado” confirma que o servidor recebeu as escolhas. Entregar sem conexão encerra a edição e deixa a tentativa aguardando correção. Com o site aberto novamente e a rede disponível, o navegador envia a entrega e o servidor calcula a nota.
- Após a finalização, a revisão reúne todas as questões, alternativas, resposta escolhida e gabarito em uma única página rolável.
- As provas são mantidas como arquivos JSON no repositório Git `https://github.com/Guiribei/tests.git`.

## Capabilities and Constraints

- O escopo inclui vários concursos. Ainda está em aberto quais editais entrarão no catálogo e quais formatos adicionais serão necessários. A distribuição de questões, os pesos e o cálculo descritos aqui correspondem às provas atuais.
- Nenhuma questão informa acerto ou erro antes da finalização da tentativa.
- Questões em branco podem ser entregues e contam como erro.
- O formato canônico de prova contém metadados, textos-base, conhecimentos gerais, conhecimentos específicos e gabarito.
- O validador atual exige dois blocos, 180 questões numeradas de 1 a 180, cinco alternativas por questão e um gabarito completo de A a E.
- Uma tentativa guarda um snapshot do JSON para não ser alterada por atualizações posteriores do Git.
- O servidor Go usa `net/http` e SQLite. Os formulários funcionam online sem JavaScript. Para uso offline, o service worker guarda a interface e o IndexedDB guarda a tentativa e as respostas pendentes. Limpar os dados do navegador antes da sincronização apaga as respostas que só existem no aparelho.
- Nas provas atuais, a nota bruta ponderada chega a 280 pontos. O aplicativo só apresenta uma estimativa da nota padronizada quando há pelo menos duas tentativas concluídas no grupo de comparação e variação nos dois blocos. Quando faltam dados, mostra os pontos brutos e informa a limitação.

## Brand Commitments

- **Computadora** é a marca. O proprietário já possui o domínio correspondente. **Caderno** é o nome do aplicativo.
- A identificação pública do aplicativo deve servir a diferentes concursos. A diretriz de marca é manter as referências à SEFAZ fora dos textos gerais da interface. Para as ações de estudo, usar “prova”, “questão”, “resposta”, “gabarito” e “resultado”.

## Evidence on Hand

- `option-one.json`, `option-two.json` e `option-three.json` contêm as três provas completas. Cada arquivo referencia o edital de origem em `metadados.edital_url`.
- O `README.md` documenta o uso offline, a sincronização e a pontuação.
- `internal/web/templates/` contém o índice de provas, a resolução das questões e a revisão dos resultados. O índice identifica o concurso pelos dados de cada prova.
- Não existem depoimentos, resultados oficiais, métricas públicas ou outras provas sociais; trabalhos futuros não devem fabricá-los.

## Product Principles

1. **Simular antes de ensinar:** durante a prova, preservar o contexto de avaliação e não antecipar o gabarito.
2. **Retomar sem atrito:** salvar respostas e permitir que uma sessão longa continue do ponto em que parou.
3. **Ser preciso sobre a nota:** separar pontuação bruta, estimativas e resultados oficiais, expondo qualquer limitação estatística.
4. **Manter o conteúdo auditável:** tratar os arquivos JSON e seus snapshots como fonte verificável de cada tentativa.
5. **Servir ao uso pessoal:** priorizar clareza, velocidade e confiabilidade para um único usuário, sem complexidade de comunidade ou operação multiusuário.
