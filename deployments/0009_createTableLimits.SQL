DROP TABLE IF EXISTS "limitsdict" CASCADE;
DROP SEQUENCE IF EXISTS limitsdict_id_seq;
CREATE SEQUENCE limitsdict_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1;

CREATE TABLE "public"."limitsdict" (
    "idlmtdct" INTEGER DEFAULT nextval('limitsdict_id_seq') NOT NULL,
    "namelmtdct" TEXT NOT NULL UNIQUE,
    "desclmtdct" TEXT NOT NULL,
    "stdvallmt" INTEGER DEFAULT 0,
    CONSTRAINT "limitsdict_pkey" PRIMARY KEY ("idlmtdct")
) WITH (oids = false);

MERGE INTO public.limitsdict as t
USING (VALUES (1, 'LMT001', 'Лимит количества запросов в день к API', 10),
              (2, 'LMT002', 'Лимит количества запросов в день к боту', 100),
              (3, 'LMT003', 'Лимит количества уведомлений после срабатывания триггера отслеживания КВ', 3))
    s (idlmtdct, namelmtdct, desclmtdct, stdvallmt)
    ON t.idlmtdct = s.idlmtdct
WHEN MATCHED THEN 
    UPDATE SET (idlmtdct, namelmtdct, desclmtdct, stdvallmt) =
            (s.idlmtdct, s.namelmtdct, s.desclmtdct, s.stdvallmt)
WHEN NOT MATCHED THEN
    INSERT (idlmtdct, namelmtdct, desclmtdct, stdvallmt)
    VALUES (s.idlmtdct, s.namelmtdct, s.desclmtdct, s.stdvallmt);

DROP TABLE IF EXISTS "limits" CASCADE;
DROP SEQUENCE IF EXISTS limits_id_seq;
CREATE SEQUENCE limits_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1;

CREATE TABLE "public"."limits" (
    "idlmt" integer DEFAULT nextval('limits_id_seq') NOT NULL,
    "valavaillmt" INTEGER NOT NULL DEFAULT 0,
    "valusedlmt" INTEGER NOT NULL DEFAULT 0,
    "activelmt" BOOLEAN DEFAULT FALSE,
    "tslmton" TIMESTAMP NULL,
    "tslmtoff" TIMESTAMP NULL,
    "userid" INTEGER NOT NULL,
    "ltmdctid" INTEGER NOT NULL,
    CONSTRAINT "limits_pkey" PRIMARY KEY ("idlmt"),
    CONSTRAINT "limits_users_fkey" FOREIGN KEY("userid")
        REFERENCES "public"."users"("idusr"),
    CONSTRAINT "limits_limitsdict_fkey" FOREIGN KEY("ltmdctid")
        REFERENCES "public"."limitsdict"("idlmtdct")
) WITH (oids = false);

MERGE INTO public.limits as t
USING (VALUES (1, 3, 0, true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 786751823, 3))
    s (idlmt, valavaillmt, valusedlmt, activelmt, tslmton, tslmtoff, userid, ltmdctid)
    ON t.idlmt = s.idlmt
WHEN MATCHED THEN 
    UPDATE SET (idlmt, valavaillmt, valusedlmt, activelmt, tslmton, tslmtoff, userid, ltmdctid) =
            (s.idlmt, s.valavaillmt, s.valusedlmt, s.activelmt, s.tslmton, s.tslmtoff, s.userid, s.ltmdctid)
WHEN NOT MATCHED THEN
    INSERT (idlmt, valavaillmt, valusedlmt, activelmt, tslmton, tslmtoff, userid, ltmdctid)
    VALUES (s.idlmt, s.valavaillmt, s.valusedlmt, s.activelmt, s.tslmton, s.tslmtoff, s.userid, s.ltmdctid);
