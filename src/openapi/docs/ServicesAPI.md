# \ServicesAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ServicesGet**](ServicesAPI.md#ServicesGet) | **Get** /services | Retrieve all services and their status
[**ServicesNameGet**](ServicesAPI.md#ServicesNameGet) | **Get** /services/{name} | Retrieve the status and versions of a service
[**ServicesNameVersionDelete**](ServicesAPI.md#ServicesNameVersionDelete) | **Delete** /services/{name}/{version} | Delete a specific version of a service
[**ServicesNameVersionGet**](ServicesAPI.md#ServicesNameVersionGet) | **Get** /services/{name}/{version} | Retrieve the status of a specific version of a service
[**ServicesNameVersionPost**](ServicesAPI.md#ServicesNameVersionPost) | **Post** /services/{name}/{version} | Enable, disable or build a specific version of a service in the pipeline
[**ServicesPost**](ServicesAPI.md#ServicesPost) | **Post** /services | Upload a new service or new version to the rover by uploading a ZIP file



## ServicesGet

> []ServicesGet200ResponseInner ServicesGet(ctx).Execute()

Retrieve all services and their status

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
	resp, r, err := apiClient.ServicesAPI.ServicesGet(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ServicesAPI.ServicesGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ServicesGet`: []ServicesGet200ResponseInner
	fmt.Fprintf(os.Stdout, "Response from `ServicesAPI.ServicesGet`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiServicesGetRequest struct via the builder pattern


### Return type

[**[]ServicesGet200ResponseInner**](ServicesGet200ResponseInner.md)

### Authorization

[BasicAuth](../README.md#BasicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ServicesNameGet

> ServicesNameGet200Response ServicesNameGet(ctx, name).Execute()

Retrieve the status and versions of a service

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
	name := "imaging" // string | The name of the service

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ServicesAPI.ServicesNameGet(context.Background(), name).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ServicesAPI.ServicesNameGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ServicesNameGet`: ServicesNameGet200Response
	fmt.Fprintf(os.Stdout, "Response from `ServicesAPI.ServicesNameGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**name** | **string** | The name of the service | 

### Other Parameters

Other parameters are passed through a pointer to a apiServicesNameGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**ServicesNameGet200Response**](ServicesNameGet200Response.md)

### Authorization

[BasicAuth](../README.md#BasicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ServicesNameVersionDelete

> ServicesNameVersionDelete(ctx, name, version).Execute()

Delete a specific version of a service

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
	name := "imaging" // string | The name of the service
	version := "1.0.0" // string | The version of the service

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.ServicesAPI.ServicesNameVersionDelete(context.Background(), name, version).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ServicesAPI.ServicesNameVersionDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**name** | **string** | The name of the service | 
**version** | **string** | The version of the service | 

### Other Parameters

Other parameters are passed through a pointer to a apiServicesNameVersionDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



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


## ServicesNameVersionGet

> ServicesNameVersionGet200Response ServicesNameVersionGet(ctx, name, version).Execute()

Retrieve the status of a specific version of a service

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
	name := "imaging" // string | The name of the service
	version := "1.0.0" // string | The version of the service

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ServicesAPI.ServicesNameVersionGet(context.Background(), name, version).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ServicesAPI.ServicesNameVersionGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ServicesNameVersionGet`: ServicesNameVersionGet200Response
	fmt.Fprintf(os.Stdout, "Response from `ServicesAPI.ServicesNameVersionGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**name** | **string** | The name of the service | 
**version** | **string** | The version of the service | 

### Other Parameters

Other parameters are passed through a pointer to a apiServicesNameVersionGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**ServicesNameVersionGet200Response**](ServicesNameVersionGet200Response.md)

### Authorization

[BasicAuth](../README.md#BasicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ServicesNameVersionPost

> ServicesNameVersionPost(ctx, name, version).Action(action).Execute()

Enable, disable or build a specific version of a service in the pipeline

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
	name := "imaging" // string | The name of the service
	version := "1.0.0" // string | The version of the service
	action := "enable" // string | The action to perform on the service version

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.ServicesAPI.ServicesNameVersionPost(context.Background(), name, version).Action(action).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ServicesAPI.ServicesNameVersionPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**name** | **string** | The name of the service | 
**version** | **string** | The version of the service | 

### Other Parameters

Other parameters are passed through a pointer to a apiServicesNameVersionPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **action** | **string** | The action to perform on the service version | 

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


## ServicesPost

> ServicesPost200Response ServicesPost(ctx).Content(content).Execute()

Upload a new service or new version to the rover by uploading a ZIP file

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
	content := os.NewFile(1234, "some_file") // *os.File | The content of the ZIP file to upload (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ServicesAPI.ServicesPost(context.Background()).Content(content).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ServicesAPI.ServicesPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ServicesPost`: ServicesPost200Response
	fmt.Fprintf(os.Stdout, "Response from `ServicesAPI.ServicesPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiServicesPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **content** | ***os.File** | The content of the ZIP file to upload | 

### Return type

[**ServicesPost200Response**](ServicesPost200Response.md)

### Authorization

[BasicAuth](../README.md#BasicAuth)

### HTTP request headers

- **Content-Type**: multipart/form-data
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

