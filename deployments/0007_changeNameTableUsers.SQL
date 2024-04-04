BEGIN;
    ALTER TABLE public.users
    RENAME TO logmsg;
END;