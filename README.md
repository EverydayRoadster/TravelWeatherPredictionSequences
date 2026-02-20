# Travel Weather Prediction Sequences

Travel Weather Prediction Sequences downloads images from metereological services like meteociel, to have them later combined into an animation.

## Usage

Travel Weather Prediction Sequences requires go installed on a computer.

go run github.com/EverydayRoadster/TravelWeatherPredictionSequences@latest

If no arguments are specified, the program will download images into folder ".meteociel/" .

Program arguments available:

- output - Output folder, defaults to .meteociel/"

- model - defaults to cfs

- run - selects the computation run 1-00h (default), 2-06h, 3-12h, 4-18h
- date - computation run date in format ("20060102"), defaults to current date. There is a transparent daily fallback to dates up to on month backwards, as to model in itself may publish model data from previous runs
- mode - 0 - Geopotential Height at 500 hPa, 1 - Temperature at 850 hPa, 2 - Precipitation, 5 - Jet Stream, 9 - Temperature at 2 meters,
- max maximum hours to predict, defaults to 7296. Download will stop earlier if no more images are available

The program will cache files previously downlaoded.
