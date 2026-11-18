import {attempt, listAttempts, forgetAttempt} from "./local-store.js?v=1";

const $ = selector => document.querySelector(selector);
const layout = $(".exam-layout");
const library = $(".local-attempts");
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

function status(message) {
  const el = $("#offline-status");
  if (el) {
    el.hidden = false;
    el.textContent = message;
  }
}

function storageError() {
  if (current) render(current);
  if (offlineShell && !current && library) {
    library.hidden = false;
    $(".local-message").textContent = "Armazenamento local indisponível. Reabra a prova com conexão.";
  }
  status("Falha ao gravar no aparelho. Confira o espaço disponível e tente novamente.");
}

function answers(row) {
  return {...row.answers, ...row.pending};
}

function delivery(row) {
  if (row.finished) return "Entrega recebida";
  if (warnings.has(row.id)) return warnings.get(row.id);
  if (row.finishing) return "Aguardando conexão para corrigir";
  if (Object.keys(row.pending).length) return "Salvo neste aparelho · sincronização pendente";
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
  $(".sidebar-title span").textContent = row.role;
  $(".progress-copy strong").textContent = count + " de " + total;
  $(".progress-track span").style.width = (100 * count / total) + "%";
  $(".mobile-progress span").textContent = "Questão " + number + " de " + total;
  $(".mobile-progress span:last-child").textContent = count + " respondidas";
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
    link.classList.toggle("answered", Boolean(choices[n]));
    link.setAttribute("aria-label", "Questão " + n + (choices[n] ? ", respondida" : ""));
  }
  for (const [index, label] of document.querySelectorAll(".alternative").entries()) {
    const input = label.querySelector("input");
    input.checked = choices[number] === input.value;
    label.lastElementChild.textContent = q.alternatives[index];
  }
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
    link.textContent = row.title + " · " + delivery(row);
    return link;
  }));
  $(".local-message").textContent = message || (offlineShell && !rows.length
    ? "Conecte-se para abrir e baixar uma prova." : "");
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
    if (sheet.finished || sheet.answers[n] === value) delete row.pending[n];
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
    if (res.status === 404) warnings.set(row.id, "Sessão indisponível. Respostas guardadas neste aparelho.");
    else if (res.status < 500) warnings.set(row.id, "Envio recusado. Respostas guardadas neste aparelho.");
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
      status("Modo online. Para baixar a prova, abra o site por HTTPS.");
      return;
    }
    status("Preparando prova para uso offline…");
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
    if (offlineShell) await renderLibrary("Abra esta tentativa com conexão para concluir o download.");
    else status("Download offline incompleto. Você pode continuar online.");
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
  window.scrollTo(0, 0);
  rememberPosition().catch(storageError);
}

document.addEventListener("click", event => {
  const link = event.target.closest(".exam-layout a");
  if (!current || !link || event.button !== 0 || event.ctrlKey || event.metaKey || event.shiftKey || event.altKey) return;
  if (!link.pathname.startsWith(prefix(id) + "/questoes/")) return;
  event.preventDefault();
  navigate(Number(link.pathname.split("/").pop()), true);
});

window.addEventListener("popstate", () => navigate(Number(location.pathname.split("/").pop()), false));

document.addEventListener("change", event => {
  const input = event.target;
  if (!input.matches(".answer-form input[name=choice]")) return;
  if (!current) {
    input.form.requestSubmit();
    return;
  }
  const n = number;
  const value = input.value;
  saveLocal(row => { row.pending[n] = value; });
});

document.addEventListener("submit", event => {
  if (current && event.target.matches(".answer-form")) {
    event.preventDefault();
    return;
  }
  const form = event.target.closest("[data-confirm-finish]");
  if (!form) return;
  if (current) event.preventDefault();
  const blank = current ? current.questions.length - Object.keys(answers(current)).length
    : Number(form.dataset.total) - Number(form.dataset.answered);
  const prompt = (blank ? "Questões em branco: " + blank + ". " : "") +
    "Finalizar a prova e encerrar a edição das respostas?";
  if (!window.confirm(prompt)) event.preventDefault();
  else if (current) saveLocal(row => { row.finishing = true; });
});

channel?.addEventListener("message", event => { updateView(event.data).catch(storageError); });

async function resume() {
  if (!window.indexedDB) {
    status("Armazenamento local indisponível. Use a prova online.");
    return;
  }
  const confirmed = document.querySelector("[data-result-attempt]")?.dataset.resultAttempt;
  if (confirmed) await forgetAttempt(confirmed);
  if (id) await prepare();
  else await renderLibrary(offlineShell && location.pathname.endsWith("/resultado")
    ? "Conecte-se para consultar a correção." : "");
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
