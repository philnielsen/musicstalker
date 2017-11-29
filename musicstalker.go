package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/ChimeraCoder/anaconda"
	"github.com/zmb3/spotify"
)

var twitapi *anaconda.TwitterApi

var spotifyClient *spotify.Client

const redirectURI = "http://localhost:8080/callback"

var (
	auth  = spotify.NewAuthenticator(redirectURI, spotify.ScopePlaylistModifyPublic)
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

	authUrl := auth.AuthURL(state)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", authUrl)

	resp, err := http.Get(authUrl)
	fmt.Println("HTTP RESPONSE", resp)

	// wait for auth to complete
	spotifyClient = <-ch

	// use the client to make calls that require authorization
	user, err := spotifyClient.CurrentUser()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("You are logged in as:", user.ID)

	//pull arg for playlist name
	playlistName := os.Args[2]

	var playlistToEdit spotify.SimplePlaylist

	//Get Playlists for User that we are logged in as
	playlists, err := spotifyClient.GetPlaylistsForUser(user.ID)
	if err != nil {
		log.Fatal(err)
	}
	for _, item := range playlists.Playlists {
		if item.Name == playlistName {
			playlistToEdit = item
		}
	}

	// //Add Specific Track to the playlist
	// newPlaylist, err := spotifyClient.AddTracksToPlaylist(user.ID, playlistToEdit.ID, "6LGabqtvan3SGYcL4guT0o")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("New Playlist ID", newPlaylist)

	twitterUsers := os.Args[1]

	userObj, err := twitapi.GetUsersShow(twitterUsers, nil)
	if err != nil {
		log.Println("Error while querying twitter API", err)
		return
	}

	twitterStream := twitapi.PublicStreamFilter(url.Values{"follow": []string{userObj.IdStr}})

	fmt.Println("Stream started, let the stalking commence")

	//twitterStream := api.PublicStreamSample(nil)
	for {
		x := <-twitterStream.C
		switch tweet := x.(type) {
		case anaconda.Tweet:
			//Add Specific Track to the playlist
			println("TWEET: ", tweet.ExtendedTweet.FullText)
			newPlaylist, err := spotifyClient.AddTracksToPlaylist(user.ID, playlistToEdit.ID, "6LGabqtvan3SGYcL4guT0o")
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("New Playlist ID", newPlaylist)

			return
		case anaconda.StatusDeletionNotice:
			// pass
		default:
			fmt.Printf("unknown type(%T) : %v \n", x, x)
		}
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
