# GARM external provider for Equinix Metal

The Equinix Metal external provider allows [garm](https://github.com/cloudbase/garm) to create runners on top of Equinix Metal bare metal servers.

## Build

Clone the repo:

```bash
git clone https://github.com/cloudbase/garm-provider-equinix
```

Build the binary:

```bash
cd garm-provider-equinix
go build .
```

Copy the binary on the same system where ```garm``` is running, and [point to it in the config](https://github.com/cloudbase/garm/blob/main/doc/providers.md#the-external-provider).

## Configure

The config file for this external provider is a simple toml used to configure the credentials needed to connect to your Equinix Metal account.

A sample config file can be found [in the testdata folder](./testdata/garm-provider-equinix.toml).

## Tweaking the provider

Garm supports sending opaque json encoded configs to the IaaS providers it hooks into. This allows the providers to implement some very provider specific functionality that doesn't necessarily translate well to other providers. Features that may exists on Azure, may not exist on AWS or OpenStack and vice versa.

To this end, this provider supports the following extra specs schema:

```json
{
    "$schema": "http://cloudbase.it/garm-provider-equinix/schemas/extra_specs#",
    "type": "object",
    "description": "Schema defining supported extra specs for the Garm Equinix Metal Provider",
    "properties": {
        "metro_code": {
            "type": "string",
            "description": "The metro in which this pool will create runners."
        },
        "hardware_reservation_id": {
            "type": "string",
            "description": "The hardware reservation ID to use."
        }
    }
}
```

An example extra specs json would look like this:

```json
{
    "metro_code": "AM",
    "hardware_reservation_id": "UUID_GOES_HERE"
}
```

To set it on an existing pool, simply run:

```bash
garm-cli pool update --extra-specs='{"metro_code": "AM"}' <POOL_ID>
```

You can also set a spec when creating a new pool, using the same flag.

Workers in that pool will be created taking into account the specs you set on the pool.