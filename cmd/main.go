package main

import (
	"os"

	"github.com/curtisnewbie/goauth"
)

func main() {
	goauth.BootstrapServer(os.Args)
}
