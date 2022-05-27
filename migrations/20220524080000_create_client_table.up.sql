CREATE TABLE clients (
  id SERIAL PRIMARY KEY,
  app_id INTEGER NOT NULL,
  variant VARCHAR(127) NULL DEFAULT NULL,
  uuid varchar(40) NOT NULL,
  "name" VARCHAR(255) NULL DEFAULT NULL, -- optional, maybe internal client name
  active BOOL NOT NULL DEFAULT TRUE,
  "locked" BOOL NOT NULL DEFAULT FALSE, -- true -> gets no updates anymore
  created_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,

  CONSTRAINT clients_ux UNIQUE (uuid),
  CONSTRAINT fk_clients_apps FOREIGN KEY (app_id) REFERENCES apps (id)
);

--bun:split

CREATE INDEX clients_app_uuid_idx ON clients (app_id, uuid);
