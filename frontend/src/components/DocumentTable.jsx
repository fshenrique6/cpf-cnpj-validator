import { useState } from "react";
import { deleteDocument, updateBlocklist } from "../api/Documents.js";

function DocumentTable({ documents, onChanged }) {
  const [error, setError] = useState("");

  async function handleToggleBlocklist(doc) {
    setError("");

    try {
      await updateBlocklist(doc.id, !doc.blocklisted);
      onChanged();
    } catch (err) {
      setError(err.message);
    }
  }

  async function handleDelete(doc) {
    if (!window.confirm(`Excluir o documento ${doc.number}?`)) {
      return;
    }

    setError("");

    try {
      await deleteDocument(doc.id);
      onChanged();
    } catch (err) {
      setError(err.message);
    }
  }

  if (documents.length === 0) {
    return <p>Nenhum documento cadastrado.</p>;
  }

  return (
    <>
      {error && <p role="alert">{error}</p>}

      <table>
        <thead>
          <tr>
            <th>Tipo</th>
            <th>Número</th>
            <th>Status</th>
            <th>Ações</th>
          </tr>
        </thead>
        <tbody>
          {documents.map((doc) => (
            <tr key={doc.id}>
              <td>{doc.type.toUpperCase()}</td>
              <td>{doc.number}</td>
              <td>{doc.blocklisted ? "Bloqueado" : "Liberado"}</td>
              <td>
                <button onClick={() => handleToggleBlocklist(doc)}>
                  {doc.blocklisted ? "Desbloquear" : "Bloquear"}
                </button>
                <button onClick={() => handleDelete(doc)}>Excluir</button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </>
  );
}

export default DocumentTable;