# Make Meme Text Searchable

I have an extensive set of memes I've been collecting since the early days of Flickr. #icanhascheeseburger

It's a pain in the ass to not be able to search my memes (now stored in Apple’s _Photos.app_) to find what I'm looking for when I need it.

<div align="center"><img src="images/had-something2.jpg" alt="I had something for this."></div>

This project uses [Amazon Rekognition] to read the text from the images, then [go-exif](https://github.com/dsoprea/go-exif) to write the text into the image metadata as a caption/description. Photo apps and services _should_ be able to parse and index this data, making your images (memes, really) searchable by the text that's in the image.

## General Program Flow

1. Read the binary image data into memory.

1. Rekognition only supports PNG and JPEG formats, so…

    1. If the image is _already_ a PNG or JPEG, skip to the next step.

    1. If the image is a GIF, WEBP, or HEIF format, convert the in-memory representation of the image to JPEG format.

1. Submit the (PNG or JPEG) bytes of the image to Rekognition.

1. Get back the results. Merge, deduplicate, and munge the resulting text matches, whatever they are. Words don't always come back in the right order, so think of them less as a _sentence_ and more of a _collection of keywords_.

    > Wait a minute, I had something for this.

    …might become…

    > a for had i minute something this wait

1. Using the EXIF library, write these words into the file (the one we read into memory), into the `ImageDescription` EXIF field of the metadata.

1. Optionally, we can:

    1. Overwrite the original file with an updated description (destructive).

    1. Read the file from one location, and write an updated copy to a new location (non-destructive).

        1. This new location can even be a different format, such as JPEG or PNG. Whatever Go's standard library supports.

## Why Go?

[Go] (aka, "Golang") compiles down to a static binary, and is stupidly fast. It can also compile to [WebAssembly], which means it can run in [Node.js] or web browsers. It also has the fastest boot time on AWS Lambda (more-or-less tied for first place with Rust), so creating a `POST` endpoint should be easy as well.

Someday, I want to learn how to develop mobile apps so that I can solve this user problem. [Compiled Go code can be called from Android and iOS](https://github.com/golang/go/wiki/Mobile).

## Progress

### Library

* [X] Importable with `go get`.
* [X] Handles image bytestreams and decoding.
* [X] Handles converting the image to JPEG before passing to Rekognition.
* [X] Sends data to Rekognition.
* [X] Parses the results from Rekognition into words.
* [X] Converts the image bytestream to JPEG before sending to Rekognition.
* [ ] Preserves any existing EXIF data.
* [ ] Writes the words into the EXIF data.

### CLI Tool

* [X] Supports reading a file.
* [X] Supports reading a directory.
* [X] Supports reading a glob.
* [ ] Supports verbose logging.
* [X] Supports AWS credentials as environment variables.
* [ ] Supports AWS credentials as a profile reference.
* [X] Supports `-v`.
* [ ] Supports `-vv` and `-vvv`.
* [ ] Supports `-q`.
* [ ] Supports outputting a copy to a new directory.
* [ ] Supports outputting a copy in a new format.
* [ ] Supports writing the Rekognition results into EXIF data at all.
* [ ] Supports status updates for jobs.
* [ ] Supports an index of already-processed images to facilitate restarting a failed queue.

## Usage

### Library

Incomplete example. Error handling removed for brevity.

```go
import "github.com/skyzyx/make-meme-text-searchable/meme"

func main() {
    // Open the file as am io.Reader.
    r, _ := os.Open("./images/paris-airport.heic")

    // Read the io.Reader, decode the image, then re-encode the image data as JPEG format.
    buf, _, _, _ := meme.ReadImage(r, meme.DefaultJPEGQuality)

    // Pass the image data to AWS Rekognition.
    results, _ := meme.DetectText(ctx, &awsConfig, buf)

    // Sanitize, de-dupe, remove punctuation, and sort the resulting words.
    words := meme.GetSanitizedText(results)

    // Write the string of words back to the image into the EXIF ImageDescription field.
    _ := meme.WriteImageDescription(r, words)
}
```

### CLI

(Brainstorming) Something like…

```bash
meme-text [--report=TEXT|JSON] [--out=FILE] [--outdir=DIR] [--outformat=GIF|HEIC|JPG|PNG|WEBP] [--quiet] [--verbose] [--force] INPUT...
```

* `INPUT` is one or more files, directories of files, or globs of files. Supports: GIF, HEIC, JPEG, PNG, WEBP. Also works with `STDIN`.
* `--report` will write data to `STDOUT` in the specified format.
* `--quiet` will silence all output.
* `--verbose` maybe be specified up to 3 times, with increasing levels of verbosity. The default value is equivalent to WARNING. -v, -vv, and -vvv are equivalent to INFO, DEBUG, and TRACE (respectively).
* `--force` disables any interactive prompts.

### Web UI

This should be relatively simple to write as long as people can drag-and-drop/upload their images into the webpage, then provide an email address to send the results to (asychronously). A simple desktop app is also possible — maybe with [Electron](https://www.electronjs.org), [Wails](https://wails.io), or [Tauri](https://tauri.studio)?

## Things to read

Things I need to read and understand. Apparently _writing_ EXIF data can be non-trivial.

* <https://en.wikipedia.org/wiki/Exif>
* <https://pkg.go.dev/github.com/dsoprea/go-jpeg-image-structure#example-SegmentList-SetExif>
* <https://github.com/dsoprea/go-jpeg-image-structure/blob/502ec55f4b9caf576fc0ba8429450022a9fc3285/jpeg_test.go#L179>
* <https://github.com/dsoprea/go-jpeg-image-structure/blob/eac4d3269d730aa516c7b18d4d396ba304d4bfd8/jpeg.go#L411>
* <https://go.dev/blog/defer-panic-and-recover>

  [Amazon Rekognition]: https://aws.amazon.com/rekognition/
  [Go]: https://go.dev
  [Node.js]: https://nodejs.org
  [WebAssembly]: https://webassembly.org
