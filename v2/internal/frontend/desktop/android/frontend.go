package android

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"text/template"

	"github.com/wailsapp/wails/v2/internal/binding"
	"github.com/wailsapp/wails/v2/internal/frontend"
	wailsruntime "github.com/wailsapp/wails/v2/internal/frontend/runtime"
	"github.com/wailsapp/wails/v2/internal/logger"
	"github.com/wailsapp/wails/v2/pkg/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options"
	android "github.com/wailsapp/wails/v2/pkg/wailsdroid"
)

const startURL = "wails://wails/"

type Frontend struct {

	// Context
	ctx context.Context

	frontendOptions *options.App
	logger          *logger.Logger
	debug           bool

	// Assets
	assets   *assetserver.AssetServer
	startURL *url.URL

	// main window handle
	mainWindow android.AppPortal
	bindings   *binding.Bindings
	dispatcher frontend.Dispatcher

	interceptHandler chan *android.ThreadSafeIntercept
	messageHandler   chan string
}

func init() {
	runtime.LockOSThread()
}

func NewFrontend(ctx context.Context, appoptions *options.App, myLogger *logger.Logger, appBindings *binding.Bindings, dispatcher frontend.Dispatcher) *Frontend {
	result := &Frontend{
		frontendOptions:  appoptions,
		logger:           myLogger,
		bindings:         appBindings,
		dispatcher:       dispatcher,
		interceptHandler: make(chan *android.ThreadSafeIntercept),
		messageHandler:   make(chan string),
		ctx:              ctx,
	}
	result.mainWindow = android.ReceiveAppPortal()
	if result.mainWindow == nil {
		log.Fatal("no AppPortal has been offered, cannot create frontend")
	}
	result.startURL, _ = url.Parse(startURL)

	if _starturl, _ := ctx.Value("starturl").(*url.URL); _starturl != nil {
		result.startURL = _starturl
	} else {
		bindingsJSON, err := appBindings.ToJSON()
		if err != nil {
			log.Fatal(err)
		}

		assets, err := assetserver.NewAssetServerMainPage(bindingsJSON, appoptions, ctx.Value("assetdir") != nil, myLogger, wailsruntime.RuntimeAssetsBundle)
		if err != nil {
			log.Fatal(err)
		}
		result.assets = assets

		// Start 10 processors to handle requests in parallel
		for i := 0; i < 10; i++ {
			go result.startRequestProcessor()
		}
	}

	go result.startMessageProcessor()

	var _debug = ctx.Value("debug")
	if _debug != nil {
		result.debug = _debug.(bool)
	}

	return result
}

func (f *Frontend) startMessageProcessor() {
	for message := range f.messageHandler {
		f.processMessage(message)
	}
}

func (f *Frontend) processMessage(message string) {
	if message == "DomReady" {
		if f.frontendOptions.OnDomReady != nil {
			f.frontendOptions.OnDomReady(f.ctx)
		}
		return
	}

	if message == "drag" {
		if !f.mainWindow.IsFullScreen() {
			f.startDrag()
		}
		return
	}

	if message == "runtime:ready" {
		cmd := fmt.Sprintf("window.wails.setCSSDragProperties('%s', '%s');", f.frontendOptions.CSSDragProperty, f.frontendOptions.CSSDragValue)
		f.ExecJS(cmd)
		return
	}

	go func() {
		result, err := f.dispatcher.ProcessMessage(message, f)
		if err != nil {
			f.logger.Error(err.Error())
			f.Callback(result)
			return
		}
		if result == "" {
			return
		}

		switch result[0] {
		case 'c':
			// Callback from a method call
			f.Callback(result[1:])
		default:
			f.logger.Info("Unknown message returned from dispatcher: %+v", result)
		}
	}()
}

func (f *Frontend) Callback(message string) {
	f.ExecJS(`window.wails.Callback(` + strconv.Quote(message) + `);`)
}

func (f *Frontend) ExecJS(js string) {
	f.mainWindow.ExecJS(js)
}

func (f *Frontend) startRequestProcessor() {
	for request := range f.interceptHandler {
		f.processRequest(request)
	}
}

func (f *Frontend) processRequest(request *android.ThreadSafeIntercept) {
	rw := &portalResponseWriter{request: request}
	defer rw.Close()
	f.assets.ProcessHTTPRequestLegacy(
		rw,
		func() (*http.Request, error) {
			req, err := http.NewRequest(http.MethodGet, request.Request.URL, nil)
			if err != nil {
				return nil, err
			}

			if req.URL.Host != f.startURL.Host {
				if req.Body != nil {
					req.Body.Close()
				}

				return nil, fmt.Errorf("Expected host '%s' in request, but was '%s'", f.startURL.Host, req.URL.Host)
			}

			return req, nil
		},
	)
}

func (f *Frontend) Run(ctx context.Context) error {
	f.ctx = context.WithValue(ctx, "frontend", f)

	go func() {
		if f.frontendOptions.OnStartup != nil {
			f.frontendOptions.OnStartup(f.ctx)
		}
	}()

	f.mainWindow.SetWebViewClientPortal(android.NewWebviewClientPortal(f.interceptHandler, f.messageHandler))
	f.mainWindow.Run(f.startURL.String())

	return nil
}

func (f *Frontend) RunMainLoop() {
	//select {}
}

func (f *Frontend) Hide() {
	f.mainWindow.Hide()
}

func (f *Frontend) Show() {
	f.mainWindow.Show()
}

func (f *Frontend) Quit() {
	f.mainWindow.Quit()
}

func (f *Frontend) WindowClose() {
	//f.mainWindow.Quit()
}

type EventNotify struct {
	Name string        `json:"name"`
	Data []interface{} `json:"data"`
}

func (f *Frontend) Notify(name string, data ...interface{}) {
	notification := EventNotify{
		Name: name,
		Data: data,
	}
	payload, err := json.Marshal(notification)
	if err != nil {
		f.logger.Error(err.Error())
		return
	}
	f.mainWindow.ExecJS(`window.wails.EventsNotify('` + template.JSEscapeString(string(payload)) + `');`)
}

func (f *Frontend) startDrag() {
	f.mainWindow.StartDrag()
}
