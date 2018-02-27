# gnmi

`gnmi` is a command-line client for interacting with a
[gNMI service](https://github.com/openconfig/reference/tree/master/rpc/gnmi).

# Installation

After installing [Go](https://golang.org/dl/) run:

```
go get github.com/aristanetworks/goarista/cmd/gnmi
```

This will install the `gnmi` binary in the `bin` directory
under [GOPATH](https://golang.org/doc/code.html#GOPATH).

# Usage

## Options

* `-addr ADDR:PORT`  
Address of the gNMI endpoint (REQUIRED)
* `-username USERNAME`  
Username to authenticate with
* `-password PASSWORD`  
Password to authenticate with
* `-tls`  
Enable TLS
* `-cafile PATH`  
Path to server TLS certificate file
* `-certfile PATH`  
Path to client TLS certificate file
* `-keyfile PATH`  
Path to client TLS private key file

## Operations

`gnmi` supports the following operations: `capabilites`, `get`,
`subscribe`, `update`, `replace`, and `delete`.

### capabilities

`capabilities` prints the result of calling the
[Capabilities gNMI RPC](https://github.com/openconfig/reference/blob/master/rpc/gnmi/gnmi-specification.md#32-capability-discovery).

Example:

```
$ gnmi [OPTIONS] capabilities
```

### get

`get` requires a path and calls the
[Get gNMI RPC](https://github.com/openconfig/reference/blob/master/rpc/gnmi/gnmi-specification.md#222-paths).

Example:

Get all configuration in the default network instance:
```
$ gnmi [OPTIONS] get '/network-instances/network-instance[name=default]'
```

### subscribe

`subscribe` requires a path and calls the
[Subscribe gNMI RPC](https://github.com/openconfig/reference/blob/master/rpc/gnmi/gnmi-specification.md#35-subscribing-to-telemetry-updates).
This command will continuously print out results until signalled to
exit, for example by typing `Ctrl-C`.

Example:

Subscribe to interface counters:
```
$ gnmi [OPTIONS] subscribe '/interfaces/interface[name=*]/state/counters'
```

### update/replace/delete

`update`, `replace`, and `delete` are used to
[modify the configuration of a gNMI endpoint](https://github.com/openconfig/reference/blob/master/rpc/gnmi/gnmi-specification.md#34-modifying-state).
All of these operations take a path that must specify a single node
element. In other words all list members must be fully-specified.

`delete` takes a path and will delete that path.

Example:

Delete BGP configuration in the default network instance:
```
$ gnmi [OPTIONS] delete '/network-instances/network-instance[name=default]/protocols/protocol[name=BGP][identifier=BGP]/'
```

`update` and `replace` both take a path and a value in JSON
format. See
[here](https://github.com/openconfig/reference/blob/master/rpc/gnmi/gnmi-specification.md#344-modes-of-update-replace-versus-update)
for documentation on the differences between `update` and `replace`.

Examples:

Disable interface Ethernet3/42:
```
gnmi [OPTIONS] update '/interfaces/interface[name=Ethernet3/42]/config/enabled' 'false'
```

Replace the BGP global configuration:
```
gnmi [OPTIONS] replace '/network-instances/network-instance[name=default]/protocols/protocol[name=BGP][identifier=BGP]/bgp/global' '{"config":{"as": 1234, "router-id": "1.2.3.4"}}'
```

Note: String values need to be quoted if they look like JSON. For example, setting the login banner to `tor[13]`:
```
gnmi [OPTIONS] update '/system/config/login-banner '"tor[13]"'
```

### CLI requests
`gnmi` offers the ability to send CLI text inside an `update` or
`replace` operation. This is achieved by doing an `update` or
`replace` and using `"cli"` as the path and a set of configure-mode
CLI commands separated by `\n`.

Example:

Configure the idle-timeout on SSH connections
```
gnmi [OPTIONS] update 'cli' 'management ssh
idle-timeout 300'
```

## Paths

Paths in `gnmi` use a simplified xpath style. Path elements are
separated by `/`. Selectors may be used on list to select certain
members. Selectors are of the form `[key-leaf=value]`. All members of a
list may be selected by not specifying any selectors, or by using a
`*` as the value in a selector. The following are equivalent:

* `/interfaces/interface`
* `/interfaces/interface[name=*]`

All characters, including `/` are allowed inside a selector value. The
character `]` must be escaped, for example `[key=[\]]` selects the
element in the list whose `key` leaf is value `[]`.

See more examples of paths in the examples above.

See
[here](https://github.com/openconfig/reference/blob/master/rpc/gnmi/gnmi-specification.md#222-paths)
for more information.
