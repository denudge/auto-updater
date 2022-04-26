CREATE UNIQUE INDEX IF NOT EXISTS releases_unique_idx ON releases (version, arch, os, variant, product, vendor);
