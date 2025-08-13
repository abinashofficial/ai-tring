package app

import (
	"aitring/handlers"
	"aitring/store"
	"aitring/services"
	audioHandler "aitring/handlers/audios"
	audioStore "aitring/store/audiostore"
	audioServ "aitring/services/audios"
)

var h handlers.Store
var repos store.Store
var serv services.Store



func setupRepos() {
	repos = store.Store{
		AudioStore: audioStore.New(),
	}
}

func setupHandlers() {
	h = handlers.Store{
		AudioHandler: audioHandler.New(serv),
	}
}

func setupService() {
	serv = services.Store{
		AudioServ: audioServ.New(repos),
	}
}


func Start() {
	envPort := "8080"
	setupRepos()
	setupService()
	setupHandlers()
	runServer(envPort, h)
}
