const API_URL = import.meta.env.VITE_API_URL ?? "http://localhost:8080";

export async function request(path, options = {}) {
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