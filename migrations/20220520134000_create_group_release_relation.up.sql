CREATE TABLE releases_groups (
  id SERIAL PRIMARY KEY,
  release_id INTEGER NOT NULL,
  group_id INTEGER NOT NULL,
  created_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,

  CONSTRAINT releases_groups_ux UNIQUE (release_id, group_id),
  CONSTRAINT fk_releases_groups_releases FOREIGN KEY (release_id) REFERENCES releases (id),
  CONSTRAINT fk_releases_groups_groups FOREIGN KEY (group_id) REFERENCES groups (id)
);

--bun:split

CREATE INDEX releases_groups_groups_idx ON releases_groups (group_id);
