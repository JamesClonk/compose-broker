### provisioning


- terminalizer record provisioning.yml
- terminalizer play provisioning.yml
- terminalizer render provisioning.yml -o provisioning.gif -q 100
- https://gifcompressor.com/


cf create-service scylladb default my-scylla-db

cf service my-scylla-db

cf bind-service my-app my-scylla-db

cf env my-app

cqlsh -u username -p password -t -e 'SHOW VERSION'

