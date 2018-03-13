package main

import (
	"strconv"
)

func main() {
	cfg := getConfig()
	a := App{}
	a.Initialize(
		cfg.DB.DBusername,
		cfg.DB.DBPassword,
		cfg.DB.DBname)
	a.Run(":" + strconv.Itoa(cfg.App.Port))
}
