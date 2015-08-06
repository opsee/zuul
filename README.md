# Zuul

Dr. Peter Venkman: What I'd really like to do is talk to Dana. Dana? It's Peter.

Dana Barrett: There is no Dana, there is only Zuul.

Dr. Peter Venkman: Oh, Zuulie, you nut, now c'mon. Just relax, c'mon. I want to talk to Dana. Dana, Dana. Can I talk to Dana?

Dana Barrett: [in an inhuman demonic voice] There is no Dana, only Zuul!

Dr. Peter Venkman: What a lovely singing voice you must have.

## Base container

The Zuul base container contains everything needed to run the OpenSSH server,
a multiplexing SSH Control Master, and SSH clients.

## Connector

In order to forward local services to the Zuul SSH server, we need a few things.
First, we want the ability to have per-connection or per-host client keys used
for authentication. These are stored in S3 and encrypted with KMS.
[s3kms](https://github.com/opsee/vinz-clortho/tree/master/README.md) is used to manage these keys as well as
retrieve them at runtime.

Running a zuul connector requires a lot of environment variables. Required unless
otherwise specified:

* `KEY_ALIAS` - This is the alias to the KMS key (alias/zuul)
* `KEY_BUCKET` - The bucket where you keep your keys (my-keys)
* `HOST_KEY_OBJECT` - This is the path within the bucket to your host key (keys/some_host_key)
* `CLIENT_KEY_OBJECT` - Same as above, but for client keys (keys/some_identity_file)
* `AWS_DEFAULT_REGION` - (Optional) AWS region (us-east-1)
* `AWS_ACCESS_KEY_ID` - (Optional) AWS access key id
* `AWS_SECRET_ACCESS_KEY` - (Optional) AWS secret key

You can then docker run your connector with this absurdly long command:

```
docker run -e "KEY_ALIAS=alias/zuul" \
           -e "KEY_BUCKET=my-keys" \
           -e "HOST_KEY_OBJECT=keys/ssh_host_rsa_key" \
           -e "CLIENT_KEY_OBJECT=keys/id_rsa" \
           -e "AWS_DEFAULT_REGION=us-west-1" \
           -e "AWS_ACCESS_KEY_ID=blahblahblah" \
           -e "AWS_SECRET_ACCESS_KEY=blahbittyblahblah" \
           --name some_host_name_connector_22_9022 \
           --rm \
           quay.io/opsee/zuul \
           connect -H some.host.name -p 22 -l 9022 -u username
```

See `connect -h` for help.

This will establish a connection to some.host.name via the [multiplexer](#multiplexer)
and forward connections made to some.host.com port 9022 to localhost (the host
running from the connector) port 22.

We embed a lot of information into the container name to avoid collisions, but
in practice they probably won't happen, because generally you don't need more
than one multiplexer. YMMV.
