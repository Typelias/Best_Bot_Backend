package subreddit

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

var hundredAmount int = 2

// Subreddit is a struct that can handle posts from multiple subbreddits
type Subreddit struct {
	subreddits    map[string][]string
	subredditList []string
	m             sync.Mutex
}

type postData struct {
	URL string `json:"url"`
}

type post struct {
	Data postData `json:"data"`
}

type resultData struct {
	Posts []post `json:"children"`
	After string `json:"after"`
}

type result struct {
	Data resultData `json:"data"`
}

func endsWith(s string, end string) bool {
	return strings.HasSuffix(s, end)
}

func (s *Subreddit) populateSubreddits() {
	s.m.Lock()
	for _, sub := range s.subredditList {
		s.subreddits[sub] = redditGetter(sub, hundredAmount)
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
	s.subredditList = []string{"funny"}
	s.subreddits = make(map[string][]string)
	s.populateSubreddits()
	go s.updator()
}

//GetAllSubreddits returns a list of all subbreddits added to the getter
func (s *Subreddit) GetAllSubreddits() []string {
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
			s.subreddits[sub] = redditGetter(sub, hundredAmount)
			fmt.Println("Added subbreddit: ", sub, "With post amount:", len(s.subreddits[sub]))
		}
	}
	fmt.Println("Added all new subbreddits")
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

func getFirstPosts(subreddit string) resultData {
	client := http.Client{
		Timeout: time.Second * 2,
	}

	req, err := http.NewRequest(http.MethodGet, "https://www.reddit.com/r/"+subreddit+"/hot.json?limit=100", nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("User-Agent", "your bot 0.1")
	res, getErr := client.Do(req)

	if getErr != nil {
		fmt.Println("Get Error", getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		fmt.Println("Read Error", readErr)
	}

	var redditResult result

	jsonErr := json.Unmarshal(body, &redditResult)
	if jsonErr != nil {
		fmt.Println("Json Error", jsonErr)
	}

	return redditResult.Data

}

func getMore(subreddit string, after string) resultData {
	client := http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest(http.MethodGet, "https://www.reddit.com/r/"+subreddit+"/hot.json?limit=100&after"+after, nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("User-Agent", "your bot 0.1")
	res, getErr := client.Do(req)

	if getErr != nil {
		fmt.Println("Get Error", getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		fmt.Println("Read Error", readErr)
	}

	var redditResult result

	jsonErr := json.Unmarshal(body, &redditResult)
	if jsonErr != nil {
		fmt.Println("Json Error", jsonErr)
	}

	return redditResult.Data

}

func redditGetter(subreddit string, numberOfHundred int) []string {
	temp := getFirstPosts(subreddit)
	after := temp.After
	posts := temp.Posts

	if numberOfHundred > 1 {
		for i := numberOfHundred; i > 1; i-- {
			temp = getMore(subreddit, after)
			posts = append(posts, temp.Posts...)
			after = temp.After
		}
	}

	var ret []string
	for _, post := range posts {
		url := post.Data.URL
		if endsWith(url, "jpg") || endsWith(url, "gif") || endsWith(url, "png") || endsWith(url, "jpeg") || endsWith(url, "gifv") {
			ret = append(ret, url)
		} else if strings.Contains(url, "gfycat") || strings.Contains(url, "redgifs") || strings.Contains(url, "imgur.") {
			ret = append(ret, url)
		}
	}

	return ret
}
