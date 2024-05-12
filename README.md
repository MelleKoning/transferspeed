# Transferspeed

Transferspeed is a test repository to see if we can benchmark transferring an arbitrary binary data over grpc http/2.

## Working with the code

The software is written in golang, and uses default libraries for generating the grpc communication. It is assumed you have the appropriate `protoc` and `protoc-gen-go` libraries installed. Generating the golang grpc code can be done from `make`.

## Pre-commit

You can run pre-commit to verify the code with default golang linters. It is assumed you have pre-commit installed.

`pre-commit run -a`
