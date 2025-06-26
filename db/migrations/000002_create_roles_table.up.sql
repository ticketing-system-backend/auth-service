CREATE TYPE role_level AS ENUM ('admin', 'superadmin', 'staff', 'hrd', 'finance', 'customer');

CREATE TABLE IF NOT EXISTS roles (
  id SERIAL PRIMARY KEY,
  nama VARCHAR NOT NULL,
  deskripsi VARCHAR NOT NULL,
  level role_level,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
)