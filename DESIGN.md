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
  header-help:
    textColor: "{colors.muted}"
  answer-option:
    padding: ".9rem .7rem"
    typography: "{typography.body}"
  answer-option-selected:
    backgroundColor: "{colors.rose}"
  map-item:
    backgroundColor: "{colors.white}"
    minHeight: "2.75rem"
  exam-details:
    textColor: "{colors.ink}"
  score-note-panel:
    backgroundColor: "{colors.mist}"
    padding: "{spacing.panel}"
  score-explanation:
    textColor: "{colors.muted}"
---

# Sistema de design do Caderno

## Overview

**Creative North Star: "Caderno de provas"**

Plataforma de resolução de questões, para não deixar a vista cansada e permitir horas a fio de estudo concentrado

A referência Mecha aparece no papel de tom frio, no contraste entre serifa e grotesca e no desenho técnico da encadernação. A ilustração acompanha a leitura do índice. Já a prova e o resultado dão o espaço central às perguntas, às alternativas e aos números.

**Key Characteristics:**

- Superfícies claras com filetes que alinham e separam registros.
- Baskervville para conteúdo de leitura e grotesca para controles e metadados.


## Colors

Cada fundo colorido de estado acompanha uma marca, uma letra ou um texto legível.

- **Cores principais** identificam a ação principal e sua passagem do ponteiro. O [cor] também marca foco, questão atual, alternativa escolhida e pontos ponderados.
- **Cor secundária** acompanha a escolha durante a prova, os links de retomada e a resposta incorreta na revisão.
- **Verde e verde claro** identificam o gabarito e os acertos por disciplina no resultado.
- **Papel e branco** formam a página e os controles. **Tinta** serve ao texto principal. **Cinza de leitura** serve aos metadados. **Cinza de filete** separa itens e **névoa** reúne texto-base e nota padronizada disponível.

**Regra da correção.** Reserve o verde dos acertos para a prova entregue. Durante a tentativa, a letra preenchida em [cor] confirma apenas a escolha.

## Typography

A Baskervville regular é servida pelo aplicativo, com Georgia e a família serifa como alternativas. Ela distingue o nome da prova, o enunciado e o valor da nota do texto que explica cada um. Helvetica Neue, Helvetica e Arial compõem a pilha sem serifa dos controles, alternativas e metadados.

O título da página tem escala fluida. O enunciado usa uma entrelinha mais aberta e uma medida máxima de 63ch, preservada na revisão. O texto-base usa até 72ch. A numeração do índice, do mapa e da pontuação usa algarismos tabulares. O nome do JSON na ficha usa monoespaçada e conserva a grafia do arquivo.

## Layout

O cabeçalho, o conteúdo e o rodapé compartilham a margem lateral fluida. O conjunto ocupa até 100rem. No índice amplo, cada linha alinha número, identificação e ação. A ilustração ocupa a faixa lateral, de até 29% da grade, diretamente sobre o papel. Em telas de até 58rem, a ação passa para a coluna do título. Até 44rem, a ilustração diminui e acompanha o título. As orientações passam para depois da lista, separadas por um filete.


## Elevation & Depth

A mudança de superfície vem da cor de fundo e de filetes de 1px. A linha escura abre seções e a linha clara separa suas entradas. No mapa, um traço interno reforça a borda da questão atual.

## Shapes

Botões, alternativas, células do mapa e painéis têm cantos retos. A letra de cada alternativa fica em um quadrado de 1,8rem. O ícone da marca representa folhas sobrepostas em tinta escura. As setas de navegação usam traço fino e acompanham a cor do texto.

## Components

- **Cabeçalho.** O ícone, o nome Caderno  compartilham o link inicial, que substitui a aba “Provas”. No celular, a assinatura passa para baixo do nome. “Uso offline” permanece como link de ajuda discreto.
- **Linha de prova.** O título, o subtítulo e o nome do arquivo permitem identificar a prova antes de abrir a ficha. A ficha expansível reúne os metadados e conserva um alvo mínimo de 44px de altura. No celular, o nome do arquivo ocupa sua própria linha. A quantidade de questões e a duração ficam junto à ação de abrir.
- **Botões.** Abrir e avançar usam preenchimento [cor]. Salvar e voltar usam branco com filete cinza. “Entregar prova” usa contorno [cor] e aparece junto ao aviso de encerramento da edição e da contagem de questões em branco como erro. A seta do botão principal avança 2px em 160ms ao receber foco ou ponteiro. A preferência por movimento reduzido elimina a transição.
- **Alternativas.** A linha inteira aciona o rádio. A escolha colore a linha de rosa e preenche o quadrado da letra em [cor]. Clicar fora das alternativas ou pressionar Esc limpa a escolha e o destaque, deixando a questão em branco; links, botões e controles de navegação preservam a resposta. A remoção é salva no aparelho e sincronizada com o site. O duplo clique descarta ou restaura uma alternativa: a letra e o texto ficam riscados em cinza de leitura. Descartar a resposta escolhida deixa a questão em branco; outra resposta escolhida é preservada. Um clique simples escolhe e restaura a alternativa. O descarte fica salvo por tentativa no mesmo navegador, inclusive offline. A tecla Delete oferece a mesma ação com o rádio em foco, e tecnologias assistivas recebem o estado “Descartada”. O foco de teclado contorna a linha. Na revisão, os rótulos “Sua resposta” e “Gabarito” aparecem abaixo do texto da alternativa para preservar sua largura de leitura.
- **Mapa de questões.** Células de pelo menos 44px por 44px formam uma grade rolável, com espaço para o foco de teclado. Uma questão respondida ganha peso e fundo cinza. A atual recebe fundo rosa, contorno [cor] e identificação de passo atual para tecnologias assistivas. Ao navegar ou reabrir o mapa, a grade revela a questão atual sem deslocar a página. A barra de progresso informa a contagem de questões respondidas às tecnologias assistivas.
- **Progresso e gravação.** A contagem de respostas acompanha a barra em todos os tamanhos. O número da questão aparece somente junto ao enunciado. A gravação e a disponibilidade offline compartilham uma faixa após as alternativas, mantendo os dois estados distintos. O código do cargo continua acessível na ficha da prova.
- **Retomada.** Cada link em “Provas neste aparelho” separa o título da prova do estado de gravação ou entrega em duas linhas, sobre fundo rosa.
- **Resultado.** Filetes escuros abrem as medidas de pontuação e a primeira recebe destaque em [cor]. A nota padronizada estimada, quando disponível, ocupa um painel cinza com sua ressalva. Quando faltam dados, “Nota padronizada indisponível” abre a explicação e a fórmula sob demanda. Barras verdes e contagens acompanham o desempenho por disciplina.

## Terminologia

| Termo | Significado na interface |
| --- | --- |
| Prova | Caderno de questões escolhido no índice. |
| Tentativa | Uma realização da prova, do início à entrega. |
| Entregar prova | Encerrar a edição das respostas e solicitar a correção. Usar o mesmo verbo no botão, no aviso e na confirmação. |
| Em branco | Questão sem alternativa escolhida. Conta como erro na entrega. |
| Salvo neste aparelho | Gravação local confirmada. “Aguardando envio” indica que as respostas ainda precisam ser enviadas ao site. |
| Sincronizado | O site recebeu as respostas. Não é sinônimo de apenas salvar no aparelho. |
| Disponível offline | A interface e as questões foram preparadas para uso sem conexão. |
| Aguardando conexão para corrigir | A entrega foi registrada no aparelho e a edição está encerrada; a correção depende do envio ao site. |
| Pontos ponderados | Pontuação bruta que aplica os pesos da prova aos acertos. Não é a nota padronizada. |
| Nota padronizada estimada | Cálculo com as tentativas concluídas neste site. O mínimo atingido pela estimativa não é uma habilitação oficial. |

Mensagens de erro descrevem o que falhou e a próxima ação possível. Não garantem que a última resposta foi salva quando o sistema não pode confirmar a gravação.

## Do's and Don'ts

### Aplicar

- Mantenha o nome do arquivo legível com a ficha fechada.
- Preserve o alinhamento entre número, título e ação quando a linha quebrar.
- Acompanhe as cores de correção com os rótulos de gabarito e resposta escolhida.

### Evitar

- Evite dar pistas de acerto ao destacar a alternativa durante a prova.
- Evite usar o espaço do enunciado para instruções secundárias quando elas cabem no texto-base expansível ou no mapa.
