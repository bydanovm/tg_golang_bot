ALTER TABLE "public"."dictcrypto" ADD CONSTRAINT constraint_name UNIQUE ("cryptoid");

DROP TABLE IF EXISTS "TypeTrackingCrypto" CASCADE;
DROP SEQUENCE IF EXISTS TypeTrackingCrypto_id_seq;
CREATE SEQUENCE TypeTrackingCrypto_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1;

CREATE TABLE "public"."TypeTrackingCrypto" (
    "idTypTrkCrp" integer DEFAULT nextval('TypeTrackingCrypto_id_seq') NOT NULL,
    "nameTypeTrkCrp" TEXT NOT NULL,
    "modeTypTrkCrp" TEXT NOT NULL,
    CONSTRAINT "TypeTrackingCrypto_pkey" PRIMARY KEY ("idTypTrkCrp")
) WITH (oids = false);

DROP TABLE IF EXISTS "TrackingCrypto" CASCADE;
DROP SEQUENCE IF EXISTS TrackingCrypto_id_seq;
CREATE SEQUENCE TrackingCrypto_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1;

CREATE TABLE "public"."TrackingCrypto" (
    "idTrkCrp" integer DEFAULT nextval('TrackingCrypto_id_seq') NOT NULL,
    "valTrkCrp" NUMERIC(15,9) DEFAULT 0,
    "typTrkCrpId" INTEGER, 
    "dctCrpId" INTEGER,
    "userId" INTEGER,
    CONSTRAINT "TrackingCrypto_pkey" PRIMARY KEY ("idTrkCrp"),
    CONSTRAINT "TrackingCrypto_TypeTrackingCrypto_fkey" FOREIGN KEY("typTrkCrpId")
        REFERENCES "public"."TypeTrackingCrypto"("idTypTrkCrp"),
    CONSTRAINT "TrackingCrypto_DictCrypto_fkey" FOREIGN KEY("dctCrpId")
        REFERENCES "public"."dictcrypto"("cryptoid"),
    CONSTRAINT "TrackingCrypto_Users_fkey" FOREIGN KEY("userId")
        REFERENCES "public"."Users"("idUsr")
) WITH (oids = false);