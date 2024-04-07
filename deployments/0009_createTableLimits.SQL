DROP TABLE IF EXISTS "limitsdict" CASCADE;
DROP SEQUENCE IF EXISTS limitsdict_id_seq;
CREATE SEQUENCE limitsdict_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1;

CREATE TABLE "public"."limitsdict" (
    "idlmtdct" INTEGER DEFAULT nextval('limitsdict_id_seq') NOT NULL,
    "namelmtdct" TEXT NOT NULL UNIQUE,
    "stdvallmt" INTEGER DEFAULT 0,
    CONSTRAINT "limitsdict_pkey" PRIMARY KEY ("idlmtdct")
) WITH (oids = false);


DROP TABLE IF EXISTS "limits" CASCADE;
DROP SEQUENCE IF EXISTS limits_id_seq;
CREATE SEQUENCE limits_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1;

CREATE TABLE "public"."limits" (
    "idlmt" integer DEFAULT nextval('limits_id_seq') NOT NULL,
    "valavaillmt" INTEGER not NULL DEFAULT 0,
    "valusedlmt" INTEGER not NULL DEFAULT 0,
    "activelmt" BOOLEAN DEFAULT FALSE,
    "tslmton" TIMESTAMP NULL,
    "tslmtoff" TIMESTAMP NULL,
    "userid" INTEGER not NULL,
    CONSTRAINT "limits_pkey" PRIMARY KEY ("idlmt"),
    CONSTRAINT "limits_users_fkey" FOREIGN KEY("userid")
        REFERENCES "public"."users"("idusr")
) WITH (oids = false);