MERGE INTO public.coinmarkets as t
USING (VALUES (1, 'https://api.livecoinwatch.com', 'Live Coin Watch', 'x-api-key', '', 'content-type', 60),
              (2, 'https://pro-api.coinmarketcap.com', 'CMC', 'X-CMC_PRO_API_KEY', 'Accepts', '', 300))
    s (idmrkt, url, description, headapikey, headaccept, headcontenttype, timeout)
    ON t.idmrkt = s.idmrkt
WHEN MATCHED THEN 
    UPDATE SET (idmrkt, url, description, headapikey, headaccept, headcontenttype, timeout) =
            (s.idmrkt, s.url, s.description, s.headapikey, s.headaccept, s.headcontenttype, s.timeout)
WHEN NOT MATCHED THEN
    INSERT (idmrkt, url, description, headapikey, headaccept, headcontenttype, timeout)
    VALUES (s.idmrkt, s.url, s.description, s.headapikey, s.headaccept, s.headcontenttype, s.timeout);

SELECT setval('coinmarkets_id_seq', 
(select max(idmrkt) from coinmarkets), true);

MERGE INTO public.coinmarketsendpoint as t
USING (VALUES (1, '/coins/list', 'POST', 'Assorted information for a list of coins', 2))
    s (idmrktend, endpoint, method, description, coinmrktid)
    ON t.idmrktend = s.idmrktend
WHEN MATCHED THEN 
    UPDATE SET (idmrktend, endpoint, method, description, coinmrktid) =
            (s.idmrktend, s.endpoint, s.method, s.description, s.coinmrktid)
WHEN NOT MATCHED THEN
    INSERT (idmrktend, endpoint, method, description, coinmrktid)
    VALUES (s.idmrktend, s.endpoint, s.method, s.description, s.coinmrktid);

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
    s (idmrkthand, key, type, description, coinmrktendid)
    ON t.idmrkthand = s.idmrkthand
WHEN MATCHED THEN 
    UPDATE SET (idmrkthand, key, type, description, coinmrktendid) =
            (s.idmrkthand, s.key, s.type, s.description, s.coinmrktendid)
WHEN NOT MATCHED THEN
    INSERT (idmrkthand, key, type, description, coinmrktendid)
    VALUES (s.idmrkthand, s.key, s.type, s.description, s.coinmrktendid);

SELECT setval('coinmarketshand_id_seq', 
(select max(idmrkthand) from coinmarketshand), true);