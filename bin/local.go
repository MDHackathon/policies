package main

import (
	"flag"
	"github.com/MDHackathon/policies/pkg"
	"io/ioutil"
)

// This will load the policies from a json configuration
// and apply a choosen operation on all object who match
// the policies rules.
//
//	S3 storage:
//		- In  -> the storage source
//		- Out -> the storage destination
//
// Rules:
// 		- match(name) -> allow a reggex on the key
//		- match(date) -> allow to operate an a set of object define by date
//						 Format "2013-06-18T02:08:03+00:00"
//
// Operation:
//		- Copy -> Copy an object from a s3 storage to an other ( actif <-> passif )
//		- Move -> Moce an object from a s3 storage to an other (  hot  <->  cold  )
//

var (
	path = flag.String("path", "./config/policies.json", "path of the json configuration")
)

func main() {

	var (
		data []byte
		err  error
	)

	if data, err = ioutil.ReadFile(*path); err != nil {
		panic(err)
	}
	policies.Execute(data)
}
