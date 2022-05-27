CREATE TABLE clients_groups (
  id SERIAL PRIMARY KEY,
  client_id INTEGER NOT NULL,
  group_id INTEGER NOT NULL,
  created_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,

  CONSTRAINT clients_groups_ux UNIQUE (client_id, group_id),
  CONSTRAINT fk_clients_groups_clients FOREIGN KEY (client_id) REFERENCES clients (id),
  CONSTRAINT fk_clients_groups_groups FOREIGN KEY (group_id) REFERENCES groups (id)
);

--bun:split

CREATE INDEX clients_groups_groups_idx ON clients_groups (group_id);
