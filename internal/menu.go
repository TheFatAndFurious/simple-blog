package internal

import (
	settings "simpleBlog/config"
)



type MenuLink struct {
    Name string
    Path string
}
var linksList []MenuLink

func InitMenu () []MenuLink {
	// first we need to concatenate pages and folders fronm the settings package
	var links []string
	links = append(links, settings.Pages...)
	links = append(links, settings.Folders...)

	for _, link := range links {
		newLink := MenuLink{
			Name: link,
			Path: link + ".html",
		}
	linksList = append(linksList, newLink)	
	}
	return linksList
}
