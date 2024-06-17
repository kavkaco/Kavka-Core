package e2e

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/kavkaco/Kavka-Core/config"
	"github.com/kavkaco/Kavka-Core/utils/net"
)

var BaseUrl string

func TestMain(m *testing.M) {
	configs := config.Read()
	baseAddr := fmt.Sprintf("%s:%d", configs.HTTP.Host, configs.HTTP.Port)

	if net.IsHostReachable(baseAddr) {
		BaseUrl = fmt.Sprintf("http://%s", baseAddr)

		os.Exit(m.Run())
		return
	}

	log.Fatal("Connection refused!")
}
