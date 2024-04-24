BEGIN;
    ALTER TABLE public.users
    RENAME TO logmsg;
    ALTER SEQUENCE "users_id_seq" RENAME TO "logmsg_id_seq";
END;