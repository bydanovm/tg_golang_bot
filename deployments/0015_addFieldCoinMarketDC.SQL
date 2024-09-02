BEGIN;
    ALTER TABLE "public"."dictcrypto"
    ADD "coinmrktid" INTEGER NOT NULL DEFAULT 2;

    ALTER TABLE "public"."dictcrypto" 
    ADD CONSTRAINT "dictcrypto_coinmarkets_fkey"
     FOREIGN KEY ("coinmrktid") REFERENCES "public"."coinmarkets"("idmrkt");
END;