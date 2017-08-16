package daemon

import (
	"fmt"
	"net/http"

	"github.com/LeReverandNox/GuessWhat/src/routing"
)

// Config is a structure holding the daemon's config
type Config struct {
	Host string
	Port string
}

// Run is the method that launch the server
func Run(cfg *Config) error {
	router := routing.NewRouter()

	fmt.Printf("The Guess What app is now running on %v:%v", cfg.Host, cfg.Port)

	err := http.ListenAndServe(cfg.Host+":"+cfg.Port, router)

	if err != nil {
		return err
	}

	return nil
}
