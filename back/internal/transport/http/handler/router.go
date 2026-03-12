package handler

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"

	transport "github.com/vkotsiuba99/NeoHome/back/internal/transport/http"
	"github.com/vkotsiuba99/NeoHome/back/internal/transport/http/middleware"
)

func NewRouter(hand *Handler, cfg transport.Config, log *slog.Logger, jwtSecret []byte) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	router.Use(middleware.Recover(cfg, log))
	router.Use(middleware.AccessLog(cfg, log))
	router.Use(middleware.CORS(cfg))

	router.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodGet, http.MethodPost)

	api := router.PathPrefix("/api/v1").Subrouter()

	auth := api.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/register", hand.CreateUser).Methods(http.MethodPost)
	auth.HandleFunc("/login/email", hand.LoginEmail).Methods(http.MethodPost)

	telemetry := api.PathPrefix("/telemetry").Subrouter()
	telemetry.HandleFunc("", hand.Telemetry).Methods(http.MethodPost)
	telemetry.HandleFunc("/mqtt", hand.TelemetryMQTT).Methods(http.MethodPost)

	private := api.NewRoute().Subrouter()
	private.Use(middleware.Auth(jwtSecret, log))

	users := private.PathPrefix("/users").Subrouter()
	users.HandleFunc("/me", hand.User).Methods(http.MethodGet)
	users.HandleFunc("/me", hand.UserUpdate).Methods(http.MethodPut)

	devices := private.PathPrefix("/devices").Subrouter()
	devices.HandleFunc("", hand.Devices).Methods(http.MethodGet)
	devices.HandleFunc("", hand.DeviceCreate).Methods(http.MethodPost)
	devices.HandleFunc("/{deviceId}", hand.DeviceUpdate).Methods(http.MethodPatch)
	devices.HandleFunc("/{deviceId}/thresholds", hand.DeviceThresholdsList).Methods(http.MethodGet)
	devices.HandleFunc("/{deviceId}/thresholds", hand.DeviceThresholdsUpdate).Methods(http.MethodPut)
	devices.HandleFunc("/{deviceId}/telemetry", hand.TelemetryDevice).Methods(http.MethodGet)
	devices.HandleFunc("/{deviceId}/latest", hand.DeviceLatest).Methods(http.MethodGet)

	alerts := private.PathPrefix("/alerts").Subrouter()
	alerts.HandleFunc("", hand.Alerts).Methods(http.MethodGet)
	alerts.HandleFunc("/{alertId}/resolve", hand.AlertResolve).Methods(http.MethodPut)

	return router
}
