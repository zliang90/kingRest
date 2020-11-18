package errors

var Env string

func SetEnv(env string) {
	Env = env
}

func InternalServerError(err error) *APIError {
	return NewAPIError("INTERNAL_SERVER_ERROR", Params{"error": err.Error()})
}

func NotFound(resource string) *APIError {
	return NewAPIError("NOT_FOUND", Params{"resource": resource})
}

func Unauthorized(err error) *APIError {
	return NewAPIError("UNAUTHORIZED", Params{"error": err})
}

func BadRequest(err error) *APIError {
	return NewAPIError("BAD_REQUEST", Params{"error": err.Error()})
}
