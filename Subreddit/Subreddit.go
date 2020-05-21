package subreddit

import (
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/turnage/graw/reddit"
)

// Subreddit is a struct that can handle posts from multiple subbreddits
type Subreddit struct {
	bot reddit.Bot
	subreddits    map[string][]string
	subredditList []string
	m             sync.Mutex
}

func endsWith(s string, end string) bool {
	return strings.HasSuffix(s, end)
}

func (s *Subreddit) getPosts(subreddit string) []string {
	url := "/r/" + subreddit
	harvest, _ := s.bot.Listing(url, "")
	var ret []string
	for _, post := range harvest.Posts {
		url := post.URL
		if endsWith(url, "jpg") || endsWith(url, "gif") || endsWith(url, "png") || endsWith(url, "jpeg") || endsWith(url, "gifv") {
			ret = append(ret, url)
		} else if strings.Contains(url, "gfycat") || strings.Contains(url, "redgifs") || strings.Contains(url, "imgur.") {
			ret = append(ret, url)
		}
	}

	return ret
}

func (s *Subreddit) populateSubreddits() {
	s.m.Lock()
	for _, sub := range s.subredditList {
		s.subreddits[sub] = s.getPosts(sub)
	}
	s.m.Unlock()
}

func (s *Subreddit) updator() {
	for true {
		time.Sleep(8 * time.Hour)
		s.populateSubreddits()

	}
}

// Init Creates the subreddit object
func (s *Subreddit) Init() {
	s.bot, _ = reddit.NewBotFromAgentFile("agent.txt", 0)
	s.subredditList = []string{"funny"}
	s.subreddits = make(map[string][]string)
	s.populateSubreddits()
	go s.updator()
}

func (s *Subreddit) GetAllSubreddits()([]string) {
	return s.subredditList
}

// UpdateSubbredditList updates the list of subreddits
func (s *Subreddit) UpdateSubbredditList(newList []string) {
	s.m.Lock()
	for _, sub := range newList {
		if search(s.subredditList, sub) {
			continue
		} else {
			s.subredditList = append(s.subredditList, sub)
			s.subreddits[sub] = s.getPosts(sub)
		}
	}
	s.m.Unlock()
}

// GetRandomPost gets random post from a subreddit
func (s *Subreddit) GetRandomPost(subreddit string) string {
	s.m.Lock()
	rand.Seed(time.Now().Unix())
	ret := s.subreddits[subreddit][(rand.Int()%len(s.subreddits[subreddit]) - 1)]
	s.m.Unlock()
	return ret
}

func search(list []string, key string) bool {
	for i := range list {
		if list[i] == key {
			return true
		}
	}

	return false
}
