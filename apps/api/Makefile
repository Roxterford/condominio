check-unix:
	@case $$(uname -s) in \
		Linux*|Darwin*) ;; \
		*) echo "Este script solo se ejecuta en Linux o macOS"; exit 1;; \
	esac

check-golines: check-unix
	@command -v golines >/dev/null 2>&1 || { \
		echo "No se encuentra instalado golines. Instalación: https://github.com/segmentio/golines"; \
		exit 1; \
	}

cov:
	go test ./... -v -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

