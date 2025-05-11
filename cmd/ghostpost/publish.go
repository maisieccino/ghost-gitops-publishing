// cmd/ghostpost/publish.go

package main

import (
	"bytes"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/rodchristiansen/ghost-gitops-publishing/internal/api"
	"github.com/rodchristiansen/ghost-gitops-publishing/internal/frontmatter"
	"github.com/rodchristiansen/ghost-gitops-publishing/internal/images"

	"github.com/spf13/cobra"
	"github.com/yuin/goldmark"
)

var httpClient = &http.Client{Timeout: 30 * time.Second}

func defaultStatus(s string) string {
	if s == "" {
		return "draft"
	}
	return s
}

func publishCmd() *cobra.Command {
	var file string
	var openEditor bool

	cmd := &cobra.Command{
		Use:   "publish",
		Short: "Push Markdown → Ghost",
		RunE: func(_ *cobra.Command, _ []string) error {
			meta, md, err := frontmatter.ParseFile(file)
			if err != nil {
				return err
			}

			imgSvc := images.New(cfg.APIURL, cfg.AdminJWT, httpClient)
			md, _ = imgSvc.Rewrite(md, filepath.Dir(file))

			var html bytes.Buffer
			if err := goldmark.Convert(md, &html); err != nil {
				return err
			}

			post := api.Post{
				Title:         meta.Title,
				Slug:          meta.Slug,
				Status:        defaultStatus(meta.Status),
				HTML:          html.String(),
				FeatureImage:  meta.FeatureImage,
				Tags:          api.WrapTags(meta.Tags),
				CustomExcerpt: meta.CustomExcerpt,
				PublishedAt:   meta.PublishedAt,
			}
			client := api.New(cfg.APIURL, cfg.AdminJWT)
			newID, err := api.Upsert(client, post, meta.PostID)
			if err != nil {
				return err
			}

			if meta.PostID == "" { // first publish → persist ID
				meta.PostID = newID
				if err := frontmatter.WriteFile(file, meta, md); err != nil {
					return err
				}
			}

			if openEditor {
				// strip trailing "/ghost/api/admin/" → siteRoot
				siteRoot := strings.Split(cfg.APIURL, "/ghost/")[0]
				url := fmt.Sprintf("%s/ghost/#/editor/post/%s", siteRoot, meta.PostID)
				_ = launchBrowser(url)
			}
			return nil

		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "", "Markdown file")
	cmd.MarkFlagRequired("file")
	cmd.Flags().BoolVarP(&openEditor, "editor", "e", false, "Open post in Ghost editor")
	return cmd
}
