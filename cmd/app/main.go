package main

import (
	"github.com/ra1nz0r/parcel_service/internal/testdata"

	_ "modernc.org/sqlite"
)

func main() {
	testdata.RunTestParcelServices()
}
