### provisioning


- terminalizer record provisioning.yml
- terminalizer play provisioning.yml
- terminalizer render provisioning.yml -o provisioning.gif -q 100
- https://gifcompressor.com/


cf create-service scylla default my-scylla-db

cf service my-scylla-db

cf bind-service my-app my-scylla-db

cf env my-app

cqlsh -u scylla -p EJOGEHGPDWSACVXE -t -e 'SHOW VERSION' portal0291-2.febdca84-52ef-4g5f-f232-da92c432944b.1213820901.composedb.com 18184

