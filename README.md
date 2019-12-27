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

### Default Datacenter

By default the service broker will provision new database deployments in the configured datacenter `COMPOSE_API_DEFAULT_DATACENTER` (see `manifest.yml`) or if none configured at all it will use `aws:eu-central-1` as default value.
When issuing service provisioning requests to the service broker it is also possible to provide the datacenter as an additional parameter.
###### Example:
```
$ cf create-service etcd default my-etcd -c '{"datacenter": "gce:europe-west1"}'
```
