package main

import (
	"fmt"
)

func main() {
	fmt.Println("Notification service is successfully running!")

	// Пустой select блокирует горутину навечно,
	// имитируя работу демона/консьюмера, чтобы контейнер не завершил работу.
	select {}
}
