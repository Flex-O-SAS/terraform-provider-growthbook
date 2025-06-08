<!-- markdown-toc start - Don't edit this section. Run M-x markdown-toc-refresh-toc -->
**Table of Contents**

- [GrowthBook Terraform Provider](#growthbook-terraform-provider)
- [Documentation](#documentation)
- [Contributing](#contributing)
  - [Acceptance tests](#acceptance-tests)

<!-- markdown-toc end -->


# GrowthBook Terraform Provider

This Terraform provider allows you to manage GrowthBook resources and retrieve data from
your GrowthBook instance using Terraform.

This is poorly written using AI as a bootstrap and covers only basic use cases.


# Documentation

See the [docs](./docs/) directory for detailed documentation.

# Contributing

## Acceptance tests

The repository runs acceptance tests against a running instance of growthbook.
This can be either cloud (not recommended) or a local instance spawned with docker.

```
#start instance
export GROWTHBOOK_API_KEY=$(docker-compose -f ./acceptance/docker-compose.yml up -d > /dev/null && sleep 2 && ./acceptance/startup.sh)

# run tests
GROWTHBOOK_API_URL=http://localhost:3100/api/v1  TF_ACC=1 go test ./...

# run test with debug output from the provider
TF_LOG=debug GROWTHBOOK_API_URL=http://localhost:3100/api/v1  TF_ACC=1 go test ./...

# cleanup growthbook instance (with volumes)
docker-compose -f ./acceptance/docker-compose.yml down -v
```
