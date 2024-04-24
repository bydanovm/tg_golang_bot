ALTER TABLE "public"."dictcrypto" ADD CONSTRAINT constraint_name UNIQUE ("cryptoid");

DROP TABLE IF EXISTS "typetrackingcrypto" CASCADE;
DROP SEQUENCE IF EXISTS typetrackingcrypto_id_seq;
CREATE SEQUENCE typetrackingcrypto_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1;

CREATE TABLE "public"."typetrackingcrypto" (
    "idtyptrkcrp" integer DEFAULT nextval('typetrackingcrypto_id_seq') NOT NULL,
    "nametypetrkcrp" TEXT NOT NULL UNIQUE,
    "desctyptrkcrp" TEXT NOT NULL,
    "risingtyptrkcrp" BOOLEAN NOT NULL DEFAULT FALSE,
    "calcproctyptrkcrp" BOOLEAN NOT NULL DEFAULT FALSE,
    CONSTRAINT "typetrackingcrypto_pkey" PRIMARY KEY ("idtyptrkcrp")
) WITH (oids = false);

-- Добавление начальных типов отслеживания КВ
MERGE INTO public.typetrackingcrypto as t
USING (VALUES (1,'RAISE_V','Превышение свыше %.3f %s (вычисление по значению)',  true, false),
              (2,'FALL_V','Падение ниже %.3f %s (вычисление по значению)',      false,true),
              (3,'RAISE_P','Превышение свыше %.0f %s (вычисление по проценту)',true, false),
              (4,'FALL_P','Падение ниже %.0f %s (вычисление по проценту)',    false,true))
    s (idtyptrkcrp, nametypetrkcrp, desctyptrkcrp, risingtyptrkcrp, calcproctyptrkcrp)
    ON t.idtyptrkcrp = s.idtyptrkcrp
WHEN MATCHED THEN 
    UPDATE SET (idtyptrkcrp, nametypetrkcrp, desctyptrkcrp, risingtyptrkcrp, calcproctyptrkcrp) =
            (s.idtyptrkcrp, s.nametypetrkcrp, s.desctyptrkcrp, s.risingtyptrkcrp, s.calcproctyptrkcrp)
WHEN NOT MATCHED THEN
    INSERT (idtyptrkcrp, nametypetrkcrp, desctyptrkcrp, risingtyptrkcrp, calcproctyptrkcrp)
    VALUES (s.idtyptrkcrp, s.nametypetrkcrp, s.desctyptrkcrp, s.risingtyptrkcrp, s.calcproctyptrkcrp);

DROP TABLE IF EXISTS "trackingcrypto" CASCADE;
DROP SEQUENCE IF EXISTS trackingcrypto_id_seq;
CREATE SEQUENCE trackingcrypto_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1;

CREATE TABLE "public"."trackingcrypto" (
    "idtrkcrp" integer DEFAULT nextval('trackingcrypto_id_seq') NOT NULL,
    "valtrkcrp" NUMERIC(19,9) DEFAULT 0,
    "ontrkcrp" BOOLEAN NOT NULL DEFAULT FALSE,
    "lmtid" INTEGER,
    "typtrkcrpid" INTEGER, 
    "dctcrpid" INTEGER,
    "userid" INTEGER,
    CONSTRAINT "trackingcrypto_pkey" PRIMARY KEY ("idtrkcrp"),
    CONSTRAINT "trackingcrypto_limits_fkey" FOREIGN KEY("lmtid")
        REFERENCES "public"."limits"("idlmt"),
    CONSTRAINT "trackingcrypto_typetrackingcrypto_fkey" FOREIGN KEY("typtrkcrpid")
        REFERENCES "public"."typetrackingcrypto"("idtyptrkcrp"),
    CONSTRAINT "trackingcrypto_dictcrypto_fkey" FOREIGN KEY("dctcrpid")
        REFERENCES "public"."dictcrypto"("cryptoid"),
    CONSTRAINT "trackingcrypto_users_fkey" FOREIGN KEY("userid")
        REFERENCES "public"."users"("idusr")
) WITH (oids = false);

-- MERGE INTO public.trackingcrypto as t
-- USING (VALUES (1, 68300.0, true, 1, 1, 786751823))
--     s (idtrkcrp, valtrkcrp, ontrkcrp, typtrkcrpid, dctcrpid, userid)
--     ON t.idtrkcrp = s.idtrkcrp
-- WHEN MATCHED THEN 
--     UPDATE SET (idtrkcrp, valtrkcrp, ontrkcrp, typtrkcrpid, dctcrpid, userid) =
--             (s.idtrkcrp, s.valtrkcrp, s.ontrkcrp, s.typtrkcrpid, s.dctcrpid, s.userid)
-- WHEN NOT MATCHED THEN
--     INSERT (idtrkcrp, valtrkcrp, ontrkcrp, typtrkcrpid, dctcrpid, userid)
--     VALUES (s.idtrkcrp, s.valtrkcrp, s.ontrkcrp, s.typtrkcrpid, s.dctcrpid, s.userid);