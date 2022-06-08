ALTER TABLE releases ADD COLUMN "link" varchar(511) NULL DEFAULT NULL;

--bun:split

ALTER TABLE releases ADD COLUMN "format" varchar(63) NULL DEFAULT NULL;
