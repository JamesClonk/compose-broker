package broker

import (
	"encoding/json"
	"net/http"

	"github.com/JamesClonk/compose-broker/api"
	"github.com/JamesClonk/compose-broker/config"
	"github.com/JamesClonk/compose-broker/log"
)

type Broker struct {
	Username       string
	Password       string
	API            config.API
	Client         *api.Client
	ServiceCatalog *ServiceCatalog
}

func NewBroker(c *config.Config) *Broker {
	b := &Broker{
		Username:       c.Username,
		Password:       c.Password,
		API:            c.API,
		Client:         api.NewClient(c),
		ServiceCatalog: LoadServiceCatalog(c.CatalogFilename),
	}
	return b
}

func (b *Broker) write(rw http.ResponseWriter, req *http.Request, code int, content interface{}) {
	log.InfoWithFields(log.Fields{
		"remote_addr":      req.RemoteAddr,
		"method":           req.Method,
		"request_uri":      req.URL.RequestURI(),
		"http_status_code": code,
		"user_agent":       req.UserAgent(),
	})

	data, err := json.MarshalIndent(content, "", "  ")
	if err != nil {
		log.Errorf("could not marshal content into json: %#v", content)
		log.Errorln(err)
	}

	rw.WriteHeader(code)
	rw.Header().Set("X-Compose-Broker", "compose-broker")
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(data)
}

func (b *Broker) Error(rw http.ResponseWriter, req *http.Request, code int, err, desc string) {
	if code < 100 {
		code = 500
	}
	if len(err) == 0 {
		err = "UnknownError"
	}
	if len(desc) == 0 {
		err = "An unknown error has occured"
	}
	b.write(rw, req, code, map[string]string{"error": err, "description": desc})
}

func (b *Broker) Health(rw http.ResponseWriter, req *http.Request) {
	b.write(rw, req, 200, map[string]string{"status": "ok"})
}
