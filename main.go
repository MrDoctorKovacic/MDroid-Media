package main

import (
	"encoding/json"
	"flag"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/MrDoctorKovacic/MDroid-Core/settings"
	"github.com/MrDoctorKovacic/MDroid-Media/bluetooth"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Config controls program settings and general persistent settings
type Config struct {
	BluetoothAddress string
}

// Reboot the machine
func reboot(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("OK")
	exec.Command("reboot", "now")
}

func init() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.CallerMarshalFunc = func(file string, line int) string {
		fileparts := strings.Split(file, "/")
		filename := strings.Replace(fileparts[len(fileparts)-1], ".go", "", -1)
		return filename + ":" + strconv.Itoa(line)
	}
	output := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "Mon Jan 2 15:04:05"}
	log.Logger = zerolog.New(output).With().Caller().Timestamp().Logger()
}

// define our router and subsequent routes here
func main() {
	flag.StringVar(&settings.Settings.File, "settings-file", "", "File to recover the persistent settings.")
	debug := flag.Bool("debug", false, "sets log level to debug")
	flag.Parse()

	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// Parse settings file
	settings.ReadFile(settings.Settings.File)

	// Check settings
	if _, err := json.Marshal(settings.GetAll()); err != nil {
		panic(err)
	}

	// Parse through config if found in settings file
	configMap, err := settings.GetComponent("MDROID")
	if err != nil {
		log.Warn().Msg("No config found in settings file, not parsing through config")
	}
	bluetooth.Setup(&configMap)

	// Init router
	router := mux.NewRouter()

	//
	// Main routes
	//
	router.HandleFunc("/restart", reboot).Methods("GET")

	//
	// Bluetooth routes
	//
	router.HandleFunc("/bluetooth", bluetooth.GetDeviceInfo).Methods("GET")
	router.HandleFunc("/bluetooth/getDeviceInfo", bluetooth.GetDeviceInfo).Methods("GET")
	router.HandleFunc("/bluetooth/getMediaInfo", bluetooth.GetMediaInfo).Methods("GET")
	router.HandleFunc("/bluetooth/connect", bluetooth.Connect).Methods("GET")
	router.HandleFunc("/bluetooth/disconnect", bluetooth.HandleDisconnect).Methods("GET")
	router.HandleFunc("/bluetooth/prev", bluetooth.Prev).Methods("GET")
	router.HandleFunc("/bluetooth/next", bluetooth.Next).Methods("GET")
	router.HandleFunc("/bluetooth/pause", bluetooth.Pause).Methods("GET")
	router.HandleFunc("/bluetooth/play", bluetooth.Play).Methods("GET")
	router.HandleFunc("/bluetooth/refresh", bluetooth.ForceRefresh).Methods("GET")

	err = http.ListenAndServe(":5353", router)
	log.Error().Msg(err.Error())
}
