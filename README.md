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

## Usage

No code is written yet, but I intend to have both an importable Go library, as well as a standalone CLI tool.

### Library

[Error handling removed for brevity.]

(Brainstorming) Something like…

```go
import "github.com/skyzyx/make-meme-text-searchable/meme"

func main() {
    files, _ := os.ReadDir(".")

    for _, filename := range files {
        contentAsBytes, _ := os.ReadFile(filename)
        wordSalad, _ := meme.Rekognize(meme.FormatAsJPEG, contentAsBytes)
        keywords := meme.DedupeAndSort(wordSalad)
        _ := meme.WriteKeywords(filename)
    }
}
```

### CLI

(Brainstorming) Something like…

```bash
make-meme-text-searchable \
    --in-dir . \
    --include *.png \
    --out-dir ~/Desktop/searchable-memes \
    --out-format jpg
```

### Web UI

This should be relatively simple to write as long as people can drag-and-drop/upload their images into the webpage, then provide an email address to send the results to (asychronously). A simple desktop app is also possible — maybe with [Electron](https://www.electronjs.org), [Wails](https://wails.io), or [Tauri](https://tauri.studio)?

## Links

For when I sit down to write actual code.

* <https://aws.amazon.com/rekognition/>
* <https://aws.amazon.com/rekognition/pricing/>
* <https://github.com/dsoprea/go-exif>
* <https://github.com/sfomuseum/go-exif-update>
* <https://github.com/jdeng/goheif>
* <https://pkg.go.dev/image/gif#Decode>
* <https://pkg.go.dev/image/png#Decode>
* <https://pkg.go.dev/image/jpeg#Decode>
* <https://pkg.go.dev/golang.org/x/image/webp#Decode>

  [Amazon Rekognition]: https://aws.amazon.com/rekognition/
  [Go]: https://go.dev
  [Node.js]: https://nodejs.org
  [WebAssembly]: https://webassembly.org
