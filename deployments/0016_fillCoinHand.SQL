MERGE INTO public.coinmarketsendpoint as t
USING (VALUES (3, '/v1/cryptocurrency/map', 'GET', 'Assorted information for a list of coins', 2),
              (4, '/v2/cryptocurrency/quotes/latest', 'GET', 'Quotes Latest v2', 2))
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
USING (VALUES (14, 'listing_status', 'string', 'active', 'Only active cryptocurrencies are returned by default', 3),
              (15, 'start', 'number', '1', 'Optionally offset the start (1-based index) of the paginated list of items to return', 3),
              (16, 'limit', 'number', '200', 'Optionally specify the number of results to return', 3),
              (17, 'sort', 'string', 'cmc_rank', 'What field to sort the list of cryptocurrencies by', 3),
              (18, 'symbol', 'string', '', 'Optionally pass a comma-separated list of cryptocurrency symbols to return CoinMarketCap IDs for', 3),
              (19, 'aux', 'string', 'platform,is_active', 'Optionally specify a comma-separated list of supplemental data fields to return', 3),
              (20, 'symbol', 'array', 'BTC,ETH', 'Optionally specify a comma-separated list of supplemental data fields to return', 4),
              (21, 'convert', 'string', 'USD', 'Optionally specify a comma-separated list of supplemental data fields to return', 4)
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