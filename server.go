package killgrave

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gorilla/mux"
)

// Server definition of mock server
type Server struct {
	impostersPath string
	router        *mux.Router
}

// NewServer initialize the mock server
func NewServer(p string, r *mux.Router) *Server {
	return &Server{
		impostersPath: p,
		router:        r,
	}
}

// Run read all the files on the impostersPath and creates different
// handlers for each imposter
func (s *Server) Run() error {
	if _, err := os.Stat(s.impostersPath); os.IsNotExist(err) {
		return invalidDirectoryError(fmt.Sprintf("the directory %s doesn't exists", s.impostersPath))
	}
	if err := s.buildImposters(); err != nil {
		return err
	}

	return nil
}

func (s *Server) buildImposters() error {
	files, _ := ioutil.ReadDir(s.impostersPath)

	for _, f := range files {
		var imposter Imposter
		if err := s.buildImposter(f.Name(), &imposter); err != nil {
			return err
		}

		if imposter.Request.Endpoint == "" {
			continue
		}
		s.router.HandleFunc(imposter.Request.Endpoint, ImposterHandler(imposter)).Methods(imposter.Request.Method)
	}

	return nil
}

func (s *Server) buildImposter(imposterFileName string, imposter *Imposter) error {
	f := s.impostersPath + "/" + imposterFileName
	imposterFile, _ := os.Open(f)
	defer imposterFile.Close()

	bytes, _ := ioutil.ReadAll(imposterFile)
	if err := json.Unmarshal(bytes, imposter); err != nil {
		return malformattedImposterError(fmt.Sprintf("error while unmarshall imposter file %s", f))
	}
	return nil
}
