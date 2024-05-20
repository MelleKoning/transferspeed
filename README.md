# Transferspeed

Transferspeed is a test repository to see if we can benchmark transferring an arbitrary binary data over grpc http/2.

## Working with the code

The software is written in golang, and uses default libraries for generating the grpc communication. It is assumed you have the appropriate `protoc` and `protoc-gen-go` libraries installed. Generating the golang grpc code can be done from `make`.

## Pre-commit

You can run pre-commit to verify the code with default golang linters. It is assumed you have pre-commit installed.

`pre-commit run -a`

# Project folders

## gRPC

As we want to see how fast grpc is for serving an image, there is a `/grpcserver` package created. Within this package there is already a dummy call to a `GetImage` remote procedure call made from a unit test. This is mimicking a server in the unit test meaning that the actual image is transferred from the `GetImage` func call all right, but it is not really using the network yet. The idea of this `server_test.go` is to verify the full image that is read from disk is indeed transferred with all the stream response chunks.

As a grpc-stream has a maximum capacity, based on the grpc and http/2 protocols, an image is broken up into 1Mb chunks to send over from server to client.

The `grpcclient` folder is meant to really setup the grpc communication over the network. To compare the "dummy" call to the actual grpc call, you can execute the Benchmark tests that are available in the `/grpcserver` and `/grpcclient` packages.

## HTTP/1.1

To compare gRPC with HTTP/1.1 or HTTP/2, we also have folders for a `httpserver`, `httpclient` and some scaffolding in place to do a request on a dummy http `/getimage` endpoint... to be filled in with actual image transport.
