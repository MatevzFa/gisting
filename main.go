package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/google/go-github/github"
	simplehttp "github.com/matevzfa/simplehttp"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var gistsAPI = simplehttp.SimpleHTTP{
	BaseURL: "https://api.github.com"}

func addGistFiles(gist *github.Gist, files *[]string) {
	for _, fpath := range *files {
		fname := github.GistFilename(filepath.Base(fpath))
		fcontent, err := ioutil.ReadFile(fpath)
		if err != nil {
			log.Fatal(err)
		}
		content := string(fcontent)
		gist.Files[fname] = github.GistFile{
			Content: &content,
		}
	}

}

func getGist(id string) github.Gist {
	gist := github.Gist{}
	gistsAPI.Get("/gists/"+id, &gist)
	return gist
}

func gistsIds() []github.Gist {
	gistsList := make([]github.Gist, 0)
	gistsAPI.Get("/gists", &gistsList)
	return gistsList
}

func downloadGist(id string) {
	gist := getGist(id)
	for fileName, file := range gist.Files {
		log.Println("Writing file " + fileName)
		error := ioutil.WriteFile(string(fileName), []byte(*file.Content), 0644)
		if error != nil {
			log.Fatal(error)
		}
	}
}

func createGist(newGist github.Gist) github.Gist {
	body, err := json.Marshal(&newGist)
	if err != nil {
		log.Fatal(err)
	}
	createdGist := github.Gist{}
	gistsAPI.Post("/gists", simplehttp.ToReader(body), &createdGist)
	return createdGist
}

func updateGist(id string, updatedGist github.Gist) github.Gist {
	body, err := json.Marshal(&updatedGist)
	if err != nil {
		log.Fatal(err)
	}
	respGist := github.Gist{}
	gistsAPI.Patch("/gists/"+id, simplehttp.ToReader(body), &updatedGist)
	return respGist
}

func setToken(ctx *kingpin.ParseContext) error {
	gistsAPI.SetHeaders([]simplehttp.HTTPHeader{
		{Key: "Authorization", Value: "token " + *ctx.Elements[0].Value},
	})
	return nil
}

// command line interface
var (
	download   = kingpin.Command("download", "Download a gist.").Alias("d")
	downloadID = download.Arg("id", "Id of that gist.").Required().String()

	create            = kingpin.Command("create", "Create a gist.").Alias("c")
	createDescription = create.Flag("description", "Description of the gist").Short('d').String()
	private           = create.Flag("private", "If creating a private gist").Bool()
	createFiles       = create.Arg("files", "Files in the gist").Required().ExistingFiles()

	update            = kingpin.Command("update", "Update a gist.").Alias("u")
	updateID          = update.Arg("id", "ID of that gist.").Required().String()
	updateDescription = update.Flag("description", "Description of the gist").Short('d').String()
	updateFiles       = update.Arg("files", "Files to update or add.").Required().ExistingFiles()

	token = kingpin.Flag("token", "OAuth token for accessing the gist API").Short('t').Action(setToken).String()
)

func main() {
	switch kingpin.Parse() {

	case "download":
		downloadGist(*downloadID)

	case "create":
		public := !*private
		gist := github.Gist{
			Public:      &public,
			Description: createDescription,
			Files:       map[github.GistFilename]github.GistFile{},
		}
		addGistFiles(&gist, createFiles)
		createdGist := createGist(gist)
		fmt.Println(*createdGist.HTMLURL)

	case "update":
		public := !*private
		gist := github.Gist{
			Public:      &public,
			Description: createDescription,
			Files:       map[github.GistFilename]github.GistFile{},
		}
		addGistFiles(&gist, updateFiles)
		updateGist(*updateID, gist)
	}
}
