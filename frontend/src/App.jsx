import { useEffect, useState } from "react";
import { listDocuments } from "./api/Documents.js";
import DocumentForm from "./components/DocumentForm";

function App() {
  const [documents, setDocuments] = useState([]);
  const [error, setError] = useState("");
  const [reloadKey, setReloadKey] = useState(0);

  useEffect(() => {
    listDocuments()
      .then((data) => {
        setDocuments(data);
        setError("");
      })
      .catch((err) => setError(err.message));
  }, [reloadKey]);

  function reload() {
    setReloadKey((key) => key + 1);
  }

  return (
    <main>
      <h1>CPF/CNPJ Validator</h1>

      <DocumentForm onCreated={reload} />

      {error && <p>{error}</p>}

      <ul>
        {documents.map((doc) => (
          <li key={doc.id}>
            {doc.type.toUpperCase()} — {doc.number}
            {doc.blocklisted && " (bloqueado)"}
          </li>
        ))}
      </ul>
    </main>
  );
}

export default App;