package main

import (
	"errors"
	"os"
	"sync"

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

	freedays, err := silae.GetFreedays(ud)
	if err != nil {
		panic(err)
	}

	wg := sync.WaitGroup{}

	for _, cf := range freedays.CollaboratorFreedays {
		for _, f := range cf.Freedays {
			wg.Add(1)

			go func() {
				defer wg.Done()

				dateStart, err := f.DateStartForOutlook()
				if err != nil {
					panic(err)
				}

				dateEnd, err := f.DateEndForOutlook()
				if err != nil {
					panic(err)
				}

				subject := f.Abbr + " " + ud.Trigram
				exists, err := ms.FindOutlookEvent(accessToken, subject, dateStart, dateEnd)
				if err != nil {
					panic(err)
				}

				if !exists {
					err = ms.CreateOutlookEvent(accessToken, subject, dateStart, dateEnd, f.IsAllDay())
					if err != nil {
						panic(err)
					}
				}
			}()
		}
	}

	wg.Wait()
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
