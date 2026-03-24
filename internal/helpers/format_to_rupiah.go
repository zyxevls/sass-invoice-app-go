package helpers

import "fmt"

func FormatRupiah(amount int64) string {
	return fmt.Sprintf("Rp %d", amount)
}
