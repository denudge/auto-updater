CREATE TABLE apps_default_groups (
  id SERIAL PRIMARY KEY,
  app_id INTEGER NOT NULL,
  group_id INTEGER NOT NULL,
  created_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,

  CONSTRAINT apps_default_groups_ux UNIQUE (app_id, group_id),
  CONSTRAINT fk_apps_default_groups_apps FOREIGN KEY (app_id) REFERENCES apps (id),
  CONSTRAINT fk_apps_default_groups_groups FOREIGN KEY (group_id) REFERENCES groups (id)
);

--bun:split

CREATE INDEX apps_default_groups_groups_idx ON apps_default_groups (group_id);
