package main

import (
	"fmt"
	"math/rand"
	"regexp"
	"slices"
	"time"
)

type giveaway map[string][]string

type contenders []string

var regex = regexp.MustCompile("@[^ ]*")

func (c *contenders) Random(totalWinners int, blockList ...string) ([]string, error) {
	if c == nil || len(*c) == 0 {
		return nil, fmt.Errorf("empty giveaway contenders")
	}

	winners := []string{}

	selectedIndexes := map[int]bool{}

	for len(winners) < totalWinners {
		if len(selectedIndexes) >= len(*c) {
			break
		}

		randomIndex := rand.Intn(len(*c))

		for selectedIndexes[randomIndex] {
			randomIndex = rand.Intn(len(*c))
		}

		selectedUser := (*c)[randomIndex]

		if slices.Contains(blockList, selectedUser) {
			continue
		}

		winners = append(winners, selectedUser)
	}

	return winners, nil
}

type startGiveawayInput struct {
	userName        string
	postCode        string
	token           string
	totalMentions   int
	totalWinners    int
	blockList       []string
	multipleEntries bool
}

func (m *model) startGiveaway(input startGiveawayInput) {
	userID, err := getUserInfo(input.userName)
	if err != nil {
		panic(err)
	}

	m.percent += 0.1

	postID := ""

	nextURL := ""

	for {
		posts, err := getPostsData(userID, input.token, nextURL)
		if err != nil {
			panic(err)
		}

		nextURL = posts.Paging.Next

		for _, post := range posts.Data {
			if post.ShortCode == input.postCode {
				postID = post.ID

				break
			}
		}

		if postID != "" {
			break
		}

		if posts.Paging.Next == "" {
			break
		}
	}

	m.percent += 0.1

	var commentsFinal []commentsData

	nextURL = ""

	for {
		comments, err := getCommentsData(postID, input.token, nextURL)
		if err != nil {
			panic(err)
		}

		nextURL = comments.Paging.Next

		commentsFinal = append(commentsFinal, comments.Data...)

		if comments.Paging.Next == "" {
			break
		}
	}

	m.percent += 0.1

	usersGiveaway := giveaway{}

	for _, comment := range commentsFinal {
		mentions := regex.FindAllString(comment.Text, -1)

		usersGiveaway[comment.Username] = append(usersGiveaway[comment.Username], mentions...)
	}

	finalList := contenders{}

	for userName, mentions := range usersGiveaway {
		uniqueMentions := slices.Compact(mentions)

		entries := 1
		if input.multipleEntries {
			entries = len(uniqueMentions) / input.totalMentions
		}

		if len(uniqueMentions) >= input.totalMentions {
			for i := 0; i < entries; i++ {
				finalList = append(finalList, userName)
			}
		}
	}

	for i := 0; i < 6; i++ {
		m.percent += 0.1

		time.Sleep(time.Second)
	}

	winners, err := finalList.Random(input.totalWinners)
	if err != nil {
		panic(err)
	}

	winnersGiveaway := giveaway{}

	for _, winner := range winners {
		winnersGiveaway[winner] = usersGiveaway[winner]
	}

	m.percent += 0.1

	time.Sleep(300 * time.Millisecond)

	m.winners = winnersGiveaway
	m.finish = true
}
