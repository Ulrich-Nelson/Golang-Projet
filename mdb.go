package main

import (
	mdb "./mdb"
	_ "github.com/lib/pq"
)

func main() {
	mdb.ExecuteCommand()
}
