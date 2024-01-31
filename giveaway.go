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

	if totalWinners > len(*c) {
		return *c, nil
	}

	for len(winners) < totalWinners {
		randomIndex := rand.Intn(len(*c))

		for selectedIndexes[randomIndex] {
			randomIndex = rand.Intn(len(*c))
		}

		winners = append(winners, (*c)[randomIndex])
	}

	return winners, nil
}

func (c *contenders) RandomAllContenders(totalWinners int, model *model, blockList ...string) ([]string, error) {
	if c == nil || len(*c) == 0 {
		return nil, fmt.Errorf("empty giveaway contenders")
	}

	winners := []string{}

	selectedIndexes := map[int]bool{}

	if totalWinners > len(*c) {
		return *c, nil
	}

	step := 0.6 / float64(totalWinners*30)

	for winnerPos := 0; winnerPos < totalWinners; winnerPos++ {
		var randomIndex int

		for i := 0; i < 30; i++ {
			randomIndex = rand.Intn(len(*c))

			for selectedIndexes[randomIndex] {
				randomIndex = rand.Intn(len(*c))
			}

			model.selectedContenders[winnerPos] = randomIndex

			model.percent += step

			time.Sleep(100 * time.Millisecond)
		}

		selectedIndexes[randomIndex] = true

		selectedUser := (*c)[randomIndex]

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
	allContenders   bool
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
		if slices.Contains(input.blockList, comment.Username) {
			continue
		}

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

	if !m.allContenders {
		for i := 0; i < 6; i++ {
			m.percent += 0.1

			time.Sleep(time.Second)
		}
	}

	var winners []string

	if !input.allContenders {
		winners, err = finalList.Random(input.totalWinners)
		if err != nil {
			panic(err)
		}
	} else {
		m.contenders = finalList

		winners, err = finalList.RandomAllContenders(input.totalWinners, m)
		if err != nil {
			panic(err)
		}
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
