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

## Creating a pool

After you [add the Equinix metal provider to garm](https://github.com/cloudbase/garm/blob/main/doc/providers.md#the-external-provider), you need to create a pool that uses it. Assuming you named your external provider as ```equinix``` in the garm config, the following command should create a new pool:

```bash
garm-cli pool create \
    --enabled=true \
    --flavor c3.small.x86 \
    --image ubuntu_22_04 \
    --min-idle-runners 0 \
    --repo f0b1c1c8-b605-4560-adb7-79b95e2e470c \
    --tags equinix,ubuntu,c3-small \
    --provider-name equinix
```

This will create a new runner pool for the repo with ID `f0b1c1c8-b605-4560-adb7-79b95e2e470c` on Equinix Metal, using the image `ubuntu_22_04` and a size of `c3.small.x86`.

The pool is set to have a minimum of zero idle runners, which means it will not attempt to spin up any runners, unless a job comes in that matches the lables of the pool.

You can, of course, tweak the values in the above command to suit your needs.

Here an example for a Windows pool:

```bash
garm-cli pool create \
   --os-type=windows \
   --enabled=true \
   --flavor=c3.small.x86 \
   --min-idle-runners=0 \
   --image windows_2022 \
   --repo f0b1c1c8-b605-4560-adb7-79b95e2e470c \
   --tags=equinix,windows2022 \
   --provider-name equinix
```

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