package http

import (
	"github.com/alioth-center/infrastructure/trace"
	"github.com/alioth-center/infrastructure/utils/values"
	"github.com/gin-gonic/gin"
)

func DefaultPreprocessors[request any, response any](processors ...EndpointPreprocessor[request, response]) []EndpointPreprocessor[request, response] {
	return append(NewPreprocessors[request, response](
		CheckRequestMethodPreprocessor[request, response],
		CheckRequestHeadersPreprocessor[request, response],
		LoadNormalRequestHeadersPreprocessor[request, response],
		CheckRequestQueriesPreprocessor[request, response],
		CheckRequestParamsPreprocessor[request, response],
		CheckRequestCookiesPreprocessor[request, response],
		CheckRequestBodyPreprocessor[request, response],
	), processors...)
}

func NewPreprocessors[request any, response any](processors ...EndpointPreprocessor[request, response]) []EndpointPreprocessor[request, response] {
	if len(processors) > 0 {
		return processors
	}

	return []EndpointPreprocessor[request, response]{}
}

type EndpointPreprocessor[request any, response any] func(endpoint *EndPoint[request, response], origin *gin.Context, dest PreprocessedContext[request, response])

func CheckRequestMethodPreprocessor[request any, response any](endpoint *EndPoint[request, response], origin *gin.Context, dest PreprocessedContext[request, response]) {
	// checking chain is aborted, no need to check
	if origin.IsAborted() {
		return
	}

	if existMethod, allowMethod := endpoint.allowMethods.isAllowed(origin.Request.Method); !existMethod {
		origin.AbortWithStatusJSON(StatusMethodNotAllowed, &FrameworkResponse{
			ErrorCode:    ErrorCodeMethodNotSupported,
			ErrorMessage: values.BuildStringsWithJoin(" ", "method", origin.Request.Method, "is not supported"),
			RequestID:    trace.GetTid(dest),
		})
	} else if !allowMethod {
		origin.AbortWithStatusJSON(StatusMethodNotAllowed, &FrameworkResponse{
			ErrorCode:    ErrorCodeMethodNotSupported,
			ErrorMessage: values.BuildStringsWithJoin(" ", "method", origin.Request.Method, "is not allowed"),
			RequestID:    trace.GetTid(dest),
		})
	}
}

func CheckRequestHeadersPreprocessor[request any, response any](endpoint *EndPoint[request, response], origin *gin.Context, dest PreprocessedContext[request, response]) {
	// checking chain is aborted, no need to check
	if origin.IsAborted() {
		return
	}

	headers := Params{}
	for key, necessary := range endpoint.parsingHeaders {
		if necessary && origin.GetHeader(key) == "" {
			origin.AbortWithStatusJSON(StatusBadRequest, &FrameworkResponse{
				ErrorCode:    ErrorCodeMissingRequiredHeader,
				ErrorMessage: values.BuildStringsWithJoin(" ", "header", key, "is required"),
				RequestID:    trace.GetTid(dest),
			})
			return
		}

		headers[key] = origin.GetHeader(key)
	}

	dest.SetHeaderParams(headers)
}

func LoadNormalRequestHeadersPreprocessor[request any, response any](endpoint *EndPoint[request, response], origin *gin.Context, dest PreprocessedContext[request, response]) {
	// checking chain is aborted, no need to check
	if origin.IsAborted() {
		return
	}

	headers := RequestHeader{
		Accept:         origin.GetHeader("Accept"),
		AcceptEncoding: origin.GetHeader("Accept-Encoding"),
		AcceptLanguage: origin.GetHeader("Accept-Language"),
		UserAgent:      origin.GetHeader("User-Agent"),
		ContentType:    origin.GetHeader("Content-Type"),
		ContentLength:  values.StringToInt(origin.GetHeader("Content-Length"), 0),
		Origin:         origin.GetHeader("Origin"),
		Referer:        origin.GetHeader("Referer"),
		Authorization:  origin.GetHeader("Authorization"),
		ApiKey:         origin.GetHeader("X-Api-Key"),
	}

	dest.SetRequestHeader(headers)
}

func CheckRequestQueriesPreprocessor[request any, response any](endpoint *EndPoint[request, response], origin *gin.Context, dest PreprocessedContext[request, response]) {
	// checking chain is aborted, no need to check
	if origin.IsAborted() {
		return
	}

	queries := Params{}
	for key, necessary := range endpoint.parsingQueries {
		if necessary && origin.Query(key) == "" {
			origin.AbortWithStatusJSON(StatusBadRequest, &FrameworkResponse{
				ErrorCode:    ErrorCodeMissingRequiredQuery,
				ErrorMessage: values.BuildStringsWithJoin(" ", "query", key, "is required"),
				RequestID:    trace.GetTid(dest),
			})
			return
		}

		queries[key] = origin.Query(key)
	}

	dest.SetQueryParams(queries)
}

func CheckRequestParamsPreprocessor[request any, response any](endpoint *EndPoint[request, response], origin *gin.Context, dest PreprocessedContext[request, response]) {
	// checking chain is aborted, no need to check
	if origin.IsAborted() {
		return
	}

	params := Params{}
	for key, necessary := range endpoint.parsingParams {
		if necessary && origin.Param(key) == "" {
			origin.AbortWithStatusJSON(StatusBadRequest, &FrameworkResponse{
				ErrorCode:    ErrorCodeMissingRequiredParam,
				ErrorMessage: values.BuildStringsWithJoin(" ", "param", key, "is required"),
				RequestID:    trace.GetTid(dest),
			})
			return
		}

		params[key] = origin.Param(key)
	}

	dest.SetPathParams(params)
}

func CheckRequestCookiesPreprocessor[request any, response any](endpoint *EndPoint[request, response], origin *gin.Context, dest PreprocessedContext[request, response]) {
	// checking chain is aborted, no need to check
	if origin.IsAborted() {
		return
	}

	cookies := Params{}
	for key, necessary := range endpoint.parsingCookies {
		value, err := origin.Cookie(key)
		if necessary && err != nil {
			origin.AbortWithStatusJSON(StatusBadRequest, &FrameworkResponse{
				ErrorCode:    ErrorCodeMissingRequiredCookie,
				ErrorMessage: values.BuildStringsWithJoin(" ", "cookie", key, "is required"),
				RequestID:    trace.GetTid(dest),
			})
			return
		}
		cookies[key] = value
	}

	dest.SetCookieParams(cookies)
}

func CheckRequestBodyPreprocessor[request any, response any](endpoint *EndPoint[request, response], origin *gin.Context, dest PreprocessedContext[request, response]) {
	// checking chain is aborted, no need to check
	if origin.IsAborted() {
		return
	}

	// read request body
	payload, readErr := origin.GetRawData()
	if readErr != nil {
		origin.AbortWithStatusJSON(StatusBadRequest, &FrameworkResponse{
			ErrorCode:    ErrorCodeInvalidRequestBody,
			ErrorMessage: values.BuildStringsWithJoin(" ", "invalid request body:", readErr.Error()),
			RequestID:    trace.GetTid(dest),
		})
		return
	}

	// unmarshal request body
	requestBody, unmarshalErr := defaultPayloadProcessor[request](origin.ContentType(), payload, nil)
	if unmarshalErr != nil {
		origin.AbortWithStatusJSON(StatusBadRequest, &FrameworkResponse{
			ErrorCode:    ErrorCodeInvalidRequestBody,
			ErrorMessage: values.BuildStringsWithJoin(" ", "invalid request body:", unmarshalErr.Error()),
			RequestID:    trace.GetTid(dest),
		})
		return
	}

	// check request body
	if checkResult := values.CheckStruct(requestBody); checkResult != "" {
		origin.AbortWithStatusJSON(StatusBadRequest, &FrameworkResponse{
			ErrorCode:    ErrorCodeBadRequestBody,
			ErrorMessage: values.BuildStringsWithJoin(" ", "missing required field:", checkResult),
			RequestID:    trace.GetTid(dest),
		})
		return
	}

	dest.SetRequest(requestBody)
}
