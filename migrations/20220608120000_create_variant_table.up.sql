CREATE TABLE variants (
  id SERIAL PRIMARY KEY,
  app_id INTEGER NOT NULL,
  "name" VARCHAR(255) NOT NULL,
  active BOOL NOT NULL DEFAULT TRUE,
  "locked" BOOL NOT NULL DEFAULT FALSE,
  upgrade_target VARCHAR(63) NULL DEFAULT NULL,
  "allow_register" boolean NOT NULL DEFAULT FALSE, -- For any non-empty variant, the default is false
  created_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,

  CONSTRAINT variants_ux UNIQUE ("name", app_id),
  CONSTRAINT fk_variants_apps FOREIGN KEY (app_id) REFERENCES apps (id)
);

--bun:split

CREATE INDEX variants_app ON variants (app_id);
