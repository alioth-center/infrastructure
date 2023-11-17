package http

type Method = string

const (
	GET     Method = "GET"
	POST    Method = "POST"
	OPTIONS Method = "OPTIONS"
	HEAD    Method = "HEAD"
	PUT     Method = "PUT"
	DELETE  Method = "DELETE"
	TRACE   Method = "TRACE"
	CONNECT Method = "CONNECT"
	PATCH   Method = "PATCH"
)

type UserAgent = string

const (
	Curl      UserAgent = "curl/7.64.1"
	Postman   UserAgent = "PostmanRuntime/7.26.8"
	ChromeOSX UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36"
	Safari    UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Safari/605.1.15"
	Firefox   UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:109.0) Gecko/20100101 Firefox/119.0"
)

type ContentType = string

const (
	ContentTypeJson       ContentType = "application/json"
	ContentTypeForm       ContentType = "application/x-www-form-urlencoded"
	ContentTypeFileStream ContentType = "application/octet-stream"
	ContentTypeMultipart  ContentType = "multipart/form-data"
	ContentTypeTextPlain  ContentType = "text/plain"
	ContentTypeTextHtml   ContentType = "text/html"
	ContentTypeTextXml    ContentType = "text/xml"
	ContentTypeTextYaml   ContentType = "text/yaml"
	ContentTypeTextCsv    ContentType = "text/csv"
	ContentTypeImagePng   ContentType = "image/png"
	ContentTypeImageJpeg  ContentType = "image/jpeg"
	ContentTypeImageGif   ContentType = "image/gif"
)
