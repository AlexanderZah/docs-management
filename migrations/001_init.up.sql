CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    login VARCHAR(255) UNIQUE NOT NULL,
    password TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS documents (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    owner_login VARCHAR(255) NOT NULL,
    mime TEXT NOT NULL,
    is_file BOOLEAN NOT NULL DEFAULT TRUE,
    is_public BOOLEAN NOT NULL DEFAULT FALSE,
    json_data JSONB,
    created_at TIMESTAMP DEFAULT now()
);


CREATE TABLE IF NOT EXISTS document_grants (
    doc_id UUID REFERENCES documents(id) ON DELETE CASCADE,
    login VARCHAR(255),
    PRIMARY KEY (doc_id, login)
);


CREATE INDEX idx_documents_owner ON documents(owner_login);
CREATE INDEX idx_documents_created ON documents(created_at);
CREATE INDEX idx_documents_name ON documents(name);