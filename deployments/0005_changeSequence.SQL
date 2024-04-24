SELECT setval('cryptoprices_id_seq', 
(select max(id) from cryptoprices), true);

SELECT setval('dictcrypto_id_seq', 
(select max(id) from dictcrypto), true);

SELECT setval('settingsproject_id_seq', 
(select max(id) from settingsproject), true);

SELECT setval('users_id_seq', 
(select max(idusr) from users), true);

SELECT setval('limits_id_seq', 
(select max(idlmt) from limits), true);

SELECT setval('trackingcrypto_id_seq', 
(select max(idtrkcrp) from trackingcrypto), true);
