# Notes API

Rest API that allows users to:

- Add, update, read and delete single notes
- Add(yet to implement), read and delete multiple notes
- Filter notes by dates, categories, and tags.

A note contains a title(**required, unique**), description(**required**), category(**optional**), date(**populated by the API**), and tags(**optional**).

## Overview

The API is written in [Golang](https://go.dev) and it is using the [Gin Web Framework](https://github.com/gin-gonic/gin). It's using [MongoDB](https://www.mongodb.com) as a database.

The tests are written with the [Ginkgo](https://github.com/onsi/ginkgo) and [Gomega](https://github.com/onsi/gomega) testing frameworks. For mocks the [Gomock](https://github.com/golang/mock) mocking framework is used.

## Contents

- [Prerequisites](#prerequisites)
- [Build](#build)
- [Configure](#configure)
- [Deploy](#deploy)
- [Run](#run)

## Prerequisites

- [Golang](https://go.dev)
- [MongoDB](https://www.mongodb.com) instance
- [Docker](https://www.docker.com)
- [Kubernetes](https://kubernetes.io)

## Build

You can build the application binary locally using `go build -o <binary-name>` from within the root directory of the project.

Building a docker image with the Dockerfile can be done with `docker build -f ./Dockerfile -t <image-name>` from within the root directory of the project.

## Configure

__Port `3040` is hardcoded as the port for the health server for the kubernetes probes.__

There are 5 environment variables you need to set to configure the application:

- **[required]** DATABASE_URI - the connection URI for the database
    - You can enable TSL communication to the database if you pass the parameters to the connection URI. For example `mongodb+srv://CLUSTER_LOCATION/?authSource=%24external&authMechanism=MONGODB-X509&retryWrites=true&w=majority&tlsCertificateKeyFile=./certs/tls.pem`, where **tlsCertificateKeyFile** points to the file where the certificate and its private key are
- **[required]** DATABASE_NAME - the name of the database
- **[required]** DATABASE_COLLECTION - the name of the database collection that will be used
- **[required]** SERVER_PORT - the server port for the HTTP traffic
- **[optional]** SERVER_TLS_PORT - the server port for the HTTPS traffic
    - Additionally for enabling TLS, you need to provide 2 flags:
        - tlsCertLocation - the location of the TLS certificate and if the certificate is signed by a certificate authority, the file should be the concatenation of the server's certificate, any intermediates, and the CA's certificate
        - tlsKeyLocation - the location of the private key of the certificate

### On Kubernetes

When deploying the application to kubernetes, the DB configuration is provided via a secret `./deployments/0db-config.yaml` and the configuration for enabling TLS is provided via a secret as well `./deployments/1tls.yaml`, where the secret keys
- tls.crt - contains the API client certificate and the CA's certificate if signed by any
- tls.key - contains the API client certificate's private key
- tls.pem - contains both the API client certificate and the certificate's private key

The additional environment variables and application arguments are provided inside the deployment file `./deployments/2deployment.yaml`