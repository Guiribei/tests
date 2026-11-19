import {attempt, listAttempts, forgetAttempt} from "./local-store.js?v=1";

const $ = selector => document.querySelector(selector);
const layout = $(".exam-layout");
const library = $(".local-attempts");
const questionMap = $(".map-disclosure");
const compact = window.matchMedia("(max-width: 44rem)");
if (questionMap) {
  questionMap.open = !compact.matches;
  compact.addEventListener("change", () => { questionMap.open = !compact.matches; });
  questionMap.addEventListener("toggle", revealCurrentQuestion);
}

function revealCurrentQuestion() {
  if (!questionMap?.open) return;
  const map = $(".question-map");
  const item = map.querySelector('[aria-current="step"]');
  if (!item) return;
  const bounds = map.getBoundingClientRect();
  const position = item.getBoundingClientRect();
  if (position.top < bounds.top + 6 || position.bottom > bounds.bottom - 6) {
    map.scrollTop += position.top - bounds.top - (map.clientHeight - position.height) / 2;
  }
}
function revealOfflineHelp() {
  if (location.hash === "#uso-offline") $("#uso-offline").open = true;
}
window.addEventListener("hashchange", revealOfflineHelp);
document.querySelector('.mast-nav a[href="#uso-offline"]').addEventListener("click", () => {
  $("#uso-offline").open = true;
});
revealOfflineHelp();
const offlineShell = layout?.dataset.attempt === "";
const route = location.pathname.match(/^\/tentativas\/([a-f0-9]+)\/questoes\/(\d+)$/);
const id = route?.[1];
const prefix = key => "/tentativas/" + key;
const questionURL = (key, n) => prefix(key) + "/questoes/" + n;
const resultURL = key => prefix(key) + "/resultado";
const channel = "BroadcastChannel" in window ? new BroadcastChannel("caderno-attempts") : null;
const warnings = new Map();
let current;
let number = Number(route?.[2]);
let syncing = false;
let rerun = false;
let preparing = false;
let localWrites = 0;
const discardKey = id ? "caderno-discarded:" + id : null;
let discarded = readDiscarded();
let alternativeClick;
let onlineSaveTimer;
let onlineChoice = $(".answer-form input:checked")?.value || "";

function readDiscarded() {
  if (!discardKey) return {};
  try {
    const value = JSON.parse(localStorage.getItem(discardKey));
    return value && typeof value === "object" && !Array.isArray(value) ? value : {};
  } catch {
    return {};
  }
}

function isDiscarded(letter) {
  return Array.isArray(discarded[number]) && discarded[number].includes(letter);
}

function renderDiscarded() {
  for (const label of document.querySelectorAll(".answer-form .alternative")) {
    const excluded = isDiscarded(label.querySelector("input").value);
    label.classList.toggle("discarded", excluded);
    label.querySelector(".alternative-state").textContent = excluded ? "Descartada" : "";
  }
}

function setDiscarded(letter, excluded) {
  if (isDiscarded(letter) === excluded) return;
  const letters = new Set(Array.isArray(discarded[number]) ? discarded[number] : []);
  if (excluded) letters.add(letter);
  else letters.delete(letter);
  if (letters.size) discarded[number] = [...letters];
  else delete discarded[number];
  renderDiscarded();
  try {
    localStorage.setItem(discardKey, JSON.stringify(discarded));
  } catch {
    status("Não foi possível salvar os descartes neste aparelho. Eles serão perdidos ao fechar a página.");
  }
}

if (layout && id) {
  $("#alternative-help").hidden = false;
  renderDiscarded();
}

window.addEventListener("storage", event => {
  if (discardKey && (event.key === discardKey || event.key === null)) {
    discarded = readDiscarded();
    renderDiscarded();
  }
});

function status(message) {
  const el = $("#offline-status");
  if (el) {
    el.hidden = false;
    el.textContent = message;
  }
}

function cancelOnlineSave() {
  clearTimeout(onlineSaveTimer);
  onlineSaveTimer = undefined;
}

async function flushOnlineSave() {
  cancelOnlineSave();
  const form = $(".answer-form");
  const choice = form.querySelector("input:checked")?.value || "";
  try {
    const response = await fetch(form.action, {
      method: "POST", body: new FormData(form), signal: AbortSignal.timeout(12000)
    });
    if (!response.ok) throw new Error("Resposta não salva");
    const finish = $(".finish-panel");
    finish.dataset.answered = Number(finish.dataset.answered) + Number(Boolean(choice)) - Number(Boolean(onlineChoice));
    onlineChoice = choice;
    return true;
  } catch {
    status("Não foi possível salvar a resposta. Confira a conexão e use Salvar resposta para tentar novamente.");
    return false;
  }
}

function storageError() {
  if (current) render(current);
  if (offlineShell && !current && library) {
    library.hidden = false;
    $(".local-message").textContent = "O navegador não conseguiu abrir as provas salvas. Conecte-se à internet e abra a prova pelo índice.";
  }
  status("Não foi possível salvar neste aparelho. Confira o espaço e as permissões de armazenamento do navegador e tente novamente.");
}

function answers(row) {
  return Object.fromEntries(Object.entries({...row.answers, ...row.pending})
    .filter(([, value]) => value !== ""));
}

function delivery(row) {
  if (row.finished) return "Resultado disponível";
  if (warnings.has(row.id)) return warnings.get(row.id);
  if (row.finishing) return "Aguardando conexão para corrigir";
  if (Object.keys(row.pending).length) return "Salvo neste aparelho · aguardando envio";
  return "Sincronizado";
}

function render(row) {
  current = row;
  const q = row.questions[number - 1];
  if (!q) return;
  layout.hidden = false;
  if (library) library.hidden = true;
  const choices = answers(row);
  const count = Object.keys(choices).length;
  const total = row.questions.length;
  const locked = row.finishing || row.finished;
  document.title = "Questão " + number + " · Caderno";
  $(".sidebar-title h2").textContent = row.title;
  $(".progress-copy strong").textContent = count + " de " + total;
  $(".progress-copy span").textContent = count === 1 ? "respondida" : "respondidas";
  $(".progress-track").setAttribute("aria-valuemax", total);
  $(".progress-track").setAttribute("aria-valuenow", count);
  $(".progress-track span").style.width = (100 * count / total) + "%";
  $(".number-label").textContent = "Questão " + number;
  $(".question-meta span:last-child").textContent = q.paper + " · " + q.discipline;
  $(".base-text").hidden = !q.baseText;
  $(".base-text summary").textContent = "Texto-base " + (q.baseID || "");
  $(".base-text p").textContent = q.baseText || "";
  $(".question-card h1").textContent = q.statement;

  const map = $(".question-map");
  if (map.children.length !== total) {
    map.replaceChildren(...row.questions.map(item => {
      const link = document.createElement("a");
      link.className = "map-item";
      link.href = questionURL(id, item.number);
      link.textContent = item.number;
      return link;
    }));
  }
  for (const link of map.children) {
    const n = Number(link.textContent);
    link.classList.toggle("current", n === number);
    if (n === number) link.setAttribute("aria-current", "step");
    else link.removeAttribute("aria-current");
    link.classList.toggle("answered", Boolean(choices[n]));
    link.setAttribute("aria-label", choices[n] ? `Questão ${n}, respondida` : `Questão ${n}, em branco`);
  }
  revealCurrentQuestion();
  for (const [index, label] of document.querySelectorAll(".alternative").entries()) {
    const input = label.querySelector("input");
    input.checked = choices[number] === input.value;
    label.querySelector(".alternative-text").textContent = q.alternatives[index];
  }
  renderDiscarded();
  $("#alternative-help").hidden = locked;
  $(".answer-form").action = prefix(id) + "/respostas/" + number;
  $(".finish-panel").action = prefix(id) + "/finalizar";
  $(".answer-form fieldset").disabled = locked;
  $(".answer-form button").hidden = true;
  $(".finish-panel button").disabled = locked;
  for (const [selector, target] of [["[data-previous]", number - 1], ["[data-next]", number + 1]]) {
    const link = $(selector);
    link.hidden = target < 1 || target > total;
    if (!link.hidden) link.href = questionURL(id, target);
  }
  status("Disponível offline");
  $("#save-state").textContent = localWrites ? "Gravando no aparelho…" : delivery(row);
}

async function renderLibrary(message = "") {
  if (!library || current) return;
  const rows = await listAttempts();
  library.hidden = !offlineShell && !rows.length;
  $(".local-links").replaceChildren(...rows.map(row => {
    const link = document.createElement("a");
    link.href = row.finished ? resultURL(row.id) : questionURL(row.id, row.position);
    const title = document.createElement("span");
    title.className = "local-attempt-title";
    title.textContent = row.title;
    const state = document.createElement("span");
    state.className = "local-attempt-state";
    state.textContent = delivery(row);
    link.append(title, state);
    return link;
  }));
  $(".local-message").textContent = message || (offlineShell && !rows.length
    ? "Nenhuma prova salva neste aparelho. Conecte-se à internet e abra uma prova pelo índice para usar offline." : "");
}

async function updateView(key) {
  if (key === id && current) {
    const row = await attempt(id);
    if (!row || row.finished) location.replace(resultURL(id));
    else render(row);
  } else if (library && !id) {
    await renderLibrary();
  }
}

async function saveLocal(change) {
  localWrites++;
  $("#save-state").textContent = "Gravando no aparelho…";
  let row;
  try {
    row = await attempt(id, saved => {
      if (saved && !saved.finishing && !saved.finished) change(saved);
      return saved;
    });
    if (!row) throw new Error("Tentativa indisponível");
  } catch {
    localWrites--;
    storageError();
    return;
  }
  localWrites--;
  render(row);
  channel?.postMessage(id);
  sync().catch(storageError);
}

function mergeSheet(row, sheet) {
  if (sheet.revision < row.revision) return row;
  row.revision = sheet.revision;
  row.answers = sheet.answers;
  row.finished = sheet.finished;
  // An acknowledgement can arrive after the user has changed the answer again.
  for (const [n, value] of Object.entries(row.pending)) {
    if (sheet.finished || (sheet.answers[n] || "") === value) delete row.pending[n];
  }
  return row;
}

async function send(row) {
  const res = await fetch(prefix(row.id) + "/sincronizar", {
    method: "POST", headers: {"Content-Type": "application/json"},
    signal: AbortSignal.timeout(12000),
    body: JSON.stringify({revision: row.revision, finish: row.finishing,
      answers: row.finishing ? answers(row) : row.pending})
  });
  if (!res.ok && res.status !== 409) {
    if (res.status === 404) warnings.set(row.id, "Não foi possível enviar as respostas: esta tentativa não foi encontrada no site. As respostas continuam neste aparelho. Volte ao índice para abrir uma prova.");
    else if (res.status < 500) warnings.set(row.id, "O site não aceitou o envio. As respostas continuam neste aparelho. Confira a conexão e reabra a prova para tentar novamente.");
    await updateView(row.id);
    return;
  }
  const sheet = await res.json();
  warnings.delete(row.id);
  const saved = await attempt(row.id, latest => latest && mergeSheet(latest, sheet));
  if (!saved) return;
  channel?.postMessage(row.id);
  await updateView(row.id);
  // Rebase an outdated batch on the server's revision, keeping newer local edits.
  if (!saved.finished && (saved.finishing || Object.keys(saved.pending).length)) rerun = true;
}

async function sync() {
  if (!navigator.onLine) return;
  if (syncing) {
    rerun = true;
    return;
  }
  syncing = true;
  try {
    do {
      rerun = false;
      for (const row of await listAttempts()) {
        if (row.finished || (!row.finishing && !Object.keys(row.pending).length)) continue;
        try {
          await send(row);
        } catch {
          // The request or its acknowledgement was lost. The durable batch
          // is retried on reconnection, on reopening, and while the tab is open.
        }
      }
    } while (rerun && navigator.onLine);
  } finally {
    syncing = false;
  }
}

async function downloadPacket() {
  const res = await fetch(prefix(id) + "/pacote", {cache: "no-store", signal: AbortSignal.timeout(15000)});
  if (!res.ok) throw new Error("Pacote indisponível");
  return res.json();
}

async function prepare() {
  if (!id || !layout || preparing || current) return;
  preparing = true;
  try {
    if (!("serviceWorker" in navigator) || !window.isSecureContext) {
      status("Responda com conexão. Para usar offline, abra o site por HTTPS.");
      return;
    }
    status("Preparando a prova para responder offline…");
    if (!navigator.serviceWorker.controller) {
      await navigator.serviceWorker.register("/sw.js");
      await Promise.race([navigator.serviceWorker.ready, new Promise((_, reject) =>
        setTimeout(() => reject(new Error("Tempo de instalação excedido")), 15000))]);
    }
    let row = await attempt(id);
    if (!row) {
      const packet = await downloadPacket();
      row = await attempt(id, existing => existing || {
        ...packet, id, position: number, pending: {}, finishing: false
      });
    }
    if (row.finished) {
      location.replace(resultURL(id));
      return;
    }
    if (new URLSearchParams(location.search).has("retomar")) number = row.position;
    if (!row.questions[number - 1]) number = row.position;
    history.replaceState(null, "", questionURL(id, number));
    render(row);
    if (row.position !== number) await rememberPosition();
    // An online HTML form can also have saved answers since this download.
    if (!offlineShell && Number(layout.dataset.revision) > row.revision) {
      try {
        const packet = await downloadPacket();
        await attempt(id, saved => saved && mergeSheet(saved, packet));
        await updateView(id);
      } catch {
        // Keep the downloaded version if the signal drops during reconciliation.
      }
    }
    sync().catch(storageError);
  } catch {
    if (offlineShell) await renderLibrary("Conecte-se à internet e reabra esta prova para terminar de baixá-la.");
    else status("A prova ainda não está pronta para uso offline. Continue com conexão e reabra a prova para tentar novamente.");
  } finally {
    preparing = false;
  }
}

async function rememberPosition() {
  const position = number;
  await attempt(id, row => {
    if (row) row.position = position;
    return row;
  });
}

function navigate(n, push) {
  if (!current?.questions[n - 1]) return;
  number = n;
  if (push) history.pushState(null, "", questionURL(id, n));
  render(current);
  $(".base-text").open = true;
  if (compact.matches) questionMap.open = false;
  $(".question-card h1").focus({preventScroll: true});
  window.scrollTo(0, 0);
  rememberPosition().catch(storageError);
}

document.addEventListener("click", async event => {
  const link = event.target.closest(".exam-layout a");
  if (!link || event.button !== 0 || event.ctrlKey || event.metaKey || event.shiftKey || event.altKey) return;
  if (!current) {
    if (onlineSaveTimer !== undefined) {
      event.preventDefault();
      if (await flushOnlineSave()) location.assign(link.href);
    }
    return;
  }
  if (!link.pathname.startsWith(prefix(id) + "/questoes/")) return;
  event.preventDefault();
  navigate(Number(link.pathname.split("/").pop()), true);
});

window.addEventListener("popstate", () => navigate(Number(location.pathname.split("/").pop()), false));

function clearAnswer() {
  const input = $(".answer-form input[name=choice]:checked");
  if (!input || input.matches(":disabled")) return false;
  input.checked = false;
  input.dispatchEvent(new Event("change", {bubbles: true}));
  return true;
}

// Remember the answer before the first click changes the radio. A double click
// must preserve a different answer, even if its first click has already saved.
document.addEventListener("mousedown", event => {
  const label = event.target.closest(".answer-form .alternative");
  if (!label || event.button !== 0) return;
  if (event.detail > 1) {
    event.preventDefault();
    return;
  }
  alternativeClick = {
    number, letter: label.querySelector("input").value,
    previous: $(".answer-form input:checked")?.value || "",
    discarded: isDiscarded(label.querySelector("input").value)
  };
});

function toggleDiscard(label, before) {
  const input = label.querySelector("input");
  if (input.matches(":disabled") || current?.finishing || current?.finished) return;
  cancelOnlineSave();
  const excluded = !(before ? before.discarded : isDiscarded(input.value));
  const previous = before ? before.previous : ($(".answer-form input:checked")?.value || "");
  const choice = excluded && previous === input.value ? "" : previous;
  setDiscarded(input.value, excluded);
  for (const radio of document.querySelectorAll(".answer-form input[name=choice]")) {
    radio.checked = radio.value === choice;
  }
  $("#discard-state").textContent = `Alternativa ${input.value} ${excluded ? "descartada" : "restaurada"}.`;
  if (current) {
    const n = number;
    saveLocal(row => { row.pending[n] = choice; });
  } else if (choice !== previous) {
    input.form.requestSubmit();
  }
}

document.addEventListener("dblclick", event => {
  const label = event.target.closest(".answer-form .alternative");
  if (!label || event.button !== 0) return;
  event.preventDefault();
  const before = alternativeClick?.number === number && alternativeClick.letter === label.querySelector("input").value
    ? alternativeClick : null;
  toggleDiscard(label, before);
  alternativeClick = null;
});

document.addEventListener("click", event => {
  if (event.defaultPrevented || event.button !== 0 || event.ctrlKey || event.metaKey || event.shiftKey || event.altKey) return;
  if (event.target.closest(".alternative, a, button, input, select, textarea, summary, [role=button], [contenteditable], .finish-panel")) return;
  if (window.getSelection()?.isCollapsed === false) return;
  clearAnswer();
});

document.addEventListener("keydown", event => {
  if (event.key === "Escape" && !event.defaultPrevented && clearAnswer()) event.preventDefault();
  if (event.key === "Delete" && event.target.matches(".answer-form input[name=choice]")) {
    event.preventDefault();
    toggleDiscard(event.target.closest(".alternative"));
  }
});

document.addEventListener("change", event => {
  const input = event.target;
  if (!input.matches(".answer-form input[name=choice]")) return;
  if (input.checked) setDiscarded(input.value, false);
  if (!current) {
    // Give the second click time to discard instead of submitting the form.
    cancelOnlineSave();
    onlineSaveTimer = setTimeout(() => {
      onlineSaveTimer = undefined;
      input.form.requestSubmit();
    }, 500);
    return;
  }
  const n = number;
  const value = input.checked ? input.value : "";
  saveLocal(row => { row.pending[n] = value; });
});

document.addEventListener("submit", async event => {
  if (event.target.matches(".answer-form")) {
    cancelOnlineSave();
    if (current) event.preventDefault();
    return;
  }
  const form = event.target.closest("[data-confirm-finish]");
  if (!form) return;
  const pendingOnline = !current && onlineSaveTimer !== undefined;
  if (current || pendingOnline) event.preventDefault();
  if (pendingOnline && !await flushOnlineSave()) return;
  const blank = current ? current.questions.length - Object.keys(answers(current)).length
    : Number(form.dataset.total) - Number(form.dataset.answered);
  const prompt = blank === 0
    ? "Entregar a prova? Você não poderá alterar as respostas após a entrega."
    : blank === 1
      ? "Entregar a prova com 1 questão em branco? Ela contará como erro e você não poderá alterar as respostas após a entrega."
      : `Entregar a prova com ${blank} questões em branco? Elas contarão como erro e você não poderá alterar as respostas após a entrega.`;
  if (!window.confirm(prompt)) event.preventDefault();
  else if (current) saveLocal(row => { row.finishing = true; });
  else if (pendingOnline) form.submit();
});

channel?.addEventListener("message", event => { updateView(event.data).catch(storageError); });

async function resume() {
  if (!window.indexedDB) {
    if (offlineShell && library) {
      library.hidden = false;
      $(".local-message").textContent = "Não é possível salvar provas neste navegador. Conecte-se à internet e abra uma prova pelo índice para responder com conexão.";
    }
    status("Não é possível salvar a prova neste navegador. Responda com conexão.");
    return;
  }
  const confirmed = document.querySelector("[data-result-attempt]")?.dataset.resultAttempt;
  if (confirmed) await forgetAttempt(confirmed);
  if (id) await prepare();
  else await renderLibrary(offlineShell && location.pathname.endsWith("/resultado")
    ? "Conecte-se à internet para receber a correção e ver o resultado." : "");
  if (current) await updateView(id);
  await sync();
}

window.addEventListener("online", () => resume().catch(storageError));
document.addEventListener("visibilitychange", () => {
  if (document.visibilityState === "visible") resume().catch(storageError);
});
window.addEventListener("pageshow", event => { if (event.persisted) resume().catch(storageError); });
setInterval(() => { if (!document.hidden) sync().catch(storageError); }, 15000);
resume().catch(storageError);
