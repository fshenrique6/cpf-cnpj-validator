import { useEffect, useState } from "react";
import { listDocuments } from "./api/Documents";
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
    <main className="app">
      <h1 className="app__title">CPF/CNPJ Validator</h1>
      <p className="app__subtitle">
        Cadastre, valide e gerencie CPFs e CNPJs.
      </p>

      <section className="card">
        <DocumentForm onCreated={reload} />
      </section>

      <section className="card">
        {error && (
          <p className="alert" role="alert">
            {error}
          </p>
        )}
        <DocumentTable documents={documents} onChanged={reload} />
      </section>
    </main>
  );
}

export default App;