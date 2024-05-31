package internal

import (
	"os"
	"path/filepath"
	settings "simpleBlog/config"
)

func InitPublic() {

for _, page := range settings.Pages {
	os.Create(filepath.Join("public", page + ".html"));
}

for _, subMenus := range settings.Folders {
	os.Chdir("public")
	os.Mkdir(subMenus, os.ModePerm)
}
os.Chdir("..")

}