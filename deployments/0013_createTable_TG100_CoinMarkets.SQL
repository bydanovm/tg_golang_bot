DROP TABLE IF EXISTS "coinmarkets" CASCADE;
DROP SEQUENCE IF EXISTS coinmarkets_id_seq;
CREATE SEQUENCE coinmarkets_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1;

CREATE TABLE "public"."coinmarkets" (
    "idmrkt" integer DEFAULT nextval('coinmarkets_id_seq') NOT NULL,
    "url" TEXT NOT NULL,
    "desc" TEXT NOT NULL,
    "headapikey" TEXT,
    "headaccept" TEXT,
    "headcontenttype" TEXT,
    "timeout" INTEGER DEFAULT 0,
    CONSTRAINT "coinmarkets_pkey" PRIMARY KEY ("idmrkt")
) WITH (oids = false);

DROP TABLE IF EXISTS "coinmarketsendpoint" CASCADE;
DROP SEQUENCE IF EXISTS coinmarketsendpoint_id_seq;
CREATE SEQUENCE coinmarketsendpoint_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1;

CREATE TABLE "public"."coinmarketsendpoint" (
    "idmrktend" integer DEFAULT nextval('coinmarketsendpoint_id_seq') NOT NULL,
    "endpoint" TEXT NOT NULL,
    "method" TEXT NOT NULL,
    "desc" TEXT NOT NULL,
    "coinmrktid" INTEGER DEFAULT 0,
    CONSTRAINT "coinmarketsendpoint_pkey" PRIMARY KEY ("idmrktend"),
    CONSTRAINT "coinmarketsendpoint_coinmarkets_fkey" FOREIGN KEY("coinmrktid")
        REFERENCES "public"."coinmarkets"("idmrkt")
) WITH (oids = false);

DROP TABLE IF EXISTS "coinmarketshand" CASCADE;
DROP SEQUENCE IF EXISTS coinmarketshand_id_seq;
CREATE SEQUENCE coinmarketshand_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1;

CREATE TABLE "public"."coinmarketshand" (
    "idmrkthand" integer DEFAULT nextval('coinmarketshand_id_seq') NOT NULL,
    "key" TEXT NOT NULL,
    "type" TEXT NOT NULL,
    "desc" TEXT NOT NULL,
    "coinmrktendid" INTEGER DEFAULT 0,
    CONSTRAINT "coinmarketshand_pkey" PRIMARY KEY ("idmrkthand"),
    CONSTRAINT "coinmarketshand_coinmarketsendpoint_fkey" FOREIGN KEY("coinmrktendid")
        REFERENCES "public"."coinmarketsendpoint"("idmrktend")
) WITH (oids = false);
