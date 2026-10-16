const API_URL = import.meta.env.VITE_API_URL ?? "http://localhost:8080";

async function request(path, options = {}) {
  const response = await fetch(`${API_URL}${path}`, {
    headers: { "Content-Type": "application/json" },
    ...options,
  });

  const data = await response.json();

  if (!response.ok) {
    throw new Error(data.error ?? "Erro inesperado");
  }

  return data;
}

export function listDocuments(params = {}) {
  const query = new URLSearchParams(
    Object.entries(params).filter(([, value]) => value !== "" && value != null)
  ).toString();

  return request(`/documents${query ? `?${query}` : ""}`);
}

export function createDocument(number) {
  return request("/documents", {
    method: "POST",
    body: JSON.stringify({ number }),
  });
}

export function updateBlocklist(id, blocklisted) {
  return request(`/documents/${id}/blocklist`, {
    method: "PATCH",
    body: JSON.stringify({ blocklisted }),
  });
}

export function deleteDocument(id) {
  return request(`/documents/${id}`, { method: "DELETE" });
}