** WARNING!!!! Alpha software **

Consul-externalservice implements external services in consul that follow
the status of their associated check.  This feature is not implemented in consul itself because it would
break consul's node level design. Checks are run and associated at one or many arbitrary consul nodes.

Every service must have at least a check associated with a node. External services nodes are arbitrary names and do not 
need to follow consul's node names. But checks associated with an external service must be run and defined in a real consul node.
This is due to the inner designs of consul that run antientropy at the node level.

Usage
=====

Run

```
consul-externalservice start --node <nodename>
```

This command starts an external service watcher for any service defined at nodename. Nodename is an arbitrary name. All checks are defined and run
from the consul node attached to the running consul-externalservice instance (currently only attaches to localhost:8500).

Consul external services are defined by creating a key-value pair in consul's database. Example:

```
curl -XPUT http://localhost:8500/v1/kv/ExternalServices/<nodename>/<servicename> \
  -d '{ "Address":"localhost", "Port":80, "Interval":"1s", "Command":"ping -c 1 localhost", "TargetState":"running" }'
```

TargetState must be one of:
"stopped" (if you currently do not want the service to be watched), "running" (if
you DO want the service to be watched), and "deleted" if you want the service
definition to be deleted by the service node watcher. Delete is not currently implemented 
but can be easily added.

Thus external services can be defined from any consul accessible point but will not be instantiated if there is no consul-externalservice watcher running for
the defined external service node.

You can run as many consul-externalservice processes per nodename as you want to allow for redundancy. Only one instance at a time will be the leader and honour
nodename service definitions.

You can check service registration with

```
dig @localhost -p 8600 <servicename>.service.consul
```

NOTE that if the service watcher dies and there are no other watchers for the same external services node, checks will remain active as long as the consul
agent is alive but will not activate or deactivate service when changing their status.

Install
=======

You only need to download the consul-externalservice executable and have a properly configured consul network accesible through localhost:8500.

Executables can be found here: https://github.com/jmcarbo/consul-externalservice/releases/tag/v0.0.1

Docker
======

If you want to run a dockerized consul-externalservice:

```
docker run -d jmcarbo/consul-externalservice /bin/start.sh <consul address to join to>
```

Development
===========

Clone repository. Use `make test` to run tests and `make build` to create binaries.
