DROP TABLE IF EXISTS "levelssecure" CASCADE;
DROP SEQUENCE IF EXISTS levelssecure_id_seq;
CREATE SEQUENCE levelssecure_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1;

CREATE TABLE "public"."levelssecure" (
    "idlvlsec" integer DEFAULT nextval('levelssecure_id_seq') NOT NULL,
    "namelvlsec" TEXT NOT NULL UNIQUE,
    "desclvlsec" TEXT NOT NULL,
    CONSTRAINT "levelssecure_pkey" PRIMARY KEY ("idlvlsec")
) WITH (oids = false);


DROP TABLE IF EXISTS "levelssecureadd" CASCADE;
DROP SEQUENCE IF EXISTS levelssecureadd_id_seq;
CREATE SEQUENCE levelssecureadd_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1;

CREATE TABLE "public"."levelssecureadd" (
    "idlvlsecadd" integer DEFAULT nextval('levelssecureadd_id_seq') NOT NULL,
    "namelvlsecadd" TEXT not NULL,
    "activelvlsecadd" BOOLEAN DEFAULT FALSE,
    "lvlsecid" INTEGER not NULL,
    CONSTRAINT "levelssecureadd_pkey" PRIMARY KEY ("idlvlsecadd"),
    CONSTRAINT "levelssecureadd_levelssecure_fkey" FOREIGN KEY("lvlsecid")
        REFERENCES "public"."levelssecure"("idlvlsec")
) WITH (oids = false);

DROP TABLE IF EXISTS "groups" CASCADE;
DROP SEQUENCE IF EXISTS groups_id_seq;
-- CREATE SEQUENCE groups_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1;

-- CREATE TABLE "public"."groups" (
--     "idgrp" integer DEFAULT nextval('groups_id_seq') NOT NULL,
--     "namegrp" TEXT NOT NULL UNIQUE,
--     "lvlsecid" INTEGER DEFAULT 5,
--     CONSTRAINT "groups_pkey" PRIMARY KEY ("idgrp"),
--     CONSTRAINT "groups_levelssecure_fkey" FOREIGN KEY("lvlsecid")
--         REFERENCES "public"."levelssecure"("idlvlsec")
-- ) WITH (oids = false);

DROP TABLE IF EXISTS "users" CASCADE;
DROP SEQUENCE IF EXISTS users_id_seq;
CREATE SEQUENCE users_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1;

CREATE TABLE "public"."users" (
    "idusr" integer DEFAULT nextval('users_id_seq') NOT NULL,
    "tsusr" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "nameusr" TEXT NOT NULL,
    "firstname" TEXT NOT NULL,
    "lastname" TEXT not NULL,
    "langcode" TEXT NOT NULL,
    "isbot" BOOLEAN NOT NULL DEFAULT FALSE,
    "isbanned" BOOLEAN NOT NULL DEFAULT FALSE,
    "chatidusr" NUMERIC(15,0) NOT NULL DEFAULT 0,
    "idlvlsec" INTEGER NOT NULL,
    CONSTRAINT "userst_pkey" PRIMARY KEY ("idusr"),
    CONSTRAINT "users_lvlsec_fkey" FOREIGN KEY("idlvlsec")
        REFERENCES "public"."levelssecure"("idlvlsec")
) WITH (oids = false);

-- Добавление начальных уровней доступа
MERGE INTO public.levelssecure as t
USING (VALUES (1,'LVL1','Уровень 1 (наивысший)'),
              (2,'LVL2','Уровень 2 (администратор)'),
              (3,'LVL3','Уровень 3 (модератор)'),
              (4,'LVL4','Уровень 4 (пользователь)'),
              (5,'LVL5','Уровень 5 (проходимец)'))
    s (idlvlsec, namelvlsec, desclvlsec)
    ON t.idlvlsec = s.idlvlsec
WHEN MATCHED THEN 
    UPDATE SET (idlvlsec, namelvlsec, desclvlsec) =
            (s.idlvlsec, s.nameLvlSec, s.desclvlsec)
WHEN NOT MATCHED THEN
    INSERT (idlvlsec, namelvlsec, desclvlsec)
    VALUES (idlvlsec, s.namelvlsec, s.desclvlsec);

-- Добавление начальных надстроек уровней доступа
MERGE INTO "public"."levelssecureadd" as t
USING (VALUES (1,'READ',true,1),
              (2,'WRITE',true,1),
              (3,'DELETE',true,1),
              (4,'UPDATE',true,1),
              (5,'APPROVE',true,1),
              (6,'UNAPPROVE',true,1),
              (7,'BANNING',true,1))
    s (idlvlsecadd, namelvlsecadd, activelvlsecadd, lvlsecid)
    ON t.idlvlsecadd = s.idlvlsecadd
WHEN MATCHED THEN 
    UPDATE SET (idlvlsecadd, namelvlsecadd, activelvlsecadd, lvlsecid) =
            (s.idlvlsecadd, s.namelvlsecadd, s.activelvlsecadd, s.lvlsecid)
WHEN NOT MATCHED THEN
    INSERT (idlvlsecadd, namelvlsecadd, activelvlsecadd, lvlsecid)
    VALUES (s.idlvlsecadd, s.namelvlsecadd, s.activelvlsecadd, s.lvlsecid);

-- -- Добавление начальных групп пользователей
-- MERGE INTO "public"."groups" as t
-- USING (VALUES (1,'Создатель',1),
--               (2,'Создатель',2),
--               (3,'Создатель',3),
--               (4,'Создатель',4),
--               (5,'Создатель',5))
--     s (idgrp, namegrp, lvlsecid)
--     ON t.idgrp = s.idgrp
-- WHEN MATCHED THEN 
--     UPDATE SET (idgrp, namegrp, lvlsecid) =
--             (s.idgrp, s.namegrp, s.lvlsecid)
-- WHEN NOT MATCHED THEN
--     INSERT (idgrp, namegrp, lvlsecid)
--     VALUES (s.idgrp, s.namegrp, s.lvlsecid);