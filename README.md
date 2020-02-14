## muikku

---

"muikku" is the Finnish name of a [tasty little fish](https://en.wikipedia.org/wiki/Coregonus_albula) and apparently what Finnish people say when [they smile on photos.](https://en.wikipedia.org/wiki/Say_cheese)

---

`muikku` is a configurable tool that will allow import of photos into custom workspaces and batch upload to photo communities.

This is my playground project to learn the [Go programming language.](https://github.com/golang/go)

It runs on `go1.13.8` and was built with the help of
[`cobra`.](https://github.com/spf13/cobra)

### Setup

It currently expects a config file at `$HOME/.muikku.yaml` like so:

```
import_path: "<path to import from (usually SD card)"
workspace: "<path to import to>"
collection_dirs:
  - "JPG"
  - "RAW"
  - "<optionable additional dirs>"
```

After that, just run `go build` and play.
