import { useState } from "react";
import { deleteDocument, updateBlocklist } from "../api/Documents";
import { formatDocument } from "../utils/Format";

function DocumentTable({ 
    documents, 
    onChanged,
    emptyMessage = "Nenhum documento encontrado."}) {
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
    if (!window.confirm(`Excluir o documento ${formatDocument(doc.number, doc.type)}?`)) {
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
    return <p className="empty">{emptyMessage}</p>;
  }

  return (
    <>
      {error && (
        <p className="alert" role="alert">
          {error}
        </p>
      )}

      <div className="table-wrapper">
        <table className="table">
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
                <td>
                  <span className="badge badge--type">{doc.type.toUpperCase()}</span>
                </td>
                <td className="table__number">
                  {formatDocument(doc.number, doc.type)}
                </td>
                <td>
                  <span
                    className={`badge ${doc.blocklisted ? "badge--blocked" : "badge--ok"}`}
                  >
                    {doc.blocklisted ? "Bloqueado" : "Liberado"}
                  </span>
                </td>
                <td>
                  <div className="table__actions">
                    <button
                      className="btn btn--ghost btn--small"
                      onClick={() => handleToggleBlocklist(doc)}
                    >
                      {doc.blocklisted ? "Desbloquear" : "Bloquear"}
                    </button>
                    <button
                      className="btn btn--danger btn--small"
                      onClick={() => handleDelete(doc)}
                    >
                      Excluir
                    </button>
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </>
  );
}

export default DocumentTable;