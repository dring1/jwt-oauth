package errors

import (
	"net/http"

	"github.com/dring1/jwt-oauth/lib/contextkeys"
)

const RecordNotFound = "Record Not Found"
const UnauthorizedUser = "Unauthorized User"
const FailedToAuthenticate = "Failed to authenticate user"
const InternalServer = "Internal Server Error"
const InvalidToken = "Invalid Token"

func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	existingError := r.Context().Value(contextkeys.Error)
	if existingError == nil {
		existingError = "An error occurred"
	}
	w.Write([]byte(existingError.(string)))
}
