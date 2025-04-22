# Image Compare Tool

This project is a command-line tool written in Go for comparing two images (PNG, JPEG, SVG) and highlighting their differences. The tool allows users to specify a tolerance level for comparison, enabling flexible image analysis. It also supports cleaning up intermediate files after comparison.

## Features

- Compare two images and determine if they are equal.
- Support for PNG, JPEG, and SVG formats (SVGs are converted to PNGs using a headless browser).
- Set a tolerance level (in percent) for comparison.
- Highlight differences between images in red.
- Clean up intermediate files (e.g., converted PNGs) using the `--clean` flag.
- Output verdicts based on comparison results.

## Installation

1. Clone the repository:

   ```
   git clone https://github.com/microsoft/vscode-remote-try-go.git
   ```

2. Navigate to the project directory:

   ```
   cd image-compare-tool
   ```

3. Build the tool:

   ```
   go build -o image-compare-tool ./cmd
   ```

## Usage

Run the tool from the terminal with the following command:

```
./image-compare-tool <path_to_base_image> <path_to_compare_image> <tolerance_level>
```

### Parameters

- `--base`: The file path to the base image.
- `--compare`: The file path to the image to compare against the base image.
- `--tolerance`: The tolerance level for comparison (in percent).
- `--clean`: Optional flag to remove intermediate files (e.g., converted PNGs images) after comparison.

### Example

```
./image-compare-tool --base image1.png --compare image2.png --tolerance 5
```

## Verdicts

- If images are equal:
  - Output: `Images are equal.`

- If images are equal within the tolerance:
  - Output: `Images are equal with tolerance {passed tolerance value}%.`

- If images are not equal and exceed the tolerance:
  - A new image will be generated highlighting the differences.
  - Output: `Images are not equal. Differences highlighted in: <path_to_highlighted_image>.`

- If images are not equal but within the tolerance:
  - Output: `Images are equal with tolerance {passed tolerance value}%, but not equal in general.`

## SVG Support

- SVG files are automatically converted to PNG using a headless browser (via `chromedp`).
- Ensure that Google Chrome or Chromium is installed on your system for SVG support.
- If the `--clean` flag is passed, the converted PNG files will be removed after comparison.

## Cleaning Up

- If the `--clean` flag is passed:
  - Converted PNG files (from SVGs) will be removed after comparison.

## License

This project is licensed under the MIT License. See the LICENSE file for details.