# OZON Internship Test Task

Microservice for HTTP requests status codes tracking.

## Installation
First up, clone the repository. For deploying with Docker, run :
```bash
docker-compose up
```
After that, microservice can be requested by 50001 port. You can use Postman or one of the clients inside project for testing.
```bash
 # for adding url to check
 go run client/add_url.go
 # retrieve last few checks for url 
 go run client/get_checks_by_url.go 
 # remove url from checking
 go run client/remove_url.go 
 # retrieve url that have n successful checks after date 
 go run client/get_urls_by_date_and_n.go 
```

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

### Getting URLs by date and n
Input: date of format "2020-Sep-19" and integer n.

```
message UrlGetRequest{
  string date = 1;
  int32 n = 2;
}
```
Output: Urls that have n or more successful (2xx) checks after date.
```
message UrlGetResponse{
  repeated string urls = 1;
}
```


## Development
For cleaning .proto files, run
```bash
make clean
```
For building (with pre-installed brotobuf and protoc), run
```bash
make gen
```