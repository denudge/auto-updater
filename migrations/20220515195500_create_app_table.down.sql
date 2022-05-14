ALTER TABLE releases ADD COLUMN vendor VARCHAR(255) NULL DEFAULT NULL;

--bun:split

ALTER TABLE releases ADD COLUMN product VARCHAR(255) NULL DEFAULT NULL;

--bun:split

ALTER TABLE releases ADD COLUMN "name" VARCHAR(255) NULL DEFAULT NULL;

--bun:split

UPDATE releases SET vendor=a.vendor, product=a.product,"name"=a."name" FROM apps a WHERE a.id=releases.app_id;

--bun:split

ALTER TABLE releases ALTER COLUMN vendor SET NOT NULL;

--bun:split

ALTER TABLE releases ALTER COLUMN product SET NOT NULL;

--bun:split

DROP INDEX releases_app;

--bun:split

ALTER TABLE releases DROP COLUMN app_id;

--bun:split

ALTER TABLE releases ADD CONSTRAINT releases_ux UNIQUE (version, arch, os, variant, product, vendor);

--bun:split

CREATE INDEX releases_vendor_product ON releases (vendor, product);

--bun:split

DROP TABLE apps;