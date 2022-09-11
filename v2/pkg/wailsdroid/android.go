// this package is called "wailsdroid" for saner importing in Android projects
// (android ns taken)
package wailsdroid

import (
	"io"
	"log"
	"strings"

	"github.com/wailsapp/wails/v2/internal/frontend"
)

const protocol = "wails://"

// was using a channel here, but received nil values? not sure why
var appPortal AppPortal

func OfferAppPortal(portal AppPortal) {
	if portal == nil {
		log.Fatal("refusing to accept offered AppPortal: it is nil")
	}
	appPortal = portal
}

func ReceiveAppPortal() AppPortal {
	return appPortal
}

// loose mapping of Java WebResourceRequest
type Request struct {
	Method string
	URL    string
}

// loose mapping of Java WebResourceResponse
type Response struct {
	// these values are nilable, but *string not supported by gomobile or I'm doing something wrong
	HasMimeType  bool
	MimeType     string
	HasEncoding  bool
	Encoding     string
	StatusCode   int32
	ReasonPhrase string
	// not supported by gomobile or I'm doing something wrong
	// Headers map[string]string
	Headers MapStream
	Data    InputStream
}

// loose mapping of Java InputStream
type InputStream interface {
	// return (-1, nil) on EOF: https://docs.oracle.com/javase/6/docs/api/java/io/InputStream.html#read()
	Read() (int32, error)
	Close() error
}

// wrap golang reader to mostly-Java-compat InputStream
type readerInputStream struct {
	reader io.Reader
}

func ReaderToInputStream(reader io.Reader) InputStream {
	return &readerInputStream{reader}
}

func (ris *readerInputStream) Read() (int32, error) {
	buf := []byte{0}
	_, err := ris.reader.Read(buf)
	if err == io.EOF {
		return -1, nil
	}
	if err != nil {
		return 0, err
	}
	return int32(buf[0]), nil
}

func (ris *readerInputStream) Close() error {
	var err error
	if closer, ok := ris.reader.(io.Closer); ok {
		err = closer.Close()
	}
	return err
}

type OpenDialogOptions interface {
	DefaultDirectory() string
	DefaultFilename() string
	Title() string
	// Filters() FileFilterIterator
	ShowHiddenFiles() bool
	CanCreateDirectories() bool
	ResolvesAliases() bool
	TreatPackagesAsDirectories() bool
}

type openDialogOptions struct {
	options frontend.OpenDialogOptions
}

func (o *openDialogOptions) DefaultDirectory() string {
	return o.options.DefaultDirectory
}
func (o *openDialogOptions) DefaultFilename() string {
	return o.options.DefaultFilename
}
func (o *openDialogOptions) Title() string {
	return o.options.Title
}

/*
	func (o *openDialogOptions) Filters() FileFilterIterator {
		return NewFileFilterIterator(o.options.Filters)
	}
*/
func (o *openDialogOptions) ShowHiddenFiles() bool {
	return o.options.ShowHiddenFiles
}
func (o *openDialogOptions) CanCreateDirectories() bool {
	return o.options.CanCreateDirectories
}
func (o *openDialogOptions) ResolvesAliases() bool {
	return o.options.ResolvesAliases
}
func (o *openDialogOptions) TreatPackagesAsDirectories() bool {
	return o.options.TreatPackagesAsDirectories
}

func NewOpenDialogOptions(o frontend.OpenDialogOptions) OpenDialogOptions {
	return &openDialogOptions{o}
}

type StringOrError struct {
	Result string
	Error  string
}

type StringIteratorOrError struct {
	Result StringIterator
	Error  string
}

/*
type ScreenIteratorOrError struct {
	Result ScreenIterator
	Error  string
}
*/

type SaveDialogOptions interface {
	DefaultDirectory() string
	DefaultFilename() string
	Title() string
	// Filters() FileFilterIterator
	ShowHiddenFiles() bool
	CanCreateDirectories() bool
	TreatPackagesAsDirectories() bool
}

type saveDialogOptions struct {
	options frontend.SaveDialogOptions
}

func (o *saveDialogOptions) DefaultDirectory() string {
	return o.options.DefaultDirectory
}
func (o *saveDialogOptions) DefaultFilename() string {
	return o.options.DefaultFilename
}
func (o *saveDialogOptions) Title() string {
	return o.options.Title
}

/*
	func (o *saveDialogOptions) Filters() FileFilterIterator {
		return NewFileFilterIterator(o.options.Filters)
	}
*/
func (o *saveDialogOptions) ShowHiddenFiles() bool {
	return o.options.ShowHiddenFiles
}
func (o *saveDialogOptions) CanCreateDirectories() bool {
	return o.options.CanCreateDirectories
}
func (o *saveDialogOptions) TreatPackagesAsDirectories() bool {
	return o.options.TreatPackagesAsDirectories
}

func NewSaveDialogOptions(o frontend.SaveDialogOptions) SaveDialogOptions {
	return &saveDialogOptions{o}
}

type MessageDialogOptions interface {
	Type() string
	Title() string
	Message() string
	Buttons() StringIterator
	DefaultButton() string
	CancelButton() string
	Icon() []byte
}

type messageDialogOptions struct {
	options frontend.MessageDialogOptions
}

func (o *messageDialogOptions) Type() string {
	return string(o.options.Type)
}
func (o *messageDialogOptions) Title() string {
	return o.options.Title
}
func (o *messageDialogOptions) Message() string {
	return o.options.Message
}
func (o *messageDialogOptions) Buttons() StringIterator {
	return NewStringIterator(o.options.Buttons)
}
func (o *messageDialogOptions) DefaultButton() string {
	return o.options.DefaultButton
}
func (o *messageDialogOptions) CancelButton() string {
	return o.options.CancelButton
}
func (o *messageDialogOptions) Icon() []byte {
	return o.options.Icon
}

func NewMessageDialogOptions(o frontend.MessageDialogOptions) MessageDialogOptions {
	return &messageDialogOptions{o}
}

// AppPortal allows interacting with the Android activity
type AppPortal interface {
	SetWebViewClientPortal(WebviewClientPortal)

	Run(startURL string)
	Hide()
	Show()
	Quit()

	BrowserOpenURL(url string)

	// OpenFileDialog(OpenDialogOptions) StringOrError
	// OpenMultipleFilesDialog(OpenDialogOptions) StringIteratorOrError
	// OpenDirectoryDialog(OpenDialogOptions) StringOrError
	// SaveFileDialog(SaveDialogOptions) StringOrError
	// MessageDialog(MessageDialogOptions)

	// ScreenGetAll() ScreenIteratorOrError

	IsFullScreen() bool
	StartDrag()

	ExecJS(string)
}

type WebviewClientPortal interface {
	// Cancel a request by returning true
	ShouldOverrideUrlLoading(request *Request) bool
	// Handle a request by returning a non-nil response
	ShouldInterceptRequest(request *Request) *Response
	// Handle a WailsInvoke message
	ReceiveMessage(message string)
}

func NewWebviewClientPortal(interceptHandler chan *ThreadSafeIntercept, messageHandler chan string) WebviewClientPortal {
	return &ourBeautifulPortal{
		interceptHandler,
		messageHandler,
	}
}

type ThreadSafeIntercept struct {
	Request      *Request
	responseChan chan *Response
}

func (tsi *ThreadSafeIntercept) Send(resp *Response) {
	tsi.responseChan <- resp
	close(tsi.responseChan)
}

type ourBeautifulPortal struct {
	interceptHandler chan *ThreadSafeIntercept
	messageHandler   chan string
}

func (p *ourBeautifulPortal) ShouldOverrideUrlLoading(request *Request) bool {
	// prevent non-wails, non-javascript URLs from loading in our embedded webview
	return !strings.HasPrefix(request.URL, protocol) && !strings.HasPrefix(request.URL, "javascript:")
}

func (p *ourBeautifulPortal) ShouldInterceptRequest(request *Request) *Response {
	tsi := &ThreadSafeIntercept{
		request,
		make(chan *Response),
	}

	p.interceptHandler <- tsi
	return <-tsi.responseChan
}

func (p *ourBeautifulPortal) ReceiveMessage(message string) {
	p.messageHandler <- message
}
