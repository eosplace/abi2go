# abi2go

# usage

The function getAPI must be provided. One way is similar to:

```
package eosgame1

import (
	"crypto/sha256"
	"fmt"
	"log"

	"github.com/eoscanada/eos-go"
)

//EOSServer store nam of server used
var EOSServer = "http://192.168.0.1:8888"

func getAPI() *eos.API {
	api := eos.New(EOSServer) //bytes.Repeat([]byte{0}, 32))

	///api.Debug = true
	///eos.Debug = true

	keyBag := eos.NewKeyBag()
	for _, key := range []string{
		///Private keys here
	} {
		if err := keyBag.Add(key); err != nil {
			log.Fatalln("Couldn't load private key:", err)
		}
	}
	api.SetSigner(keyBag)
	return api
}
```