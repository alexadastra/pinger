# OZON Internship Test Task

Microservice for HTTP requests status codes tracking.

## Installation
First up, clone the repository.
For local usage run following commands:

```bash
make postgres
make run
```
This will enable PostgreSQL on localhost:5432 and gRPC endpoints on localhost:50001. For deploying with Docker, run `Dockerfile` and then `docker-compose up`.

## Usage

Here are inputs and outputs for each use-case:
### Adding new URL to check
Input: url to check (obligatory), HTTP request method (optional, "get" by default), time interval to check in seconds (optional, 24 hours by default).
```
message UrlPostRequest{
  string url = 1;
  string method = 2;
  string time_interval = 3;
}
```
Output: full struct for checking url, including id and timestamp of creation
```
message Url {
  string id = 1;
  string url = 2;
  string method = 3;
  string time_created = 4;
  string time_interval = 5;
}
```
### Getting checks for URL
Input: URL to select for (obligatory) and limit for query (optional, 5 by default)

```
message CheckGetRequest{
  string url = 1;
  int32 limit = 2;
}
```
Output: List of checks with not empty status code and creating time fields.
```
message Check{
  string id = 1;
  string url = 2;
  int32 status_code = 3;
  string time_checked = 4;
}
message CheckGetResponse{
  repeated Check checks = 1;
}
```
### Deleting URL to check
Both input and output are the same CheckGetRequest with only URL accepting as argument.


## Development
For cleaning .proto files, run
```bash
make clean
```
For building (with pre-installed brotobuf and protoc), run
```bash
make gen
```