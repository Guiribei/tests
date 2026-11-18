const CACHE = "caderno-offline-v2";
const SHELL = "/modo-offline";
const assets = [SHELL, "/static/app.css?v=5", "/static/app.js?v=4",
  "/static/local-store.js?v=1", "/static/baskervville.ttf", "/static/favicon.svg?v=2"];

self.addEventListener("install", event => {
  event.waitUntil(caches.open(CACHE).then(cache => cache.addAll(assets)));
});

self.addEventListener("activate", event => {
  event.waitUntil((async () => {
    for (const key of await caches.keys()) {
      if (key.startsWith("caderno-offline-") && key !== CACHE) await caches.delete(key);
    }
    await self.clients.claim();
  })());
});

self.addEventListener("fetch", event => {
  const {request} = event;
  const url = new URL(request.url);
  if (request.method !== "GET" || url.origin !== self.location.origin) return;
  if (assets.includes(url.pathname + url.search)) {
    event.respondWith(caches.open(CACHE).then(async cache => (await cache.match(request)) || fetch(request)));
  } else if (request.mode === "navigate" && (url.pathname === "/" ||
      /^\/tentativas\/[a-f0-9]+\/(questoes\/\d+|resultado)$/.test(url.pathname))) {
    event.respondWith((async () => {
      try {
        const response = await fetch(request, {signal: AbortSignal.timeout(6000)});
        if (response.status < 500) return response;
        return (await caches.match(SHELL, {cacheName: CACHE})) || response;
      } catch {
        return (await caches.match(SHELL, {cacheName: CACHE})) || Response.error();
      }
    })());
  }
});
