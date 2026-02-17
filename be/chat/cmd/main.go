package main

import (
	"github.com/versoit/diploma/be/chat/internal/app"
	"go.uber.org/fx"
)

func main() {
	fx.New(app.Module).Run()
}
