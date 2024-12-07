# Nuclei gRPC Application

This application serves nuclei engine as GRPC service.

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

```json
{
  "targets": ["https://hotstar.com"],
  "tags": ["dns", "ssl", "tech"]
}
```

### Example Response

The response is a server-side streaming response. Example JSON responses can be found in the `examples/jsons` directory.

## Contributing

Feel free to submit issues, fork the repository and send pull requests!

## License

This project is licensed under the MIT License.
