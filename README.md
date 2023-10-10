# standard-deviation-service

## How to run

### service executable
```
go run cmd/main.go
```
or
```
go run cmd/main.go -reqs 15 -port 8081
```

### container image
```
docker build --tag stddev .
docker run -v /etc/ssl/certs:/etc/ssl/certs --net host stddev
```
or
```
docker build --tag stddev .
docker run -v /etc/ssl/certs:/etc/ssl/certs --net host stddev -reqs 15 -port 8081
```

### parameters
```
-reqs    number of requests per second
-port    port number
```

## Task

### Description
Create a simple REST service in Go supporting the following GET operation:
```
/random/mean?requests={r}&length={l}
```
which performs `{r}` concurrent requests to [random.org](https://random.org) API asking for `{l}` number of random integers.

For each of `{r}` requests, the service must calculate standard deviation of the drawn integers and additionally standard deviation of sum of all sets.
Results must be presented in JSON format.

### Example
`GET /random/mean?requests=2&length=5`

Response:
```json
[
  {
     "stddev": 1,
     "data": [1, 2, 3, 4, 5]
  },
  {
     "stddev": 1,
     "data": [1, 2, 3, 4, 5]
  },
  { // stddev of sum
     "stddev": 1,
     "data": [1, 1, 2, 2, 3, 3, 4, 4, 5, 5]
  }
]
```

### Requirements
1. Proper error handling when communicating with external service (timeouts, invalid statuses).
2. Usage of contexts.
3. Solution should be delivered as a git repository.
4. Provide a `.Dockerfile` that builds the image with the application.
5. Application should run flawlessly after running `docker build ...` & `docker run ...` commands.