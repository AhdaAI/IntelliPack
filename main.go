package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/AhdaAI/IntelliPack/system"
)

func main() {
	start := time.Now()
	extract := flag.Bool("e", false, "Extract mode (Default: false)")
	input := flag.String("i", "", "Folder or File to compress.")
	output := flag.String("o", "", "Output compressed file (e.g. output.tar.zst).")
	outputDir := flag.String("dir", "", "Output directory (e.g. E:/Testing).")

	flag.Parse()
	if *input == "" {
		fmt.Println("Input is empty.")
		flag.Usage()
		return
	}

	if *outputDir == "" && *output == "" {
		fmt.Println("Output is empty, please specified either -o or -dir.")
		flag.Usage()
		return
	}

	if !*extract {
		out := *output
		if filepath.Ext(*output) != "" {
			out = *output + ".tar.zst"
		} else if *outputDir != "" {
			out = *outputDir + "/" + filepath.Base(*input) + ".tar.zst"
		}
		err := system.CompressFolder(*input, out)
		if err != nil {
			log.Fatal(err)
		}
		elapsed := time.Since(start)
		log.Printf("Compression took : %s", elapsed.Round(time.Second).String())
	} else {
		out := *output
		if filepath.Ext(*output) != "" {
			out = *output + ".tar.zst"
		} else if *outputDir != "" {
			out = *outputDir + "/" + strings.TrimSuffix(filepath.Base(*input), ".tar.zst")
		}
		err := system.ExtractArchive(*input, out)
		if err != nil {
			log.Fatalf("Extract failed: %v", err)
		}
		elapsed := time.Since(start)
		log.Printf("Extraction took : %s", elapsed.Round(time.Second).String())
	}
	log.Println("Done!")
}
