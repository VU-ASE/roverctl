# \PipelineAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**PipelineGet**](PipelineAPI.md#PipelineGet) | **Get** /pipeline | Retrieve pipeline status and process execution information
[**PipelineNameGet**](PipelineAPI.md#PipelineNameGet) | **Get** /pipeline/{name} | Retrieve the status of a service running as a process in the pipeline
[**PipelinePost**](PipelineAPI.md#PipelinePost) | **Post** /pipeline | Start or stop the pipeline of all enabled services



## PipelineGet

> PipelineGet200Response PipelineGet(ctx).Execute()

Retrieve pipeline status and process execution information

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.PipelineAPI.PipelineGet(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PipelineAPI.PipelineGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `PipelineGet`: PipelineGet200Response
	fmt.Fprintf(os.Stdout, "Response from `PipelineAPI.PipelineGet`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiPipelineGetRequest struct via the builder pattern


### Return type

[**PipelineGet200Response**](PipelineGet200Response.md)

### Authorization

[BasicAuth](../README.md#BasicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PipelineNameGet

> PipelineNameGet200Response PipelineNameGet(ctx, name).LogLines(logLines).Execute()

Retrieve the status of a service running as a process in the pipeline

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	name := "imaging" // string | The name of the service running as a process in the pipeline
	logLines := int32(100) // int32 | The number of log lines to retrieve (optional) (default to 50)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.PipelineAPI.PipelineNameGet(context.Background(), name).LogLines(logLines).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PipelineAPI.PipelineNameGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `PipelineNameGet`: PipelineNameGet200Response
	fmt.Fprintf(os.Stdout, "Response from `PipelineAPI.PipelineNameGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**name** | **string** | The name of the service running as a process in the pipeline | 

### Other Parameters

Other parameters are passed through a pointer to a apiPipelineNameGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **logLines** | **int32** | The number of log lines to retrieve | [default to 50]

### Return type

[**PipelineNameGet200Response**](PipelineNameGet200Response.md)

### Authorization

[BasicAuth](../README.md#BasicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PipelinePost

> PipelinePost(ctx).Action(action).Execute()

Start or stop the pipeline of all enabled services

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	action := "start" // string | The action to perform on the pipeline

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.PipelineAPI.PipelinePost(context.Background()).Action(action).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PipelineAPI.PipelinePost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiPipelinePostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **action** | **string** | The action to perform on the pipeline | 

### Return type

 (empty response body)

### Authorization

[BasicAuth](../README.md#BasicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

