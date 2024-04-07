ALTER TABLE "public"."dictcrypto" ADD CONSTRAINT constraint_name UNIQUE ("cryptoid");

DROP TABLE IF EXISTS "typetrackingcrypto" CASCADE;
DROP SEQUENCE IF EXISTS typetrackingcrypto_id_seq;
CREATE SEQUENCE typetrackingcrypto_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1;

CREATE TABLE "public"."typetrackingcrypto" (
    "idtyptrkcrp" integer DEFAULT nextval('typetrackingcrypto_id_seq') NOT NULL,
    "nametypetrkcrp" TEXT NOT NULL,
    "modetyptrkcrp" TEXT NOT NULL,
    CONSTRAINT "typetrackingcrypto_pkey" PRIMARY KEY ("idtyptrkcrp")
) WITH (oids = false);

DROP TABLE IF EXISTS "trackingcrypto" CASCADE;
DROP SEQUENCE IF EXISTS trackingcrypto_id_seq;
CREATE SEQUENCE trackingcrypto_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1;

CREATE TABLE "public"."trackingcrypto" (
    "idtrkcrp" integer DEFAULT nextval('trackingcrypto_id_seq') NOT NULL,
    "valtrkcrp" NUMERIC(15,0) DEFAULT 0,
    "typtrkcrpid" INTEGER, 
    "dctcrpid" INTEGER,
    "userid" INTEGER,
    CONSTRAINT "trackingcrypto_pkey" PRIMARY KEY ("idtrkcrp"),
    CONSTRAINT "trackingcrypto_typetrackingcrypto_fkey" FOREIGN KEY("typtrkcrpid")
        REFERENCES "public"."typetrackingcrypto"("idtyptrkcrp"),
    CONSTRAINT "trackingcrypto_dictcrypto_fkey" FOREIGN KEY("dctcrpid")
        REFERENCES "public"."dictcrypto"("cryptoid"),
    CONSTRAINT "trackingcrypto_users_fkey" FOREIGN KEY("userid")
        REFERENCES "public"."users"("idusr")
) WITH (oids = false);