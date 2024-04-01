DROP TABLE IF EXISTS "cryptoprices";
DROP SEQUENCE IF EXISTS cryptoprices_id_seq;
CREATE SEQUENCE cryptoprices_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 37363 CACHE 1;

CREATE TABLE "public"."cryptoprices" (
    "id" integer DEFAULT nextval('cryptoprices_id_seq') NOT NULL,
    "timestamp" timestamp DEFAULT CURRENT_TIMESTAMP,
    "cryptoid" integer,
    "cryptoprice" numeric(15,9),
    "cryptoupdate" timestamp,
    CONSTRAINT "cryptoprices_pkey" PRIMARY KEY ("id")
) WITH (oids = false);


DROP TABLE IF EXISTS "dictcrypto";
DROP SEQUENCE IF EXISTS dictcrypto_id_seq;
CREATE SEQUENCE dictcrypto_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 127 CACHE 1;

CREATE TABLE "public"."dictcrypto" (
    "id" integer DEFAULT nextval('dictcrypto_id_seq') NOT NULL,
    "timestamp" timestamp DEFAULT CURRENT_TIMESTAMP,
    "cryptoid" integer,
    "cryptoname" text,
    "cryptolastprice" numeric(15,9),
    "cryptoupdate" timestamp,
    CONSTRAINT "dictcrypto_pkey" PRIMARY KEY ("id")
) WITH (oids = false);


DROP TABLE IF EXISTS "settingsproject";
DROP SEQUENCE IF EXISTS settingsproject_id_seq;
CREATE SEQUENCE settingsproject_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1;

CREATE TABLE "public"."settingsproject" (
    "id" integer DEFAULT nextval('settingsproject_id_seq') NOT NULL,
    "name" text,
    "description" text,
    "active" boolean,
    "type" text,
    "value" text,
    "timestart" timestamp DEFAULT CURRENT_TIMESTAMP,
    "timelast" timestamp DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT "settingsproject_pkey" PRIMARY KEY ("id")
) WITH (oids = false);


DROP TABLE IF EXISTS "users";
DROP SEQUENCE IF EXISTS users_id_seq;
CREATE SEQUENCE users_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 190 CACHE 1;

CREATE TABLE "public"."users" (
    "id" integer DEFAULT nextval('users_id_seq') NOT NULL,
    "timestamp" timestamp DEFAULT CURRENT_TIMESTAMP,
    "username" text,
    "chat_id" integer,
    "message" text,
    "answer" text,
    CONSTRAINT "users_pkey" PRIMARY KEY ("id")
) WITH (oids = false);