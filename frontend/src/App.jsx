import { useEffect, useState } from "react";
import { listDocuments } from "./api/Documents.js";
import DocumentForm from "./components/DocumentForm";
import DocumentTable from "./components/DocumentTable";

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

      <DocumentTable documents={documents} onChanged={reload} />
    </main>
  );
}

export default App;