package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"golang.org/x/term"
)

func main() {

	var output string
	var model string
	var run int
	var dateS string = time.Now().Format("20060102")
	var mode string
	var maxCount int

	flag.StringVar(&output, "output", ".meteociel/", "Output folder")
	flag.StringVar(&model, "model", "cfs", "Model name (e.g. cfs)")
	flag.IntVar(&run, "run", 1, "Model run 1-4")
	flag.StringVar(&dateS, "date", time.Now().Format("20060102"), "Run date (e.g. 20260219, default: current date)")
	flag.StringVar(&mode, "mode", "0", "Mode (e.g. 0,1,2,5,9) for subject of calculation")
	flag.IntVar(&maxCount, "max", 7296, "Max hours to download (default 7296)")
	flag.Parse()

	if output == "" || model == "" || run == 0 || dateS == "" || mode == "" {
		fmt.Println("Usage:")
		fmt.Println("  go run main.go -output <folder> -model <model> -date <YYYYMMDDHH> -mode <mode>")
		return
	}

	dateS = dateS + fmt.Sprintf("%02d", (run-1)*6) // append 00,06,12,18 run ...

	baseURL := "https://modeles12.meteociel.fr/modeles"

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	saveDir := filepath.Join(output, model, dateS, mode)

	err := os.MkdirAll(saveDir, 0755)
	if err != nil {
		fmt.Println("Failed creating directory:", err)
		return
	}

	totalFiles := (maxCount / 6)
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

		filename := fmt.Sprintf("%04d.png", hour)
		savePath := filepath.Join(saveDir, filename)

		if _, err := os.Stat(savePath); err == nil {
			continue
		}

		url := fmt.Sprintf(
			"%s/%s/runs/%s/run%d/%s-%s-%d.png",
			baseURL,
			model,
			currentDate,
			run,
			model,
			mode,
			hour,
		)

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

	fmt.Printf("Download complete. %d Files downloaded", current-1)
	fmt.Println()

	switch mode {
	case "9":
		for i := 0; i < 4; i++ {
			err = renderVideo(output, saveDir, model, dateS, mode, "3", i)
			if err != nil {
				break
			}
		}

	default:
		err = renderVideo(output, saveDir, model, dateS, mode, "12", -1)
	}
	if err != nil {
		fmt.Println("Video generation complete.")
	}
}

func previousRunDate(date string) (string, error) {
	t, err := time.Parse("2006010215", date)
	if err != nil {
		return "", err
	}
	t = t.Add(-24 * time.Hour)
	return t.Format("2006010215"), nil
}

func renderVideo(output, saveDir, model, dateS, mode, framerate string, interleave int) error {
	ffmpeg := "ffmpeg"
	_, err := exec.LookPath(ffmpeg)
	if err != nil {
		fmt.Println(ffmpeg, " not installed, skipping video generation")
		return nil
	}
	cmd := exec.Command(
		ffmpeg,
		"-framerate", framerate,
		"-pattern_type", "glob",
		"-i", saveDir+"/*.png",
		"-c:v", "libx264",
		"-pix_fmt", "yuv420p",
		"-y")
	if interleave > -1 {
		cmd.Args = append(cmd.Args, "-vf", "select='eq(mod(n\\,4)\\,"+strconv.Itoa(interleave)+")',setpts=N/FRAME_RATE/TB")
		cmd.Args = append(cmd.Args, filepath.Join(output, model+"_"+dateS+"_"+mode+"_"+strconv.Itoa(interleave)+".mp4"))
	} else {
		cmd.Args = append(cmd.Args, filepath.Join(output, model+"_"+dateS+"_"+mode+".mp4"))
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
