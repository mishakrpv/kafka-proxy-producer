# Configuration

## Suppported Downstream Message Sources

You can map data from different parts of the incoming HTTP request to form the Kafka message. The following placeholders are supported:

- **`[FromBody]`**: Data from the HTTP request body.
- **`[FromRoute]`**: Data from the URL path (e.g., dynamic segments like `/items/{id}`).
- **`[FromQuery]`**: Data from the query string parameters (e.g., `/items?id=123`).
- **`[FromForm]`**: Data from the form-encoded body of the HTTP request.
- **`[FromHeader]`**: Data from the HTTP request headers (e.g., `Authorization`, `Content-Type`).
