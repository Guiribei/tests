# Product

<!-- impeccable:product-schema 1 -->

## Platform

web

## Users

O usuário principal e único é o proprietário do produto, que o utiliza para estudar e realizar simulados completos do concurso SEFAZ/SC 2026.

## Product Purpose

Computadora é uma aplicação pessoal para realizar provas simuladas completas, preservar o progresso entre sessões e revisar o desempenho somente depois da entrega. O produto deve permitir uma rotina de treino próxima da prova: escolher um caderno, responder ou pular questões livremente, finalizar e então conferir respostas, gabarito e pontuação.

O sucesso do produto significa conseguir fazer e revisar as provas com continuidade, sem pistas prematuras e sem depender de uma plataforma comercial de concursos.

## Positioning

O produto transforma arquivos JSON versionados em Git em provas completas e retomáveis. Cada tentativa mantém um snapshot do caderno utilizado, oculta o gabarito até a entrega e calcula os resultados conforme a estrutura e os pesos do edital.

## Operating Context

- O uso é individual e pode acontecer em sessões longas ou interrompidas.
- A prova contém 180 questões: 80 de conhecimentos gerais, com peso 1, e 100 de conhecimentos específicos, com peso 2.
- O usuário pode responder em sequência ou navegar diretamente para qualquer questão.
- O progresso é retomado no mesmo navegador por meio de uma sessão anônima.
- Após a finalização, a revisão reúne todas as questões, alternativas, resposta escolhida e gabarito em uma única página rolável.
- As provas são mantidas como arquivos JSON no repositório Git `https://github.com/Guiribei/tests.git`.

## Capabilities and Constraints

- O escopo de conteúdo é exclusivo ao concurso SEFAZ/SC 2026; não há intenção confirmada de atender outros concursos ou editais.
- A interface pública deve usar apenas o termo “prova” e não deve exibir referências textuais à SEFAZ.
- Nenhuma questão informa acerto ou erro antes da finalização da tentativa.
- Questões em branco podem ser entregues e contam como erro.
- O formato canônico de prova contém metadados, textos-base, conhecimentos gerais, conhecimentos específicos e gabarito.
- Cada arquivo deve conter exatamente 180 questões numeradas continuamente, cinco alternativas por questão e um gabarito completo de A a E.
- Uma tentativa guarda um snapshot do JSON para não ser alterada por atualizações posteriores do Git.
- A persistência usa SQLite e a aplicação não exige cadastro ou autenticação.
- A nota bruta ponderada tem máximo de 280 pontos.
- A nota padronizada depende da média e do desvio-padrão do grupo. Quando não existem dados suficientes, o produto deve informar a limitação em vez de apresentar uma nota oficial ou inventada.
- A aplicação é implementada em Go 1.27.4 com `net/http`, HTMX e SQLite.

## Brand Commitments

- O nome do produto é **Computadora**.
- O nome foi escolhido porque o proprietário já possui o domínio correspondente; isso não constitui uma promessa funcional ou temática.
- A marca e a interface não devem usar SEFAZ como identificação pública. A terminologia visível deve se concentrar em “prova”, “questão”, “resposta”, “gabarito” e “resultado”.

## Evidence on Hand

- `option-one.json`: prova completa no formato canônico.
- `option-two.json`: segunda prova convertida para o mesmo formato canônico.
- `sefsc126_edital_de_abertura_final_publicar(3).pdf`: edital usado como fonte para estrutura, pesos e regra de pontuação.
- `README.md`: decisões técnicas e operacionais atuais.
- `internal/web/templates/`: fluxos implementados de seleção, resolução e revisão.
- Não existem depoimentos, resultados oficiais, métricas públicas ou outras provas sociais; trabalhos futuros não devem fabricá-los.

## Product Principles

1. **Simular antes de ensinar:** durante a prova, preservar o contexto de avaliação e não antecipar o gabarito.
2. **Retomar sem atrito:** salvar respostas e permitir que uma sessão longa continue do ponto em que parou.
3. **Ser preciso sobre a nota:** separar pontuação bruta, estimativas e resultados oficiais, expondo qualquer limitação estatística.
4. **Manter o conteúdo auditável:** tratar os arquivos JSON e seus snapshots como fonte verificável de cada tentativa.
5. **Servir ao uso pessoal:** priorizar clareza, velocidade e confiabilidade para um único usuário, sem complexidade de comunidade ou operação multiusuário.
