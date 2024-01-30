package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

type posts struct {
	Data   []postsData `json:"data"`
	Paging postsPaging `json:"paging"`
}

type postsData struct {
	ID        string `json:"id"`
	ShortCode string `json:"shortcode"`
}

type postsPaging struct {
	Next string `json:"next"`
}

func getPostsData(userID, token, url string) (posts, error) {
	if url == "" {
		url = "https://graph.facebook.com/v19.0/%s/media?fields=shortcode&access_token=%s"
		url = fmt.Sprintf(url, userID, token)
	}

	resp, err := http.Get(url)
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

type comments struct {
	Data   []commentsData `json:"data"`
	Paging commentsPaging `json:"paging"`
}

type commentsData struct {
	Text     string `json:"text"`
	Username string `json:"username"`
}

type commentsPaging struct {
	Next string `json:"next"`
}

func getCommentsData(postID, token, url string) (comments, error) {
	if url == "" {
		url = "https://graph.facebook.com/v19.0/%s/comments?fields=text,username&access_token=%s"
		url = fmt.Sprintf(url, postID, token)
	}

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}

	var response comments
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return comments{}, err
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
