package controllers

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"

	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	I "github.com/IBM/fp-go/identity"
)

// withMultipartFileResource is a generic utility function that manages the lifecycle
// of a multipart.File resource. It ensures that the file is properly opened and closed
// while executing a provided function.
//
// Type Parameters:
//   - T: The type of the value returned by the provided function.
//
// Parameters:
//   - file (*multipart.FileHeader): The file header representing the uploaded file.
//
// Returns:
//   - A function that takes another function as input. This inner function operates
//     on the opened multipart.File and returns an Either[error, T] result. The outer
//     function ensures that the file is closed after the inner function completes.
//
// Usage:
//
//	This function is useful for safely handling multipart file uploads in a resource-
//	managed way, ensuring that the file is closed even if an error occurs during processing.
func withMultipartFileResource[T any](file *multipart.FileHeader) func(func(multipart.File) E.Either[error, T]) E.Either[error, T] {
	return E.WithResource[error, multipart.File, T](
		E.Eitherize0(file.Open),
		func(file multipart.File) E.Either[error, any] {
			return E.Eitherize0(func() (any, error) {
				return F.ToAny("multipart.File was closed"), file.Close()
			})()
		},
	)
}

// withOSFileResource is a generic utility function that manages the lifecycle of an *os.File
// resource for a given multipart.FileHeader. It ensures that the file is properly created
// and closed, while allowing the caller to perform operations on the file.
//
// The function takes a *multipart.FileHeader as input and returns a higher-order function
// that accepts a callback. The callback is a function that performs operations on the
// *os.File and returns an Either[error, T] result.
//
// The resource management is handled using the E.WithResource function, which ensures
// that the file is created and closed safely, even in the presence of errors.
//
// Type Parameters:
//   - T: The type of the result produced by the callback function.
//
// Parameters:
//   - file (*multipart.FileHeader): The file header representing the uploaded file.
//
// Returns:
//   - func(func(*os.File) E.Either[error, T]) E.Either[error, T]:
//     A higher-order function that takes a callback to operate on the *os.File and
//     returns an Either[error, T] result.
func withOSFileResource[T any](file *multipart.FileHeader) func(func(*os.File) E.Either[error, T]) E.Either[error, T] {
	return E.WithResource[error, *os.File, T](
		E.Eitherize0(func() (*os.File, error) { return os.Create(file.Filename) }),
		func(dst *os.File) E.Either[error, any] {
			return E.Eitherize0(func() (any, error) { return F.ToAny("os.File was closed"), dst.Close() })()
		},
	)
}

// This controller only writes the uploaded file to local file system
// TODO - how to test?
func Upload(c echo.Context) error {
	// Read form fields
	name := c.FormValue("name")
	email := c.FormValue("email")

	//-----------
	// Read file
	//-----------

	// Source
	result := F.Pipe2(
		"file",
		E.Eitherize1(c.FormFile),
		E.Chain(
			func(file *multipart.FileHeader) E.Either[error, *multipart.FileHeader] {
				return withMultipartFileResource[*multipart.FileHeader](file)(
					func(src multipart.File) E.Either[error, *multipart.FileHeader] {
						return withOSFileResource[*multipart.FileHeader](file)(
							func(dst *os.File) E.Either[error, *multipart.FileHeader] {
								return E.Map[error](func(_ int64) *multipart.FileHeader {
									return file
								})(E.Eitherize2(io.Copy)(dst, src))
							})
					})
			}),
	)

	return E.Fold(
		I.Of,
		func(file *multipart.FileHeader) error {
			return c.HTML(http.StatusOK, fmt.Sprintf("<p>File %s uploaded successfully with fields name=%s and email=%s.</p>", file.Filename, name, email))
		},
	)(result)
}
