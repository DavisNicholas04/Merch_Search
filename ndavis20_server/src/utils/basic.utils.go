// this package is to be removed upon making the https://github.com/DavisNicholas04/Merch_Search public
// in its place both: github.com/DavisNicholas04/Merch_Search/blob/main/polling_agent/src/utils/basic.utils.go
// will be imported

package utils

import (
	"errors"
	"github.com/jamespearly/loggly"
	"github.com/joho/godotenv"
	"log"
	"net"
	"os"
)

func InstantiateClient(tag string) *loggly.ClientType {
	return loggly.New(tag)
}

func FileExist(file string) bool {
	_, err := os.Stat(file)
	if errors.Is(err, os.ErrNotExist) {
		return false
	} else {
		return true
	}
}

func LoadDotEnv(filenames ...string) {
	for _, name := range filenames {
		if FileExist(name) {
			err1 := godotenv.Load(filenames...)
			if err1 != nil {
				log.Fatalln("Error loading .env file")
			}
		}
	}
}

// Get preferred outbound ip of this machine
func GetOutboundIP() net.IP {
	connection, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer func(conn *net.Conn) {
		err := (*conn).Close()
		if err != nil {

		}
	}(&connection)

	localAddr := connection.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
