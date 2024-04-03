BEGIN;
    ALTER TABLE public.dictcrypto
    ADD "cryptocounter" integer NOT NULL DEFAULT 0;
END;