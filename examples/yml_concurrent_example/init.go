package main

import "os"

func init() {
	// you can read from .env file, this is just an example
	_ = os.Setenv("BASE_URL", "https://dummyjson.com")
}
