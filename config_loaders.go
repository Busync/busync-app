package busylight_sync

import "fmt"

func LoadConfigInGivenFormat(fileFormat string) error {
	switch fileFormat {
	default:
		return fmt.Errorf("%s is not implemented", fileFormat)
	}
}
