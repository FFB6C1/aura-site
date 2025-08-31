package builder

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ffb6c1/aura-site/internal/config"
	"github.com/ffb6c1/aura-site/internal/file"
	"github.com/ffb6c1/aura-site/internal/markdown"
)

type gallery struct {
	path    string
	imgData []img
}

type img struct {
	name        string
	filename    string
	thumb       string
	alt         string
	description string
}

func checkGallery(page string) (string, []gallery) {
	if !strings.Contains(page, "!!!gallery") {
		return markdown.Convert(page, "builder"), nil
	}

	splitPage := strings.Split(page, "!!!gallery[")
	galleries := []gallery{}
	newPage := markdown.Convert(splitPage[0], "builder")
	for _, splitItem := range splitPage[1:] {
		galleryDir, restOfItem, _ := strings.Cut(splitItem, "]")
		if galleryDir == "" {
			newPage += markdown.Convert(restOfItem, "builder")
			continue
		}
		galHTML, gal := makeGallery(galleryDir)
		newPage += galHTML + markdown.Convert(restOfItem, "builder")
		galleries = append(galleries, gal)
	}

	return newPage, galleries
}

func makeGallery(galleryDir string) (string, gallery) {
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

	galHTML := imagesToGallery(galleryDir, fullImgData)

	return galHTML, gallery{
		path:    galleryDir,
		imgData: fullImgData,
	}
}

func finishGalleries(galleries []gallery) {
	for _, gal := range galleries {
		if err := makePagesForImages(gal.path, gal.imgData); err != nil {
			fmt.Printf("could not make pages/could not copy images for %s", gal.path)
		}
	}
}

func sortImages(allImages [][]string) (map[string]string, map[string]string) {
	images := make(map[string]string)
	thumbs := make(map[string]string)
	for _, img := range allImages {
		if name, thumb := strings.CutSuffix(img[0], "-thumb"); thumb {
			thumbs[name] = img[0] + "." + img[1]
			continue
		}
		images[img[0]] = img[0] + "." + img[1]
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

func makePagesForImages(galleryName string, imgData []img) error {
	config := config.GetConfig()
	galleryPath := filepath.Join(config.GetExportPath(), galleryName)
	srcPath := filepath.Join(config.GetImportPath(), galleryName)
	if err := file.MakeDirectory(galleryPath); err != nil {
		return err
	}

	pages := map[string]string{}

	for _, image := range imgData {
		page := fmt.Sprintf("<img src=\"images/%s\" alt=\"%s\" class=\"galleryImg\">\n", image.filename, image.alt)
		if image.description != "" {
			page += markdown.Convert(image.description, "builder")
		}
		pages[image.name] = page
	}
	if err := createPages(pages, galleryPath); err != nil {
		return err
	}

	if err := copyImages(srcPath, galleryPath, imgData); err != nil {
		return err
	}

	css := config.GetCSS()

	if err := file.WriteFileFromString(filepath.Join(galleryPath, "styles.css"), css); err != nil {
		return err
	}

	return nil
}

func copyImages(srcPath, dstPath string, imgData []img) error {
	imgPath, thumbsPath, err := makeImgDirs(dstPath)
	if err != nil {
		return err
	}

	for _, image := range imgData {
		if image.thumb != "" {
			thumb := filepath.Join(srcPath, image.thumb)
			newThumb := filepath.Join(thumbsPath, image.thumb)
			if err := file.CopyFile(thumb, newThumb); err != nil {
				fmt.Println(err)
				return err
			}
		}
		imagePath := filepath.Join(srcPath, image.filename)
		newImagePath := filepath.Join(imgPath, image.filename)
		if err := file.CopyFile(imagePath, newImagePath); err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}

func makeImgDirs(path string) (string, string, error) {
	thumbsPath := filepath.Join(path, "thumbs")
	imgPath := filepath.Join(path, "images")

	if err := file.MakeDirectory(thumbsPath); err != nil {
		fmt.Println("error making directory")
		return "", "", err
	}
	if err := file.MakeDirectory(imgPath); err != nil {
		fmt.Println("error making directory")
		return "", "", err
	}
	return imgPath, thumbsPath, nil
}
