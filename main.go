package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/viper"

	"silae2calendar/ms"
	"silae2calendar/silae"
)

func main() {
	if err := loadConfig(); err != nil {
		panic(err)
	}

	// Connect Silae
	ud, err := silae.GetUserData(viper.GetString("silae_username"), viper.GetString("silae_password"))
	if err != nil {
		panic(err)
	}

	// Connect Microsoft
	accessToken, err := ms.GetAccessToken()
	if err != nil {
		panic(err)
	}
	fmt.Println(accessToken)

	freedays, err := silae.GetFreedays(ud)
	if err != nil {
		panic(err)
	}

	for _, cf := range freedays.CollaboratorFreedays {
		for _, f := range cf.Freedays {
			fmt.Println(f)

			exists, err := ms.FindOutlookEvent(accessToken, f.Abbr+" "+ud.Trigram, f.DateStart, f.DateEnd)
			if err != nil {
				panic(err)
			}

			fmt.Println(exists)
		}
	}

	err = ms.CreateOutlookEvent(accessToken)
	if err != nil {
		panic(err)
	}
}

func loadConfig() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	viper.AutomaticEnv()
	viper.SetConfigName(".silae2calendar")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(home)

	err = viper.ReadInConfig()
	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			err := viper.SafeWriteConfig()
			if err != nil {
				return err
			}
		}
	}

	return nil
}
