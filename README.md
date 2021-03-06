# :sparkles: compose-broker :game_die:

[![CircleCI](https://circleci.com/gh/JamesClonk/compose-broker.svg?style=svg)](https://circleci.com/gh/JamesClonk/compose-broker)
[![License](https://img.shields.io/badge/license-Apache--2.0-blue)](https://github.com/JamesClonk/compose-broker/blob/master/LICENSE)
[![Platform](https://img.shields.io/badge/platform-Cloud%20Foundry-lightgrey)](https://developer.swisscom.com/)

> #### Conquer the Data Layer
> Performance and reliable data layers for developers who'd rather spend their time building apps than managing databases.

**compose-broker** is a [Compose.io](https://www.compose.com/) [service broker](https://www.openservicebrokerapi.org/) for [Cloud Foundry](https://www.cloudfoundry.org/) and [Kubernetes](https://kubernetes.io/)

It supports databases and services such as:
- [🐘 PostgreSQL](https://www.compose.com/databases/postgresql)
- [🐬 MySQL](https://www.compose.com/databases/mysql)
- [👻 RethinkDB](https://www.compose.com/databases/rethinkdb)
- [🐙 ScyllaDB (Cassandra)](https://www.compose.com/databases/scylladb)
- [🕷 Elasticsearch](https://www.compose.com/databases/elasticsearch)
- [🐦 Redis](https://www.compose.com/databases/redis)
- [🐟 etcd](https://www.compose.com/databases/etcd)
- [🐇 RabbitMQ](https://www.compose.com/databases/rabbitmq)

## Usage

#### Deploy service broker to Cloud Foundry

1. create an [API Token](https://app.compose.io/oauth/api_tokens) on your Compose.io [account](https://app.compose.io/account)
2. pick a Cloud Foundry provider.
   I'd suggest the [Swisscom AppCloud](https://developer.swisscom.com/)
3. push the app, providing the API key and a username/password to secure the service broker with
4. register the service broker in your space (`--space-scoped`)
5. check `cf marketplace` to see your new available service plans

![create service broker](https://raw.githubusercontent.com/JamesClonk/compose-broker/recordings/setup-min.gif "create service broker")

#### Provision new databases

1. create a new service instance (`cf cs`)
2. bind the service instance to your app (`cf bs`), or create a service key (`cf csk`)
3. inspect the service binding/key, have a look at the credentials (`cf env`/`cf sk`)
4. use the given credentials to connect to your new database
5. enjoy!

![provision service](https://raw.githubusercontent.com/JamesClonk/compose-broker/recordings/provisioning-min.gif "provision service")

## Configuration

All configuration of the service broker is done through environment variables (provided by `manifest.yml` during a `cf push`) and the included `catalog.yml`.

### manifest.yml

Possible configuration values are:
```yaml
BROKER_LOG_LEVEL: info # optional, can be set to debug, info, warning, error or fatal, defaults to info
BROKER_LOG_TIMESTAMP: false # optional, add timestamp to logging messages (not needed when deployed on Cloud Foundry), defaults to false
BROKER_SKIP_SSL_VALIDATION: false, # optional, disables SSL certificate verification for API calls, defaults to false
BROKER_AUTH_USERNAME: broker-username # required, HTTP basic auth username to secure service broker with
BROKER_AUTH_PASSWORD: broker-password # required, HTTP basic auth password to secure service broker with
BROKER_CATALOG_FILENAME: catalog.yml # optional, filename containing all catalog information, defaults to catalog.yml
COMPOSE_API_URL: https://api.compose.io/2016-07/ # optional, Base URL of Compose.io API, defaults to https://api.compose.io/2016-07
COMPOSE_API_TOKEN: e7fb89a0-26f8-4ee5-890e-3c68079b15ea # required, Compose.io API Token
COMPOSE_API_DEFAULT_DATACENTER: gce:europe-west1 # optional, defaults to aws:eu-central-1
COMPOSE_API_DEFAULT_ACCOUNT_ID: 586eab527c65836dde5533e8 # optional, service broker will try to read it from Compose.io API if not set
```

### catalog.yml

The service brokers provided `catalog.yml` configures all available services offerings and plans. (See [OSBAPI spec](https://github.com/openservicebrokerapi/servicebroker/blob/v2.15/spec.md#service-offering-object))

Adjust this file to define your own different service plans with different deployment sizes / units. A plan's metadata allows for optionally configuring custom [unit](https://apidocs.compose.com/docs/scaling) sizes, enabling/disabling Redis' [cache mode](https://help.compose.com/docs/redis-cache-and-storage-modes), specifying explicit software versions or the [datacenter](https://apidocs.compose.com/docs/datacenters) to host the deployment. 

Review the included Redis example plans for these properties:
https://github.com/JamesClonk/compose-broker/blob/f7331ef8cc1a18c7fc4b060931e0cb35e7580f5e/catalog.yml#L19-L59

###### Plan metadata example:
```yaml
metadata:
  # Number of resource units to allocate to the deployment (optional, defaults to 1)
  units: 2
  # Whether to optimize the deployment to be used as a cache (optional, Redis only)
  cache_mode: true
  # Version of the software to deploy (optional)
  version: "4.0.14"
  # Datacenter to use for deployment (optional, defaults to $COMPOSE_API_DEFAULT_DATACENTER)
  datacenter: aws:eu-central-1
```

#### Account ID & Datacenter

By default the service broker will provision new database deployments with the configured account id `COMPOSE_API_DEFAULT_ACCOUNT_ID` and datacenter `COMPOSE_API_DEFAULT_DATACENTER` (see `manifest.yml`).
If no account id is configured it will try to read the value over the Compose.io API and take the first account it finds.
Similarly if no datacenter is configured it will use `aws:eu-central-1` as default value.

When issuing service provisioning requests to the service broker it is possible to provide the account id and/or the datacenter as additional parameters.
###### Example:
```bash
cf create-service etcd default my-etcd -c '{
  "account_id": "454f3deb8cad236ffb3452e9", 
  "datacenter": "gce:europe-west1"
}'
```

#### Units

When issuing service provisioning requests to the service broker it is also possible to request a specific [unit](https://apidocs.compose.com/docs/scaling) size (which basically determines the scaling for your Compose.io deployments, _memory_ and _disk space_) instead of the configured value in the service brokers `catalog.yml`.

It is also possible to update existing service instances with `cf update-service`.
###### Example:
```bash
cf create-service scylla default my-scylla-db -c '{ "units": 10 }'
# or
cf update-service my-postgres-db -c '{ "units": 4 }'
```

#### Version

Last but not least it is also possible during service provisioning to request specific software versions for your database deployments.
You can get all available versions from the [Compose.io API](https://apidocs.compose.com/reference#2016-07-get-databases): https://api.compose.io/2016-07/databases

###### Example:
```bash
cf create-service rethink default my-rethinkdb -c '{ "version": "2.3.7" }'
```
