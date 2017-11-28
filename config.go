package main

import "os"

// ConsumerKey - Twitter API Consumer Key
var ConsumerKey = os.Getenv("ConsumerKey")

// ConsumerSecret - Twitter API Consumer Secret
var ConsumerSecret = os.Getenv("ConsumerSecret")

// Token - Twitter API Access Token
var Token = os.Getenv("Token")

// TokenSecret - Twitter API Access Token Secret
var TokenSecret = os.Getenv("TokenSecret")

// SpotifyClientID - Spotify Client ID For API Access
var SpotifyClientID = os.Getenv("SPOTIFY_ID")

//SpotifyClientSecret - Spotify Client Secret for API Access
var SpotifyClientSecret = os.Getenv("SPOTIFY_SECRET")
