function Filters({ filters, onChange, onReset }) {
  function handleChange(event) {
    const { name, value } = event.target;
    onChange({ ...filters, [name]: value });
  }

  return (
    <div className="filters">
      <label className="filters__field">
        Tipo
        <select
          className="filters__select"
          name="type"
          value={filters.type}
          onChange={handleChange}
        >
          <option value="">Todos</option>
          <option value="cpf">CPF</option>
          <option value="cnpj">CNPJ</option>
        </select>
      </label>

      <label className="filters__field">
        Status
        <select
          className="filters__select"
          name="blocklisted"
          value={filters.blocklisted}
          onChange={handleChange}
        >
          <option value="">Todos</option>
          <option value="false">Liberados</option>
          <option value="true">Bloqueados</option>
        </select>
      </label>

      <label className="filters__field">
        Ordenar por
        <select
          className="filters__select"
          name="sort_by"
          value={filters.sort_by}
          onChange={handleChange}
        >
          <option value="created_at">Data de cadastro</option>
          <option value="updated_at">Última atualização</option>
          <option value="number">Número</option>
          <option value="type">Tipo</option>
        </select>
      </label>

      <label className="filters__field">
        Ordem
        <select
          className="filters__select"
          name="order"
          value={filters.order}
          onChange={handleChange}
        >
          <option value="desc">Decrescente</option>
          <option value="asc">Crescente</option>
        </select>
      </label>

      <button
        className="btn btn--ghost btn--small filters__reset"
        type="button"
        onClick={onReset}
      >
        Limpar
      </button>
    </div>
  );
}

export default Filters;