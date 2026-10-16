import { useState } from "react";
import { createDocument } from "../api/Documents.js";

function DocumentForm({ onCreated }) {
  const [number, setNumber] = useState("");
  const [error, setError] = useState("");
  const [submitting, setSubmitting] = useState(false);

  async function handleSubmit(event) {
    event.preventDefault();
    setError("");
    setSubmitting(true);

    try {
      await createDocument(number);
      setNumber("");
      onCreated();
    } catch (err) {
      setError(err.message);
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <form onSubmit={handleSubmit}>
      <input
        type="text"
        placeholder="CPF ou CNPJ"
        value={number}
        onChange={(event) => setNumber(event.target.value)}
        disabled={submitting}
      />
      <button type="submit" disabled={submitting || number.trim() === ""}>
        Cadastrar
      </button>
      {error && <p role="alert">{error}</p>}
    </form>
  );
}

export default DocumentForm;