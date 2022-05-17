CREATE TABLE groups (
  id SERIAL PRIMARY KEY,
  app_id INTEGER NOT NULL,
  "name" VARCHAR(255) NOT NULL,
  created_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,

  CONSTRAINT groups_ux UNIQUE ("name", app_id),
  CONSTRAINT fk_groups_apps FOREIGN KEY (app_id) REFERENCES apps (id)
);

--bun:split

CREATE INDEX groups_app ON groups (app_id);
