package huebridgex

import (
	"github.com/julienschmidt/httprouter"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var handlerMap map[string]huestate
var correctedNameMap map[string]string

func init() {
	log.SetOutput(ioutil.Discard)
	handlerMap = make(map[string]huestate)
	correctedNameMap = make(map[string]string)
	upnpTemplateInit()
}

func SetLogger(w io.Writer) {
	log.SetOutput(w)
}

func ListenAndServe(addr string) error {
	router := httprouter.New()
	router.GET(upnp_uri, upnpSetup(addr))

	router.GET("/api/:userId", getLightsList)
	router.PUT("/api/:userId/lights/:lightId/state", setLightState)
	router.GET("/api/:userId/lights/:lightId", getLightInfo)

	go upnpResponder(addr, upnp_uri)
	return http.ListenAndServe(addr, requestLogger(router))
}

// Handler state is the state of the "light" after the handler function if error is set to true echo will reply with "sorry the device is not responding"
type Handler func(Request, *Response)

func isAlpha(r rune) bool {
	switch {
	case 'a' <= r && r <= 'z':
		return true
	case 'A' <= r && r <= 'Z':
		return true
	default:
		return false
	}
}

func getHueStateByName(name string) (huestate, bool) {
	if correctedName, corrected := correctedNameMap[name]; corrected {
		name = correctedName
	}
	hstate, ok := handlerMap[name]
	return hstate, ok
}

func setHueStateByName(name string, state huestate) {
	if correctedName, corrected := correctedNameMap[name]; corrected {
		name = correctedName
	}
	handlerMap[name] = state
}

// Handle will set the device state
func Handle(deviceName string, h Handler) {
	log.Println("[HANDLE]", deviceName)
	removeIllegalRunes := func(r rune) rune {
		if isAlpha(r) {
			return r
		}
		return -1
	}
	correctedName := strings.Map(removeIllegalRunes, deviceName)
	correctedNameMap[correctedName] = deviceName
	handlerMap[deviceName] = huestate{
		Handler: h,
		OnState: false,
	}
}

func requestLogger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//log.Println("[WEB]", r.RemoteAddr, r.Method, r.URL)
		//		log.Printf("\t%+v\n", r)
		h.ServeHTTP(w, r)
	})
}
