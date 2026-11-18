let connection;

function openDatabase() {
  if (!connection) {
    connection = new Promise((resolve, reject) => {
      const req = indexedDB.open("caderno-offline", 1);
      req.onupgradeneeded = () => req.result.createObjectStore("attempts", {keyPath: "id"});
      req.onsuccess = () => resolve(req.result);
      req.onerror = () => reject(req.error);
    });
  }
  return connection;
}

export async function listAttempts() {
  const db = await openDatabase();
  return new Promise((resolve, reject) => {
    const req = db.transaction("attempts").objectStore("attempts").getAll();
    req.onsuccess = () => resolve(req.result);
    req.onerror = () => reject(req.error);
  });
}

// Read and modify in the same IndexedDB transaction, including across tabs.
// Resolve only after commit, so the UI can accurately report a local save.
export async function attempt(id, change) {
  const db = await openDatabase();
  return new Promise((resolve, reject) => {
    const tx = db.transaction("attempts", change ? "readwrite" : "readonly");
    const store = tx.objectStore("attempts");
    const req = store.get(id);
    let value;
    req.onsuccess = () => {
      value = req.result;
      if (change) {
        value = change(value);
        if (value) store.put(value);
      }
    };
    tx.oncomplete = () => resolve(value);
    tx.onabort = () => reject(tx.error || new Error("Gravação interrompida"));
  });
}

export async function forgetAttempt(id) {
  const db = await openDatabase();
  return new Promise((resolve, reject) => {
    const tx = db.transaction("attempts", "readwrite");
    tx.objectStore("attempts").delete(id);
    tx.oncomplete = resolve;
    tx.onabort = () => reject(tx.error);
  });
}
