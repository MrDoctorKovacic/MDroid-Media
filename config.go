package bluetooth

import (
	"regexp"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

// Regex expressions for parsing dbus output
var (
	BluetoothAddress string
	tmuxStarted      bool
	replySerialRegex *regexp.Regexp
	findStringRegex  *regexp.Regexp
	cleanRegex       *regexp.Regexp
)

func init() {
	replySerialRegex = regexp.MustCompile(`(.*reply_serial=2\n\s*variant\s*)array`)
	findStringRegex = regexp.MustCompile(`string\s"(.*)"|uint32\s(\d)+`)
	cleanRegex = regexp.MustCompile(`(string|uint32|\")+`)
}

// Setup bluetooth with address
func Setup(configAddr *map[string]string) {
	configMap := *configAddr
	bluetoothAddress, usingBluetooth := configMap["BLUETOOTH_ADDRESS"]
	if !usingBluetooth {
		log.Warn().Msg("No bluetooth address found in config, using empty address")
		BluetoothAddress = ""
		return
	}

	SetAddress(bluetoothAddress)
	log.Info().Msg("Enabling auto refresh of BT address")
	go startAutoRefresh()
}

// SetRoutes handles module routing
func SetRoutes(router *mux.Router) {
	//
	// Bluetooth routes
	//
	router.HandleFunc("/bluetooth", GetDeviceInfo).Methods("GET")
	router.HandleFunc("/bluetooth/getDeviceInfo", GetDeviceInfo).Methods("GET")
	router.HandleFunc("/bluetooth/getMediaInfo", GetMediaInfo).Methods("GET")
	router.HandleFunc("/bluetooth/connect", HandleConnect).Methods("GET")
	router.HandleFunc("/bluetooth/disconnect", HandleConnect).Methods("GET")
	router.HandleFunc("/bluetooth/prev", Prev).Methods("GET")
	router.HandleFunc("/bluetooth/next", Next).Methods("GET")
	router.HandleFunc("/bluetooth/pause", HandlePause).Methods("GET")
	router.HandleFunc("/bluetooth/play", HandlePlay).Methods("GET")
	router.HandleFunc("/bluetooth/refresh", ForceRefresh).Methods("GET")
}
