CREATE TABLE releases (
  id SERIAL PRIMARY KEY,
  vendor VARCHAR(255) NULL,
  product VARCHAR(255) NOT NULL,
  "name" VARCHAR(255) NULL DEFAULT NULL,
  variant VARCHAR(127) NULL DEFAULT NULL,
  description VARCHAR(1022) NULL DEFAULT NULL,
  os VARCHAR(127) NULL DEFAULT NULL,
  arch VARCHAR(127) NULL DEFAULT NULL,
  released_at TIMESTAMP(0) NOT NULL,
  version VARCHAR(63) NOT NULL,
  unstable BOOL NOT NULL DEFAULT FALSE,
  alias VARCHAR(127) NULL DEFAULT NULL,
  signature VARCHAR(255) NULL DEFAULT NULL,
  tags VARCHAR(127) ARRAY DEFAULT array[]::varchar[],
  upgrade_target VARCHAR(63) NULL DEFAULT NULL,
  should_upgrade INT NOT NULL DEFAULT 0,
  updated_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,

  CONSTRAINT releases_ux UNIQUE (version, arch, os, variant, product, vendor)
);

--bun:split

CREATE INDEX IF NOT EXISTS releases_vendor_product ON releases (vendor, product);
