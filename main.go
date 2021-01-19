// Declaration of the main package
package main

// Importing packages
import (
	"crypto/tls"
	"encoding/json"
	"html/template"
	// ADDED "io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	log "unknwon.dev/clog/v2"

	api "gogs"
	"webauthn/protocol"
)

func init() {
	err := log.NewConsole()
	if err != nil {
		panic("unable to create new logger: " + err.Error())
	}
}

func handleIndexHelper(client *api.Client, template_file string) func(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(template_file))

	return func(w http.ResponseWriter, r *http.Request) {
		var data struct {
			UserName        string
			Repos           []string
			WebauthnOptions string
		}

		// Fetch the `UserName`
		user, err := client.GetSelfInfo()
		if err != nil {
			log.Error("%v", err)
			return
		}
		data.UserName = user.UserName

		// Fetch the repositories
		repos, err := client.ListMyRepos()
		if err != nil {
			log.Error("%v", err)
			return
		}
		data.Repos = make([]string, len(repos))

		for i, repo := range repos {
			data.Repos[i] = repo.Name
		}

		// TODO: Need to check if webauthn is enabled, maybe in the `Repo_GenericWebauthnBegin`
		// function, returns an empty `options`
		//
		// Get the webauthn assertion options to pre-load into the web-page
		options, err := client.Repo_GenericWebauthnBegin()
		if err != nil {
			log.Error("%v", err)
			return
		}

		// Fill in the txAuthn text with a format placeholder
		options.Response.Extensions = protocol.AuthenticationExtensions{
			"txAuthSimple": "Confirm deletion of repository: {0}/{1}",
		}

		// JSON encode the `options` and place them in the template `data` struct
		repo_options, err := json.Marshal(options.Response)
		if err != nil {
			log.Error("%v", err)
			return
		}
		data.WebauthnOptions = string(repo_options)

		tmpl.Execute(w, data)
	}
}

func handleDeleteRepoHelper(client *api.Client) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Read the form entries
		repo := r.PostFormValue("repo_name")
		webauthnData := r.PostFormValue("webauthn_data")

		// Fetch our `UserName`, presumably the owner of `repo`
		user, err := client.GetSelfInfo()
		if err != nil {
			log.Error("%v", err)
			return
		}

		// Populate the webauthn container
		opt := api.WebauthnContainer{
			WebauthnData: webauthnData,
		}

		// Delete the repository
		err = client.DeleteRepo(user.UserName, repo, opt)
		if err != nil {
			log.Error("%v", err)

			// Redirect back to the index
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		log.Info("Deleted repository: %s", repo)

		// Redirect back to the index
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// Main function
func main() {
	// Server configurations
	serverAddress := ":8080"

	// Connect to the gogs API
	url := "https://localhost:3000"
	token := "48f07353f272b9166450eba14b7576ffa7104cce"
	client := api.NewClient(url, token)

	// The HTTPS certificate is self-signed, skip verifying it
	http_client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	client.SetHTTPClient(http_client)

	r := mux.NewRouter()

	// Serve some basic templated HTML
	r.HandleFunc("/", handleIndexHelper(client, "index.tmpl")).Methods("GET")
	r.HandleFunc("/delete_repo", handleDeleteRepoHelper(client)).Methods("POST")

	// Serve the javascript parts of this app
	jsDir := http.FileServer(http.Dir("./js"))
	r.PathPrefix("/").Handler(http.StripPrefix("/js", jsDir))

	log.Info("Starting server at %s", serverAddress)
	log.Fatal("%v", http.ListenAndServeTLS(serverAddress, "server.crt", "server.key", r))
}
