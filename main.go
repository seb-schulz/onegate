package main

import (
	"github.com/seb-schulz/onegate/cmd"
	_ "github.com/seb-schulz/onegate/cmd/config"
	_ "github.com/seb-schulz/onegate/cmd/user"
)

func main() {
	cmd.Execute()
}
