CREATE TABLE apps (
  id SERIAL PRIMARY KEY,
  vendor VARCHAR(255) NOT NULL,
  product VARCHAR(255) NOT NULL,
  "name" VARCHAR(255) NULL DEFAULT NULL,
  is_active BOOL NOT NULL DEFAULT TRUE,
  is_locked BOOL NOT NULL DEFAULT FALSE,
  upgrade_target VARCHAR(63) NULL DEFAULT NULL,
  created_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,

  CONSTRAINT apps_ux UNIQUE (vendor, product)
);

--bun:split

INSERT INTO apps (vendor, product, "name") SELECT DISTINCT vendor, product, "name" FROM releases;

--bun:split

ALTER TABLE releases ADD COLUMN app_id INTEGER NULL DEFAULT NULL;

--bun:split

ALTER TABLE releases ADD CONSTRAINT fk_releases_apps FOREIGN KEY (app_id) REFERENCES apps (id);

--bun:split

UPDATE releases SET app_id=a.id FROM apps a WHERE a.vendor=releases.vendor AND a.product=releases.product;

--bun:split

ALTER TABLE releases ALTER COLUMN app_id SET NOT NULL;

--bun:split

ALTER TABLE releases DROP CONSTRAINT releases_ux;

--bun:split

ALTER TABLE releases ADD CONSTRAINT releases_ux UNIQUE (version, arch, os, variant, app_id);

--bun:split

DROP INDEX releases_vendor_product;

--bun:split

ALTER TABLE releases DROP COLUMN vendor;

--bun:split

ALTER TABLE releases DROP COLUMN product;

--bun:split

ALTER TABLE releases DROP COLUMN "name";

--bun:split

CREATE INDEX releases_app ON releases (app_id);
