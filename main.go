package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/asticode/go-astisub"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "vtt2ass"
	app.Version = "0.0.1"
	app.Author = "Raincal"
	app.Email = "cyj94228@gmail.com"
	app.Usage = "Convert vtt to ass"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "output, o",
			Usage: "options: output path",
		},
	}

	app.Action = func(c *cli.Context) error {
		inputPath := c.Args().Get(0)
		outputPath := c.String("output")

		if inputPath == "" {
			fmt.Println("Please enter a input path!")
			return nil
		}

		dir, file := path.Split(inputPath)
		fileSuffix := path.Ext(file)
		filename := strings.TrimSuffix(file, fileSuffix)

		if fileSuffix != ".vtt" {
			fmt.Println("Only support .vtt file!")
			return nil
		}

		vttFile, err := astisub.OpenFile(inputPath)
		if err != nil {
			fmt.Printf("err was %v", err)
			return nil
		}

		// Remove the styling from YouTube
		pattern := "</c>|<c.color\\S{6}>|<\\d+:\\d+:\\d+.\\d+><c>"
		re := regexp.MustCompile(pattern)

		for _, item := range vttFile.Items {
			for idxLine, line := range item.Lines {
				for idxLineItem, lineItem := range line.Items {
					item.Lines[idxLine].Items[idxLineItem].Text = re.ReplaceAllString(lineItem.Text, "")
				}
			}
		}

		if outputPath == "" {
			outputPath = dir + filename + ".ass"
		}

		err = vttFile.Write(outputPath)
		if err != nil {
			fmt.Printf("err was %v", err)
			return nil
		}

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
