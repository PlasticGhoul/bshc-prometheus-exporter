# bshc-prometheus-exporter
*Prometheus Exporter for Bosch Smart Home Controller*

This Prometheus Exporter collects data from the Bosch Smart Home Controller API and publishes them for Prometheus to collect.

## Installation
### Docker
In this directory run the follwing command to build a ready-to-start Docker image:
```
docker build . -t bshc-prometheus-exporter:latest
```

After that you can run the following command to start the container:
```
docker run -it -d --name bshc-prometheus-exporter -p <host-port>:<container-port> -v <path to config>:/app/config/config.yaml bshc-prometheus-exporter:latest
```

### Build Binary
To build a binary just run:
```
go build -v -o ./bin/ ./...
```
After that you can run the follwing command to start the exporter:
```
./bshc-prometheus-exporter -c <path to config file>
```

### Run from directory
To run the exporter without building it, run the following commands:
```
go get
go run -c <path to config file>
```

## Parameters
The following parameters are supported:  
| Parameter | Purpose |
|-----------|---------|
| -c/--config | Path to config file |
| -b/--bind | HTTP bind IP address or hostname |
| -p/--port | HTTP port |
| -bh/--bshchost | BSHC hostname or IP address |
| -bp/--bshcport | BSHC API port |
| -cc/--clientcert | Client certificate for authentication |
| -ck/--clientkey | Client key for authentication |
| -d/--debug | Enable debug log output |

***Hint***  
Every parameter can be also set as a config value inside the config file except `-c/--congfig`.  
Please note that configuration ist read with the follwing priority (the further up in the list, the higher the priority):  
1. CLI parameters
2. Configuration file
3. (Potential) Default values

## Configuration file
The configuration file is written in YAML and contains the following sections/keys:
- http
  - bind --> Hostname or IP address to bind the HTTP server to
  - port --> Port to bind the HTTP server to
- bshc
  - host --> Hostname or IP address of the BSHC
  - port --> API port of the BSHC
  - client_cert --> Client certificate for authentication
  - client_key --> Client key for authentication
- services
  - temperature_level --> Enable temperature_level
  - humidity_level --> Enable humidity_level of devices
  - valve_tappet --> Enable valve_tappet (valve positiona) for thermostats

***An example/template configuration can be found in the `config` folder of this repository***