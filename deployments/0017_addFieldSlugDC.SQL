BEGIN;
    ALTER TABLE public.dictcrypto
    ADD "slug" TEXT NOT NULL DEFAULT '';
    
    ALTER TABLE public.dictcrypto
    ADD "cryptofullname" TEXT NOT NULL DEFAULT '';
END;