---
name: Questões de concurso
description: Cores, tipografia e componentes da interface de provas.
colors:
  paper: "#ffffff"
  white: "#fff"
  ink: "#293239"
  muted: "#5e6b74"
  line: "#dce1e5"
  control-line: "#82929e"
  mist: "#eef1f3"
  accent: "#365b78"
  accent-hover: "#28475f"
  selected: "#edf3f8"
  selected-line: "#839fb5"
  correct: "#2b654b"
  correct-wash: "#edf5ef"
  incorrect: "#944b3c"
  incorrect-wash: "#fbefeb"
typography:
  headline:
    fontFamily: "-apple-system, BlinkMacSystemFont, \"Segoe UI\", \"Helvetica Neue\", Arial, sans-serif"
    fontSize: "1.75rem"
    fontWeight: 500
    lineHeight: 1.3
    letterSpacing: "-.025em"
  section:
    fontFamily: "-apple-system, BlinkMacSystemFont, \"Segoe UI\", \"Helvetica Neue\", Arial, sans-serif"
    fontSize: "1.25rem"
    fontWeight: 500
    lineHeight: 1.3
    letterSpacing: "-.015em"
  title:
    fontFamily: "-apple-system, BlinkMacSystemFont, \"Segoe UI\", \"Helvetica Neue\", Arial, sans-serif"
    fontSize: "1.125rem"
    fontWeight: 500
    lineHeight: 1.3
    letterSpacing: "-.01em"
  statement:
    fontFamily: "-apple-system, BlinkMacSystemFont, \"Segoe UI\", \"Helvetica Neue\", Arial, sans-serif"
    fontSize: "1.125rem"
    fontWeight: 400
    lineHeight: 1.7
  body:
    fontFamily: "-apple-system, BlinkMacSystemFont, \"Segoe UI\", \"Helvetica Neue\", Arial, sans-serif"
    fontSize: "1rem"
    fontWeight: 400
    lineHeight: 1.6
  option:
    fontFamily: "-apple-system, BlinkMacSystemFont, \"Segoe UI\", \"Helvetica Neue\", Arial, sans-serif"
    fontSize: ".9375rem"
    fontWeight: 400
    lineHeight: 1.65
  label:
    fontFamily: "-apple-system, BlinkMacSystemFont, \"Segoe UI\", \"Helvetica Neue\", Arial, sans-serif"
    fontSize: ".8125rem"
    fontWeight: 400
    lineHeight: 1.6
  action:
    fontFamily: "-apple-system, BlinkMacSystemFont, \"Segoe UI\", \"Helvetica Neue\", Arial, sans-serif"
    fontSize: ".875rem"
    fontWeight: 500
    lineHeight: 1.3
  filename:
    fontFamily: "ui-monospace, monospace"
    fontSize: ".75rem"
    fontWeight: 400
    lineHeight: 1.6
rounded:
  progress: "2px"
  sm: "0"
  control: "3px"
  surface: "0"
  circle: "50%"
spacing:
  compact: ".5rem"
  inline: ".75rem"
  base: "1rem"
  panel: "1.25rem"
  roomy: "1.5rem"
  reading: "1.75rem"
  gutter: "2rem"
components:
  button-primary:
    backgroundColor: "{colors.accent}"
    textColor: "{colors.white}"
    typography: "{typography.action}"
    rounded: "{rounded.control}"
    padding: ".7rem 1rem"
  button-primary-hover:
    backgroundColor: "{colors.accent-hover}"
  button-primary-active:
    backgroundColor: "#1f394e"
  button-secondary:
    backgroundColor: "{colors.white}"
    textColor: "{colors.ink}"
    typography: "{typography.action}"
    rounded: "{rounded.control}"
    padding: ".7rem 1rem"
  button-text:
    backgroundColor: "transparent"
    textColor: "{colors.muted}"
    typography: "{typography.action}"
    rounded: "{rounded.control}"
    padding: ".7rem 1rem"
  answer-option:
    typography: "{typography.option}"
    rounded: "{rounded.control}"
    padding: ".875rem .5rem"
  answer-option-selected:
    backgroundColor: "{colors.selected}"
  map-item:
    backgroundColor: "{colors.white}"
    rounded: "{rounded.sm}"
  map-item-current:
    backgroundColor: "{colors.selected}"
    textColor: "{colors.accent}"
  catalog:
    backgroundColor: "{colors.white}"
    rounded: "{rounded.surface}"
  question-card:
    backgroundColor: "{colors.white}"
    rounded: "{rounded.surface}"
    padding: "0"
  local-attempt:
    backgroundColor: "{colors.selected}"
    textColor: "{colors.accent}"
    rounded: "{rounded.control}"
    padding: ".85rem 1rem"
  score-panel:
    backgroundColor: "{colors.white}"
    rounded: "{rounded.surface}"
  base-text:
    backgroundColor: "{colors.mist}"
    rounded: "{rounded.control}"
    padding: ".5rem 1rem"
  review-correct:
    backgroundColor: "{colors.correct-wash}"
    textColor: "{colors.correct}"
  review-incorrect:
    backgroundColor: "{colors.incorrect-wash}"
    textColor: "{colors.incorrect}"
---
# Design System: Questões de concurso

## Overview

**Creative North Star: "Responder questões com clareza"**

A interface prioriza a leitura e a resolução de questões. “Questões de concurso” descreve a função do aplicativo. Os controles usam texto direto, tipografia do sistema e cores discretas.

A página é branca e o conteúdo é organizado por alinhamento, espaço e divisórias finas. Azul identifica avanço, foco e escolhas; verde e terracota aparecem na correção após a entrega. Não há caixas envolvendo catálogo, questão ou resultado, nem ilustração decorativa ou assinatura de marca.

**Key Characteristics:**

- Tipografia sem serifa, com enunciados de entrelinha aberta.
- Leitura direta sobre fundo branco, com divisórias no lugar de caixas.
- Estados de escolha, gravação e correção identificados também por texto.

## Colors

Os valores normativos estão no frontmatter e correspondem às propriedades da folha de estilos.

- **Primary:** azul discreto (`accent`) para abrir e avançar, foco e questão atual; `accent-hover` para o ponteiro. `selected` marca uma resposta escolhida e a questão atual. A retomada usa links e divisórias, sem fundo colorido.
- **Correção:** `correct` e `correct-wash` mostram acerto e gabarito; `incorrect` e `incorrect-wash` mostram a resposta escolhida incorreta. Esses estados aparecem depois da entrega.
- **Neutral:** `paper` é o fundo da página; `white`, a área de leitura. `ink` serve ao texto principal, `muted` ao apoio. `line` separa seções, `control-line` delimita controles e `mist` reúne textos-base e estados neutros.

## Typography

A pilha sem serifa do sistema serve a títulos, enunciados e controles. A única família distinta é a monoespaçada dos nomes de arquivo; não há fonte externa.

O frontmatter registra a hierarquia de página, seção, prova, enunciado, corpo, alternativa, rótulo e ação. Em telas de até 44rem, o título de página passa a 1.625rem, o título da prova a 1.125rem e o enunciado a 1.0625rem com entrelinha de 1.7. O enunciado e o texto-base têm medida máxima de 72ch. Contagens, mapa e pontuação usam algarismos tabulares.

## Layout

Cabeçalho, rodapé e conteúdo compartilham um quadro de até 74rem. A margem lateral usa `gutter`, passa a 1.5rem até 60rem e a 1rem até 44rem. O catálogo alinha número, identificação e ação em três colunas; no celular, a ação passa para baixo da identificação. Até 22rem, o cabeçalho da lista pode quebrar em duas linhas.

A prova usa uma barra lateral de 13.5rem e uma área de leitura flexível, com intervalo de 3rem. Até 60rem, a barra passa a 12rem e o intervalo a 1.25rem. Até 44rem, a composição vira uma coluna, a barra perde a fixação e o mapa fica recolhido inicialmente. A questão não tem preenchimento lateral próprio: enunciado, alternativas, navegação e avisos compartilham o mesmo alinhamento.

A pontuação usa quatro colunas e passa a duas no celular; disciplinas passam de duas colunas a uma. A revisão ocupa até 52rem. Linhas e textos longos quebram sem empurrar os controles para fora da tela.

## Elevation & Depth

A hierarquia vem do espaço, dos pesos de texto e de divisórias de 1px. Catálogo, questão, revisão e pontuação ficam diretamente na página, sem contêineres com borda ou sombra. Texto-base e estimativa de nota usam um fundo neutro quando precisam reunir uma explicação. O foco tem contorno azul de 2px.

## Shapes

Os botões têm cantos discretos de 3px. Células do mapa, alternativas e fundos de contexto têm cantos retos; as letras das alternativas permanecem circulares. As setas usam o mesmo traço simples em todos os controles.

## Components

- **Cabeçalho e ajuda:** “Questões de concurso” leva ao índice. “Uso offline” abre as orientações expansíveis no rodapé. Os links e resumos mantêm alvos de pelo menos 2.75rem.
- **Lista de provas:** linhas abertas com divisórias e identificação numerada. Cada entrada reúne título, subtítulo, detalhes expansíveis, nome do JSON, quantidade de questões, duração sugerida e ação de abrir. O arquivo permanece visível com os detalhes fechados.
- **Botões:** abrir, salvar e voltar usam branco com borda de controle; avançar usa azul com texto branco. Entregar usa o botão de texto junto ao aviso de encerramento da edição. Todos têm altura mínima de 2.75rem, foco visível e estados de ponteiro, pressão e desabilitado. As setas permanecem estáticas; a preferência por movimento reduzido também é respeitada.
- **Alternativas:** a linha inteira aciona o rádio. A seleção colore a linha e preenche o círculo da letra em azul. O foco contorna a linha. Clicar fora das alternativas ou pressionar Esc limpa a resposta; links, botões e outros controles preservam a escolha. Duplo clique ou Delete com o rádio em foco descarta ou restaura a alternativa, riscando letra e texto em cinza. Descartar a escolhida deixa a questão em branco; outra resposta escolhida permanece. Um clique simples escolhe e restaura. O descarte persiste por tentativa no mesmo navegador, inclusive offline, e é anunciado às tecnologias assistivas.
- **Mapa e progresso:** células de pelo menos 44px por 44px formam uma grade rolável. Questões em branco são links numéricos sem caixa; respondidas têm fundo neutro e peso maior. A atual tem fundo de seleção, contorno azul e `aria-current="step"`. Abrir ou navegar no mapa revela a célula atual sem deslocar a página. A barra informa a contagem de respondidas também às tecnologias assistivas. No celular, a navegação fecha o mapa e leva o foco ao enunciado.
- **Gravação e retomada:** uma faixa após as alternativas separa o estado de gravação da disponibilidade offline. “Salvo neste aparelho” confirma gravação local; “aguardando envio” indica respostas pendentes; “Sincronizado” confirma recebimento pelo site. “Disponível offline” confirma que a prova foi preparada para uso sem rede. “Provas neste aparelho” mostra título e estado em duas linhas. Entregar sem conexão bloqueia a edição e informa “Aguardando conexão para corrigir”.
- **Resultado:** pontos ponderados recebem destaque azul. A nota padronizada estimada, quando disponível, aparece em painel neutro com a ressalva; quando indisponível, uma seção expansível explica os dados necessários e a fórmula. Disciplinas mostram barras verdes e contagens. Na revisão, “Sua resposta” e “Gabarito” ficam abaixo do texto da alternativa, preservando a largura de leitura.
- **Termos e erros:** “Prova” é o conjunto escolhido no índice; “tentativa” é uma realização até a entrega. “Entregar prova” encerra a edição e solicita a correção; questões “em branco” contam como erro. “Pontos ponderados” e “nota padronizada estimada” são medidas distintas. Erros descrevem a falha e uma ação possível; não afirmam que a resposta foi salva quando a gravação não foi confirmada.

## Do's and Don'ts

### Do

- Do manter “prova”, “questão”, “resposta”, “gabarito” e “resultado” nos controles de estudo.
- Do acompanhar os estados de correção com os rótulos “Sua resposta”, “Gabarito”, “Acertou”, “Errou” ou “Em branco”.
- Do preservar foco visível, alvos de pelo menos 44px e a leitura das alternativas em telas estreitas.
- Do distinguir gravação local, sincronização e disponibilidade offline nas mensagens.

### Don't

- Don't reintroduzir o nome rejeitado, vinho, ilustrações de encadernação ou outra metáfora de marca.
- Don't indicar acerto, erro ou gabarito antes da entrega.
- Don't usar a cor como única identificação de uma escolha ou de um resultado.
- Don't apresentar uma estimativa de nota como resultado oficial ou confirmar gravação sem sucesso comprovado.
