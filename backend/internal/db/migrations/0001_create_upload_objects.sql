-- Migration inicial: criação da tabela upload_objects
-- Esta migration NÃO é executada pela aplicação.
-- Ela deve ser rodada por um processo externo (Makefile/CI/etc.).

CREATE TABLE IF NOT EXISTS upload_objects (
  id         VARCHAR(255) PRIMARY KEY,
  file_name  TEXT        NOT NULL,
  file_size  BIGINT      NOT NULL,
  mime_type  TEXT        NOT NULL,
  object_key TEXT        NOT NULL,
  status     TEXT        NOT NULL,
  expires_at TIMESTAMPTZ NOT NULL
);

