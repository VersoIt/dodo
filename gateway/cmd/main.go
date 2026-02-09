package main

import (
	"github.com/versoit/diploma/gateway/internal/app"
	"go.uber.org/fx"
)

func main() {
	fx.New(app.Module).Run()
}
