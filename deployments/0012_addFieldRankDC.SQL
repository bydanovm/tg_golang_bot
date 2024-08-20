BEGIN;
    ALTER TABLE public.dictcrypto
    ADD "cryptorank" integer NOT NULL DEFAULT 0;
END;