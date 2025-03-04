// Copyright  The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mockserver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/vcenterreceiver/internal/mockserver"

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	xj "github.com/basgys/goxml2json"
	"github.com/stretchr/testify/require"
)

const (
	// MockUsername is the correct user for authentication to the Mock Server
	MockUsername = "otelu"
	// MockPassword is the correct password for authentication to the Mock Server
	MockPassword = "otelp"
)

var errNotFound = errors.New("not found")

type soapRequest struct {
	Envelope soapEnvelope `json:"Envelope"`
}

type soapEnvelope struct {
	Body map[string]interface{} `json:"Body"`
}

// MockServer has access to recorded SOAP responses and will serve them over http based off the scraper's API calls
func MockServer(t *testing.T) *httptest.Server {
	vsphereMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// converting to JSON in order to iterate over map keys
		jsonified, err := xj.Convert(r.Body)
		require.NoError(t, err)
		sr := &soapRequest{}
		err = json.Unmarshal(jsonified.Bytes(), sr)
		require.NoError(t, err)
		require.Len(t, sr.Envelope.Body, 1)

		var requestType string
		for k := range sr.Envelope.Body {
			requestType = k
		}
		require.NotEmpty(t, requestType)

		body, err := routeBody(t, requestType, sr.Envelope.Body)
		if errors.Is(err, errNotFound) {
			w.WriteHeader(404)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/xml")
		_, _ = w.Write(body)
	}))
	return vsphereMock
}

func routeBody(t *testing.T, requestType string, body map[string]interface{}) ([]byte, error) {
	switch requestType {
	case "RetrieveServiceContent":
		return loadResponse("service-content.xml")
	case "Login":
		return loadResponse("login.xml")
	case "Logout":
		return loadResponse("logout.xml")
	case "RetrieveProperties":
		return routeRetreiveProperties(t, body)
	case "QueryPerf":
		return routePerformanceQuery(t, body)
	}

	return []byte{}, errNotFound
}

func routeRetreiveProperties(t *testing.T, body map[string]interface{}) ([]byte, error) {
	rp, ok := body["RetrieveProperties"].(map[string]interface{})
	require.True(t, ok)
	specSet := rp["specSet"].(map[string]interface{})

	var objectSetArray = false
	objectSet, ok := specSet["objectSet"].(map[string]interface{})
	if !ok {
		objectSetArray = true
	}

	var propSetArray = false
	propSet, ok := specSet["propSet"].(map[string]interface{})
	if !ok {
		propSetArray = true
	}

	var obj map[string]interface{}
	var content string
	var contentType string
	if !objectSetArray {
		obj = objectSet["obj"].(map[string]interface{})
		content = obj["#content"].(string)
		contentType = obj["-type"].(string)
	}

	switch {
	case content == "group-d1" && contentType == "Folder":
		return loadResponse("datacenter.xml")

	case content == "datacenter-3" && contentType == "Datacenter":
		return loadResponse("datacenter-properties.xml")

	case content == "domain-c8" && contentType == "ClusterComputeResource":
		if propSetArray {
			pSet := specSet["propSet"].([]interface{})
			for _, prop := range pSet {
				spec := prop.(map[string]interface{})
				specType := spec["type"].(string)
				if specType == "ResourcePool" {
					return loadResponse("resource-pool.xml")
				}
			}
		}
		path := propSet["pathSet"].(string)
		switch path {
		case "datastore":
			return loadResponse("cluster-datastore.xml")
		case "summary":
			return loadResponse("cluster-summary.xml")
		case "host":
			return loadResponse("host-list.xml")
		}

	case content == "PerfMgr" && contentType == "PerformanceManager":
		return loadResponse("perf-manager.xml")

	case content == "group-h5" && contentType == "Folder":
		if propSetArray {
			arr := specSet["propSet"].([]interface{})
			for _, i := range arr {
				m, ok := i.(map[string]interface{})
				require.True(t, ok)
				if m["type"] == "ClusterComputeResource" {
					return loadResponse("host-cluster.xml")
				}
			}
		}
		return loadResponse("host-parent.xml")

	case content == "datastore-1003" && contentType == "Datastore":
		if objectSetArray {
			return loadResponse("datastore-list.xml")
		}
		return loadResponse("datastore-summary.xml")

	case contentType == "HostSystem":
		if ps, ok := propSet["pathSet"].([]interface{}); ok {
			for _, v := range ps {
				if v == "summary.hardware" {
					return loadResponse("host-properties.xml")
				}
			}
		} else {
			ps, ok := propSet["pathSet"].(string)
			require.True(t, ok)
			if ps == "name" {
				return loadResponse("host-names.xml")
			}

		}

	case content == "group-v4" && contentType == "Folder":
		if propSetArray {
			return loadResponse("vm-group.xml")
		}
		if propSet == nil {
			return loadResponse("vm-folder.xml")
		}
		return loadResponse("vm-folder-parent.xml")

	case content == "vm-1040" && contentType == "VirtualMachine":
		if propSet["pathSet"] == "summary.runtime.host" {
			return loadResponse("vm-host.xml")
		}
		return loadResponse("vm-properties.xml")

	case (content == "group-v1034" || content == "group-v1001") && contentType == "Folder":
		return loadResponse("vm-empty-folder.xml")

	case contentType == "ResourcePool":
		if ps, ok := propSet["pathSet"].([]interface{}); ok {
			for _, prop := range ps {
				if prop == "summary" {
					return loadResponse("resource-pool-summary.xml")
				}
			}
		}

		if ss, ok := objectSet["selectSet"].(map[string]interface{}); ok && ss["path"] == "resourcePool" {
			return loadResponse("resource-pool-group.xml")
		}

	case objectSetArray:
		objectArray := specSet["objectSet"].([]interface{})
		for _, i := range objectArray {
			m, ok := i.(map[string]interface{})
			require.True(t, ok)
			mObj := m["obj"].(map[string](interface{}))
			typeString := mObj["-type"]
			if typeString == "HostSystem" {
				return loadResponse("host-names.xml")
			}
		}
	}

	return []byte{}, errNotFound
}

func routePerformanceQuery(t *testing.T, body map[string]interface{}) ([]byte, error) {
	queryPerf := body["QueryPerf"].(map[string]interface{})
	require.NotNil(t, queryPerf)
	querySpec := queryPerf["querySpec"].(map[string]interface{})
	entity := querySpec["entity"].(map[string]interface{})
	switch entity["-type"] {
	case "HostSystem":
		return loadResponse("host-performance-counters.xml")
	case "VirtualMachine":
		return loadResponse("vm-performance-counters.xml")
	}
	return []byte{}, errNotFound
}

func loadResponse(filename string) ([]byte, error) {
	return os.ReadFile(filepath.Join("internal", "mockserver", "responses", filename))
}
