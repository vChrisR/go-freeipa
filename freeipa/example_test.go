package freeipa_test

import (
	"crypto/tls"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/vchrisr/go-freeipa/freeipa"
)

func Example_addUser() {
	tspt := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // WARNING DO NOT USE THIS OPTION IN PRODUCTION
		},
	}
	c, e := freeipa.Connect("dc1.test.local", tspt, "admin", "walrus123")
	if e != nil {
		log.Fatal(e)
	}

	rand.Seed(time.Now().UTC().UnixNano())
	uid := fmt.Sprintf("jdoe%v", rand.Int())

	res, e := c.UserAdd(&freeipa.UserAddArgs{
		Givenname: "John",
		Sn:        "Doe",
	}, &freeipa.UserAddOptionalArgs{
		UID: freeipa.String(uid),
	})
	if e != nil {
		log.Fatal(e)
	}

	fmt.Printf("Added user %v", *res.Result.Cn)
	// Output: Added user John Doe
}

func Example_errorHandling() {
	tspt := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // WARNING DO NOT USE THIS OPTION IN PRODUCTION
		},
	}
	c, e := freeipa.Connect("dc1.test.local", tspt, "admin", "walrus123")
	if e != nil {
		log.Fatal(e)
	}

	_, e = c.UserShow(&freeipa.UserShowArgs{}, &freeipa.UserShowOptionalArgs{
		UID: freeipa.String("somemissinguid"),
	})
	if e == nil {
		fmt.Printf("No error")
	} else if ipaE, ok := e.(*freeipa.Error); ok {
		fmt.Printf("FreeIPA error %v: %v\n", ipaE.Code, ipaE.Message)
		if ipaE.Code == freeipa.NotFoundCode {
			fmt.Println("(matched expected error code)")
		}
	} else {
		fmt.Printf("Other error: %v", e)
	}

	// Output: FreeIPA error 4001: somemissinguid: user not found
	// (matched expected error code)
}
