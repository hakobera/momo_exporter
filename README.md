# WebRTC Native Client Momo Exporter for Prometheus

This is a simple server that scrapes WebRTC Native Client Momo stats and exports them via HTTP for Prometheus consumption.

## Prerequisities

This exporter is only works Momo 2020.11 and later.

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

```
Copyright 2020, Kazuyuki Honda

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
