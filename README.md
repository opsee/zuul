# Zuul

Dr. Peter Venkman: What I'd really like to do is talk to Dana. Dana? It's Peter.

Dana Barrett: There is no Dana, there is only Zuul.

Dr. Peter Venkman: Oh, Zuulie, you nut, now c'mon. Just relax, c'mon. I want to talk to Dana. Dana, Dana. Can I talk to Dana?

Dana Barrett: [in an inhuman demonic voice] There is no Dana, only Zuul!

Dr. Peter Venkman: What a lovely singing voice you must have.

## Base container

The Zuul base container contains everything needed to run the OpenSSH server,
a multiplexing SSH Control Master, and SSH clients.

## Keys

In order to forward local services to the Zuul SSH server, we need a few things.
First, we want the ability to have per-connection or per-host client keys used
for authentication. These are stored in S3 and encrypted with KMS.
[s3kms](https://github.com/opsee/vinz-clortho/tree/master/README.md) is used to manage these keys as well as
retrieve them at runtime.

## Environment Variables

All components need the following set of environment variables to run.

* KEY_BUCKET
* AWS_DEFAULT_REGION
* AWS_ACCESS_KEY_ID
* AWS_SECRET_ACCESS_KEY

When running this in production, it is recommended that you use
[ec2-env](https://github.com/opsee/ec2-env)
or something like it to setup the environment. If running locally,
export these variables in your environment for more convenient execution
when using `docker run`.

## Server

### Environment Varialbes

* SERVER_PRIVATE_KEY_OBJECT
* CLIENT_PUBLIC_KEY_OBJECT

### Running

For the sake of convenience, we assume you are running things locally. The
following will run a server, export the listening SSHd on port 4022 and name
the container "zuul-server".

```
docker run --rm --name "zuul-server" -P 22:4022 quay.io/opsee/zuul server
```

## Multiplexer

The multiplexer can be used to multiplex connections to each running Zuul server.
This is largely to keep down the number of SSH connections between hosts.

### Environment Variables



### Running

## Connector

### Environment Variables

* CLIENT_PRIVATE_KEY_OBJECT
* SERVER_PUBLIC_KEY_OBJECT

### Running

In order to run Zuul, you'll need the hostname/IP address and port of the running
multiplexer. Assuming you're running this all locally, you can link the
containers and use the container name.

```
docker run --rm quay.io/opsee/zuul connect \
  -H some.host.name -p 22 -l 9022 -u username
```

See `connect -h` for help.

This will establish a connection to some.host.name via the [multiplexer](#multiplexer)
and forward connections made to some.host.com port 9022 to localhost (the host
running from the connector) port 22.

We embed a lot of information into the container name to avoid collisions, but
in practice they probably won't happen, because generally you don't need more
than one multiplexer. YMMV.

## Development

### Contributing

Fork, branch, pull-request.

### Conventions

UPPER_CASE is reserved for variables that come from the environment--all of which
should be given a default value via parameter expansion:

```
UPPER_CASE=${UPPER_CASE:-"some value"}
```

lower_case is for variables used within zuul scripts.

TODO: Convention for "exported" variables in common.sh/client.sh ?
