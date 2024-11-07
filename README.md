# Kafka Proxy Producer

![Build Status][build-shield]
![Go version][go-shield]
[![Apache License][license-shield]][license-url]

Kafka Proxy Producer is a web server that facilitates seamless message production to Kafka broker based on specified routes defined in a `configuration.json` file. This project bridges your applications and Kafka by mapping incoming HTTP requests to Kafka messages, enhancing integration capabilities.

## Features

- **Dynamic Route Configuration**: Define routes in a json file that specify how incoming HTTP requests should be transformed into Kafka messages.
- **Supports Multiple HTTP Methods**: Configure routes for various HTTP methods (e.g., POST) to handle different types of requests.
- **Flexible Message Mapping**: Customize the Kafka messages produced based on the incoming request data.

## Configuration

Routes are defined in a `configuration.json` file with the following structure:

```json
{
  "Routes": [
    {
      "DownstreamTopicPartition": {},
      "DownstreamMessage": {},
      "UpstreamPathTemplate": "/",
      "UpstreamHttpMethod": []
    }
  ],
  "LaunchSettings": {}
}
```

- **DownstreamTopicPartition**: Specifies the Kafka topic and partition to which the message will be sent.
  - `Topic`: The name of the Kafka topic.
  - `Partition`: The partition number for the topic.
  - `Offset`: The offset for the message (use null to let Kafka handle it).
  - `Metadata`: Additional metadata (optional).
- **DownstreamMessage**: Defines the structure of the Kafka message that will be produced. Use placeholders like `[FromRoute]`, `[FromBody]`, and `[FromHeader]` to dynamically map data from the HTTP request.

> [!NOTE]
> When you define nested json structures in the `DownstreamMessage` section of the configuration, this structure will be mirrored in both the incoming HTTP request and the resulting Kafka message.

- **UpstreamPathTemplate**: The URL path that the web server will listen to for incoming requests. You can use path parameters for dynamic routing like this `.../{name}/...`.
- **UpstreamHttpMethod**: An array of allowed HTTP methods for this route (e.g., `["Post", "Put"]`).

See [configuration documentation](./docs/configuration.md) for more details.

## Getting Started

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/mishakrpv/kafka-proxy-producer.git
   cd kafka-proxy-producer
   ```

2. Install dependencies:

   ```bash
   go mod download && go mod verify
   ```

3. Configure your `configuration.json` file as per your requirements.

### Running the Server

You can run the server either directly or as a Docker container.

#### Running Directly

Start the server by executing:

```bash
make run
```

#### Running in Docker

To run the server as a Docker container, use the following command:

```bash
docker build -t kafka-proxy-producer .
```

Then, run the container while specifying the configuration file path:

```bash
docker run -e CONFIG_PATH="/path/to/your/configuration.json" -p 5465:5465 kafka-proxy-producer
```

### Example Usage

1. Execute following command in the [example](./example) directory:

   ```bash
   docker-compose up
   ```

2. Send a POST request to the defined endpoint:

   ```bash
   POST http://localhost:5465/items/123
   Content-Type: application/json
   Authorization: Bearer your_token

   {
      "principal": {
        "name": "John Doe"
      }
   }
   ```

## License

This project is licensed under the Apache License 2.0. See the [LICENSE](LICENSE) file for details.

<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->

[license-shield]: https://img.shields.io/badge/license-Apache%202.0-red?style=flat-square
[license-url]: https://github.com/mishakrpv/kafka-proxy-producer/blob/main/LICENSE
[go-shield]: https://img.shields.io/github/go-mod/go-version/mishakrpv/kafka-proxy-producer
[build-shield]: https://github.com/mishakrpv/kafka-proxy-producer/actions/workflows/go.yml/badge.svg
