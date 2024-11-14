package main

import (
	"fmt"
	"os"

	"silae2calendar/ms"
	"silae2calendar/silae"
)

func main() {
	userData, err := silae.GetUserData(os.Getenv("SILAE_USERNAME"), os.Getenv("SILAE_PASSWORD"))
	if err != nil {
		panic(err)
	}

	freedays, err := silae.GetFreedays(userData)
	if err != nil {
		panic(err)
	}

	for _, cf := range freedays.CollaboratorFreedays {
		for _, f := range cf.Freedays {
			fmt.Println(f)
		}
	}

	accessToken, err := ms.GetAccessToken()
	if err != nil {
		panic(err)
	}

	err = ms.CreateOutlookEvent(accessToken)
	if err != nil {
		panic(err)
	}
}
