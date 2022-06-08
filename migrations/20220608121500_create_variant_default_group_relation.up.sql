CREATE TABLE variants_default_groups (
  id SERIAL PRIMARY KEY,
  variant_id INTEGER NOT NULL,
  group_id INTEGER NOT NULL,
  created_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,

  CONSTRAINT variants_default_groups_ux UNIQUE (variant_id, group_id),
  CONSTRAINT fk_variants_default_groups_variants FOREIGN KEY (variant_id) REFERENCES variants (id),
  CONSTRAINT fk_variants_default_groups_groups FOREIGN KEY (group_id) REFERENCES groups (id)
);

--bun:split

CREATE INDEX variants_groups_groups_idx ON variants_default_groups (group_id);
