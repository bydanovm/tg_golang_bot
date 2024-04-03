SELECT setval('cryptoprices_id_seq', 
(select max(id) from cryptoprices), true);

SELECT setval('dictcrypto_id_seq', 
(select max(id) from dictcrypto), true);

SELECT setval('settingsproject_id_seq', 
(select max(id) from settingsproject), true);

SELECT setval('users_id_seq', 
(select max(id) from users), true);