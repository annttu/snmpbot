[![Build Status](https://travis-ci.com/qmsk/snmpbot.svg?branch=master)](https://travis-ci.com/qmsk/snmpbot)

# github.com/qmsk/snmpbot

SNMP client (manager) library for Go with support for SMI MIBs.

REST (HTTP/JSON) API for writing SNMP applications.

## Quickstart

The easiest way to run `snmpbot` is using the [pre-built Docker image](https://hub.docker.com/r/qmsk/snmpbot/):

```
$ docker run --net=host qmsk/snmpbot
INFO mibs: Load MIBs from directory: /opt/qmsk/snmpbot/mibs
...
INFO web: Listen on :8286...
```

The pre-build Docker images contain the [most common MIB definitions](https://github.com/qmsk/snmpbot-mibs).

### Querying a scalar object

You can then query the API at `:8286`:

#### `$ curl -s 'http://localhost:8286/api/hosts/public@localhost/objects/SNMPv2-MIB::sysDescr' | jq`
```json
{
  "ID": "SNMPv2-MIB::sysDescr",
  "Instances": [
    {
      "HostID": "public@localhost",
      "Value": "Linux tehobari 4.15.0-33-generic #36~16.04.1-Ubuntu SMP Wed Aug 15 17:21:05 UTC 2018 x86_64"
    }
  ]
}
```

### Querying multi-instance objects

You can query an object with multiple instances, and `snmpbot` will parse the OID indexes:

#### `$ curl -s 'http://localhost:8286/api/hosts/public@localhost/objects/IF-MIB::ifName' | jq`
```
{
  "ID": "IF-MIB::ifName",
  "IndexKeys": [
    "IF-MIB::ifIndex"
  ],
  "Instances": [
    {
      "HostID": "public@localhost",
      "Index": {
        "IF-MIB::ifIndex": 1
      },
      "Value": "lo"
    },
    {
      "HostID": "public@localhost",
      "Index": {
        "IF-MIB::ifIndex": 2
      },
      "Value": "enp0s31f6"
    },
    {
      "HostID": "public@localhost",
      "Index": {
        "IF-MIB::ifIndex": 4
      },
      "Value": "wlp4s0"
    },
    {
      "HostID": "public@localhost",
      "Index": {
        "IF-MIB::ifIndex": 5
      },
      "Value": "docker0"
    },
    {
      "HostID": "public@localhost",
      "Index": {
        "IF-MIB::ifIndex": 95
      },
      "Value": "wwp0s20f0u3c2"
    }
  ]
}
```

### Querying tables

You can also query tables, and `snmpbot` will group the objects into entries with indexes:

#### `curl -s 'http://localhost:8286/api/hosts/public@localhost/tables/IF-MIB::ifTable' | jq`
```json
{
  "ID": "IF-MIB::ifTable",
  "IndexKeys": [
    "IF-MIB::ifIndex"
  ],
  "ObjectKeys": [
    "IF-MIB::ifIndex",
    "IF-MIB::ifDescr",
    "IF-MIB::ifType",
    "IF-MIB::ifMtu",
    "IF-MIB::ifSpeed",
    "IF-MIB::ifPhysAddress",
    "IF-MIB::ifAdminStatus",
    "IF-MIB::ifOperStatus",
    "IF-MIB::ifLastChange",
    "IF-MIB::ifInOctets",
    "IF-MIB::ifInUcastPkts",
    "IF-MIB::ifInNUcastPkts",
    "IF-MIB::ifInDiscards",
    "IF-MIB::ifInErrors",
    "IF-MIB::ifInUnknownProtos",
    "IF-MIB::ifOutOctets",
    "IF-MIB::ifOutUcastPkts",
    "IF-MIB::ifOutNUcastPkts",
    "IF-MIB::ifOutDiscards",
    "IF-MIB::ifOutErrors",
    "IF-MIB::ifOutQLen",
    "IF-MIB::ifSpecific"
  ],
  "Entries": [
    {
      "HostID": "public@localhost",
      "Index": {
        "IF-MIB::ifIndex": 1
      },
      "Objects": {
        "IF-MIB::ifAdminStatus": "up",
        "IF-MIB::ifDescr": "lo",
        "IF-MIB::ifInDiscards": 0,
        "IF-MIB::ifInErrors": 0,
        "IF-MIB::ifInNUcastPkts": 0,
        "IF-MIB::ifInOctets": 1692789727,
        "IF-MIB::ifInUcastPkts": 1086499,
        "IF-MIB::ifInUnknownProtos": 0,
        "IF-MIB::ifIndex": 1,
        "IF-MIB::ifLastChange": 0,
        "IF-MIB::ifMtu": 65536,
        "IF-MIB::ifOperStatus": "up",
        "IF-MIB::ifOutDiscards": 0,
        "IF-MIB::ifOutErrors": 0,
        "IF-MIB::ifOutNUcastPkts": 0,
        "IF-MIB::ifOutOctets": 1692789727,
        "IF-MIB::ifOutQLen": 0,
        "IF-MIB::ifOutUcastPkts": 1086499,
        "IF-MIB::ifPhysAddress": "",
        "IF-MIB::ifSpecific": ".0.0",
        "IF-MIB::ifSpeed": 10000000,
        "IF-MIB::ifType": "softwareLoopback"
      }
    },
    ...
  ]
}
```

### Advanced usage

See the remainder of the README for more advanced usage, including:

* Querying multiple objects in parallel
* Querying multiple tables in parallel
* Configuring hosts (static TOML config file, or dynamic HTTP `POST/PUT/DELETE` API)
* Querying multiple hosts in parallel

## Requirements

### Go version 1.10

* [encoding/asn1: add MarshalWithParams](https://github.com/golang/go/commit/c32626a4ce9293979c407c4e6a799d1bec37aa18#diff-3c740598b49abcc8e5f74f801dc75255)

### Go version 1.9

* [encoding/asn1: add NullBytes and NullRawValue for working with ASN.1 NULL](https://github.com/golang/go/commit/d9b1f9e85ee097ebc95c5904cee921ba7be4f732)

### SNMP MIBs

SNMP MIBs must be pre-processed into a custom JSON format for use with `snmpbot`.

Common pre-processed MIBs can be found at [github.com/qmsk/snmpbot-mibs](https://github.com/qmsk/snmpbot-mibs):

    git clone https://github.com/qmsk/snmpbot-mibs
    export SNMPBOT_MIBS=$PWD/snmpbot-mibs

Custom MIBs can be imported using the [`mib-import.py` script](./scripts).

## Go Libraries

### `github.com/qmsk/snmpbot/snmp`

[![](https://godoc.org/github.com/qmsk/snmpbot/snmp?status.svg)](http://godoc.org/github.com/qmsk/snmpbot/snmp)

Low-level SNMP protocol support.

* Using [github.com/geoffgarside/ber](github.com/geoffgarside/ber) for BER decoding.

### `github.com/qmsk/snmpbot/client`

[![](https://godoc.org/github.com/qmsk/snmpbot/client?status.svg)](http://godoc.org/github.com/qmsk/snmpbot/client)

SNMP client with support for UDP queries

* Multiple parallel requests (goroutine-safe)
* Request timeout and retry
* Get request splitting (large numbers of OIDs)

### `github.com/qmsk/snmpbot/mibs`

[![](https://godoc.org/github.com/qmsk/snmpbot/mibs?status.svg)](http://godoc.org/github.com/qmsk/snmpbot/mibs)

SMI support for MIBs

* Initialize using `Load(path string)` to load the `.json` MIB files
* Resolving strings like `"interfaces::ifDescr"` to `*Object`
* Resolving OIDs like `ParseOID(".1.3.6.1.2.1.2.2.1.2")` to `*Object`
* Decoding SMI object `SYNTAX` to `interface{}`, including `encoding/json` support
* Decoding SMI table `INDEX` syntax from OIDs

### `github.com/qmsk/snmpbot/server`

[![](https://godoc.org/github.com/qmsk/snmpbot/server?status.svg)](http://godoc.org/github.com/qmsk/snmpbot/server)

SNMP manager with REST API.

## CLI Commands

These are mostly meant for testing.

All of the commands support the same common flags:

```
  -debug
        Log debug
  -quiet
        Do not log warnings
  -snmp-community string
        Default SNMP community (default "public")
  -snmp-maxvars uint
        Maximum request VarBinds (default 10)
  -snmp-mibs string
        Load MIBs from path (default $SNMPBOT_MIBS)
  -snmp-retry int
        SNMP request retry
  -snmp-timeout duration
        SNMP request timeout (default 1s)
  -snmp-udp-size uint
        Maximum UDP recv size (default 1500)
  -verbose
        Log info
```

Apart from the `snmpbot` command, the first argument is a SNMP address of the form `[COMMUNITY@]HOST[:PORT]`, and the remainder are SMI objects of the form `[MIB[::OBJECT]][.INDEX]`.

### `github.com/qmsk/snmpbot/cmd/snmpget`

Testing `GetRequest`; note that non-tabular SNMP objects cannot be fetched without the `.0` instance suffix.

#### `snmpget public@edgeswitch-098730 system::sysDescr`
```
2017/11/18 22:02:49 VarBind[.1.3.6.1.2.1.1.1](system::sysDescr): SNMP VarBind Error: NoSuchInstance
```

#### `snmpget public@edgeswitch-098730 system::sysDescr.0`
```
system::sysDescr.0 = EdgeSwitch 24-Port Lite, 1.7.0.4922887, Linux 3.6.5-f4a26ed5, 0.0.0.0000000
```

### `github.com/qmsk/snmpbot/cmd/snmpwalk`

Testing `GetNextRequest`

#### `snmpwalk public@edgeswitch-098730 system::sysDescr.0`
```
```

The output is empty, because there is nothing in the `system::sysDescr` tree after `.0`.

#### `snmpwalk public@edgeswitch-098730 system::sysDescr`
```
system::sysDescr.0 = EdgeSwitch 24-Port Lite, 1.7.0.4922887, Linux 3.6.5-f4a26ed5, 0.0.0.0000000
```

#### `snmpwalk public@edgeswitch-098730 system`
```
system::sysDescr.0 = EdgeSwitch 24-Port Lite, 1.7.0.4922887, Linux 3.6.5-f4a26ed5, 0.0.0.0000000
system::sysObjectID.0 = .1.3.6.1.4.1.4413
system::sysUpTime.0 = 2352h54m5s
system::sysContact.0 =
system::sysName.0 = UBNT EdgeSwitch
system::sysLocation.0 =
system::sysServices.0 = 6
system::sysORLastChange.0 = 1m7s
...
```

#### `snmpwalk public@edgeswitch-098730 Q-BRIDGE-MIB::dot1qVlanStaticName`

```
Q-BRIDGE-MIB::dot1qVlanStaticName[1] = default
Q-BRIDGE-MIB::dot1qVlanStaticName[2] = wan
Q-BRIDGE-MIB::dot1qVlanStaticName[1002] = iot
Q-BRIDGE-MIB::dot1qVlanStaticName[1003] = guest
```

Object indexes are decoded.

### `github.com/qmsk/snmpbot/cmd/snmpobject`

Very similar to `snmpwalk`, but implemented slightly differently - only works for known objects, not arbitrary OIDs.

#### `snmpobject public@edgeswitch-098730 system::sysDescr interfaces::ifDescr`

```
system::sysDescr = EdgeSwitch 24-Port Lite, 1.7.0.4922887, Linux 3.6.5-f4a26ed5, 0.0.0.0000000
interfaces::ifDescr[1] = Slot: 0 Port: 1 Gigabit - Level
interfaces::ifDescr[2] = Slot: 0 Port: 2 Gigabit - Level
interfaces::ifDescr[3] = Slot: 0 Port: 3 Gigabit - Level
interfaces::ifDescr[4] = Slot: 0 Port: 4 Gigabit - Level
```

#### `snmpobject public@edgeswitch-098730 Q-BRIDGE-MIB::dot1qTpFdbPort`
```
Q-BRIDGE-MIB::dot1qTpFdbPort[1][44:xx:xx:xx:xx:19] = 24
Q-BRIDGE-MIB::dot1qTpFdbPort[1][50:xx:xx:xx:xx:b2] = 2
Q-BRIDGE-MIB::dot1qTpFdbPort[1][70:xx:xx:xx:xx:e4] = 23
Q-BRIDGE-MIB::dot1qTpFdbPort[1][ac:xx:xx:xx:xx:c9] = 6
Q-BRIDGE-MIB::dot1qTpFdbPort[1][f0:xx:xx:xx:xx:30] = 65
Q-BRIDGE-MIB::dot1qTpFdbPort[1][f0:xx:xx:xx:xx:3f] = 24
Q-BRIDGE-MIB::dot1qTpFdbPort[1][f0:xx:xx:xx:xx:45] = 24
Q-BRIDGE-MIB::dot1qTpFdbPort[2][00:xx:xx:xx:xx:27] = 24
Q-BRIDGE-MIB::dot1qTpFdbPort[2][70:xx:xx:xx:xx:e4] = 23
Q-BRIDGE-MIB::dot1qTpFdbPort[1002][74:xx:xx:xx:xx:c6] = 24
Q-BRIDGE-MIB::dot1qTpFdbPort[1002][f0:xx:xx:xx:xx:45] = 24
```

### `github.com/qmsk/snmpbot/cmd/snmpprobe`

Use `GetNextRequest` to probe the presence of subtrees, including entire MIBs or objects within MIBs

#### `snmpprobe public@edgeswitch-098730 BRIDGE-MIB BRIDGE-MIB::dot1dBaseBridgeAddress BRIDGE-MIB::dot1dStpPortTable`
```
BRIDGE-MIB = true
BRIDGE-MIB::dot1dBaseBridgeAddress = true
BRIDGE-MIB::dot1dStpPortTable = false
```

### `github.com/qmsk/snmpbot/cmd/snmptable`

Use `GetNextRequest` to walk and decode SMI tables

#### `snmptable public@edgeswitch-098730 interfaces::ifTable`
```
interfaces::ifIndex |       interfaces::ifIndex interfaces::ifDescr                interfaces::ifType interfaces::ifMtu interfaces::ifSpeed interfaces::ifPhysAddress interfaces::ifAdminStatus interfaces::ifOperStatus interfaces::ifLastChange
---                 |       ---                 ---                                ---                ---               ---                 ---                       ---                       ---                      ---
1                   |       1                   Slot: 0 Port: 1 Gigabit - Level    6                  1518              0                   f0:9f:c2:09:87:31         up                        down                     2257h16m19s
2                   |       2                   Slot: 0 Port: 2 Gigabit - Level    6                  1518              1000000000          f0:9f:c2:09:87:31         up                        up                       2342h30m20s
3                   |       3                   Slot: 0 Port: 3 Gigabit - Level    6                  1518              0                   f0:9f:c2:09:87:31         up                        down                     525h50m17s
4                   |       4                   Slot: 0 Port: 4 Gigabit - Level    6                  1518              0                   f0:9f:c2:09:87:31         up                        down                     0s
5                   |       5                   Slot: 0 Port: 5 Gigabit - Level    6                  1518              0                   f0:9f:c2:09:87:31         up                        down                     1200h48m53s
6                   |       6                   Slot: 0 Port: 6 Gigabit - Level    6                  1518              100000000           f0:9f:c2:09:87:31         up                        up                       2318h40m45s
7                   |       7                   Slot: 0 Port: 7 Gigabit - Level    6                  1518              0                   f0:9f:c2:09:87:31         up                        down                     0s
8                   |       8                   Slot: 0 Port: 8 Gigabit - Level    6                  1518              0                   f0:9f:c2:09:87:31         up                        down                     0s
9                   |       9                   Slot: 0 Port: 9 Gigabit - Level    6                  1518              10000000            f0:9f:c2:09:87:31         up                        up                       2331h13m45s
10                  |       10                  Slot: 0 Port: 10 Gigabit - Level   6                  1518              0                   f0:9f:c2:09:87:31         up                        down                     0s
11                  |       11                  Slot: 0 Port: 11 Gigabit - Level   6                  1518              0                   f0:9f:c2:09:87:31         up                        down                     0s
12                  |       12                  Slot: 0 Port: 12 Gigabit - Level   6                  1518              0                   f0:9f:c2:09:87:31         up                        down                     0s
13                  |       13                  Slot: 0 Port: 13 Gigabit - Level   6                  1518              0                   f0:9f:c2:09:87:31         up                        down                     0s
14                  |       14                  Slot: 0 Port: 14 Gigabit - Level   6                  1518              0                   f0:9f:c2:09:87:31         up                        down                     0s
15                  |       15                  Slot: 0 Port: 15 Gigabit - Level   6                  1518              0                   f0:9f:c2:09:87:31         up                        down                     0s
16                  |       16                  Slot: 0 Port: 16 Gigabit - Level   6                  1518              0                   f0:9f:c2:09:87:31         up                        down                     0s
17                  |       17                  Slot: 0 Port: 17 Gigabit - Level   6                  1518              100000000           f0:9f:c2:09:87:31         up                        up                       1m17s
18                  |       18                  Slot: 0 Port: 18 Gigabit - Level   6                  1518              0                   f0:9f:c2:09:87:31         up                        down                     0s
19                  |       19                  Slot: 0 Port: 19 Gigabit - Level   6                  1518              0                   f0:9f:c2:09:87:31         up                        down                     0s
20                  |       20                  Slot: 0 Port: 20 Gigabit - Level   6                  1518              0                   f0:9f:c2:09:87:31         up                        down                     0s
21                  |       21                  Slot: 0 Port: 21 Gigabit - Level   6                  1518              0                   f0:9f:c2:09:87:31         up                        down                     0s
22                  |       22                  Slot: 0 Port: 22 Gigabit - Level   6                  1518              0                   f0:9f:c2:09:87:31         up                        down                     0s
23                  |       23                  Slot: 0 Port: 23 Gigabit - Level   6                  1518              1000000000          f0:9f:c2:09:87:31         up                        up                       1551h25m59s
24                  |       24                  Slot: 0 Port: 24 Gigabit - Level   6                  1518              1000000000          f0:9f:c2:09:87:31         up                        up                       382h12m33s
25                  |       25                  Slot: 0 Port: 25 Gigabit - Level   6                  1518              0                   f0:9f:c2:09:87:31         up                        down                     0s
26                  |       26                  Slot: 0 Port: 26 Gigabit - Level   6                  1518              0                   f0:9f:c2:09:87:31         up                        down                     0s
65                  |       65                   CPU Interface for Slot: 5 Port: 1 1                  1518              0                   f0:9f:c2:09:87:30         up                        up                       0s
66                  |       66                   Link Aggregate 1                  1                  1518              0                   f0:9f:c2:09:87:31         up                        down                     0s
67                  |       67                   Link Aggregate 2                  1                  1518              0                   f0:9f:c2:09:87:31         up                        down                     0s
68                  |       68                   Link Aggregate 3                  1                  1518              0                   f0:9f:c2:09:87:31         up                        down                     0s
69                  |       69                   Link Aggregate 4                  1                  1518              0                   f0:9f:c2:09:87:31         up                        down                     0s
70                  |       70                   Link Aggregate 5                  1                  1518              0                   f0:9f:c2:09:87:31         up                        down                     0s
71                  |       71                   Link Aggregate 6                  1                  1518              0                   f0:9f:c2:09:87:31         up                        down                     0s
```

#### `snmptable public@edgeswitch-098730 LLDP-MIB::lldpRemTable`
```
LLDP-MIB::lldpRemTimeMark LLDP-MIB::lldpRemLocalPortNum LLDP-MIB::lldpRemIndex |       LLDP-MIB::lldpRemChassisIdSubtype LLDP-MIB::lldpRemChassisId LLDP-MIB::lldpRemPortIdSubtype LLDP-MIB::lldpRemPortId LLDP-MIB::lldpRemPortDesc LLDP-MIB::lldpRemSysName LLDP-MIB::lldpRemSysDesc
---                       ---                           ---                    |       ---                               ---                        ---                            ---                     ---                       ---                      ---
21h59m32.24s              24                            1                      |       macAddress                        f0 9f c2 64 6d 45          macAddress                     f0 9f c2 64 6d 45       switch0                   erx-home                 UBNT EdgeRouter X SFP 6-Port running on v1.9.1.1.4977602.170427.0113
21h59m32.24s              24                            2                      |       macAddress                        f0 9f c2 64 6d 45          macAddress                     f0 9f c2 64 6d 3f       eth0                      erx-home                 UBNT EdgeRouter X SFP 6-Port running on v1.9.1.1.4977602.170427.0113
```

## `github.com/qmsk/snmpbot/cmd/snmpbot`

This is the command using the [`server`](#server) to provide the HTTP REST API.

### Usage

```
$ $GOPATH/bin/snmpbot -help
  -config string
        Load TOML config
  -debug
        Log debug
  -http-listen string
        HTTP server listen: [HOST]:PORT (default ":8286")
  -http-static string
        HTTP sever /static path: PATH
  -quiet
        Do not log warnings
  -snmp-community string
        Default SNMP community (default "public")
  -snmp-maxvars uint
        Maximum request VarBinds (default 10)
  -snmp-mibs string
        Load MIBs from path (default $SNMPBOT_MIBS)
  -snmp-retry int
        SNMP request retry
  -snmp-timeout duration
        SNMP request timeout (default 1s)
  -snmp-udp-size uint
        Maximum UDP recv size (default 1500)
  -verbose
        Log info
```

Examples:

    $GOPATH/bin/snmpbot -config config.toml -http-listen :8286 -snmp-udp-size 9000 -snmp-maxvars 50

## Config

Loaded using `snmpbot -config`

```toml
[hosts.edgeswitch-098730]
SNMP = "public@edgeswitch-098730"
Location = "test"

[hosts.erx-home]
SNMP = "secret@erx-home"
Location = "home"
```

The configuration file is optional, dynamic hosts can be queried without any config, using `GET /hosts/...?snmp=community@host` (also `-snmp-community=...`).

***NOTE***: The mass-querying `/objects/...` and `/tables/...` endpoints only query configured objects.

## API

See the [`api`](https://godoc.org/github.com/qmsk/snmpbot/api) package docs for the exact details

#### `GET /api/`
```json
{
   "Hosts" : [
      {
         "ID" : "edgeswitch-098730",
         "SNMP" : "public@172.28.2.2:161"
      },
      {
         "ID" : "erx-home",
         "SNMP" : "secret@172.28.0.1:161"
      }
   ],
   "MIBs" : [
      {
         "ID" : "IF-MIB"
      },
      ...
   ],
   "Tables" : [
      {
         "ID" : "Q-BRIDGE-MIB::dot1qVlanStaticTable",
         "ObjectKeys" : [
            "Q-BRIDGE-MIB::dot1qVlanStaticName",
            "Q-BRIDGE-MIB::dot1qVlanStaticEgressPorts",
            "Q-BRIDGE-MIB::dot1qVlanForbiddenEgressPorts",
            "Q-BRIDGE-MIB::dot1qVlanStaticUntaggedPorts",
            "Q-BRIDGE-MIB::dot1qVlanStaticRowStatus"
         ],
         "IndexKeys" : [
            "Q-BRIDGE-MIB::dot1qVlanIndex"
         ]
      },
      ...
    ],
   "Objects" : [
      {
         "IndexKeys" : [
            "interfaces::ifIndex"
         ],
         "ID" : "interfaces::ifDescr"
      },
      {
         "ID" : "system::sysContact"
      },
      ...
   ]
}
```

#### `GET /api/hosts/`

Query configured hosts.

```json
[
   {
      "ID" : "edgeswitch-098730",
      "SNMP" : "public@172.28.2.2:161"
   },
   ...
]
```

#### `POST /api/hosts/`

Configure a new host. The new host will show up in `/hosts/` and can be queried via the `/objects/...` and `/tables/...` endpoints.

```json
{
   "ID" : "edgeswitch-098730",
   "SNMP" : "public@172.28.2.2:161",
   "Online": true,
   ...
}
```

Response is as for `GET`.

Returns HTTP 409 Conflict if the host is already configured. Use `PUT /v1/hosts/:id` instead.

##### Request `Content-Type: application/json`
```json
{
  "ID": "test",
  "SNMP": "community@test.example.com",
  "Location": "testing"
}
```

##### Request `Content-Type: application/x-www-form-urlencoded`
```
id=test&snmp=community@test.example.com&location=testing
```

#### `GET /api/hosts/:id`

Query information about a specific configured host.

```json
{
   "ID" : "edgeswitch-098730",
   "SNMP" : "public@172.28.2.2:161",
   "Location": "testing",
   "Online": true
}
```

***Note***: Does not return anything useful for a dynamic host.

#### `PUT /api/hosts/:id`

Add or replace a configured host.

```json
{
   "ID" : "edgeswitch-098730",
   "SNMP" : "public@172.28.2.2:161",
   "Online": true,
   ...
}
```

Response is as for `GET`.

##### Request `Content-Type: application/json`
```json
{
  "SNMP": "community@test.example.com",
  "Location": "testing"
}
```

##### Request `Content-Type: application/x-www-form-urlencoded`
```
snmp=community@test.example.com&location=testing
```

#### `DELETE /api/hosts/:id`

Remove a configured host.

Returns HTTP 204 No Content.

#### `GET /api/hosts/test/...?snmp=community@test.example.com`

Query a dynamic host using an arbitrary `?snmp=[<community> "@"] <host>` target.

The `[<community> "@"]` is optional, and defaults to the `-config` => `[ClientOptions] Community=` or `-snmp-community`.

```json
{
   "MIBs" : [
      {
         "ID" : "SNMPv2-MIB"
      },
      ...
   ],
   "ID" : "test",
   "SNMP" : "community@192.0.2.1:161"
}
```

The given SNMP host will be probed for supported MIBs.

***NOTE***: The `?snmp=` query parameter works for all host API paths, but is ignored for configured hosts.

#### `GET /api/hosts/:host/`

Query host with information about tables/objects for probed MIBs.

```json
{
   "ID" : "edgeswitch-098730",
   "Objects" : [
      {
         "ID" : "system::sysObjectID"
      },
      ...
   ],
   "Tables" : [
      {
         "IndexKeys" : [
            "BRIDGE-MIB::dot1dStpPort"
         ],
         "ObjectKeys" : [
            "BRIDGE-MIB::dot1dStpPort",
            "BRIDGE-MIB::dot1dStpPortPriority",
            "BRIDGE-MIB::dot1dStpPortState",
            "BRIDGE-MIB::dot1dStpPortEnable",
            "BRIDGE-MIB::dot1dStpPortPathCost",
            "BRIDGE-MIB::dot1dStpPortDesignatedRoot",
            "BRIDGE-MIB::dot1dStpPortDesignatedCost",
            "BRIDGE-MIB::dot1dStpPortDesignatedBridge",
            "BRIDGE-MIB::dot1dStpPortDesignatedPort",
            "BRIDGE-MIB::dot1dStpPortForwardTransitions"
         ],
         "ID" : "BRIDGE-MIB::dot1dStpPortTable"
      },
      ...
   ],
   "MIBs" : [
      {
         "ID" : "SNMPv2-MIB"
      },
      ...
   ],
   "SNMP" : "public@172.28.2.2:161"
}
```

#### `GET /api/hosts/:host/tables/LLDP-MIB::lldpRemTable`

Query a specific table for a specific host (dynamic or configured).

```json
{
   "ObjectKeys" : [
      ...
   ],
   "IndexKeys" : [
      ...
   ],
   "Entries" : [
      ...
   ],
   "ID" : "LLDP-MIB::lldpRemTable"
}
```

***Note***: The queried table does not necessarily need to belong to a probed MIB.

#### `GET /api/hosts/:host/tables/?table=LLDP-MIB::*`

Query multiple tables from probed mibs for a specific host (dynamic or configured).

Use multiple `?table=A&table=B` filters to query different tables.

Use (multiple) `?host=test-*` query parameters to limit queried hosts.

```json
[
   {
      "ObjectKeys" : [
         "LLDP-MIB::lldpLocPortIdSubtype",
         "LLDP-MIB::lldpLocPortId",
         "LLDP-MIB::lldpLocPortDesc"
      ],
      "Entries" : [
         {
            "HostID" : "edgeswitch-098730",
            "Index" : {
               "LLDP-MIB::lldpLocPortNum" : 1
            },
            "Objects" : {
               "LLDP-MIB::lldpLocPortIdSubtype" : "macAddress",
               "LLDP-MIB::lldpLocPortId" : "f0 9f c2 09 87 31",
               "LLDP-MIB::lldpLocPortDesc" : ""
            }
         },
         ...
      ],
      "ID" : "LLDP-MIB::lldpLocPortTable",
      "IndexKeys" : [
         "LLDP-MIB::lldpLocPortNum"
      ]
   },
   {
      "ObjectKeys" : [
         "LLDP-MIB::lldpRemChassisIdSubtype",
         "LLDP-MIB::lldpRemChassisId",
         "LLDP-MIB::lldpRemPortIdSubtype",
         "LLDP-MIB::lldpRemPortId",
         "LLDP-MIB::lldpRemPortDesc",
         "LLDP-MIB::lldpRemSysName",
         "LLDP-MIB::lldpRemSysDesc"
      ],
      "Entries" : [
         {
            "Index" : {
               "LLDP-MIB::lldpRemTimeMark" : 79172.24,
               "LLDP-MIB::lldpRemLocalPortNum" : 24,
               "LLDP-MIB::lldpRemIndex" : 1
            },
            "HostID" : "edgeswitch-098730",
            "Objects" : {
               "LLDP-MIB::lldpRemChassisIdSubtype" : "macAddress",
               "LLDP-MIB::lldpRemSysDesc" : "UBNT EdgeRouter X SFP 6-Port running on v1.9.1.1.4977602.170427.0113",
               "LLDP-MIB::lldpRemPortDesc" : "switch0",
               "LLDP-MIB::lldpRemSysName" : "erx-home",
               "LLDP-MIB::lldpRemPortIdSubtype" : "macAddress",
               "LLDP-MIB::lldpRemChassisId" : "f0 9f c2 64 6d 45",
               "LLDP-MIB::lldpRemPortId" : "f0 9f c2 64 6d 45"
            }
         },
         ...
      ],
      "ID" : "LLDP-MIB::lldpRemTable",
      "IndexKeys" : [
         "LLDP-MIB::lldpRemTimeMark",
         "LLDP-MIB::lldpRemLocalPortNum",
         "LLDP-MIB::lldpRemIndex"
      ]
   }
]
```

***Note***: The queried tables must belong to a probed MIB.

#### `GET /api/hosts/:host/objects/SNMPv2-MIB::sysDescr`

Query a specific object for a specific host (configured or dynamic).

```json
{
   "Instances" : [
      {
         "Value" : "EdgeSwitch 24-Port Lite, 1.7.0.4922887, Linux 3.6.5-f4a26ed5, 0.0.0.0000000",
         "HostID" : "edgeswitch-098730"
      }
   ],
   "ID" : "SNMPv2-MIB::sysDescr"
}
```

***Note***: The queried object does not necessarily need to belong to a probed MIB.

***Note***: Objects belonging to table entries will return multiple instances with different `Index` values.

#### `GET /api/hosts/:host/objects/?object=SNMPv2-MIB::sys*&object=IF-MIB::ifDescr`

Query matching objects from probed MIBs for a specific host (dynamic or configured).

Use multiple `?object=A&object=B` filters to query different objects. Use `?object=MIB::foo*` filters to query multiple matching objects.

Use `?table=IF-MIB::ifTable` to query for objects belonging to matching tables. These can also be combined with `?object=` filters, e.g. `GET /api/hosts/public@localhost/objects/?table=IF-MIB::ifXTable&object=*HC*Pkts` to return 64-bit per-interface packet counters as separate objects.

```json
   {
      "IndexKeys" : [
         "interfaces::ifIndex"
      ],
      "ID" : "interfaces::ifDescr",
      "Instances" : [
         {
            "Value" : "Slot: 0 Port: 1 Gigabit - Level",
            "HostID" : "edgeswitch-098730",
            "Index" : {
               "interfaces::ifIndex" : 1
            }
         },
         {
            "HostID" : "edgeswitch-098730",
            "Value" : "Slot: 0 Port: 2 Gigabit - Level",
            "Index" : {
               "interfaces::ifIndex" : 2
            }
         },
         ...
      ]
   },
   {
      "ID" : "system::sysDescr",
      "Instances" : [
         {
            "HostID" : "edgeswitch-098730",
            "Value" : "EdgeSwitch 24-Port Lite, 1.7.0.4922887, Linux 3.6.5-f4a26ed5, 0.0.0.0000000"
         }
      ]
   }
]
```

***Note***: The queried object must belong to a probed MIB.

#### `GET /api/tables/IF-MIB::ifTable?host=test-*`

Query specific table across all hosts.

Use (multiple) `?host=test-*` query parameters to limit queried hosts.

```json
{
   "ID" : "IF-MIB::ifTable",
   "IndexKeys" : [
      "IF-MIB::ifIndex"
   ],
   "ObjectKeys" : [
      "IF-MIB::ifIndex",
      "IF-MIB::ifDescr",
      "IF-MIB::ifType",
      "IF-MIB::ifMtu",
      "IF-MIB::ifSpeed",
      "IF-MIB::ifPhysAddress",
      "IF-MIB::ifAdminStatus",
      "IF-MIB::ifOperStatus",
      "IF-MIB::ifLastChange"
   ],
   "Entries" : [
      {
         "Index" : {
            "IF-MIB::ifIndex" : 1
         },
         "Objects" : {
            "IF-MIB::ifPhysAddress" : "",
            "IF-MIB::ifLastChange" : 0,
            "IF-MIB::ifMtu" : 65536,
            "IF-MIB::ifType" : 24,
            "IF-MIB::ifSpeed" : 10000000,
            "IF-MIB::ifIndex" : 1,
            "IF-MIB::ifAdminStatus" : "up",
            "IF-MIB::ifDescr" : "lo",
            "IF-MIB::ifOperStatus" : "up"
         },
         "HostID" : "erx-home"
      },
      ...
      {
         "HostID" : "edgeswitch-098730",
         "Objects" : {
            "IF-MIB::ifMtu" : 1518,
            "IF-MIB::ifLastChange" : 8433020,
            "IF-MIB::ifPhysAddress" : "f0:9f:c2:09:87:31",
            "IF-MIB::ifSpeed" : 1000000000,
            "IF-MIB::ifType" : 6,
            "IF-MIB::ifDescr" : "Slot: 0 Port: 2 Gigabit - Level",
            "IF-MIB::ifAdminStatus" : "up",
            "IF-MIB::ifIndex" : 2,
            "IF-MIB::ifOperStatus" : "up"
         },
         "Index" : {
            "IF-MIB::ifIndex" : 2
         }
      },
   ]
}
```

***Note***: Only configured hosts are queried.

***Note***: The table contains entries from all hosts for that SNMP table: Different `Entries` in the same table can have different `HostID` values.

#### `GET /api/tables/?table=LLDP-MIB::lldpRemTable&table=Q-BRIDGE-MIB::dot1qTpFdbTable`

Query multiple tables across all hosts.

Use multiple `?table=A&table=B` filters to query different tables.

Use (multiple) `?host=test-*` query parameters to limit queried hosts.


```json
[
   {
      "ID" : "LLDP-MIB::lldpRemTable",
      ...
   },
   {
      "ID" : "Q-BRIDGE-MIB::dot1qTpFdbTable",
      ...
   }
]
```

***Note***: Only configured hosts are queried.

***Note***: Each table contains entries from all hosts for that SNMP table: Different `Entries` in the same table can have different `HostID` values.

#### `GET /api/objects/IF-MIB::ifDescr`

Query specific object across all hosts. Use `?host=test-*` to filter queried hosts.

```json
{
   "IndexKeys" : [
      "IF-MIB::ifIndex"
   ],
   "ID" : "IF-MIB::ifDescr",
   "Instances" : [
      {
         "HostID" : "erx-home",
         "Value" : "lo",
         "Index" : {
            "IF-MIB::ifIndex" : 1
         }
      },
      {
         "Index" : {
            "IF-MIB::ifIndex" : 1
         },
         "Value" : "Slot: 0 Port: 1 Gigabit - Level",
         "HostID" : "edgeswitch-098730"
      },
      ...
   ]
}
```

***Note***: Only configured hosts are queried.

#### `GET /api/objects/IF-MIB::ifDescr?host=edgeswitch-*`

Query specific object across all hosts. Use `?host=...` to filter queried hosts.

```json
{
   "Instances" : [
      {
         "Index" : {
            "IF-MIB::ifIndex" : 1
         },
         "HostID" : "edgeswitch-098730",
         "Value" : "Slot: 0 Port: 1 Gigabit - Level"
      },
      {
         "HostID" : "edgeswitch-098730",
         "Index" : {
            "IF-MIB::ifIndex" : 2
         },
         "Value" : "Slot: 0 Port: 2 Gigabit - Level"
      },
      {
         "Value" : "Slot: 0 Port: 3 Gigabit - Level",
         "HostID" : "edgeswitch-098730",
         "Index" : {
            "IF-MIB::ifIndex" : 3
         }
      },
      ...
      {
         "HostID" : "edgeswitch-098730",
         "Index" : {
            "IF-MIB::ifIndex" : 71
         },
         "Value" : " Link Aggregate 6"
      }
   ],
   "ID" : "IF-MIB::ifDescr",
   "IndexKeys" : [
      "IF-MIB::ifIndex"
   ]
}
```

***Note***: Only configured hosts are queried.

#### `GET /api/objects/?object=SNMPv2-MIB::*`

Query multiple matching objects across all hosts.

Use multiple `?object=A&object=B` filters to query different objects.

Use (multiple) `?host=test-*` query parameters to limit queried hosts.

```json
[
   {
      "ID" : "SNMPv2-MIB::sysLocation",
      "Instances" : [
         {
            "Value" : "home.qmsk.net",
            "HostID" : "erx-home"
         },
         {
            "Value" : "",
            "HostID" : "edgeswitch-098730"
         }
      ]
   },
   {
      "Instances" : [
         {
            "HostID" : "erx-home",
            "Value" : "EdgeOS v1.9.1.1.4977602.170427.0113"
         },
         {
            "HostID" : "edgeswitch-098730",
            "Value" : "EdgeSwitch 24-Port Lite, 1.7.0.4922887, Linux 3.6.5-f4a26ed5, 0.0.0.0000000"
         }
      ],
      "ID" : "SNMPv2-MIB::sysDescr"
   },
   ...
]
```

***Note***: Only configured hosts are queried.
