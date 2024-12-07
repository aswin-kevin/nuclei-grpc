# Nuclei gRPC Application

This application serves nuclei engine as GRPC service. It uses nuclei engine V3.

## Prerequisites

- Go 1.16 or higher
- gRPC Go plugin

## Installation

1. Clone the repository:

   ```sh
   git clone https://github.com/aswin-kevin/nuclei-grpc.git
   cd nuclei-grpc
   ```

2. Install dependencies:
   ```sh
   go mod tidy
   ```

## Usage

### Starting the Server

To start the gRPC server, run the following command:

```sh
go run main.go start
```

### Starting the Server with Custom Address and Port

To start the gRPC server with a custom address and port, use the following command:

```sh
go run main.go start --address <custom_address> --port <custom_port>
```

Replace `<custom_address>` with the desired address (e.g., `localhost`) and `<custom_port>` with the desired port (e.g., `50051`).

### Example Request

To make a request to the gRPC server, use the following payload:

1. Filter using tags

```json
{
  "targets": ["https://hotstar.com"],
  "tags": ["dns", "ssl", "tech"]
}
```

2. Filter using templates

```json
{
  "targets": ["https://hotstar.com"],
  "templates": ["http/cves", "network/cves"]
}
```

3. Filter using templates relative paths

```json
{
  "targets": ["https://hotstar.com"],
  "templates": ["http/cves/xyz.yaml", "network/cves/new.yaml"]
}
```

4. Filter using templates ids

```json
{
  "targets": ["https://hotstar.com"],
  "template_ids": ["detect-dangling-cname", "dnssec-detection"]
}
```

### Example Response

The response is a server-side streaming response. Example JSON responses can be found in the `examples/jsons` directory.

### Protofile

Protofile is placed in `pkg/service/service.proto` , you this on your client side.

If you want to rebuild the protofile use the following command

```sh
protoc --go_out=. --go-grpc_out=. service.proto
```

## Contributing

Feel free to submit issues, fork the repository and send pull requests!

## License

This project is licensed under the MIT License.
