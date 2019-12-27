### setup


- terminalizer record setup.yml
- terminalizer play setup.yml
- terminalizer render setup.yml -o setup.gif -q 100
- https://gifcompressor.com/


cf push --var auth_username=broker --var auth_password=secret --var api_token=235302cc-34f1-4425-8584-8e1516cfdaa2

curl https://broker:secret@compose-broker.applicationcloud.io/health

cf create-service-broker compose broker secret https://compose-broker.applicationcloud.io --space-scoped

cf marketplace



