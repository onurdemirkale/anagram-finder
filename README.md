# Anagram Finder

Anagram Finder is a service designed to provide anagrams for a given set of words. The API supports multiple input methods and sort-map algorithm to find anagrams.

Sort-map algorithm was chosen due to its simplicity and effectiveness. Compared to hash and trie algorithms, sort-map offered the same time complexity. Trie algorithm had the advantage of prefix searches, which would be beneficial if the functionality was to be expanded.

## Overview

```
├─ /cmd
│ └─ /anagramfinder
│ └─ main.go - Entry point of the application.
│
├─ /pkg
│ ├─ /anagram
│ │ ├─ anagram_finder.go - Defines an interface for anagram finders.
│ │ ├─ anagram_finder_factory.go - Factory to create an instance of anagram finder.
│ │ └─ sort_map_anagram_finder.go - Implementation of anagram finder using sorted map.
│ │
│ └─ /inputsource
│ ├─ input_source.go - Defines an interface for input sources.
│ ├─ input_source_factory.go - Factory to create an instance of input source.
│ ├─ http_body_input_source.go - Implementation to handle inputs from HTTP body.
│ ├─ http_file_input_source.go - Implementation to handle inputs from HTTP files.
│ └─ http_url_input_source.go - Implementation to handle inputs from HTTP URLs.
│
└─ /api
├─ anagram_handler.go - Handles and parses requests and initiates anagram finding process.
├─ anagram_handler_test.go - Integration tests for anagram requests.
├─ anagram_request.go - Defines and validates the incoming anagram request.
├─ anagram_response.go - Defines and serves the response of the anagram request.
├─ error_handler.go - Maps and handles errors for HTTP responses.
│
├─ /k8s - Kubernetes deployments and services.


```

## Prerequisites

- Go 1.20
- Docker  (optional for containerized deployment)
- Kubernetes (optional for containerized deployment)

## Supported Algorithms

- Sort-Map: Sorts the characters of a word and then uses this sorted version as a key in a map. All words that sort to the same string are anagrams of each other.

## Deployment

### Local Development

Build the project:

```sh
go build
```

Run the resulting binary:

```sh
./anagram-finder
```

### Docker

Build the Docker image: 

```sh
docker build -t anagram-finder:latest .
```

Run the Docker container:

```sh
docker run -p 8080:8080 anagram-finder:latest
```

### Kubernetes

Ensure you have kubectl configured for your cluster.

Deploy the application:

```sh
kubectl apply -f k8s/
```

In the current Kubernetes deployment, `imagePullPolicy` field is set to `IfNotPresent`, which ensures Kubernetes uses the local image if it exists. Update this field before a deployment.

To test the Kubernetes deployment, forward a port.

```sh
kubectl port-forward service/anagram-finder-service 8080:80
```

## API Usage

1. Using a file:

```sh
curl -X POST -H "Content-Type: multipart/form-data" \
  -F "file=@anagrams.txt" \
  -F "inputType=http_file" \
  -F "algorithm=basic" \
  http://localhost:8080/anagram -o downloaded_anagrams.txt
```

2. Using JSON:

```sh
curl -X POST -H 'Content-Type: application/json' \
http://localhost:8080/anagram \
-d '{
  "inputType": "http_body",
  "inputData": "cat,tac",
  "algorithm": "sort_map"
}'
```

3. Using URL:

```sh
curl -X POST -H 'Content-Type: application/json' \
http://localhost:8080/anagram \
-d '{
  "inputType": "http_url",
  "inputData": <url>,
  "algorithm": "sort_map"
}'
```

## Test Coverage

anagram-finder/api/anagram_handler.go:18:			NewAnagramHandler	100.0%
anagram-finder/api/anagram_handler.go:25:			FindAnagrams		75.0%
anagram-finder/api/anagram_handler.go:46:			parseRequest		77.3%
anagram-finder/api/anagram_handler.go:89:			processAnagrams		70.0%
anagram-finder/api/anagram_request.go:20:			validate		100.0%
anagram-finder/api/anagram_request.go:36:			validateInputType	100.0%
anagram-finder/api/anagram_request.go:49:			validateAlgorithm	100.0%
anagram-finder/api/anagram_request.go:61:			validateInputData	83.3%
anagram-finder/api/anagram_response.go:15:			serveResponse		100.0%
anagram-finder/api/error_handler.go:34:			handleError		75.0%
anagram-finder/pkg/anagram/factory.go:11:			NewAnagramFinderFactory	0.0%
anagram-finder/pkg/anagram/factory.go:15:			CreateAnagramFinder	0.0%
anagram-finder/pkg/anagram/sort_map_anagram_finder.go:11:	NewSortMapAnagramFinder	100.0%
anagram-finder/pkg/anagram/sort_map_anagram_finder.go:19:	FindAnagrams		100.0%
anagram-finder/pkg/anagram/sort_map_anagram_finder.go:43:	sortWord		100.0%
anagram-finder/pkg/inputsource/factory.go:14:			NewInputSourceFactory	0.0%
anagram-finder/pkg/inputsource/factory.go:18:			CreateInputSource	0.0%
anagram-finder/pkg/inputsource/http_body_input_source.go:11:	NewHttpBodyInputSource	100.0%
anagram-finder/pkg/inputsource/http_body_input_source.go:16:	GetWords		100.0%
anagram-finder/pkg/inputsource/http_file_input_source.go:12:	NewHttpFileInputSource	100.0%
anagram-finder/pkg/inputsource/http_file_input_source.go:16:	GetWords		87.5%
anagram-finder/pkg/inputsource/http_url_input_source.go:15:	NewHttpUrlInputSource	100.0%
anagram-finder/pkg/inputsource/http_url_input_source.go:19:	GetWords		84.6%
total:											(statements)		79.8%

## Swagger Documentation

You can access the API documentation at docs/swagger.yaml in the root directory. It provides detailed information about request formats, responses, and possible status codes.