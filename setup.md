### setup


- terminalizer record setup.yml
- terminalizer play setup.yml
- terminalizer render setup.yml -o setup.gif -q 100
- https://gifcompressor.com/


cf push --var auth_username=broker --var auth_password=secret --var api_token=6c9e75fc-08dd-49bd-ab70-5b7c6364ab0b

curl https://broker:secret@compose-broker.applicationcloud.io/health

cf create-service-broker compose broker secret https://compose-broker.applicationcloud.io --space-scoped

cf marketplace



