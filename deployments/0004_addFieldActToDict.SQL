BEGIN;
    ALTER TABLE public.dictcrypto
    ADD "active" BOOLEAN NOT NULL DEFAULT TRUE;
END;