package controllers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	E "github.com/IBM/fp-go/either"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	U "htmx-upload/util"
)

// This test is written by ChatGPT and is based on the provided code.
// It tests the withMultipartFileResource function, which is a utility for handling
func TestWithMultipartFileResource(t *testing.T) {
	t.Run("successfully processes multipart file", func(t *testing.T) {
		e := echo.New()

		req, err := E.Unwrap(E.Chain(func(path string) E.Either[error, *http.Request] {
			return U.NewfileUploadRequest("http://localhost:8080/upload", map[string]string{"name": "Richard Chuo", "email": "good@gmail.com"}, "file", fmt.Sprintf("%s/../../test.txt", path))
		})(E.Eitherize0(os.Getwd)()))

		if err != nil {
			assert.Fail(t, fmt.Sprintf("Failed to create request: %v", err))
		}

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if assert.NoError(t, Upload(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "<p>File test.txt uploaded successfully with fields name=Richard Chuo and email=good@gmail.com.</p>", rec.Body.String())
		}
	})
}
