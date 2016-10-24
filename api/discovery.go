package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/markdaws/gohome"
)

// RegisterDiscoveryHandlers registers all of the discovery specific API REST routes
func RegisterDiscoveryHandlers(r *mux.Router, s *apiServer) {
	r.HandleFunc("/api/v1/discovery/{modelNumber}",
		apiDiscoveryHandler(s.system)).Methods("GET")

	//TODO: Implement with extensions
	/*
		r.HandleFunc("/api/v1/discovery/{modelNumber}/token",
			apiDiscoveryTokenHandler(s.system)).Methods("GET")
		r.HandleFunc("/api/v1/discovery/{modelNumber}/access",
			apiDiscoveryAccessHandler(s.system)).Methods("GET")
	*/
}

func apiDiscoveryHandler(system *gohome.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		//TODO: This shouldn't be model number, should be some unique importer generated ID
		modelNumber := vars["modelNumber"]

		//TODO: fix, This is blocking
		network := system.Extensions.FindNetwork(system, &gohome.Device{ModelNumber: modelNumber})
		if network == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		devs, err := network.Devices(system, modelNumber)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		jsonDevices := make(devices, len(devs))
		for i, device := range devs {

			//TODO: This API shouldn't be sending back client ids, the client should
			//make these values up
			jsonZones := make(zones, len(device.Zones))
			var j int
			for _, zn := range device.Zones {
				jsonZones[j] = jsonZone{
					Address:     zn.Address,
					ID:          zn.ID,
					Name:        zn.Name,
					Description: zn.Description,
					DeviceID:    device.ID,
					Type:        zn.Type.ToString(),
					Output:      zn.Output.ToString(),
					ClientID:    modelNumber + "_zone_" + strconv.Itoa(j),
				}
				j++
			}

			jsonSensors := make(sensors, len(device.Sensors))
			j = 0
			for _, sen := range device.Sensors {
				jsonSensors[j] = jsonSensor{
					ID:          sen.ID,
					Name:        sen.Name,
					Description: sen.Description,
					Address:     sen.Address,

					//TODO: Shouldn't be setting ClientID here
					ClientID: modelNumber + "_sensor_" + strconv.Itoa(j),

					//TODO: Attrs
				}
				j++
			}

			var connPoolJSON *jsonConnPool
			if device.Connections != nil {
				connCfg := device.Connections.Config
				connPoolJSON = &jsonConnPool{
					Name:     connCfg.Name,
					PoolSize: int32(connCfg.Size),
				}
			}
			jsonDevices[i] = jsonDevice{
				ID:          device.ID,
				Address:     device.Address,
				Name:        device.Name,
				Description: device.Description,
				ModelNumber: device.ModelNumber,
				Token:       "",
				ClientID:    modelNumber + "_" + strconv.Itoa(i),
				Zones:       jsonZones,
				Sensors:     jsonSensors,
				ConnPool:    connPoolJSON,
			}
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(jsonDevices); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		/*TODO: Remove
		json.NewEncoder(w).Encode(struct {
			Location string `json:"location"`
		}{Location: data["location"]})
		*/
	}
}

/*
func apiDiscoveryZoneHandler(system *gohome.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)

		//This is blocking, waits 5 seconds
		zs, err := discovery.Zones(vars["modelNumber"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		jsonZones := make(zones, len(zs))
		for i, zone := range zs {
			jsonZones[i] = jsonZone{
				Address:     zone.Address,
				Name:        zone.Name,
				Description: zone.Description,
				Type:        zone.Type.ToString(),
				Output:      zone.Output.ToString(),
			}
		}
		sort.Sort(jsonZones)
		if err := json.NewEncoder(w).Encode(jsonZones); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
*/

/*
func apiDiscoveryTokenHandler(system *gohome.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)

		//TODO: Make non-blocking: this is blocking
		token, err := discovery.DiscoverToken(vars["modelNumber"], r.URL.Query().Get("address"))
		if err != nil {
			if err == discovery.ErrUnauthorized {
				// Let the caller know this was a specific kind of error
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(struct {
					Unauthorized bool `json:"unauthorized"`
				}{Unauthorized: true})
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}

			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(struct {
			Token string `json:"token"`
		}{Token: token})
	}
}

func apiDiscoveryAccessHandler(system *gohome.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)

		//TODO: Make non-blocking: this is blocking
		err := discovery.VerifyConnection(
			vars["modelNumber"],
			r.URL.Query().Get("address"),
			r.URL.Query().Get("token"),
		)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(struct{}{})
	}
}
*/
