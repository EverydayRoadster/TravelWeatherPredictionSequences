package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/term"
)

func main() {

	var output string
	var model string
	var run int
	var dateS string = time.Now().Format("20060102")
	var parameter string
	var maxCount int

	flag.StringVar(&output, "output", ".meteociel/", "Output folder")
	flag.StringVar(&model, "model", "cfs", "Model name (e.g. cfs)")
	flag.IntVar(&run, "run", 1, "Model run 1-4")
	flag.StringVar(&dateS, "date", time.Now().Format("20060102"), "Run date (e.g. 20260219, default: current date)")
	flag.StringVar(&parameter, "parameter", "1", "Parameter (e.g. 0,1,2,5,9)")
	flag.IntVar(&maxCount, "max", 7296, "Max hours to download (default 7296)")
	flag.Parse()

	if output == "" || model == "" || run == 0 || dateS == "" || parameter == "" {
		fmt.Println("Usage:")
		fmt.Println("  go run main.go -output <folder> -model <model> -date <YYYYMMDDHH> -parameter <param>")
		return
	}

	dateS = dateS + fmt.Sprintf("%02d", (run-1)*6) // append 00,06,12,18 run ...

	baseURL := "https://modeles12.meteociel.fr/modeles"

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	saveDir := filepath.Join(output, model, dateS, parameter)

	err := os.MkdirAll(saveDir, 0755)
	if err != nil {
		fmt.Println("Failed creating directory:", err)
		return
	}

	totalFiles := maxCount / 6
	current := 0

	fallbackDate, err := time.Parse("2006010215", dateS)
	fallbackDate = fallbackDate.Add(-1 * 24 * 31 * time.Hour) // one month back
	currentDate := dateS

	for hour := 6; hour <= maxCount; hour += 6 {
		current++

		if term.IsTerminal(int(os.Stdout.Fd())) {
			fmt.Printf("\rDownloading %d / %d", current, totalFiles)
		} else {
			fmt.Printf("Downloading %d / %d\n", current, totalFiles)
		}

		url := fmt.Sprintf(
			"%s/%s/runs/%s/run%d/%s-%s-%d.png",
			baseURL,
			model,
			currentDate,
			run,
			model,
			parameter,
			hour,
		)

		filename := fmt.Sprintf("%04d.png", hour)
		savePath := filepath.Join(saveDir, filename)

		if _, err := os.Stat(savePath); err == nil {
			continue
		}

		resp, err := client.Get(url)
		if err != nil {
			fmt.Println("\nError:", err)
			continue
		}
		time.Sleep(250 * time.Millisecond) // be gentle to the server

		if resp.StatusCode != http.StatusOK {
			fmt.Println("\nSkipping (not found):", hour)
			resp.Body.Close()
			continue
		}

		out, err := os.Create(savePath)
		if err != nil {
			fmt.Println("\nError creating file:", err)
			resp.Body.Close()
			continue
		}

		written, err := io.Copy(out, resp.Body)

		out.Close()
		resp.Body.Close()

		if err != nil {
			fmt.Println("\nError saving file:", err)
			continue
		}

		if written == 0 {
			os.Remove(savePath) // delete empty file
			currentD, err := time.Parse("2006010215", currentDate)

			if !currentD.Before(fallbackDate) { // switch date to previous
				currentDate, err = previousRunDate(currentDate)
				if err != nil {
					break
				}
				hour -= 6 // retry same hour
				current--
				continue
			}

			fmt.Println("\nFallback already used but still empty file received. Aborting.")
			break
		}
	}

	fmt.Printf(" Download complete. %d Files downloaded", current-1)
	fmt.Println()
}

func previousRunDate(date string) (string, error) {
	t, err := time.Parse("2006010215", date)
	if err != nil {
		return "", err
	}
	t = t.Add(-24 * time.Hour)
	return t.Format("2006010215"), nil
}
