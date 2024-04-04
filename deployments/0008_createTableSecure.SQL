DROP TABLE IF EXISTS "LevelsSecure" CASCADE;
DROP SEQUENCE IF EXISTS LevelsSecure_id_seq;
CREATE SEQUENCE LevelsSecure_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1;

CREATE TABLE "public"."LevelsSecure" (
    "idLvlSec" integer DEFAULT nextval('LevelsSecure_id_seq') NOT NULL,
    "nameLvlSec" TEXT not NULL,
    CONSTRAINT "LevelsSecure_pkey" PRIMARY KEY ("idLvlSec")
) WITH (oids = false);


DROP TABLE IF EXISTS "LevelsSecureAdd" CASCADE;
DROP SEQUENCE IF EXISTS LevelsSecureAdd_id_seq;
CREATE SEQUENCE LevelsSecureAdd_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1;

CREATE TABLE "public"."LevelsSecureAdd" (
    "idLvlSecAdd" integer DEFAULT nextval('LevelsSecureAdd_id_seq') NOT NULL,
    "nameLvlSecAdd" TEXT not NULL,
    "activeLvlSecAdd" BOOLEAN DEFAULT FALSE,
    "lvlSecId" INTEGER not NULL,
    CONSTRAINT "LevelsSecureAdd_pkey" PRIMARY KEY ("idLvlSecAdd"),
    CONSTRAINT "LevelsSecureAdd_LevelsSecure_fkey" FOREIGN KEY("lvlSecId")
        REFERENCES "public"."LevelsSecure"("idLvlSec")
) WITH (oids = false);

DROP TABLE IF EXISTS "Groups" CASCADE;
DROP SEQUENCE IF EXISTS Groups_id_seq;
CREATE SEQUENCE Groups_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1;

CREATE TABLE "public"."Groups" (
    "idGrp" integer DEFAULT nextval('Groups_id_seq') NOT NULL,
    "nameGrp" TEXT not NULL,
    "lvlSecId" INTEGER,
    CONSTRAINT "Groups_pkey" PRIMARY KEY ("idGrp"),
    CONSTRAINT "Groups_LevelsSecure_fkey" FOREIGN KEY("lvlSecId")
        REFERENCES "public"."LevelsSecure"("idLvlSec")
) WITH (oids = false);

DROP TABLE IF EXISTS "Users" CASCADE;
DROP SEQUENCE IF EXISTS Users_id_seq;
CREATE SEQUENCE Users_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1;

CREATE TABLE "public"."Users" (
    "idUsr" integer DEFAULT nextval('Users_id_seq') NOT NULL,
    "tsUsr" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "nameUsr" TEXT not NULL,
    "chatIdUsr" INTEGER not null,
    "groupId" INTEGER NOT NULL,
    CONSTRAINT "Users_pkey" PRIMARY KEY ("idUsr"),
    CONSTRAINT "Users_Groups_fkey" FOREIGN KEY("groupId")
        REFERENCES "public"."Groups"("idGrp")
) WITH (oids = false);