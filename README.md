# Best_Bot_Backend
Backend for Best_Bot  
https://github.com/Typelias/Best_Bot  
# Usage
It runs on port 8080
- "/subreddit" method POST. It expects an array of strings.  
Used this to tell the server what Subreddits to get post from  

- "/subreddit/{id}" method GET.   
Gets a random Image(jpg,png)/gif/gifv or gfycat URL

- "/subreddits" method GET  
Returns an array with all subreddits the backend keeps track of

# Requirements
* GoLang
* It should download dependencies when building. If it somehow does not they are:
  * gorilla/mux: https://github.com/gorilla/mux

# Installation
- Clone repo
- Run go build
