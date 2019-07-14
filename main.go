package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os/exec"

	"github.com/MrDoctorKovacic/MDroid-Core/bluetooth"
	"github.com/MrDoctorKovacic/MDroid-Core/status"
	"github.com/MrDoctorKovacic/MDroid-Core/settings"
	"github.com/gorilla/mux"
)

// Config controls program settings and general persistent settings
type Config struct {
	BluetoothAddress string
}

// MainStatus will control logging and reporting of status / warnings / errors
var MainStatus = status.NewStatus("Main")

// Reboot the machine
func reboot(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("OK")
	exec.Command("reboot", "now")
}

// define our router and subsequent routes here
func main() {

	// Start with program arguments
	var (
		settingsFile string
	)
	flag.StringVar(&settingsFile, "settings-file", "", "File to recover the persistent settings.")
	flag.Parse()

	// Parse settings file
	settingsData, useVerboseOutput := settings.SetupSettings(settingsFile)

	// Log settings
	out, err := json.Marshal(settingsData)
	if err != nil {
		panic(err)
	}
	MainStatus.Log(status.OK(), "Using settings: "+string(out))

	// Parse through config if found in settings file
	config, ok := settingsData["CONFIG"]
	if ok {
		// Set up bluetooth
		bluetoothAddress, usingBluetooth := config["BLUETOOTH_ADDRESS"]
		if usingBluetooth {
			bluetooth.EnableAutoRefresh()
			bluetooth.SetAddress(bluetoothAddress)
		}
	} else {
		MainStatus.Log(status.Warning(), "No config found in settings file, not parsing through config")
	}

	// Init router
	router := mux.NewRouter()

	//
	// Main routes
	//
	router.HandleFunc("/restart", reboot).Methods("GET")

	//
	// Settings routes
	//
	router.HandleFunc("/settings", settings.GetAllSettings).Methods("GET")
	router.HandleFunc("/settings/{component}", settings.GetSetting).Methods("GET")
	router.HandleFunc("/settings/{component}/{name}", settings.GetSettingValue).Methods("GET")
	router.HandleFunc("/settings/{component}/{name}/{value}", settings.SetSettingValue).Methods("POST")

	//
	// Bluetooth routes
	//
	router.HandleFunc("/bluetooth", bluetooth.GetDeviceInfo).Methods("GET")
	router.HandleFunc("/bluetooth/getDeviceInfo", bluetooth.GetDeviceInfo).Methods("GET")
	router.HandleFunc("/bluetooth/getMediaInfo", bluetooth.GetMediaInfo).Methods("GET")
	router.HandleFunc("/bluetooth/connect", bluetooth.Connect).Methods("GET")
	router.HandleFunc("/bluetooth/disconnect", bluetooth.Connect).Methods("GET")
	router.HandleFunc("/bluetooth/prev", bluetooth.Prev).Methods("GET")
	router.HandleFunc("/bluetooth/next", bluetooth.Next).Methods("GET")
	router.HandleFunc("/bluetooth/pause", bluetooth.Pause).Methods("GET")
	router.HandleFunc("/bluetooth/play", bluetooth.Play).Methods("GET")
	router.HandleFunc("/bluetooth/refresh", bluetooth.ForceRefresh).Methods("GET")

	log.Fatal(http.ListenAndServe(":5353", router))
}
