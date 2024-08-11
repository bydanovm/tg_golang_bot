BEGIN;
    ALTER TABLE public.users
    ALTER "idusr" TYPE bigint,
    ALTER "idusr" SET DEFAULT nextval('users_id_seq'),
    ALTER "idusr" SET NOT NULL;

    ALTER TABLE public.users
    ALTER "chatidusr" TYPE bigint,
    ALTER "chatidusr" SET DEFAULT '0',
    ALTER "chatidusr" SET NOT NULL;

    ALTER TABLE public.trackingcrypto
    ALTER "userid" TYPE bigint,
    ALTER "userid" DROP DEFAULT,
    ALTER "userid" DROP NOT NULL;

    ALTER TABLE public.limits
    ALTER "userid" TYPE bigint,
    ALTER "userid" DROP DEFAULT,
    ALTER "userid" SET NOT NULL;
END;