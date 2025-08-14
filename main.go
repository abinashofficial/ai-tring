package main

import(
	"aitring/app"
"aitring/tests"
"flag"
)


func main() {
		var runUnitTests bool

	flag.BoolVar(&runUnitTests, "runUnitTests", true, "Setting to true will run unit tests")
	flag.Parse()

	if runUnitTests {
		tests.Start()
	} else {
		app.Start()
	}

}