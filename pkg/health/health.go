package health

import (
	"net/http"
	"os"
	"time"

	"starbase.ag/liquidity/pkg/logger"
)

func Health() error {
	port := os.Getenv("SRV_PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/v1/health", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})

	server := &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Info().Msgf("health start at port: %v", port)

	return server.ListenAndServe()
}
