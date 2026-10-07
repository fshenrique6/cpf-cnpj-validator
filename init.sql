CREATE TABLE documents (
    id SERIAL PRIMARY KEY,
    number TEXT NOT NULL UNIQUE,
    type TEXT NOT NULL CHECK (type IN ('cpf', 'cnpj')),
    blocklisted BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);
