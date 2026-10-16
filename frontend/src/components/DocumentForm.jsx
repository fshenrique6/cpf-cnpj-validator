import { useState } from "react";
import { createDocument } from "../api/Documents";

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
    <form className="form" onSubmit={handleSubmit}>
      <input
        className="form__input"
        type="text"
        placeholder="Digite um CPF ou CNPJ"
        value={number}
        onChange={(event) => setNumber(event.target.value)}
        disabled={submitting}
      />
      <button
        className="btn btn--primary"
        type="submit"
        disabled={submitting || number.trim() === ""}
      >
        Cadastrar
      </button>
      {error && (
        <p className="alert" role="alert">
          {error}
        </p>
      )}
    </form>
  );
}

export default DocumentForm;