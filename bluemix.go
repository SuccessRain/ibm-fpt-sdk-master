package watson

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
)

type vcapService struct {
	Name        string          `json:"name"`
	Label       string          `json:"label,omitempty"`
	Plan        string          `json:"plan,omitempty"`
	Credentials vcapCredentials `json:"credentials"`
}

type vcapCredentials struct {
	Url      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
	ApiKey string `json:"apikey"`
}

// getBluemixCredentials parses the VCAP_SERVICES environment variable, and returns the
// credential information for the service with the given name. If a non-empty
// plan is also provided, it returns credential information for the specified
// plan.
func getBluemixCredentials(name, plan string) (Credentials, error) {
	vcap_services := os.Getenv("VCAP_SERVICES")
	if len(vcap_services) == 0 {
		return Credentials{}, errors.New("VCAP_SERVICES undefined")
	}
	var vcap map[string][]vcapService
	err := json.Unmarshal([]byte(vcap_services), &vcap)
	if err != nil {
		return Credentials{}, errors.New("failed to parse VCAP_SERVICES " + err.Error())
	}
	for vname, vservice := range vcap {
		if !strings.HasPrefix(vname, name) {
			continue
		}
		for i := range vservice {
			if len(plan) == 0 || plan == vservice[i].Plan {
				creds := Credentials{
					ServiceName: name,
					ServicePlan: vservice[i].Plan,
					Url: vservice[i].Credentials.Url,
					Username: vservice[i].Credentials.Username,
					Password: vservice[i].Credentials.Password,
					ApiKey: vservice[i].Credentials.ApiKey,
				}
				return creds, nil
			}
		}
	}
	return Credentials{}, errors.New("service instance not found in VCAP_SERVICES")
}
