package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"

	"github.com/grafana/grafana-kiosk/pkg/initialize"
	"github.com/grafana/grafana-kiosk/pkg/kiosk"
)

// LoginMethod specifies the type of login to be used by the kiosk
type LoginMethod int

// Login Methods
const (
	ANONYMOUS LoginMethod = 0
	LOCAL     LoginMethod = 1
	GCOM      LoginMethod = 2
)

// KioskMode specifes the mode of the kiosk
type KioskMode int

// Kiosk Modes
const (
	// TV will hide the sidebar but allow usage of menu
	TV KioskMode = 0
	// NORMAL will disable sidebar and top navigation bar
	NORMAL KioskMode = 1
	// NONE will not enter kiosk mode
	NONE KioskMode = 2
)

var (
	loginMethod = LOCAL
	kioskMode   = NORMAL
)

func main() {
	var Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %v\n", os.Args[0])
		flag.PrintDefaults()
	}
	urlPtr := flag.String("URL", "https://play.grafana.org", "URL to Grafana server (Required)")
	methodPtr := flag.String("login-method", "anon", "login method: [anon|local|gcom]")
	usernamePtr := flag.String("username", "guest", "username (Required)")
	passwordPtr := flag.String("password", "guest", "password (Required)")
	// kiosk=tv includes sidebar menu
	// kiosk no sidebar ever
	kioskModePtr := flag.String("kiosk-mode", "default", "kiosk mode [default|tv|false]")
	autoFit := flag.Bool("autofit", true, "autofit panels in kiosk mode")
	// when the URL is a playlist, append "inactive" to the URL
	isPlayList := flag.Bool("playlist", false, "URL is a playlist: [true|false]")
	LXDEEnabled := flag.Bool("lxde", true, "initialize LXDE for kiosk mode")
	LXDEHomePtr := flag.String("lxde-home", "/home/pi", "path to home directory of LXDE user running X Server")
	flag.Parse()

	// make sure the url has content
	if *urlPtr == "" {
		Usage()
		os.Exit(1)
	}
	// validate url
	_, err := url.ParseRequestURI(*urlPtr)
	if err != nil {
		Usage()
		panic(err)
	}

	if *isPlayList == true {
		println("playlist")
	}

	if *LXDEEnabled == true {
		initialize.LXDE(*LXDEHomePtr)
	}
	switch *kioskModePtr {
	case "tv": // NO SIDEBAR ACCESS
		kioskMode = TV
	case "false": // DO NOT USE KIOSK MODE
		kioskMode = NONE
	case "default": // NO TOPNAV or SIDEBAR
		kioskMode = NORMAL
	default:
		kioskMode = NORMAL
	}

	switch *methodPtr {
	case "anon":
		loginMethod = ANONYMOUS
	case "local":
		loginMethod = LOCAL
	case "gcom":
		loginMethod = GCOM
	default:
		loginMethod = ANONYMOUS
	}

	switch loginMethod {
	case LOCAL:
		println("Launching local login kiosk")
		kiosk.GrafanaKioskLocal(urlPtr, usernamePtr, passwordPtr, *autoFit)
	case GCOM:
		println("Launching GCOM login kiosk")
		kiosk.GrafanaKioskGCOM(urlPtr, usernamePtr, passwordPtr, *autoFit)
	case ANONYMOUS:
		println("Launching ANON login kiosk")
		kiosk.GrafanaKioskAnonymous(urlPtr, *autoFit)
	}
}
