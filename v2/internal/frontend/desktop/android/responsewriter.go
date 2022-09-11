package android

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"syscall"

	android "github.com/wailsapp/wails/v2/pkg/wailsdroid"
)

type portalResponseWriter struct {
	request *android.ThreadSafeIntercept

	header      http.Header
	wroteHeader bool
	w           io.WriteCloser
	wErr        error
}

func (rw *portalResponseWriter) Header() http.Header {
	if rw.header == nil {
		rw.header = http.Header{}
	}
	return rw.header
}

func (rw *portalResponseWriter) Write(buf []byte) (int, error) {
	rw.WriteHeader(http.StatusOK)
	if rw.wErr != nil {
		return 0, rw.wErr
	}
	return rw.w.Write(buf)
}

func (rw *portalResponseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}
	rw.wroteHeader = true

	contentType := rw.header.Get("Content-Type")
	contentEncoding := rw.header.Get("Content-Encoding")

	// handle text/html; charset=utf-8
	mediaType, params, err := mime.ParseMediaType(contentType)
	if err == nil {
		if mediaType != "" {
			contentType = mediaType
		}
		if encoding, ok := params["charset"]; ok {
			contentEncoding = encoding
		}
	}

	if contentType == "" {
		contentType = "text/html" // hack: no content type was being passed, so default to html
	}
	response := &android.Response{
		HasMimeType: contentType != "",
		MimeType:    contentType,
		HasEncoding: contentEncoding != "",
		Encoding:    contentEncoding,
		StatusCode:  int32(code),
		Headers:     android.HeaderMapToMapStream(rw.header),
	}

	defer func() {
		if response.ReasonPhrase == "" {
			response.ReasonPhrase = http.StatusText(int(response.StatusCode))
		}
		rw.request.Send(response)
	}()

	if code != http.StatusOK {
		rw.w = &nopCloser{io.Discard}
		return
	}

	// We can't use os.Pipe here, because that returns files with a finalizer for closing the FD. But the control over the
	// read FD is given to the InputStream and will be closed there.
	// Furthermore we especially don't want to have the FD_CLOEXEC
	r, w, err := pipe()
	if err != nil {
		rw.wErr = fmt.Errorf("Unable opening pipe: %s", err)
		response.StatusCode = http.StatusInternalServerError
		response.ReasonPhrase = rw.wErr.Error()
		return
	}
	rw.w = w

	response.Data = android.ReaderToInputStream(r)
}

type nopCloser struct {
	io.Writer
}

func (nopCloser) Close() error { return nil }

func (rw *portalResponseWriter) Close() {
	if rw.w != nil {
		rw.w.Close()
	}
}

func pipe() (r *os.File, w *os.File, err error) {
	var p [2]int
	e := syscall.Pipe2(p[0:], 0)
	if e != nil {
		return nil, nil, fmt.Errorf("pipe2: %s", e)
	}

	return os.NewFile(uintptr(p[0]), "|0"), os.NewFile(uintptr(p[1]), "|1"), nil
}
