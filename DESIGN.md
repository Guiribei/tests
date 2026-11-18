---
name: Caderno
description: Sistema visual da Computadora para escolher, responder e revisar provas de estudo individual.
colors:
  wine: "#911a46"
  wine-dark: "#701334"
  rose: "#f7e9ed"
  correct: "#245a4a"
  correct-wash: "#e7efeb"
  paper: "#f9faf9"
  white: "#fff"
  ink: "#262933"
  muted: "#525e65"
  line: "#c6ccce"
  mist: "#e9edef"
typography:
  display:
    fontFamily: "Baskervville, Georgia, serif"
    fontSize: "clamp(3.3rem, 6.25vw, 6rem)"
    fontWeight: 400
    lineHeight: 1.03
    letterSpacing: "-.04em"
  headline:
    fontFamily: "Baskervville, Georgia, serif"
    fontSize: "2.15rem"
    fontWeight: 400
    lineHeight: 1.2
    letterSpacing: "-.032em"
  title:
    fontFamily: "Baskervville, Georgia, serif"
    fontSize: "clamp(1.55rem, 2.55vw, 2.65rem)"
    fontWeight: 400
    lineHeight: 1.12
    letterSpacing: "-.034em"
  statement:
    fontFamily: "Baskervville, Georgia, serif"
    fontSize: "1.55rem"
    fontWeight: 400
    lineHeight: 1.45
    letterSpacing: "-.012em"
  body:
    fontFamily: '"Helvetica Neue", Helvetica, Arial, sans-serif'
    fontSize: "1rem"
    fontWeight: 400
    lineHeight: 1.5
  label:
    fontFamily: '"Helvetica Neue", Helvetica, Arial, sans-serif'
    fontSize: ".78rem"
    fontWeight: 550
    lineHeight: 1.4
    letterSpacing: ".055em"
  action:
    fontFamily: '"Helvetica Neue", Helvetica, Arial, sans-serif'
    fontSize: ".89rem"
    fontWeight: 550
    lineHeight: 1.2
  filename:
    fontFamily: "ui-monospace, monospace"
    fontSize: ".76rem"
rounded:
  square: "0"
spacing:
  gutter: "clamp(1.2rem, 3.2vw, 3.5rem)"
  grid-gap: ".25rem"
  inline: "1rem"
  panel: "1.4rem"
components:
  button-primary:
    backgroundColor: "{colors.wine}"
    textColor: "{colors.white}"
    typography: "{typography.action}"
    rounded: "{rounded.square}"
    padding: ".7rem 1rem"
  button-primary-hover:
    backgroundColor: "{colors.wine-dark}"
  button-secondary:
    backgroundColor: "{colors.white}"
    textColor: "{colors.ink}"
    typography: "{typography.action}"
    rounded: "{rounded.square}"
    padding: ".7rem 1rem"
  button-finish:
    backgroundColor: "transparent"
    textColor: "{colors.wine}"
    typography: "{typography.action}"
    rounded: "{rounded.square}"
    padding: ".7rem 1rem"
  header-nav-current:
    backgroundColor: "{colors.ink}"
    textColor: "{colors.white}"
  answer-option:
    padding: ".9rem .7rem"
    typography: "{typography.body}"
  answer-option-selected:
    backgroundColor: "{colors.rose}"
  map-item:
    backgroundColor: "{colors.white}"
    height: "2.05rem"
  exam-details:
    textColor: "{colors.ink}"
  score-note-panel:
    backgroundColor: "{colors.mist}"
    padding: "{spacing.panel}"
---

# Sistema de design do Caderno

## Overview

**Creative North Star: "Caderno de provas"**

O Caderno, da Computadora, organiza provas como folhas de uma mesma publicação. A linha horizontal separa registros, a serifa dá espaço ao enunciado e a numeração orienta quem vai e volta entre questões.

A referência Mecha aparece no papel de tom frio, no contraste entre serifa e grotesca e no desenho técnico da encadernação. A ilustração acompanha a leitura do índice. Já a prova e o resultado dão o espaço central às perguntas, às alternativas e aos números.

**Key Characteristics:**

- Superfícies claras com filetes que alinham e separam registros.
- Baskervville para conteúdo de leitura e grotesca para controles e metadados.
- Vinho para ação e escolha, verde para correção.
- Ícone de folhas sobrepostas ao lado de Caderno e da assinatura Computadora.

## Colors

O vinho pontua uma página dominada por cinzas frios. Cada fundo colorido de estado acompanha uma marca, uma letra ou um texto legível.

- **Vinho e vinho escuro** identificam a ação principal e sua passagem do ponteiro. O vinho também marca foco, questão atual, alternativa escolhida e pontos ponderados.
- **Rosa claro** acompanha a escolha durante a prova, os links de retomada e a resposta incorreta na revisão.
- **Verde e verde claro** identificam o gabarito e os acertos por disciplina no resultado.
- **Papel e branco** formam a página e os controles. **Tinta** serve ao texto principal e à aba atual. **Cinza de leitura** serve aos metadados. **Cinza de filete** separa itens e **névoa** reúne instruções, texto-base e nota padronizada.

**Regra da correção.** Reserve o verde dos acertos para a prova entregue. Durante a tentativa, a letra preenchida em vinho confirma apenas a escolha.

## Typography

A Baskervville regular é servida pelo aplicativo, com Georgia e a família serifa como alternativas. Ela distingue o nome da prova, o enunciado e o valor da nota do texto que explica cada um. Helvetica Neue, Helvetica e Arial compõem a pilha sem serifa dos controles, alternativas e metadados.

O título da página tem escala fluida. O enunciado usa uma entrelinha mais aberta e uma medida máxima de 63ch, preservada na revisão. O texto-base usa até 72ch. A numeração do índice, do mapa e da pontuação usa algarismos tabulares. O nome do JSON na ficha usa monoespaçada e conserva a grafia do arquivo.

## Layout

O cabeçalho, o conteúdo e o rodapé compartilham a margem lateral fluida. O conjunto ocupa até 100rem. No índice amplo, cada linha alinha número, identificação e ação. A ilustração ocupa a faixa lateral, de até 29% da grade. Em telas de até 58rem, a ação passa para a coluna do título. Até 44rem, a ilustração diminui e acompanha o título. As orientações passam para depois da lista.

Na prova, o mapa fica em uma coluna lateral e o corpo ocupa até 51rem. Até 44rem, o mapa recolhível passa para o fluxo acima do enunciado. No resultado, as quatro medidas da pontuação passam a duas colunas e as disciplinas passam a uma. A revisão corre em uma coluna de até 56rem.

## Elevation & Depth

A mudança de superfície vem da cor de fundo e de filetes de 1px. A linha escura abre seções e a linha clara separa suas entradas. No mapa, um traço interno reforça a borda da questão atual.

## Shapes

Botões, alternativas, células do mapa e painéis têm cantos retos. A letra de cada alternativa fica em um quadrado de 1,8rem. O ícone da marca representa folhas sobrepostas em tinta escura. As setas de navegação usam traço fino e acompanham a cor do texto.

## Components

- **Cabeçalho.** O ícone, o nome Caderno e a assinatura Computadora compartilham o link inicial. No celular, a assinatura passa para baixo do nome. As abas são contíguas e a aba atual usa tinta escura com texto branco.
- **Linha de prova.** O título, o subtítulo e o nome do arquivo permitem identificar a prova antes de abrir a ficha. A ficha expansível reúne os metadados. A quantidade de questões e a duração ficam junto à ação de abrir.
- **Botões.** Abrir e avançar usam preenchimento vinho. Salvar e voltar usam branco com filete cinza. Finalizar usa contorno vinho e aparece junto ao aviso de bloqueio das respostas. A seta do botão principal avança 2px em 160ms ao receber foco ou ponteiro. A preferência por movimento reduzido elimina a transição.
- **Alternativas.** A linha inteira aciona o rádio. A escolha colore a linha de rosa e preenche o quadrado da letra em vinho. O foco de teclado contorna a linha. Na revisão, os rótulos “Sua resposta” e “Gabarito” explicam as cores.
- **Mapa de questões.** Células numeradas formam uma grade rolável. Uma questão respondida ganha peso e fundo cinza. A atual recebe contorno vinho e identificação de passo atual para tecnologias assistivas.
- **Resultado.** Filetes escuros abrem as medidas de pontuação e a primeira recebe destaque em vinho. Um painel cinza reúne a nota padronizada ou a razão de sua indisponibilidade. Barras verdes e contagens acompanham o desempenho por disciplina.

## Do's and Don'ts

### Aplicar

- Mantenha o nome do arquivo legível com a ficha fechada.
- Preserve o alinhamento entre número, título e ação quando a linha quebrar.
- Acompanhe as cores de correção com os rótulos de gabarito e resposta escolhida.

### Evitar

- Evite dar pistas de acerto ao destacar a alternativa durante a prova.
- Evite usar o espaço do enunciado para instruções secundárias quando elas cabem no texto-base expansível ou no mapa.
