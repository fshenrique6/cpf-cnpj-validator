import { useEffect, useState } from "react";
import { listDocuments } from "./api/Documents";
import DocumentForm from "./components/DocumentForm";
import DocumentTable from "./components/DocumentTable";
import Filters from "./components/Filters";
import StatusBar from "./components/StatusBar";

const INITIAL_FILTERS = {
  type: "",
  blocklisted: "",
  sort_by: "created_at",
  order: "desc",
};

function App() {
  const [documents, setDocuments] = useState([]);
  const [filters, setFilters] = useState(INITIAL_FILTERS);
  const [error, setError] = useState("");
  const [reloadKey, setReloadKey] = useState(0);

  useEffect(() => {
    let ignore = false;

    listDocuments(filters)
      .then((data) => {
        if (!ignore) {
          setDocuments(data);
          setError("");
        }
      })
      .catch((err) => {
        if (!ignore) {
          setError(err.message);
        }
      });

    return () => {
      ignore = true;
    };
  }, [filters, reloadKey]);

  function reload() {
    setReloadKey((key) => key + 1);
  }

  const hasActiveFilters = filters.type !== "" || filters.blocklisted !== "";

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
        <Filters
          filters={filters}
          onChange={setFilters}
          onReset={() => setFilters(INITIAL_FILTERS)}
        />
      </section>

      <section className="card">
        {error && (
          <p className="alert" role="alert">
            {error}
          </p>
        )}
        <DocumentTable
          documents={documents}
          onChanged={reload}
          emptyMessage={
            hasActiveFilters
              ? "Nenhum documento encontrado com esses filtros."
              : "Nenhum documento cadastrado."
          }
        />
      </section>

      <StatusBar />
    </main>
  );
}

export default App;