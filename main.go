package main

import (
	"github.com/seb-schulz/onegate/cmd"
	_ "github.com/seb-schulz/onegate/cmd/client"
	_ "github.com/seb-schulz/onegate/cmd/generate"
	_ "github.com/seb-schulz/onegate/cmd/session"
	_ "github.com/seb-schulz/onegate/cmd/user"
)

func main() {
	cmd.Execute()
}
