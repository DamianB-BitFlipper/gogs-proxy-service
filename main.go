// Declaration of the main package
package main

// Importing packages
import (
	"crypto/tls"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	log "unknwon.dev/clog/v2"

	gogs_api "gogs"
)

func init() {
	err := log.NewConsole()
	if err != nil {
		panic("unable to create new logger: " + err.Error())
	}
}

func handleIndexHelper(client *gogs_api.Client, template_file string) func(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(template_file))

	return func(w http.ResponseWriter, r *http.Request) {
		var data struct {
			UserName string
			Repos    []string
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

		tmpl.Execute(w, data)
	}
}

func handleDeleteRepoHelper(client *gogs_api.Client) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			log.Error("%v", err)
			return
		}

		// TODO: Make this work
		repo := r.Form.Get("repoName")

		// Fetch our `UserName`, presumably the owner of `repo`
		user, err := client.GetSelfInfo()
		if err != nil {
			log.Error("%v", err)
			return
		}

		// Delete the repository
		err = client.DeleteRepo(user.UserName, repo)
		if err != nil {
			log.Error("%v", err)
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
	client := gogs_api.NewClient(url, token)

	// The HTTPS certificate is self-signed, skip verifying it
	http_client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	client.SetHTTPClient(http_client)

	// Serve some basic templated HTML
	r := mux.NewRouter()

	r.HandleFunc("/", handleIndexHelper(client, "index.tmpl")).Methods("GET")
	r.HandleFunc("/delete_repo", handleDeleteRepoHelper(client)).Methods("POST")

	log.Info("Starting server at %s", serverAddress)
	log.Fatal("%v", http.ListenAndServeTLS(serverAddress, "server.crt", "server.key", r))
}
