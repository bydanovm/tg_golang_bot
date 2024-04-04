DROP TABLE IF EXISTS "LimitsDict" CASCADE;
DROP SEQUENCE IF EXISTS LimitsDict_id_seq;
CREATE SEQUENCE LimitsDict_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1;

CREATE TABLE "public"."LimitsDict" (
    "idLmtDct" integer DEFAULT nextval('LimitsDict_id_seq') NOT NULL,
    "nameLmtDct" TEXT not NULL,
    "stdValLmt" INTEGER DEFAULT 0,
    CONSTRAINT "LimitsDict_pkey" PRIMARY KEY ("idLmtDct")
) WITH (oids = false);


DROP TABLE IF EXISTS "Limits" CASCADE;
DROP SEQUENCE IF EXISTS Limits_id_seq;
CREATE SEQUENCE Limits_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1;

CREATE TABLE "public"."Limits" (
    "idLmt" integer DEFAULT nextval('Limits_id_seq') NOT NULL,
    "valAvailLmt" INTEGER not NULL DEFAULT 0,
    "valUsedLmt" INTEGER not NULL DEFAULT 0,
    "activeLmt" BOOLEAN DEFAULT FALSE,
    "tsLmtOn" TIMESTAMP NULL,
    "tsLmtOff" TIMESTAMP NULL,
    "userId" INTEGER not NULL,
    CONSTRAINT "Limits_pkey" PRIMARY KEY ("idLmt"),
    CONSTRAINT "Limits_Users_fkey" FOREIGN KEY("userId")
        REFERENCES "public"."Users"("idUsr")
) WITH (oids = false);