package main

import (
	"fmt"
	"os"

	"silae2calendar/silae"
)

func main() {
	userData := silae.GetUserData(os.Getenv("SILAE_USERNAME"), os.Getenv("SILAE_PASSWORD"))
	freedays := silae.GetFreedays(userData)

	for _, cf := range freedays.CollaboratorFreedays {
		for _, f := range cf.Freedays {
			fmt.Println(f)
		}
	}
}
