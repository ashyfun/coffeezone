package coffeezone

import (
	"log"
	"os"
)

func SetLogFileOutput(file string) (*os.File, error) {
	var (
		f   *os.File
		err error
	)
	if file != "" {
		f, err = os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return nil, err
		}

		log.SetOutput(f)
	}

	return f, nil
}
