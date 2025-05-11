# Ghost GitOps Publishing

Draft, update, and publish Ghost posts via a GitOps Markdown workflow.

Ghost gives you the blogging platform.  
Git gives you version control, history, branches, reviews, cloud pipelines.

`ghostpost` glues them together.

You keep your post as a Markdown file with front-matter.  
One command turns it into a Ghost draft—or updates an existing post.

## Features

### Front-matter driven

Define title, slug, tags, status, excerpt, schedule, visibility, authors, templates, and more.  
`ghostpost` reads and writes the `post_id` for you.

### Idempotent updates

`ghostpost` fetches the current `updated_at` lock and issues a `PUT`.  
Your edits replace the previous version so you are always in sync.

### Open in online CMS editor after publishing

```bash
ghostpost publish -f post.md --editor
```

Opens your browser at `/ghost/#/editor/post/{post_id}`

## Why `ghostpost`

What if you could manage your blog like code?

### Version control

- Track every edit in `git log`
- Run spell-check in CI
- Review via pull requests

### Stateless deploys

- No local state  
- The `post_id` lives in your front-matter  
- Everything self-contained in the `.md` article

### Automatic images

- Reference local image paths  
- `ghostpost` uploads and updates URLs

### Zero dependencies

- One static Go binary for macOS, Linux, Windows

## Setup

Create `~/.ghostpost/config.yaml` (or keep it in the repo):

```yaml
api_url:   https://your-site.ghost.io/ghost/api/admin/
admin_jwt: 123abc456def:deadbeefcafef00d...   # Admin key or signed JWT
```

- You can paste the raw **Admin API key**; `ghostpost` will auto-sign it.  
- The trailing slash in `api_url` is required.

## Your first post

`welcome.md`:

```md
---
title: Welcome to Focused Systems
slug: welcome-to-focused-systems
custom_excerpt: Insights, guides, and updates on technology and workflows
tags:
  - DevOps
  - CI/CD
feature_image: assets/hero.png      # relative path or URL to your cover image
status: draft                       # draft | published | scheduled
published_at: 2025-06-15T09:00:00Z  # ISO timestamp for scheduling (optional)
visibility: public                  # public | members | paid | specific
tiers:
  - free                            # which paid tiers can see this (optional)
featured: false                     # show as “featured” in your theme?
authors:
  - rodchristiansen                 # Ghost user slugs
custom_template: post               # choose a custom template (optional)
post_id: ""                         # filled in by ghostpost after first publish
---

## Hello world

Welcome to my very first post—managed entirely via GitOps!  
```

Publish:

```bash
ghostpost publish -f welcome.md
```

`ghostpost` will:

1. Upload any local images it finds.  
2. Convert Markdown → HTML with Goldmark.  
3. Create the post in Ghost.  
4. Write the returned `post_id` back into the front-matter.

Your file now contains:

```yaml
post_id: 681fcffa6cf6ba0001ccf0e9
```

Commit that change—now the ID tracks with your content.

## Fix a typo

Edit the file, run the same command again.

`ghostpost`:

- Pulls the current `updated_at` timestamp.  
- Sends a `PUT /posts/{id}` with that lock.  
- Ghost patches the post.

## Jump straight to the editor

```bash
ghostpost publish -f welcome.md --editor
```

Your browser opens:  
`https://your-site.ghost.io/ghost/#/editor/post/681fcffa6cf6ba0001ccf0e9`

## Front-matter keys

| Key               | Purpose                                  |
| ----------------- | ---------------------------------------- |
| `title`           | Post title                               |
| `slug`            | URL slug (optional)                      |
| `tags`            | Array or YAML list                       |
| `feature_image`   | Path or URL                              |
| `status`          | `draft`, `published`, `scheduled`        |
| `published_at`    | ISO date string for scheduling           |
| `visibility`      | `public`, `members`, `paid`, `specific`  |
| `tiers`           | Array of paid tiers (for `specific`)     |
| `featured`        | `true`/`false` to feature the post       |
| `custom_excerpt`  | Manual excerpt                           |
| `authors`         | Array of author names                    |
| `custom_template` | Template name (e.g. `Full Feature Image`)|
| `post_id`         | Populated by `ghostpost` after first push|

## CI example

```yaml
name: Publish
on:
  push:
    branches: [main]
    paths: ["posts/**/*.md"]

jobs:
  ghost:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.24.2'
      - run: go install github.com/rodchristiansen/ghost-gitops-publishing/cmd/ghostpost@latest
      - run: |
          for f in posts/*.md; do
            ghostpost publish -f "$f"
          done
        env:
          GHOST_API_URL:   ${{ secrets.GHOST_API_URL }}
          GHOST_ADMIN_JWT: ${{ secrets.GHOST_ADMIN_JWT }}
```

## Contributing

Pull requests and issues welcome!

- Run `go vet ./...` and `go test ./...` before pushing.  
- Keep sentences short; write like a friend who figured something out.

## License

MIT. See the `LICENSE` file.

> Inspired by the “Articles as Code” idea from [post2ghost](https://www.how-hard-can-it.be/post2ghost/)  
