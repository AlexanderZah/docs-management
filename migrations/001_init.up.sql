CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    login VARCHAR(255) UNIQUE NOT NULL,
    password TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS documents (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,               
    is_file BOOLEAN NOT NULL DEFAULT TRUE, 
    is_public BOOLEAN NOT NULL DEFAULT FALSE, 
    token TEXT NOT NULL,         
    mime TEXT NOT NULL,               
    grants TEXT[] NOT NULL DEFAULT '{}',
    json_data JSONB,                  
    content BYTEA,                     
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_documents_created ON public.documents USING btree (created_at);
CREATE INDEX idx_documents_name ON public.documents USING btree (name);
CREATE INDEX idx_documents_grants ON public.documents USING gin (grants);