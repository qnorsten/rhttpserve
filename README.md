# rserve [![Build Status](https://travis-ci.org/brandur/rserve.svg?branch=master)](https://travis-ci.org/brandur/rserve)

An HTTP server that can serve files out of an rclone
remote. Includes a command line utility to generate time
expiring signed URLs (in the style of S3) that will be
verified by the server.

This has the effect of allowing files to be shared simply
(and temporarily) without having to walk through web-based
sharing prompts and the like.

## Build

    $ go install

## Setup

Generate a public/private key pair:

    $ rserve generate
    RSERVE_PUBLIC_KEY=
    RSERVE_PRIVATE_KEY=

### Server

The server needs to be configured with a public key:

    $ export RSERVE_PUBLIC_KEY=

Any rclone remotes you plan on serving files from should
also be configured in the environment:

    $ export RCLONE_CONFIG_REMOTE_TYPE="amazon cloud drive"
    $ export RCLONE_CONFIG_REMOTE_CLIENT_ID=
    $ export RCLONE_CONFIG_REMOTE_CLIENT_SECRET=
    $ export RCLONE_CONFIG_REMOTE_TOKEN=

The server can then be started with:

    $ rserve serve

It defaults to listening on port 8090, but tries to read a
value out of `PORT` if one is configured.

### Client

The client needs a private key and the host that the server
is listening on:

    $ export RSERVE_PRIVATE_KEY=
    $ export RSERVE_HOST=localhost:8090

Because you'll likely be running the client locally, it
might be useful to store these values in your `.zshrc` or
equivalent.

## Usage

With both server and client set up, it's now possible to
have rserve generate a URL for a file in your remote:

    rserve sign myremote:magazines/mercantallist.the/mercantallist.the.2015-07-04.pdf

Compose with `xargs` to sign all files in a directory:

    rclone ls -q myremote:magazines/mercantallist.the/ | awk '{print "myremote:magazines/mercantallist.the/" $2}' | xargs rserve sign --curl --skip-check

## Development

## Vendoring Dependencies

Dependencies are managed with govendor. New ones can be vendored using these
commands:

    go get -u github.com/kardianos/govendor
    govendor add +external
