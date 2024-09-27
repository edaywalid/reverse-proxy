# REVERSE PROXY

## Description

a proxy server that routes http requests to differnt servers based on the round robin algorithm.

## Usage

### Prerequisites

- Go 1.22.5
- Docker (optional)
- make

### Installation

- Clone the repository

```bash
    git clone https://github.com/edaywalid/reverse-proxy
    cd reverse-proxy
```

- Create a `config.yaml` file in the root directory of the project and add the following content

```yaml
servers:
  - url: http://localhost:8001
  - url: http://localhost:8002
  - url: http://localhost:8003
https_port: 8443
http_port: 8080
cert_file: ./ssl/localhost.crt
key_file: ./ssl/localhost.key
```

- Note that the `cert_file` and `key_file` are needed for the https server to run. You can generate a self-signed certificate using the following command

- u need to have openssl installed on your machine

```bash
    openssl req -x509 -newkey rsa:4096 -keyout ssl/localhost.key -out ssl/localhost.crt -days 365 -nodes
```

- You will be prompted to enter the following information

      Country Name (2 letter code) [AU]:US
      State or Province Name (full name) [Some-State]:California
      Locality Name (eg, city) []:San Fransisco
      Organization Name (eg, company) [Internet Widgits Pty Ltd]:My Company
      Organizational Unit Name (eg, section) []:
      Common Name (e.g. server FQDN or YOUR name) []:localhost
      Email Address []:

- For testing Run the docker containers

```bash
    make runc
```

- For running the server

```bash
    make run-server
```

- To run both the server and the containers

```bash
    make run-all
```
