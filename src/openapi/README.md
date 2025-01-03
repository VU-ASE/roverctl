# Go API client for openapi

API exposed from each rover to allow process, service, source and file management

## Overview
This API client was generated by the [OpenAPI Generator](https://openapi-generator.tech) project.  By using the [OpenAPI-spec](https://www.openapis.org/) from a remote server, you can easily generate an API client.

- API version: 1.0.0
- Package version: 1.0.0
- Generator version: 7.9.0
- Build package: org.openapitools.codegen.languages.GoClientCodegen

## Installation

Install the following dependencies:

```sh
go get github.com/stretchr/testify/assert
go get golang.org/x/net/context
```

Put the package under your project folder and add the following in import:

```go
import openapi "github.com/GIT_USER_ID/GIT_REPO_ID"
```

To use a proxy, set the environment variable `HTTP_PROXY`:

```go
os.Setenv("HTTP_PROXY", "http://proxy_name:proxy_port")
```

## Configuration of Server URL

Default configuration comes with `Servers` field that contains server objects as defined in the OpenAPI specification.

### Select Server Configuration

For using other server than the one defined on index 0 set context value `openapi.ContextServerIndex` of type `int`.

```go
ctx := context.WithValue(context.Background(), openapi.ContextServerIndex, 1)
```

### Templated Server URL

Templated server URL is formatted using default variables from configuration or from context value `openapi.ContextServerVariables` of type `map[string]string`.

```go
ctx := context.WithValue(context.Background(), openapi.ContextServerVariables, map[string]string{
	"basePath": "v2",
})
```

Note, enum values are always validated and all unused variables are silently ignored.

### URLs Configuration per Operation

Each operation can use different server URL defined using `OperationServers` map in the `Configuration`.
An operation is uniquely identified by `"{classname}Service.{nickname}"` string.
Similar rules for overriding default operation server index and variables applies by using `openapi.ContextOperationServerIndices` and `openapi.ContextOperationServerVariables` context maps.

```go
ctx := context.WithValue(context.Background(), openapi.ContextOperationServerIndices, map[string]int{
	"{classname}Service.{nickname}": 2,
})
ctx = context.WithValue(context.Background(), openapi.ContextOperationServerVariables, map[string]map[string]string{
	"{classname}Service.{nickname}": {
		"port": "8443",
	},
})
```

## Documentation for API Endpoints

All URIs are relative to *http://localhost*

Class | Method | HTTP request | Description
------------ | ------------- | ------------- | -------------
*HealthAPI* | [**StatusGet**](docs/HealthAPI.md#statusget) | **Get** /status | Retrieve the health and versioning information
*HealthAPI* | [**UpdatePost**](docs/HealthAPI.md#updatepost) | **Post** /update | Self-update the roverd daemon process
*PipelineAPI* | [**LogsAuthorNameVersionGet**](docs/PipelineAPI.md#logsauthornameversionget) | **Get** /logs/{author}/{name}/{version} | Retrieve logs for any service. Logs from running or previously run services can be viewed and will be kept until rover reboot.
*PipelineAPI* | [**PipelineGet**](docs/PipelineAPI.md#pipelineget) | **Get** /pipeline | Retrieve pipeline status and process execution information
*PipelineAPI* | [**PipelinePost**](docs/PipelineAPI.md#pipelinepost) | **Post** /pipeline | Set the services that are enabled in this pipeline, by specifying the fully qualified services
*PipelineAPI* | [**PipelineStartPost**](docs/PipelineAPI.md#pipelinestartpost) | **Post** /pipeline/start | Start the pipeline
*PipelineAPI* | [**PipelineStopPost**](docs/PipelineAPI.md#pipelinestoppost) | **Post** /pipeline/stop | Stop the pipeline
*ServicesAPI* | [**FetchPost**](docs/ServicesAPI.md#fetchpost) | **Post** /fetch | Fetches the zip file from the given URL and installs the service onto the filesystem
*ServicesAPI* | [**ServicesAuthorGet**](docs/ServicesAPI.md#servicesauthorget) | **Get** /services/{author} | Retrieve the list of parsable services for a specific author
*ServicesAPI* | [**ServicesAuthorServiceGet**](docs/ServicesAPI.md#servicesauthorserviceget) | **Get** /services/{author}/{service} | Retrieve the list of parsable service versions for a specific author and service
*ServicesAPI* | [**ServicesAuthorServiceVersionDelete**](docs/ServicesAPI.md#servicesauthorserviceversiondelete) | **Delete** /services/{author}/{service}/{version} | Delete a specific version of a service
*ServicesAPI* | [**ServicesAuthorServiceVersionGet**](docs/ServicesAPI.md#servicesauthorserviceversionget) | **Get** /services/{author}/{service}/{version} | Retrieve the status of a specific version of a service
*ServicesAPI* | [**ServicesAuthorServiceVersionPost**](docs/ServicesAPI.md#servicesauthorserviceversionpost) | **Post** /services/{author}/{service}/{version} | Build a fully qualified service version
*ServicesAPI* | [**ServicesGet**](docs/ServicesAPI.md#servicesget) | **Get** /services | Retrieve the list of all authors that have parsable services. With these authors you can query further for services
*ServicesAPI* | [**UploadPost**](docs/ServicesAPI.md#uploadpost) | **Post** /upload | Upload a new service or new version to the rover by uploading a ZIP file


## Documentation For Models

 - [DaemonStatus](docs/DaemonStatus.md)
 - [FetchPost200Response](docs/FetchPost200Response.md)
 - [FetchPostRequest](docs/FetchPostRequest.md)
 - [GenericError](docs/GenericError.md)
 - [PipelineGet200Response](docs/PipelineGet200Response.md)
 - [PipelineGet200ResponseEnabledInner](docs/PipelineGet200ResponseEnabledInner.md)
 - [PipelineGet200ResponseEnabledInnerProcess](docs/PipelineGet200ResponseEnabledInnerProcess.md)
 - [PipelineGet200ResponseEnabledInnerService](docs/PipelineGet200ResponseEnabledInnerService.md)
 - [PipelinePost400Response](docs/PipelinePost400Response.md)
 - [PipelinePost400ResponseValidationErrors](docs/PipelinePost400ResponseValidationErrors.md)
 - [PipelinePostRequestInner](docs/PipelinePostRequestInner.md)
 - [PipelineStatus](docs/PipelineStatus.md)
 - [ProcessStatus](docs/ProcessStatus.md)
 - [ReferencedService](docs/ReferencedService.md)
 - [ServiceStatus](docs/ServiceStatus.md)
 - [ServicesAuthorServiceVersionDelete200Response](docs/ServicesAuthorServiceVersionDelete200Response.md)
 - [ServicesAuthorServiceVersionGet200Response](docs/ServicesAuthorServiceVersionGet200Response.md)
 - [ServicesAuthorServiceVersionGet200ResponseConfigurationInner](docs/ServicesAuthorServiceVersionGet200ResponseConfigurationInner.md)
 - [ServicesAuthorServiceVersionGet200ResponseConfigurationInnerValue](docs/ServicesAuthorServiceVersionGet200ResponseConfigurationInnerValue.md)
 - [ServicesAuthorServiceVersionGet200ResponseInputsInner](docs/ServicesAuthorServiceVersionGet200ResponseInputsInner.md)
 - [ServicesAuthorServiceVersionPost400Response](docs/ServicesAuthorServiceVersionPost400Response.md)
 - [StatusGet200Response](docs/StatusGet200Response.md)
 - [StatusGet200ResponseCpuInner](docs/StatusGet200ResponseCpuInner.md)
 - [StatusGet200ResponseMemory](docs/StatusGet200ResponseMemory.md)
 - [UnmetServiceError](docs/UnmetServiceError.md)
 - [UnmetStreamError](docs/UnmetStreamError.md)
 - [UpdatePost200Response](docs/UpdatePost200Response.md)


## Documentation For Authorization


Authentication schemes defined for the API:
### BasicAuth

- **Type**: HTTP basic authentication

Example

```go
auth := context.WithValue(context.Background(), openapi.ContextBasicAuth, openapi.BasicAuth{
	UserName: "username",
	Password: "password",
})
r, err := client.Service.Operation(auth, args)
```


## Documentation for Utility Methods

Due to the fact that model structure members are all pointers, this package contains
a number of utility functions to easily obtain pointers to values of basic types.
Each of these functions takes a value of the given basic type and returns a pointer to it:

* `PtrBool`
* `PtrInt`
* `PtrInt32`
* `PtrInt64`
* `PtrFloat`
* `PtrFloat32`
* `PtrFloat64`
* `PtrString`
* `PtrTime`

## Author



