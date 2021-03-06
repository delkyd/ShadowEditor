package editor

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/tengge1/shadoweditor/helper"
)

// CopyAssets copy the assets needed to the exported scene.
func CopyAssets(path string) error {
	// copy html files
	sourceName := helper.MapPath("/index.html")
	destName := filepath.Join(path, "editor.html")
	if err := helper.CopyFile(sourceName, destName); err != nil {
		return err
	}

	bytes, err := ioutil.ReadFile(destName)
	if err != nil {
		return err
	}

	text := strings.ReplaceAll(string(bytes), "location.origin", "'.'") // make api path to current path
	if err := ioutil.WriteFile(destName, []byte(text), 0755); err != nil {
		return err
	}

	// copy build folder
	dirName := filepath.Join(path, "build")
	if _, err := os.Stat(dirName); err != nil {
		os.MkdirAll(dirName, 0755)
	}

	sourceName = helper.MapPath("/build/ShadowEditor.js")
	destName = filepath.Join(path, "build", "ShadowEditor.js")
	if err := helper.CopyFile(sourceName, destName); err != nil {
		return err
	}

	// copy assets folder
	sourceName = helper.MapPath("/assets")
	destName = filepath.Join(path, "assets")
	if err := helper.CopyDirectory(sourceName, destName); err != nil {
		return err
	}

	// copy language pack
	sourceName = helper.MapPath("/lang")
	destName = filepath.Join(path, "lang")
	if err := helper.CopyDirectory(sourceName, destName); err != nil {
		return err
	}

	// copy website icon
	sourceName = helper.MapPath("/favicon.ico")
	destName = filepath.Join(path, "favicon.ico")
	return helper.CopyDirectory(sourceName, destName)
}
