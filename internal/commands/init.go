package commands

import (
	"fmt"

	"github.com/davidmks/got/internal/repository"
)

func Init() error {
	if err := repository.Initialize(); err != nil {
		return err
	}

	fmt.Println("Initialized empty got repository")
	return nil
}
