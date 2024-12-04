# map-cli
Mattilsynet operational cli for resources, supports authentication towards NATS.

# Plan:

1. Authentication NATS, check auth folder TODOS
    1. #TODO: make it possible to switch projects somehow in conjunction with nats credentials, don't know how yet.
2. Managed environment creation, uses map-query-api and map-command-api which facilitates creation of ME towards subject which is used by map-managed-environment wasmcloud module.
3. GCP-project creation within managed-environment wish, this needs to be formalized somehow. #TODO
