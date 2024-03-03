package main

import "os"

func init() {
	// you can read from .env file, this is just an example
	_ = os.Setenv("BASE_URL", "https://dummyjson.com")
	_ = os.Setenv("USERNAME", "kminchelle")
	_ = os.Setenv("PASSWORD", "0lelplR")
}
