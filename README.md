<p align="center">
    <img alt="Beetle Logo" src="https://raw.githubusercontent.com/clivern/Beetle/master/assets/img/gopher.png?v=1.0.1" width="150" />
    <h3 align="center">Beetle</h3>
    <p align="center">Kubernetes multi-cluster deployment service</p>
    <p align="center">
        <a href="https://travis-ci.org/Clivern/Beetle"><img src="https://travis-ci.org/Clivern/Beetle.svg?branch=master"></a>
        <a href="https://github.com/Clivern/Beetle/releases"><img src="https://img.shields.io/badge/Version-0.0.1-red.svg"></a>
        <a href="https://goreportcard.com/report/github.com/Clivern/Beetle"><img src="https://goreportcard.com/badge/github.com/clivern/Beetle?v=0.0.1"></a>
     <a href="https://hub.docker.com/r/clivern/beetle"><img src="https://img.shields.io/badge/Docker-Latest-orange"></a>
        <a href="https://github.com/Clivern/Beetle/blob/master/LICENSE"><img src="https://img.shields.io/badge/LICENSE-MIT-orange.svg"></a>
    </p>
</p>


Application deployment and management should be automated, auditable, and easy to understand and that's what beetle tries to achieve in a simple manner. Beetle automates the deployment and rollback of your applications in a multi-cluster, multi-namespaces kubernetes environments. Easy to integrate with through API endpoints & webhooks to fit a variety of workflows.


## Documentation

## Deployment

### On a Linux Server

Download [the latest Beetle binary.](https://github.com/Clivern/Beetle/releases)

```zsh
$ curl -sL https://github.com/Clivern/Beetle/releases/download/x.x.x/beetle_x.x.x_OS.tar.gz | tar xz
```

Create your config file as explained on [development part](#development) and run beetle with systemd or anything else you prefer.

```
$ ./beetle -config=/custom/path/config.prod.yml
```


## Development

Beetle uses Go Modules to manage dependencies. First Create a prod config file.

```zsh
$ git clone https://github.com/Clivern/Beetle.git
$ cp config.dist.yml config.prod.yml
```

Then add your default configs. You probably wondering how the following configs even work! let's pick one and explain.

The item mode: `${BEETLE_APP_MODE:-dev}` means that the mode is dev unless environment variable `BEETLE_APP_MODE` is defined. so you can always override the value by defining the environment variable `export BEETLE_APP_MODE=prod`. and same for others

```yaml
# App configs
app:
    # Env mode (dev or prod)
    mode: ${BEETLE_APP_MODE:-dev}
    # HTTP port
    port: ${BEETLE_API_PORT:-8080}
    # App URL
    domain: ${BEETLE_APP_DOMAIN:-http://127.0.0.1:8080}
    # TLS configs
    tls:
        status: ${BEETLE_API_TLS_STATUS:-off}
        pemPath: ${BEETLE_API_TLS_PEMPATH:-cert/server.pem}
        keyPath: ${BEETLE_API_TLS_KEYPATH:-cert/server.key}

    # Message Broker Configs
    broker:
        # Broker driver (native)
        driver: ${BEETLE_BROKER_DRIVER:-native}
        # Native driver configs
        native:
            # Queue max capacity
            capacity: ${BEETLE_BROKER_NATIVE_CAPACITY:-5000}
            # Number of concurrent workers
            workers: ${BEETLE_BROKER_NATIVE_WORKERS:-4}

    # Application Database
    database:
        # Database driver (mysql)
        driver: ${BEETLE_DATABASE_DRIVER:-mysql}
        # MySQL DB Configs
        mysql:
            # MySQL Hostname
            host: ${BEETLE_DATABASE_MYSQL_HOST:-localhost}
            # MySQL Port
            port: ${BEETLE_DATABASE_MYSQL_PORT:-3306}
            # MySQL Database
            database: ${BEETLE_DATABASE_MYSQL_DATABASE:-beetle}
            # MySQL Username
            username: ${BEETLE_DATABASE_MYSQL_USERNAME:-root}
            # MySQL Password
            password: ${BEETLE_DATABASE_MYSQL_PASSWORD:- }

    # Supported Notifications Webhooks
    webhooks:
        -
            name: http_service
            type: http
            url: http://example.com/api/listen
            method: post
            headers:
                - X-AUTH-TOKEN=1234
            events:
                - deployment
                - rollback

    # Kubernetes Clusters
    clusters:
        -
            # kubernetes cluster name
            name: ${BEETLE_DEFAULT_CLUSTER_NAME:-default}

            # kubernetes cluster kubctl config
            kubeconfig: ${BEETLE_DEFAULT_CLUSTER_CONFIG:-/app/config/default-kubctl.yaml}

            # Enabled Notifications Webhooks
            notify:
                - http_service

# Log configs
log:
    # Log level, it can be debug, info, warn, error, panic, fatal
    level: ${BEETLE_LOG_LEVEL:-info}
    # output can be stdout or abs path to log file /var/logs/beetle.log
    output: ${BEETLE_LOG_OUTPUT:-stdout}
    # Format can be json
    format: ${BEETLE_LOG_FORMAT:-json}
```

And then run the application.

```zsh
$ go build beetle.go
$ ./beetle

// OR

$ make run

// To Provide a custom config file
$ ./beetle -config=/custom/path/config.prod.yml
$ go run beetle.go -config=/custom/path/config.prod.yml
```

## Versioning

For transparency into our release cycle and in striving to maintain backward compatibility, Beetle is maintained under the [Semantic Versioning guidelines](https://semver.org/) and release process is predictable and business-friendly.

See the [Releases section of our GitHub project](https://github.com/clivern/beetle/releases) for changelogs for each release version of Beetle. It contains summaries of the most noteworthy changes made in each release.


## Bug tracker

If you have any suggestions, bug reports, or annoyances please report them to our issue tracker at https://github.com/clivern/beetle/issues


## Security Issues

If you discover a security vulnerability within Beetle, please send an email to [hello@clivern.com](mailto:hello@clivern.com)


## Contributing

We are an open source, community-driven project so please feel free to join us. see the [contributing guidelines](CONTRIBUTING.md) for more details.


## License

© 2020, clivern. Released under [MIT License](https://opensource.org/licenses/mit-license.php).

**Beetle** is authored and maintained by [@clivern](http://github.com/clivern).
