package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/mbndr/figlet4go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/withmandala/go-log"
	"gopkg.in/yaml.v3"
)

// Global variables
var (
	logger                   *log.Logger
	temperatureGauge         *prometheus.GaugeVec
	humidityGauge            *prometheus.GaugeVec
	valveTappetGauge         *prometheus.GaugeVec
	setpointTemperatureGauge *prometheus.GaugeVec
	configPath               string
	httpBind                 string
	httpPort                 string
	bshcHost                 string
	bshcPort                 string
	bshcClientCert           string
	bshcClientKey            string
	debug                    bool
	configPathDefault        = "config/config.yaml"
	httpBindDefault          = ""
	httpPortDefault          = ""
	bshcHostDefault          = ""
	bshcPortDefault          = ""
	bshcClientCertDefault    = ""
	bshcClientKeyDefault     = ""
	skipTlsVerify			bool
	skipTlsVerifyDefault	= false
	c                        conf
	devices                  = make(map[string]interface{})
	rooms                    = make(map[string]interface{})
)

// Config struct
type conf struct {
	HTTP struct {
		Bind string `yaml:"bind"`
		Port string `yaml:"port"`
	} `yaml:"http"`

	BSHC struct {
		Host       string `yaml:"host"`
		Port       string `yaml:"port"`
		ClientCert string `yaml:"client_cert"`
		ClientKey  string `yaml:"client_key"`
		SkipTLSVerify bool `yaml:"skip_tls_verify"`
	} `yaml:"bshc"`

	SERVICES struct {
		TemperatureLevel bool `yaml:"temperature_level"`
		HumidityLevel    bool `yaml:"humidity_level"`
		ValveTappet      bool `yaml:"valve_tappet"`
	} `yaml:"services"`
}

// Load config file
func (c *conf) getConf(configPath string) *conf {
	logger.Infof("Loading configuration from %s", configPath)
	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		logger.Fatalf("Failed to read config file: %v", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		logger.Fatalf("Failed to unmarshal config file: %v", err)
	}
	logger.Info("Configuration loaded successfully")
	return c
}

// Initialize logger
func initLogger(debug bool) *log.Logger {
	logger := log.New(os.Stderr).WithColor()
	if debug {
		logger = logger.WithDebug()
		logger.Info("Debug mode enabled")
	} else {
		logger.Info("Debug mode disabled")
	}
	return logger
}

// Make GET request
func makeGetRequest(url, clientCert, clientKey string) (*http.Response, error) {
	logger.Debugf("Making GET request to URL: %s", url)

	// Load client cert
	cert, err := tls.LoadX509KeyPair(clientCert, clientKey)
	if err != nil {
		logger.Errorf("Could not load client certificate: %v", err)
		return nil, err
	}
	logger.Debug("Client certificate loaded successfully")

	// Setup HTTPS client
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: transport}
	
	if skipTlsVerify {
		tlsConfig.InsecureSkipVerify = true
		logger.Debug("TLS verification skipped")
	} else {
		tlsConfig.InsecureSkipVerify = false
		logger.Debug("TLS verification enabled")
	}	
	
	logger.Debug("HTTPS client configured successfully")

	// Make GET request
	resp, err := client.Get(url)
	if err != nil {
		logger.Errorf("Could not make GET request: %v", err)
		return nil, err
	}
	logger.Infof("GET request to URL %s completed with status code: %d", url, resp.StatusCode)

	return resp, nil
}

func getDeviceNames() {
	logger.Info("Fetching device names")

	// Make GET request to devices endpoint
	devicesURL := fmt.Sprintf("https://%s:%s/smarthome/devices", bshcHost, bshcPort)
	resp, err := makeGetRequest(devicesURL, bshcClientCert, bshcClientKey)
	if err != nil {
		logger.Errorf("Failed to get devices: %v", err)
		return
	}
	defer resp.Body.Close()

	// Check if response status code is 200
	if resp.StatusCode != http.StatusOK {
		logger.Errorf("Unexpected status code: %d", resp.StatusCode)
		return
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("Failed to read response body: %v", err)
		return
	}

	// Parse JSON response
	var devicesArray []map[string]interface{}
	err = json.Unmarshal(body, &devicesArray)
	if err != nil {
		logger.Errorf("Failed to unmarshal devices response: %v", err)
		return
	}

	// Save devices with ids, names, and room ids to a map
	for _, device := range devicesArray {
		if device["deviceModel"] == "VENTILATION_SERVICE" || device["deviceModel"] == "HUE_BRIDGE_MANAGER" {
			continue
		}

		deviceID, ok := device["id"].(string)
		if !ok {
			logger.Errorf("Invalid device id format for device: %v", device)
			continue
		}
		deviceName, ok := device["name"].(string)
		if !ok {
			logger.Errorf("Invalid device name format for device: %v", device)
			continue
		}
		roomID, ok := device["roomId"].(string)
		if !ok {
			logger.Errorf("Invalid room id format for device: %v", device)
			continue
		}
		devices[deviceID] = map[string]string{
			"name":   deviceName,
			"roomId": roomID,
		}
		logger.Debugf("Device added: ID=%s, Name=%s, RoomID=%s", deviceID, deviceName, roomID)
	}
	logger.Info("Device names fetched successfully")
}

func getRoomNames() {
	logger.Info("Fetching room names")

	// Make GET request to rooms endpoint
	roomsURL := fmt.Sprintf("https://%s:%s/smarthome/rooms", bshcHost, bshcPort)
	resp, err := makeGetRequest(roomsURL, bshcClientCert, bshcClientKey)
	if err != nil {
		logger.Errorf("Failed to get rooms: %v", err)
		return
	}
	defer resp.Body.Close()

	// Check if response status code is 200
	if resp.StatusCode != http.StatusOK {
		logger.Errorf("Unexpected status code: %d", resp.StatusCode)
		return
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("Failed to read response body: %v", err)
		return
	}

	// Parse JSON response
	var roomsArray []map[string]interface{}
	err = json.Unmarshal(body, &roomsArray)
	if err != nil {
		logger.Errorf("Failed to unmarshal rooms response: %v", err)
		return
	}

	// Save rooms with ids and names to a map
	for _, room := range roomsArray {
		roomID, ok := room["id"].(string)
		if !ok {
			logger.Errorf("Invalid room id format for room: %v", room)
			continue
		}
		roomName, ok := room["name"].(string)
		if !ok {
			logger.Errorf("Invalid room name format for room: %v", room)
			continue
		}
		rooms[roomID] = roomName
		logger.Debugf("Room added: ID=%s, Name=%s", roomID, roomName)
	}
	logger.Info("Room names fetched successfully")
}

func updateMetrics() {
	logger.Debug("Updating metrics")

	// Construct URL for services endpoint
	servicesURL := fmt.Sprintf("https://%s:%s/smarthome/services", bshcHost, bshcPort)

	// Make GET request to services endpoint
	resp, err := makeGetRequest(servicesURL, bshcClientCert, bshcClientKey)
	if err != nil {
		logger.Errorf("Failed to get services: %v", err)
		return
	}
	defer resp.Body.Close()

	// Check if response status code is 200
	if resp.StatusCode != http.StatusOK {
		logger.Errorf("Unexpected status code: %d", resp.StatusCode)
		return
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("Failed to read response body: %v", err)
		return
	}

	// Parse JSON response
	var services []map[string]interface{}
	err = json.Unmarshal(body, &services)
	if err != nil {
		logger.Errorf("Failed to unmarshal services response: %v", err)
		return
	}

	// Filter services with ID "TemperatureLevel"
	if c.SERVICES.TemperatureLevel {
		logger.Debug("Processing TemperatureLevel services")
		var temperatureLevelServices []map[string]interface{}
		for _, service := range services {
			if service["id"] == "TemperatureLevel" {
				state, ok := service["state"].(map[string]interface{})
				if !ok {
					logger.Errorf("Invalid state format for service: %v", service)
					continue
				}
				filteredService := map[string]interface{}{
					"id":          service["id"],
					"deviceId":    service["deviceId"],
					"temperature": state["temperature"],
				}
				temperatureLevelServices = append(temperatureLevelServices, filteredService)
			}
		}

		// Update Prometheus metrics with temperature levels
		for _, service := range temperatureLevelServices {
			deviceID, ok := service["deviceId"].(string)
			if !ok {
				logger.Errorf("Invalid deviceId format for service: %v", service)
				continue
			}
			temperature, ok := service["temperature"].(float64)
			if !ok {
				logger.Errorf("Invalid temperature format for service: %v", service)
				continue
			}
			deviceName, ok := devices[deviceID].(map[string]string)["name"]
			if !ok {
				logger.Errorf("Invalid device name format for device: %v", deviceID)
				continue
			}
			roomID, ok := devices[deviceID].(map[string]string)["roomId"]
			if !ok {
				logger.Errorf("Invalid room ID format for device: %v", deviceID)
				continue
			}
			roomName, ok := rooms[roomID].(string)
			if !ok {
				logger.Errorf("Invalid room name format for room: %v", roomID)
				continue
			}

			logger.Debugf("Updating temperature metric for device %s in room %s", deviceName, roomName)
			temperatureGauge.WithLabelValues(deviceID, deviceName, roomName).Set(temperature)
		}

		// Filter services with ID "RoomClimateControl"
		var setpointTemperatureLevelServices []map[string]interface{}
		for _, service := range services {
			if service["id"] == "RoomClimateControl" {
				state, ok := service["state"].(map[string]interface{})
				if !ok {
					logger.Errorf("Invalid state format for service: %v", service)
					continue
				}
				filteredService := map[string]interface{}{
					"id":                  service["id"],
					"deviceId":            service["deviceId"],
					"setpointTemperature": state["setpointTemperature"],
				}
				setpointTemperatureLevelServices = append(setpointTemperatureLevelServices, filteredService)
			}
		}

		// Update Prometheus metrics with desired temperature levels
		for _, service := range setpointTemperatureLevelServices {
			deviceID, ok := service["deviceId"].(string)
			if !ok {
				logger.Errorf("Invalid deviceId format for service: %v", service)
				continue
			}
			setpointTemperature, ok := service["setpointTemperature"].(float64)
			if !ok {
				logger.Errorf("Invalid setpointTemperature format for service: %v", service)
				continue
			}
			deviceName, ok := devices[deviceID].(map[string]string)["name"]
			if !ok {
				logger.Errorf("Invalid device name format for device: %v", deviceID)
				continue
			}
			roomID, ok := devices[deviceID].(map[string]string)["roomId"]
			if !ok {
				logger.Errorf("Invalid room ID format for device: %v", deviceID)
				continue
			}
			roomName, ok := rooms[roomID].(string)
			if !ok {
				logger.Errorf("Invalid room name format for room: %v", roomID)
				continue
			}

			logger.Debugf("Updating temperature metric for device %s in room %s", deviceName, roomName)
			setpointTemperatureGauge.WithLabelValues(deviceID, deviceName, roomName).Set(setpointTemperature)
		}
	}

	// Filter services with ID "HumidityLevel"
	if c.SERVICES.HumidityLevel {
		logger.Debug("Processing HumidityLevel services")
		var humidityLevelServices []map[string]interface{}
		for _, service := range services {
			if service["id"] == "HumidityLevel" {
				state, ok := service["state"].(map[string]interface{})
				if !ok {
					logger.Errorf("Invalid state format for service: %v", service)
					continue
				}
				filteredService := map[string]interface{}{
					"id":       service["id"],
					"deviceId": service["deviceId"],
					"humidity": state["humidity"],
				}
				humidityLevelServices = append(humidityLevelServices, filteredService)
			}
		}

		// Update Prometheus metrics with humidity levels
		for _, service := range humidityLevelServices {
			deviceID, ok := service["deviceId"].(string)
			if !ok {
				logger.Errorf("Invalid deviceId format for service: %v", service)
				continue
			}
			humidity, ok := service["humidity"].(float64)
			if !ok {
				logger.Errorf("Invalid humidity format for service: %v", service)
				continue
			}

			deviceName, ok := devices[deviceID].(map[string]string)["name"]
			if !ok {
				logger.Errorf("Invalid device name format for device: %v", deviceID)
				continue
			}
			roomID, ok := devices[deviceID].(map[string]string)["roomId"]
			if !ok {
				logger.Errorf("Invalid room ID format for device: %v", deviceID)
				continue
			}
			roomName, ok := rooms[roomID].(string)
			if !ok {
				logger.Errorf("Invalid room name format for room: %v", roomID)
				continue
			}

			logger.Debugf("Updating humidity metric for device %s in room %s", deviceName, roomName)
			humidityGauge.WithLabelValues(deviceID, deviceName, roomName).Set(humidity)
		}
	}

	// Filter services with ID "ValveTappet"
	if c.SERVICES.ValveTappet {
		logger.Debug("Processing ValveTappet services")
		var valveTappetServices []map[string]interface{}
		for _, service := range services {
			if service["id"] == "ValveTappet" {
				state, ok := service["state"].(map[string]interface{})
				if !ok {
					logger.Errorf("Invalid state format for service: %v", service)
					continue
				}
				filteredService := map[string]interface{}{
					"id":       service["id"],
					"deviceId": service["deviceId"],
					"position": state["position"],
				}
				valveTappetServices = append(valveTappetServices, filteredService)
			}
		}

		// Update Prometheus metrics with valve tappet
		for _, service := range valveTappetServices {
			deviceID, ok := service["deviceId"].(string)
			if !ok {
				logger.Errorf("Invalid deviceId format for service: %v", service)
				continue
			}
			valve, ok := service["position"].(float64)
			if !ok {
				logger.Errorf("Invalid position format for service: %v", service)
				continue
			}

			deviceName, ok := devices[deviceID].(map[string]string)["name"]
			if !ok {
				logger.Errorf("Invalid device name format for device: %v", deviceID)
				continue
			}
			roomID, ok := devices[deviceID].(map[string]string)["roomId"]
			if !ok {
				logger.Errorf("Invalid room ID format for device: %v", deviceID)
				continue
			}
			roomName, ok := rooms[roomID].(string)
			if !ok {
				logger.Errorf("Invalid room name format for room: %v", roomID)
				continue
			}

			logger.Debugf("Updating valve tappet metric for device %s in room %s", deviceName, roomName)
			valveTappetGauge.WithLabelValues(deviceID, deviceName, roomName).Set(valve)
		}
	}

	logger.Debug("Metrics updated successfully")
}

func main() {
	// Splash art
	ascii := figlet4go.NewAsciiRender()
	options := figlet4go.NewRenderOptions()
	options.FontColor = []figlet4go.Color{figlet4go.ColorGreen}

	renderStr, _ := ascii.RenderOpts("PlasticGhoul", options)
	fmt.Print(renderStr)
	fmt.Println("           BSHC Prometheus Exporter")
	fmt.Println()

	// Parse flags
	flag.StringVar(&configPath, "c", configPathDefault, "Path to the config file")
	flag.StringVar(&configPath, "config", configPathDefault, "Path to the config file")
	flag.StringVar(&httpBind, "b", httpBindDefault, "Host to bind the HTTP server to")
	flag.StringVar(&httpBind, "bind", httpBindDefault, "Host to bind the HTTP server to")
	flag.StringVar(&httpPort, "p", httpPortDefault, "Port to bind the HTTP server to")
	flag.StringVar(&httpPort, "port", httpPortDefault, "Port to bind the HTTP server to")
	flag.StringVar(&bshcHost, "bh", bshcHostDefault, "BSHC host")
	flag.StringVar(&bshcHost, "bshchost", bshcHostDefault, "BSHC host")
	flag.StringVar(&bshcPort, "bp", bshcPortDefault, "BSHC port")
	flag.StringVar(&bshcPort, "bshcport", bshcPortDefault, "BSHC port")
	flag.StringVar(&bshcClientCert, "cc", bshcClientCertDefault, "BSHC client cert")
	flag.StringVar(&bshcClientCert, "clientcert", bshcClientCertDefault, "BSHC client cert")
	flag.StringVar(&bshcClientKey, "ck", bshcClientKeyDefault, "BSHC client key")
	flag.StringVar(&bshcClientKey, "clientkey", bshcClientKeyDefault, "BSHC client key")
	flag.BoolVar(&skipTlsVerify, "insecure", false, "Skip TLS verification")
	flag.BoolVar(&skipTlsVerify, "i", false, "Skip TLS verification")
	flag.BoolVar(&debug, "d", debug, "Enable debug mode")
	flag.BoolVar(&debug, "debug", debug, "Enable debug mode")
	flag.Parse()

	// Initialize logger
	logger = initLogger(debug)

	// Start application
	logger.Info("Starting BSHC Prometheus Exporter")

	// Check if config file exists
	if _, err := os.Stat(configPath); err == nil {
		// Load config file
		c.getConf(configPath)

		if (flag.Lookup("bind").Value.String() == httpBindDefault || flag.Lookup("b").Value.String() == httpBindDefault) && c.HTTP.Bind != httpBindDefault {
			httpBind = c.HTTP.Bind
		}
		if flag.Lookup("port").Value.String() == httpPortDefault || flag.Lookup("p").Value.String() == httpPortDefault && c.HTTP.Port != httpBindDefault {
			httpPort = c.HTTP.Port
		}
		if flag.Lookup("bshchost").Value.String() == bshcHostDefault || flag.Lookup("bh").Value.String() == bshcHostDefault && c.BSHC.Host != bshcHostDefault {
			bshcHost = c.BSHC.Host
		}
		if flag.Lookup("bshcport").Value.String() == bshcPort || flag.Lookup("bp").Value.String() == bshcPort && c.BSHC.Port != bshcPort {
			bshcPort = c.BSHC.Port
		}
		if flag.Lookup("clientcert").Value.String() == bshcClientCertDefault || flag.Lookup("cc").Value.String() == bshcClientCertDefault && c.BSHC.ClientCert != bshcClientCertDefault {
			bshcClientCert = c.BSHC.ClientCert
		}
		if flag.Lookup("clientkey").Value.String() == bshcClientKeyDefault || flag.Lookup("ck").Value.String() == bshcClientKeyDefault && c.BSHC.ClientKey != bshcClientKeyDefault {
			bshcClientKey = c.BSHC.ClientKey
		}
		if flag.Lookup("insecure").Value.String() == fmt.Sprint(skipTlsVerifyDefault) || flag.Lookup("i").Value.String() == fmt.Sprint(skipTlsVerifyDefault) && c.BSHC.SkipTLSVerify != skipTlsVerifyDefault {
			skipTlsVerify = c.BSHC.SkipTLSVerify
		}
	}

	// Check if required config values are set
	if httpBind == "" || httpPort == "" || bshcHost == "" || bshcPort == "" || bshcClientCert == "" || bshcClientKey == "" {
		logger.Fatal("Missing required config values")
	}

	// DEBUG: Print config values
	logger.Debug("Config Path: " + configPath)
	logger.Debug("HTTP Bind: " + httpBind)
	logger.Debug("HTTP Port: " + fmt.Sprint(httpPort))
	logger.Debug("BSHC Host: " + bshcHost)
	logger.Debug("BSHC Port: " + fmt.Sprint(bshcPort))
	logger.Debug("BSHC Client Cert: " + bshcClientCert)
	logger.Debug("BSHC Client Key: " + bshcClientKey)
	logger.Debug("Temperature Level: " + fmt.Sprint(c.SERVICES.TemperatureLevel))
	logger.Debug("Humidity Level: " + fmt.Sprint(c.SERVICES.HumidityLevel))
	logger.Debug("Valve Tappet: " + fmt.Sprint(c.SERVICES.ValveTappet))

	// Get device names
	getDeviceNames()
	getRoomNames()

	// Define Prometheus metrics
	temperatureGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "temperature_level",
			Help: "Temperature level of the devices",
		},
		[]string{"device_id", "device_name", "room_name"},
	)

	setpointTemperatureGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "setpoint_temperature_level",
			Help: "Desired temperature level of the devices",
		},
		[]string{"device_id", "device_name", "room_name"},
	)

	humidityGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "humidity_level",
			Help: "Humidity level of the devices",
		},
		[]string{"device_id", "device_name", "room_name"},
	)

	valveTappetGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "valve_tappet",
			Help: "Valve tappet of the devices",
		},
		[]string{"device_id", "device_name", "room_name"},
	)

	// Register Prometheus metrics
	logger.Info("Registering Prometheus metrics")
	prometheus.MustRegister(temperatureGauge)
	prometheus.MustRegister(setpointTemperatureGauge)
	prometheus.MustRegister(humidityGauge)
	prometheus.MustRegister(valveTappetGauge)

	// HTTP handler for Prometheus metrics
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		logger.Debug("Handling /metrics request")
		updateMetrics()
		promhttp.Handler().ServeHTTP(w, r)
	})

	// Start HTTP server for Prometheus metrics
	logger.Infof("Starting HTTP server on %s:%s", httpBind, httpPort)
	if err := http.ListenAndServe(fmt.Sprintf("%s:%s", httpBind, httpPort), nil); err != nil {
		logger.Fatalf("Failed to start HTTP server: %v", err)
	}
}
