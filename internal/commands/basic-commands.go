package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/walter2310/basic-cache-db/internal/cache"
)

var reader *bufio.Reader = bufio.NewReader(os.Stdin)

func GetInput() (string, error) {
	fmt.Print(">>> ")

	str, err := reader.ReadString('\n')
	if err != nil {
		return "", nil
	}

	str = strings.Replace(str, "\r\n", "", 1)

	return str, nil
}

func ExecuteCommands(command string, c *cache.Cache) {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return
	}

	switch strings.ToUpper(parts[0]) {
	case "SET":
		if len(parts) < 3 {
			fmt.Println("Usage: SET <key> <value>")
			return
		}

		ttl := time.Duration(time.Second * 120)

		if len(parts) == 4 {
			parsedTTL, err := time.ParseDuration(parts[3])
			if err != nil {
				fmt.Println("Invalid TTL format. Example: 10s, 5m, 1h")
				return
			}

			ttl = parsedTTL
		}

		c.Set(parts[1], parts[2], ttl)

	case "GET":
		if len(parts) < 2 {
			fmt.Println("Usage: GET <key>")
			return
		}

		value, exists := c.Get(parts[1])

		if exists {
			fmt.Println("Value:", value)
		} else {
			fmt.Println("Key not found or expired.")
		}

	case "KEYS":
		if len(parts) < 2 {
			fmt.Println("Usage: KEYS <pattern>")
			return
		}

		keys := c.ListKeys(parts[1])
		if len(keys) == 0 {
			fmt.Println("No matching keys found.")
		}

		fmt.Println("Matching keys:")
		for _, key := range keys {
			fmt.Println("-", key)
		}
	}
}
