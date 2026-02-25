# Travel Weather Prediction Sequences

Travel Weather Prediction Sequences downloads images from metereological services like meteociel, to have them combined into an animation.

## Features

* **Download** – Fetch PNGs from weather prediction calculation, for a specific model, run and date range.
* **Organise** – Files are stored under `output/<model>/<date><run>/<mode>/`.
* **Video** – Renders a video from the images; mode `9` produces four interleaved videos, any other mode produces a single video.
* **Cache** – Existing images are reused; no re‑download unless the cache is cleared.
* **Cleanup** – By default stale directories (of previous days) for the current model are removed before a new run.

## Usage

```bash
# Basic run (downloads and renders video)
go run github.com/EverydayRoadster/TravelWeatherPredictionSequences@latest \
  -output ./output \

# Skip the automatic cleanup of old directories
# use the 12'o clock run 
# produce images and animation of the prediction for JetStream
go run github.com/EverydayRoadster/TravelWeatherPredictionSequences@latest \
  -output ./output \
  -run 2 \
  -mode 5 \
  -noclean
```

| Flag      | Description |
|-----------|-------------|
| `-output` | Path to the base folder where the data will be stored. Default: `.meteociel/`. |
| `-model`  | Name of the model to query (`cfs` is the default and only model supported right now). |
| `-run`    | Run number of model calculation 1-00h (default), 2-06h, 3-12h, 4-18h |
| `-mode`   | Mode of Parameters: 0 - Geopotential Height at 500 hPa, 1 - Temperature at 850 hPa, 2 - Precipitation, 5 - Jet Stream, 9 - Temperature at 2 meters |
| `-noclean`| **New flag** – when present the tool will *not* delete old directories before starting a new run. |

## License

MIT © 2026
