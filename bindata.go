package main

import "log"

//MustAsset - Dummy MustAsset function, shuts up the linter if go-bindata hasn't been run yet. This file should be overwritten when make is run.
func MustAsset(interface{}) []byte {
	log.Fatal("No assets have been packed, please run go-bindata before building.")
	return []byte{}
}
