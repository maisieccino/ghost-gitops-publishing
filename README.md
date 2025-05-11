# ghostpost — Git-first publishing to Ghost

Ghost gives you a fast writing UI.
Git gives you real history, branches, reviews, CI.

`ghostpost` glues them together.

You keep every post as a Markdown file with front-matter.
`ghostpost` turns that file into a Ghost draft (or update) with one command.


## Why you might care

* **Version control**
  See every change in `git log`, run spell-check in CI, review with pull requests.

* **Stateless deploys**
  No local caches. The post ID lives in the Markdown front-matter.

* **Images on the fly**
  Local image paths become Ghost URLs automatically.

* **No runtime deps**
  A single static Go binary for macOS, Linux, or Windows.


## Install

```bash
go install github.com/rodchristiansen/ghost-gitops-publishing/cmd/ghostpost@latest
```

The binary lands in your `$GOPATH/bin` (often `~/go/bin`).
Add that directory to `$PATH` if needed.


## Setup

Create `~/.ghostpost/config.yaml` (or keep it in the repo).

```yaml
api_url:  https://your-site.ghost.io/ghost/api/admin/
admin_jwt:  123abc456def:deadbeefcafef00d...   # Admin key or signed JWT
```

* You can paste the raw **Admin API key**.
  ghostpost signs it for you.
* The trailing slash in `api_url` is needed.


## Your first post

`welcome.md`

```md

title: Welcome to Focused Systems
slug: welcome-to-focused-systems
tags: [DevOps, CI/CD]
status: draft            # default is draft
custom_excerpt: Why this blog exists

Hello world…
```

Publish:

```bash
ghostpost publish -f welcome.md
```

ghostpost:

1. Uploads any local images it finds.
2. Converts Markdown → HTML with Goldmark.
3. Creates the post in Ghost.
4. Writes the returned `post_id` back into the front-matter.

The file now looks like:

```yaml
post_id: 681fcffa6cf6ba0001ccf0e9
```

Commit that change—now the ID tracks with your content.



## Fix a typo

Edit the file, run the same command again.

ghostpost:

* Pulls the current `updated_at` timestamp.
* Sends a `PUT /posts/{id}` with that lock.
* Ghost patches the post.



## Jump straight to the editor

```bash
ghostpost publish -f welcome.md --editor
# or
ghostpost publish -f welcome.md -e
```

Your browser opens:

```
https://your-site.ghost.io/ghost/#/editor/post/681fcffa6cf6ba0001ccf0e9
```



## Front-matter keys

| key              | purpose                             |
| - | -- |
| `title`          | Post title                          |
| `slug`           | URL slug (optional)                 |
| `tags`           | Array or YAML list                  |
| `feature_image`  | Path or URL                         |
| `status`         | `draft`, `published`, `scheduled`   |
| `published_at`   | ISO date string for schedule        |
| `custom_excerpt` | Manual excerpt                      |
| `post_id`        | Added by ghostpost after first push |



## CI example

```yaml
name: Publish
on:
  push:
    branches: [main]
    paths: ["posts/**.md"]

jobs:
  ghost:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with: {go-version: '1.22'}
      - run: go install github.com/rodchristiansen/ghost-gitops-publishing/cmd/ghostpost@latest
      - run: |
          for f in posts/*.md; do
            ghostpost publish -f "$f"
          done
        env:
          GHOST_API_URL:  ${{ secrets.GHOST_API_URL }}
          GHOST_ADMIN_JWT: ${{ secrets.GHOST_ADMIN_JWT }}
```



## Roadmap

* Tag management (`ghostpost tags …`)
* Image garbage-collection
* Live preview server
* Windows cross-compile releases

## Contributing

Pull requests and issues welcome.

* Run `go vet ./...` and `go test ./...` before pushing.
* Keep sentences short.
  Write like a friend who figured something out.

## License

MIT. See `LICENSE` file.
