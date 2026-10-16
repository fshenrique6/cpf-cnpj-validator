import { useEffect, useState } from "react";
import { listDocuments } from "./api/Documents.js";

function App() {
  const [documents, setDocuments] = useState([]);
  const [error, setError] = useState("");

  useEffect(() => {
    listDocuments()
      .then(setDocuments)
      .catch((err) => setError(err.message));
  }, []);

  return (
    <main>
      <h1>CPF/CNPJ Validator</h1>

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