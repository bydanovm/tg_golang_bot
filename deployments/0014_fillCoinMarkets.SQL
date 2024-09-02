MERGE INTO public.coinmarkets as t
USING (VALUES (1, 'https://api.livecoinwatch.com', 'Live Coin Watch', 'LCW', 'x-api-key', '', 'content-type', 60),
              (2, 'https://pro-api.coinmarketcap.com', 'Coin Market Cup', 'CMC', 'X-CMC_PRO_API_KEY', 'Accepts', '', 300))
    s (idmrkt, url, shortdesc, description, headapikey, headaccept, headcontenttype, timeout)
    ON t.idmrkt = s.idmrkt
WHEN MATCHED THEN 
    UPDATE SET (idmrkt, url, shortdesc, description, headapikey, headaccept, headcontenttype, timeout) =
            (s.idmrkt, s.url, s.shortdesc, s.description, s.headapikey, s.headaccept, s.headcontenttype, s.timeout)
WHEN NOT MATCHED THEN
    INSERT (idmrkt, url, shortdesc, description, headapikey, headaccept, headcontenttype, timeout)
    VALUES (s.idmrkt, s.url, s.shortdesc, s.description, s.headapikey, s.headaccept, s.headcontenttype, s.timeout);

SELECT setval('coinmarkets_id_seq', 
(select max(idmrkt) from coinmarkets), true);

MERGE INTO public.coinmarketsendpoint as t
USING (VALUES (1, '/coins/list', 'POST', 'Assorted information for a list of coins', 1),
              (2, '/coins/map', 'POST', 'Assorted information for a custom map of coins.', 1))
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
USING (VALUES (1, 'currency', 'string', 'USD', 'any valid coin or fiat code', 1),
              (2, 'sort', 'string', 'rank', 'sorting parameter, rank, price, volume, code, name, age', 1),
              (3, 'order', 'string', 'ascending', 'sorting order, ascending or descending', 1),
              (4, 'offset', 'number', '0', 'offset of the list, default 0', 1),
              (5, 'limit', 'number', '100', 'limit of the list, default 10, maximum 100', 1),
              (6, 'meta', 'boolean', 'false', 'to include full coin information or not', 1),
              (7, 'codes', 'array', '["ETH","BTC","TON"]', 'any valid coin or fiat code', 2),
              (8, 'currency', 'string', 'USD', 'any valid coin or fiat code', 2),
              (9, 'sort', 'string', 'rank', 'sorting parameter, rank, price, volume, code, name, age', 2),
              (10, 'order', 'string', 'ascending', 'sorting order, ascending or descending', 2),
              (11, 'offset', 'number', '0', 'offset of the list, default 0', 2),
              (12, 'limit', 'number', '100', 'limit of the list, default 10, maximum 100', 2),
              (13, 'meta', 'boolean', 'false', 'to include full coin information or not', 2)
      )
    s (idmrkthand, key, type, value, description, coinmrktendid)
    ON t.idmrkthand = s.idmrkthand
WHEN MATCHED THEN 
    UPDATE SET (idmrkthand, key, type, value, description, coinmrktendid) =
            (s.idmrkthand, s.key, s.type, s.value, s.description, s.coinmrktendid)
WHEN NOT MATCHED THEN
    INSERT (idmrkthand, key, type, value, description, coinmrktendid)
    VALUES (s.idmrkthand, s.key, s.type, s.value, s.description, s.coinmrktendid);

SELECT setval('coinmarketshand_id_seq', 
(select max(idmrkthand) from coinmarketshand), true);