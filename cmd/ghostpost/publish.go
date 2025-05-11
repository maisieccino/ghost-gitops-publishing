// cmd/ghostpost/publish.go

package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
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

			// Compute SHA256 digest of Markdown body
			h := sha256.Sum256(md)
			nowHash := hex.EncodeToString(h[:])

			// If hash matches, skip publishing
			if meta.Hash == nowHash {
				fmt.Println("↻ no changes since last publish, skipping…")
				return nil
			}

			imgSvc := images.New(cfg.APIURL, cfg.AdminJWT, httpClient)
			md, _ = imgSvc.Rewrite(md, filepath.Dir(file))

			var html bytes.Buffer
			if err := goldmark.Convert(md, &html); err != nil {
				return err
			}

			client := api.New(cfg.APIURL, cfg.AdminJWT)

			// Map author names to IDs with error handling
			allAuthors, err := client.ListAuthors(context.Background())
			var authorIDs []string
			if err != nil {
				fmt.Println("warning: could not fetch authors from Ghost, using names as IDs")
				authorIDs = meta.Authors
			} else {
				nameToAuthorID := map[string]string{}
				for _, a := range allAuthors {
					nameToAuthorID[a.Name] = a.ID
				}
				for _, name := range meta.Authors {
					if id, ok := nameToAuthorID[name]; ok {
						authorIDs = append(authorIDs, id)
					}
				}
			}

			// Map tier names/slugs to TierRef (ID+Name+Slug)
			allTiers, err := client.ListTiers(context.Background())
			if err != nil {
				return fmt.Errorf("could not fetch tiers: %w", err)
			}
			byName := make(map[string]api.TierRef, len(allTiers))
			bySlug := make(map[string]api.TierRef, len(allTiers))
			for _, t := range allTiers {
				byName[t.Name] = t
				bySlug[t.Slug] = t
			}
			var tierRefs []api.TierRef
			for _, want := range meta.Tiers {
				if t, ok := byName[want]; ok {
					tierRefs = append(tierRefs, t)
				} else if t, ok := bySlug[want]; ok {
					tierRefs = append(tierRefs, t)
				} else {
					return fmt.Errorf("unknown tier %q (available: %v)", want, keys(byName))
				}
			}

			post := api.Post{
				Title:          meta.Title,
				Slug:           meta.Slug,
				Status:         defaultStatus(meta.Status),
				HTML:           html.String(),
				FeatureImage:   meta.FeatureImage,
				Tags:           api.WrapTags(meta.Tags),
				CustomExcerpt:  meta.CustomExcerpt,
				PublishedAt:    meta.PublishedAt,
				Visibility:     meta.Visibility,
				Tiers:          api.WrapTiers(tierRefs),
				Featured:       meta.Featured,
				Authors:        api.WrapAuthors(authorIDs),
				CustomTemplate: meta.CustomTemplate,
			}
			newID, err := api.Upsert(client, post, meta.PostID)
			if err != nil {
				return err
			}

			// Always refresh the post from Ghost so we get the real published_at + status
			ghostPost, err := client.GetPost(context.Background(), newID)
			if err != nil {
				return err
			}

			dirty := false
			if meta.PostID == "" {
				meta.PostID = newID
				dirty = true
			}
			if meta.PublishedAt != ghostPost.PublishedAt {
				meta.PublishedAt = ghostPost.PublishedAt
				dirty = true
			}
			if meta.Status != ghostPost.Status {
				meta.Status = ghostPost.Status
				dirty = true
			}
			// update meta.Authors with human-readable names from ghostPost
			var newAuthors []string
			for _, a := range ghostPost.Authors {
				newAuthors = append(newAuthors, a.Name)
			}
			if len(meta.Authors) != len(newAuthors) {
				meta.Authors = newAuthors
				dirty = true
			} else {
				for i := range meta.Authors {
					if meta.Authors[i] != newAuthors[i] {
						meta.Authors = newAuthors
						dirty = true
						break
					}
				}
			}
			// update meta.Tiers with human-readable names from ghostPost
			var newTiers []string
			for _, t := range ghostPost.Tiers {
				newTiers = append(newTiers, t.Name)
			}
			if !api.EqualStringSlices(meta.Tiers, newTiers) {
				meta.Tiers = newTiers
				dirty = true
			}
			// Always update hash after publish
			if meta.Hash != nowHash {
				meta.Hash = nowHash
				dirty = true
			}
			if dirty {
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

// Helper to list keys for error messages
func keys[K comparable, V any](m map[K]V) []K {
	out := make([]K, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
