package app

import (
	"aitring/handlers"
	"aitring/store"
	"aitring/services"
	audioHandler "aitring/handlers/audios"
	audioStore "aitring/store/audiostore"
	audioServ "aitring/services/audios"
	"log"
	"aitring/store/pipelinestore"
	"context"
)

var h handlers.Store
var repos store.Store
var serv services.Store



func setupRepos() {
	cfg := pipelinestore.DefaultConfig()
	audStore := audioStore.NewAudioStore("metadata.json")
	repos = store.Store{
		AudioStore: audioStore.New(),
		PipelineStore: pipelinestore.New(cfg, audStore),
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
		log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	
	setupRepos()
	setupService()
	setupHandlers()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	repos.PipelineStore.Start(ctx)
	envPort := "8080"

	runServer(envPort, h)
}
