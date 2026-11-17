document.addEventListener("submit", (event) => {
  const form = event.target.closest("[data-confirm-finish]");
  if (!form) return;

  const answered = Number(form.dataset.answered || 0);
  const total = Number(form.dataset.total || 180);
  const blank = total - answered;
  const message = blank > 0
    ? `Você deixou ${blank} questão${blank === 1 ? "" : "ões"} em branco. Depois de finalizar, as respostas não poderão ser alteradas. Deseja corrigir agora?`
    : "Depois de finalizar, as respostas não poderão ser alteradas. Deseja corrigir agora?";

  if (!window.confirm(message)) event.preventDefault();
});

document.addEventListener("htmx:responseError", () => {
  const state = document.querySelector("#save-state");
  if (state) state.textContent = "Não foi possível salvar. Tente novamente.";
});

