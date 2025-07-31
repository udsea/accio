package sites

// Site represents a website where a username can be checked
type Site struct {
	Name        string // Name of the site
	URL         string // URL for display purposes
	ErrorType   string // Type of error to check for (status_code, message, etc.)
	ErrorMsg    string // Error message to look for
	URLProbe    bool   // Whether to use URL probing
	URLFormat   string // URL format with {} placeholder for username
	RegexCheck  string // Regex to check in the response
	CheckMethod string // HTTP method to use (GET, HEAD, etc.)
}

// GetSites returns a list of sites to check
func GetSites() []Site {
	return []Site{
		{
			Name:        "GitHub",
			URL:         "https://github.com/{}",
			ErrorType:   "status_code",
			URLProbe:    false,
			URLFormat:   "https://github.com/{}",
			CheckMethod: "GET",
		},
		{
			Name:        "Twitter",
			URL:         "https://twitter.com/{}",
			ErrorType:   "status_code",
			URLProbe:    false,
			URLFormat:   "https://twitter.com/{}",
			CheckMethod: "GET",
		},
		{
			Name:        "Instagram",
			URL:         "https://www.instagram.com/{}",
			ErrorType:   "status_code",
			URLProbe:    false,
			URLFormat:   "https://www.instagram.com/{}",
			CheckMethod: "GET",
		},
		{
			Name:        "Facebook",
			URL:         "https://www.facebook.com/{}",
			ErrorType:   "status_code",
			URLProbe:    false,
			URLFormat:   "https://www.facebook.com/{}",
			CheckMethod: "GET",
		},
		{
			Name:        "YouTube",
			URL:         "https://www.youtube.com/{}",
			ErrorType:   "status_code",
			URLProbe:    false,
			URLFormat:   "https://www.youtube.com/{}",
			CheckMethod: "GET",
		},
		{
			Name:        "Pinterest",
			URL:         "https://www.pinterest.com/{}",
			ErrorType:   "status_code",
			URLProbe:    false,
			URLFormat:   "https://www.pinterest.com/{}",
			CheckMethod: "GET",
		},
		{
			Name:        "Reddit",
			URL:         "https://www.reddit.com/user/{}",
			ErrorType:   "status_code",
			URLProbe:    false,
			URLFormat:   "https://www.reddit.com/user/{}",
			CheckMethod: "GET",
		},
		{
			Name:        "Twitch",
			URL:         "https://www.twitch.tv/{}",
			ErrorType:   "status_code",
			URLProbe:    false,
			URLFormat:   "https://www.twitch.tv/{}",
			CheckMethod: "GET",
		},
		{
			Name:        "Medium",
			URL:         "https://medium.com/@{}",
			ErrorType:   "status_code",
			URLProbe:    false,
			URLFormat:   "https://medium.com/@{}",
			CheckMethod: "GET",
		},
		{
			Name:        "Quora",
			URL:         "https://www.quora.com/profile/{}",
			ErrorType:   "status_code",
			URLProbe:    false,
			URLFormat:   "https://www.quora.com/profile/{}",
			CheckMethod: "GET",
		},
		{
			Name:        "Flickr",
			URL:         "https://www.flickr.com/people/{}",
			ErrorType:   "status_code",
			URLProbe:    false,
			URLFormat:   "https://www.flickr.com/people/{}",
			CheckMethod: "GET",
		},
		{
			Name:        "Steam",
			URL:         "https://steamcommunity.com/id/{}",
			ErrorType:   "status_code",
			URLProbe:    false,
			URLFormat:   "https://steamcommunity.com/id/{}",
			CheckMethod: "GET",
		},
		{
			Name:        "Vimeo",
			URL:         "https://vimeo.com/{}",
			ErrorType:   "status_code",
			URLProbe:    false,
			URLFormat:   "https://vimeo.com/{}",
			CheckMethod: "GET",
		},
		{
			Name:        "SoundCloud",
			URL:         "https://soundcloud.com/{}",
			ErrorType:   "status_code",
			URLProbe:    false,
			URLFormat:   "https://soundcloud.com/{}",
			CheckMethod: "GET",
		},
		{
			Name:        "Disqus",
			URL:         "https://disqus.com/by/{}",
			ErrorType:   "status_code",
			URLProbe:    false,
			URLFormat:   "https://disqus.com/by/{}",
			CheckMethod: "GET",
		},
		{
			Name:        "Hackernews",
			URL:         "https://news.ycombinator.com/user?id={}",
			ErrorType:   "status_code",
			URLProbe:    false,
			URLFormat:   "https://news.ycombinator.com/user?id={}",
			CheckMethod: "GET",
		},
		{
			Name:        "Deviantart",
			URL:         "https://{}.deviantart.com",
			ErrorType:   "status_code",
			URLProbe:    false,
			URLFormat:   "https://{}.deviantart.com",
			CheckMethod: "GET",
		},
		{
			Name:        "Patreon",
			URL:         "https://www.patreon.com/{}",
			ErrorType:   "status_code",
			URLProbe:    false,
			URLFormat:   "https://www.patreon.com/{}",
			CheckMethod: "GET",
		},
		{
			Name:        "BitBucket",
			URL:         "https://bitbucket.org/{}",
			ErrorType:   "status_code",
			URLProbe:    false,
			URLFormat:   "https://bitbucket.org/{}",
			CheckMethod: "GET",
		},
		{
			Name:        "GitLab",
			URL:         "https://gitlab.com/{}",
			ErrorType:   "status_code",
			URLProbe:    false,
			URLFormat:   "https://gitlab.com/{}",
			CheckMethod: "GET",
		},
		{
			Name:        "Spotify",
			URL:         "https://open.spotify.com/user/{}",
			ErrorType:   "status_code",
			URLProbe:    false,
			URLFormat:   "https://open.spotify.com/user/{}",
			CheckMethod: "GET",
		},
		{
			Name:        "Behance",
			URL:         "https://www.behance.net/{}",
			ErrorType:   "status_code",
			URLProbe:    false,
			URLFormat:   "https://www.behance.net/{}",
			CheckMethod: "GET",
		},
		{
			Name:        "Goodreads",
			URL:         "https://www.goodreads.com/{}",
			ErrorType:   "status_code",
			URLProbe:    false,
			URLFormat:   "https://www.goodreads.com/{}",
			CheckMethod: "GET",
		},
		{
			Name:        "Instructables",
			URL:         "https://www.instructables.com/member/{}",
			ErrorType:   "status_code",
			URLProbe:    false,
			URLFormat:   "https://www.instructables.com/member/{}",
			CheckMethod: "GET",
		},
		{
			Name:        "Keybase",
			URL:         "https://keybase.io/{}",
			ErrorType:   "status_code",
			URLProbe:    false,
			URLFormat:   "https://keybase.io/{}",
			CheckMethod: "GET",
		},
		{
			Name:        "Kongregate",
			URL:         "https://www.kongregate.com/accounts/{}",
			ErrorType:   "status_code",
			URLProbe:    false,
			URLFormat:   "https://www.kongregate.com/accounts/{}",
			CheckMethod: "GET",
		},
		{
			Name:        "Livejournal",
			URL:         "https://{}.livejournal.com",
			ErrorType:   "status_code",
			URLProbe:    false,
			URLFormat:   "https://{}.livejournal.com",
			CheckMethod: "GET",
		},
		{
			Name:        "AngelList",
			URL:         "https://angel.co/{}",
			ErrorType:   "status_code",
			URLProbe:    false,
			URLFormat:   "https://angel.co/{}",
			CheckMethod: "GET",
		},
		{
			Name:        "Last.fm",
			URL:         "https://www.last.fm/user/{}",
			ErrorType:   "status_code",
			URLProbe:    false,
			URLFormat:   "https://www.last.fm/user/{}",
			CheckMethod: "GET",
		},
		{
			Name:        "Dribbble",
			URL:         "https://dribbble.com/{}",
			ErrorType:   "status_code",
			URLProbe:    false,
			URLFormat:   "https://dribbble.com/{}",
			CheckMethod: "GET",
		},
	}
}

// GetSiteByName returns a site by its name
func GetSiteByName(name string) (Site, bool) {
	for _, site := range GetSites() {
		if site.Name == name {
			return site, true
		}
	}
	return Site{}, false
}
