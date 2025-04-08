package util

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	R "github.com/IBM/fp-go/record"
	T "github.com/IBM/fp-go/tuple"

	AF "github.com/spf13/afero"
)

func ReadTestFile(fileName string) E.Either[error, []byte] {
	return E.Eitherize1(AF.Afero{Fs: AF.NewOsFs()}.ReadFile)(fileName)
}

func withOSFile(path string) func(func(file *os.File) E.Either[error, *http.Request]) E.Either[error, *http.Request] {
	return E.WithResource[error, *os.File, *http.Request](
		func() E.Either[error, *os.File] {
			return E.Eitherize1(os.Open)(path)
		},
		func(file *os.File) E.Either[error, any] {
			return E.Eitherize0(func() (any, error) {
				return F.ToAny("file was closed"), file.Close()
			})()
		},
	)
}

func withMultipartWriter(body *bytes.Buffer) func(func(writer *multipart.Writer) E.Either[error, *http.Request]) E.Either[error, *http.Request] {
	return E.WithResource[error, *multipart.Writer, *http.Request](
		func() E.Either[error, *multipart.Writer] {
			return E.Right[error](multipart.NewWriter(body))
		},
		func(writer *multipart.Writer) E.Either[error, any] {
			return E.Eitherize0(func() (any, error) {
				return F.ToAny("multipart.Writer was closed"), writer.Close()
			})()
		},
	)
}

func readAllE(file *os.File) E.Either[error, []byte] {
	return E.Eitherize1(io.ReadAll)(file)
}

func fileStatE(file *os.File) E.Either[error, os.FileInfo] {
	return E.Eitherize0(file.Stat)()
}

func writeField(w *multipart.Writer) func(T.Tuple2[string, string]) E.Either[error, any] {
	return func(t T.Tuple2[string, string]) E.Either[error, any] {
		return E.Eitherize2(func(key string, value string) (any, error) {
			return F.ToAny("writeField"), w.WriteField(key, value)
		})(t.F1, t.F2)
	}
}

func traverseTuple(f func(T.Tuple2[string, string]) E.Either[error, any]) func([]T.Tuple2[string, string]) E.Either[error, []any] {
	return func(ta []T.Tuple2[string, string]) E.Either[error, []any] {
		return E.TraverseArray(f)(ta)
	}
}

// Creates a new file upload http request with optional extra params
func NewfileUploadRequest(uri string, params map[string]string, paramName, path string) E.Either[error, *http.Request] {
	body := new(bytes.Buffer)

	return withOSFile(path)(func(file *os.File) E.Either[error, *http.Request] {
		return E.Chain(func(fileContents []byte) E.Either[error, *http.Request] {
			return withMultipartWriter(body)(func(writer *multipart.Writer) E.Either[error, *http.Request] {
				return E.Chain(
					func(fi os.FileInfo) E.Either[error, *http.Request] {
						return E.Chain(
							func(part io.Writer) E.Either[error, *http.Request] {
								return E.Chain(
									func(_ int) E.Either[error, *http.Request] {
										return E.Chain(func(_ []any) E.Either[error, *http.Request] {
											return E.Map[error](func(req *http.Request) *http.Request {
												req.Header.Set("Content-Type", writer.FormDataContentType())

												return req
											})(E.Eitherize3(http.NewRequest)("POST", uri, body))
										})(traverseTuple(writeField(writer))(R.ToArray(params)))
									},
								)(E.Eitherize1(part.Write)(fileContents))
							},
						)(E.Eitherize2(writer.CreateFormFile)(paramName, fi.Name()))
					},
				)(fileStatE(file))
			})
		})(readAllE(file))
	})
}
