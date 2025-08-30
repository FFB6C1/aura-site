package builder

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ffb6c1/aura-site/internal/config"
	"github.com/ffb6c1/aura-site/internal/file"
	"github.com/ffb6c1/aura-site/internal/markdown"
)

type img struct {
	name        string
	filename    string
	thumb       string
	alt         string
	description string
}

func gallery(page string) string {
	if !strings.Contains(page, "!!!gallery") {
		return markdown.Convert(page, "builder")
	}

	splitPage := strings.Split(page, "!!!gallery[")
	newPage := markdown.Convert(splitPage[0], "builder")
	for _, splitItem := range splitPage[1:] {
		galleryDir, restOfItem, _ := strings.Cut(splitItem, "]")
		if galleryDir == "" {
			newPage += markdown.Convert(restOfItem, "builder")
			continue
		}
		gallery := makeGallery(galleryDir)
		newPage += gallery + markdown.Convert(restOfItem, "builder")
	}

	return newPage
}

func makeGallery(galleryDir string) string {
	path := filepath.Join(config.GetConfig().GetImportPath(), galleryDir)
	allImages, err := file.GetImgFiles(path)
	if err != nil {
		fmt.Printf("could not get images from %s\n", path)
	}
	if len(allImages) == 0 {
		fmt.Printf("no applicable images at %s\n", path)
	}

	images, thumbs := sortImages(allImages)
	alt := getAltText(path)
	descs := getDescriptions(path)

	fullImgData := buildImages(images, thumbs, alt, descs)

	gallery := imagesToGallery(galleryDir, fullImgData)
	return gallery
}

func sortImages(allImages [][]string) (map[string]string, map[string]string) {
	images := make(map[string]string)
	thumbs := make(map[string]string)
	for _, img := range allImages {
		if name, thumb := strings.CutSuffix(img[0], "-thumb"); thumb {
			thumbs[name] = img[0] + "." + img[1]
			continue
		}
		images[img[0]] = img[0] + img[1]
	}
	return images, thumbs
}

func getAltText(path string) map[string]string {
	altMap := make(map[string]string)
	alt, err := file.FileToString(filepath.Join(path, "alt.txt"))
	if err != nil {
		alt, err = file.FileToString(filepath.Join(path, "alt.md"))
		if err != nil {
			return altMap
		}
	}
	split := strings.Split(alt, "!!!")
	for _, splitItem := range split {
		image, altText, found := strings.Cut(splitItem, ":")
		if !found {
			continue
		}
		altMap[image] = altText
	}
	return altMap
}

func getDescriptions(path string) map[string]string {
	descs, err := file.GetMDFiles(path)
	if err != nil {
		return map[string]string{}
	}
	return descs
}

func buildImages(images, thumbs, alt, descs map[string]string) []img {
	imgs := []img{}
	for name, file := range images {
		image := img{
			name:        name,
			filename:    file,
			thumb:       thumbs[name],
			alt:         alt[name],
			description: descs[name],
		}
		imgs = append(imgs, image)
	}
	return imgs
}

func imagesToGallery(galleryName string, images []img) string {
	thumbnails := ""
	noThumbnails := ""

	for _, image := range images {
		address := fmt.Sprintf("%s/%s.html", galleryName, image.name)
		if image.thumb == "" {
			link := fmt.Sprintf("<a href=\"%s\">%s</a><br />\n", address, image.name)
			noThumbnails += link
		} else {
			thumbAddress := fmt.Sprintf("%s/thumbs/%s", galleryName, image.thumb)
			link := fmt.Sprintf("<a href=\"%s\"><img src=\"%s\" alt=\"%s\" class=\"thumb\"></a>\n", address, thumbAddress, image.alt)
			thumbnails += link
		}
	}

	return "<div class=\"gallery\">\n" + noThumbnails + thumbnails + "</div>"
}
