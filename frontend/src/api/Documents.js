import { request } from "./Client";

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