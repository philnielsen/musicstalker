package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ChimeraCoder/anaconda"
	"github.com/zmb3/spotify"
)

var twitapi *anaconda.TwitterApi

const redirectURI = "http://localhost:8080/callback"

var (
	auth  = spotify.NewAuthenticator(redirectURI, spotify.ScopeUserLibraryRead)
	ch    = make(chan *spotify.Client)
	state = "abc123"
)

func main() {
	//Twitter API Setup
	anaconda.SetConsumerKey(ConsumerKey)
	anaconda.SetConsumerSecret(ConsumerSecret)
	twitapi = anaconda.NewTwitterApi(Token, TokenSecret)

	//Spotify Client
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
	go http.ListenAndServe(":8080", nil)

	url := auth.AuthURL(state)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

	// wait for auth to complete
	spotifyClient := <-ch

	// use the client to make calls that require authorization
	user, err := spotifyClient.CurrentUser()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("You are logged in as:", user.ID)

	twitterUsers := os.Args[1]

	userObjs, err := twitapi.GetUsersLookup(twitterUsers, nil)
	fmt.Println("Users:")
	for _, item := range userObjs {
		fmt.Println("   ", item.Name)
	}

	if err != nil {
		log.Println("Error while querying twitter API", err)
		return
	}

	currentUserTracks, err := spotifyClient.CurrentUsersTracks()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Saved Songs:")
	for _, item := range currentUserTracks.Tracks {
		fmt.Println("   ", item.Name)
	}

	results, err := spotifyClient.Search("holiday", spotify.SearchTypePlaylist|spotify.SearchTypeAlbum)
	if err != nil {
		log.Fatal(err)
	}

	// handle album results
	if results.Albums != nil {
		fmt.Println("Albums:")
		for _, item := range results.Albums.Albums {
			fmt.Println("   ", item.Name)
		}
	}
	// handle playlist results
	if results.Playlists != nil {
		fmt.Println("Playlists:")
		for _, item := range results.Playlists.Playlists {
			fmt.Println("   ", item.Name)
		}
	}
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}
	// use the token to get an authenticated client
	client := auth.NewClient(tok)
	fmt.Fprintf(w, "Login Completed!")
	ch <- &client
}
