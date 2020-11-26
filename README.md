# WebRTC Native Client Momo Exporter for Prometheus

This is a simple server that scrapes WebRTC Native Client Momo stats and exports them via HTTP for Prometheus consumption.

**CAUTION: Momo Exporter is only works with [customized version of Momo](https://github.com/hakobera/momo/pull/4)**

## Getting Started

To run it:

```sh
$ ./momo_exporter [flags]
```

Help on flags:

```sh
$ ./momo_exporter --help
```

## Usage

### HTTP stats URL

Specify custom URLs for the Momo stats port using the --momo.scrape-uri flag.

```sh
$ momo_exporter --momo.scrape-uri="http://localhost:8081/metrics"
```

## License

Apache License 2.0, see [LICENSE](https://github.com/hakobera/momo_exporter/blob/main/LICENSE)
