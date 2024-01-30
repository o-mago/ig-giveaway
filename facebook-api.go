package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

type posts struct {
	Data []postsData `json:"data"`
}

type postsData struct {
	Comments  postsCommentsData `json:"comments"`
	ShortCode string            `json:"shortcode"`
}

type postsCommentsData struct {
	Data []postsComments `json:"data"`
}

type postsComments struct {
	Text     string `json:"text"`
	Username string `json:"username"`
}

func getPostsData(userID, token string) (posts, error) {
	url := "https://graph.facebook.com/v19.0/%s/media?fields=comments{username,text},shortcode&limit=10&access_token=%s"

	resp, err := http.Get(fmt.Sprintf(url, userID, token))
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}

	var response posts
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return posts{}, err
	}

	return response, nil
}

type userInfo struct {
	Data userInfoData `json:"data"`
}

type userInfoData struct {
	User userInfoDataUser `json:"user"`
}

type userInfoDataUser struct {
	ID string `json:"fbid"`
}

func getUserInfo(username string) (string, error) {
	target := "https://www.instagram.com/api/v1/users/web_profile_info/?"

	params := url.Values{}

	params.Set("username", username)

	req, err := http.NewRequest("GET", target+params.Encode(), nil)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}

	req.Header = http.Header{
		"x-ig-app-id": {"936619743392459"},
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("do: could not make request: %s\n", err)
		os.Exit(1)
	}

	var response userInfo
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}

	return response.Data.User.ID, nil
}
