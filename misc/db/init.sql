-- -------------------------------------------------------------
-- TablePlus 3.7.1(332)
--
-- https://tableplus.com/
--
-- Database: gold_store
-- Generation Time: 2022-10-05 17:27:31.6470
-- -------------------------------------------------------------


DROP TABLE IF EXISTS "public"."harga";
-- This script only contains the table creation statements and does not fully represent the table in the database. It's still missing: indices, triggers. Do not use it as a backup.

-- Table Definition
CREATE TABLE "public"."harga" (
    "id" varchar NOT NULL,
    "topup_price" numeric NOT NULL,
    "buyback_price" numeric NOT NULL,
    "created_at" timestamp,
    "updated_at" timestamp,
    "user_id" varchar,
    PRIMARY KEY ("id")
);

DROP TABLE IF EXISTS "public"."rekening";
-- This script only contains the table creation statements and does not fully represent the table in the database. It's still missing: indices, triggers. Do not use it as a backup.

-- Table Definition
CREATE TABLE "public"."rekening" (
    "id" varchar NOT NULL,
    "user_id" varchar,
    "norek" varchar,
    "saldo" float4,
    PRIMARY KEY ("id")
);

DROP TABLE IF EXISTS "public"."transaksi";
-- This script only contains the table creation statements and does not fully represent the table in the database. It's still missing: indices, triggers. Do not use it as a backup.

-- Table Definition
CREATE TABLE "public"."transaksi" (
    "id" varchar NOT NULL,
    "rekening_id" varchar NOT NULL,
    "gram" float8 NOT NULL,
    "type" varchar NOT NULL,
    "topup_price" numeric,
    "buyback_price" numeric,
    "created_at" int8,
    PRIMARY KEY ("id")
);

DROP TABLE IF EXISTS "public"."user";
-- This script only contains the table creation statements and does not fully represent the table in the database. It's still missing: indices, triggers. Do not use it as a backup.

-- Table Definition
CREATE TABLE "public"."user" (
    "id" varchar NOT NULL,
    "role" varchar,
    PRIMARY KEY ("id")
);

INSERT INTO "public"."harga" ("id", "topup_price", "buyback_price", "created_at", "updated_at", "user_id") VALUES
('1', '900000', '1000000', '2022-10-05 08:44:09.586862', NULL, NULL);

INSERT INTO "public"."rekening" ("id", "user_id", "norek", "saldo") VALUES
('1', 'admin', NULL, NULL),
('2', 'customer', '3424255233', '1.2');

