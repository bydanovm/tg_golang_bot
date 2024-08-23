MERGE INTO public.coinmarkets as t
USING (VALUES (1, 'https://api.livecoinwatch.com', 'Live Coin Watch', 'x-api-key', NULL, 'content-type', 60))
    s (idmrkt, url, "desc", headapikey, headaccept, headcontenttype, timeout)
    ON t.idmrkt = s.idmrkt
WHEN MATCHED THEN 
    UPDATE SET (idmrkt, url, "desc", headapikey, headaccept, headcontenttype, timeout) =
            (s.idmrkt, s.url, s.desc, s.headapikey, s.headaccept, s.headcontenttype, s.timeout)
WHEN NOT MATCHED THEN
    INSERT (idmrkt, url, "desc", headapikey, headaccept, headcontenttype, timeout)
    VALUES (s.idmrkt, s.url, s.desc, s.headapikey, s.headaccept, s.headcontenttype, s.timeout);

SELECT setval('coinmarkets_id_seq', 
(select max(idmrkt) from coinmarkets), true);

MERGE INTO public.coinmarketsendpoint as t
USING (VALUES (1, '/coins/list', 'POST', 'Assorted information for a list of coins', 1))
    s (idmrktend, endpoint, method, "desc", coinmrktid)
    ON t.idmrktend = s.idmrktend
WHEN MATCHED THEN 
    UPDATE SET (idmrktend, endpoint, method, "desc", coinmrktid) =
            (s.idmrktend, s.endpoint, s.method, s.desc, s.coinmrktid)
WHEN NOT MATCHED THEN
    INSERT (idmrktend, endpoint, method, "desc", coinmrktid)
    VALUES (s.idmrktend, s.endpoint, s.method, s.desc, s.coinmrktid);

SELECT setval('coinmarketsendpoint_id_seq', 
(select max(idmrktend) from coinmarketsendpoint), true);

MERGE INTO public.coinmarketshand as t
USING (VALUES (1, 'currency', 'string', 'any valid coin or fiat code', 1),
              (2, 'sort', 'string', 'sorting parameter, rank, price, volume, code, name, age', 1),
              (3, 'order', 'string', 'sorting order, ascending or descending', 1),
              (4, 'offset', 'number', 'offset of the list, default 0', 1),
              (5, 'limit', 'number', 'limit of the list, default 10, maximum 100', 1),
              (6, 'meta', 'boolean', 'to include full coin information or not', 1)
      )
    s (idmrkthand, key, type, "desc", coinmrktendid)
    ON t.idmrkthand = s.idmrkthand
WHEN MATCHED THEN 
    UPDATE SET (idmrkthand, key, type, "desc", coinmrktendid) =
            (s.idmrkthand, s.key, s.type, s.desc, s.coinmrktendid)
WHEN NOT MATCHED THEN
    INSERT (idmrkthand, key, type, "desc", coinmrktendid)
    VALUES (s.idmrkthand, s.key, s.type, s.desc, s.coinmrktendid);

SELECT setval('coinmarketshand_id_seq', 
(select max(idmrkthand) from coinmarketshand), true);