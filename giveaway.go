package main

import (
	"fmt"
	"math/rand"
	"regexp"
	"time"
)

type giveaway map[string]string

var regex = regexp.MustCompile("@[^ ]*")

func (g *giveaway) Random() (string, error) {
	if g == nil || len(*g) == 0 {
		return "", fmt.Errorf("empty giveaway contenders")
	}

	randomIndex := rand.Intn(len(*g))

	index := 0
	for userName := range *g {
		if index == randomIndex {
			return userName, nil
		}

		index++
	}

	return "", fmt.Errorf("empty giveaway contenders")
}

func (m *model) startGiveaway(userName, postCode, token, totalMentions string, shouldFilter bool) {
	userID, err := getUserInfo(userName)
	if err != nil {
		panic(err)
	}

	posts, err := getPostsData(userID, token)
	if err != nil {
		panic(err)
	}

	contenders := giveaway{}

	for _, post := range posts.Data {
		if post.ShortCode != postCode {
			continue
		}

		for _, comment := range post.Comments.Data {
			mentions := regex.FindAllString(comment.Text, -1)

			mentionCounter := len(mentions)

			if mentionCounter >= 3 {
				contenders[comment.Username] = comment.Text
			}
		}
	}

	for i := 0; i < 10; i++ {
		m.percent += 0.1

		time.Sleep(time.Second)
	}

	winner, err := contenders.Random()
	if err != nil {
		panic(err)
	}

	m.winner = "@" + winner
	m.winnerText = contenders[winner]
}
