package anigate

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func (s *Service) fsStat(args map[string]any) (map[string]any, error) {
	rp, err := s.policy.resolve(stringArg(args, "workspace"), stringArgDefault(args, "path", "."))
	if err != nil {
		return nil, err
	}
	info, err := os.Lstat(rp.Abs)
	if err != nil {
		return nil, err
	}
	return statMap(rp, info), nil
}

func (s *Service) fsTree(args map[string]any) (map[string]any, error) {
	rp, err := s.policy.resolve(stringArg(args, "workspace"), stringArgDefault(args, "path", "."))
	if err != nil {
		return nil, err
	}
	depth := intArgDefault(args, "depth", 2)
	if depth < 0 || depth > 8 {
		depth = 2
	}
	maxEntries := intArgDefault(args, "max_entries", 200)
	if maxEntries <= 0 || maxEntries > 2000 {
		maxEntries = 200
	}
	count := 0
	truncated := false
	tree, err := s.buildTree(rp.Workspace.Name, rp.Abs, rp.Rel, depth, maxEntries, &count, &truncated)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"workspace": rp.Workspace.Name,
		"path":      rp.Rel,
		"tree":      tree,
		"count":     count,
		"truncated": truncated,
	}, nil
}

func (s *Service) buildTree(workspace, abs, rel string, depth, maxEntries int, count *int, truncated *bool) (map[string]any, error) {
	if *count >= maxEntries {
		*truncated = true
		return nil, nil
	}
	info, err := os.Lstat(abs)
	if err != nil {
		return nil, err
	}
	*count++
	node := map[string]any{
		"name":     filepath.Base(abs),
		"path":     filepath.ToSlash(rel),
		"is_dir":   info.IsDir(),
		"size":     info.Size(),
		"modified": info.ModTime().UTC().Format(time.RFC3339),
	}
	if !info.IsDir() || depth == 0 {
		return node, nil
	}
	entries, err := os.ReadDir(abs)
	if err != nil {
		return node, nil
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].IsDir() != entries[j].IsDir() {
			return entries[i].IsDir()
		}
		return entries[i].Name() < entries[j].Name()
	})
	var children []map[string]any
	for _, entry := range entries {
		if shouldSkipTreeEntry(entry.Name()) {
			continue
		}
		childRel := filepath.Join(rel, entry.Name())
		child, err := s.buildTree(workspace, filepath.Join(abs, entry.Name()), childRel, depth-1, maxEntries, count, truncated)
		if err != nil {
			return nil, err
		}
		if child != nil {
			children = append(children, child)
		}
		if *truncated {
			break
		}
	}
	node["children"] = children
	return node, nil
}

func (s *Service) fsWritePreview(args map[string]any) (map[string]any, error) {
	workspace := stringArg(args, "workspace")
	if err := s.workspaceAllows(workspace, "write"); err != nil {
		return nil, err
	}
	rp, err := s.policy.resolve(workspace, stringArg(args, "path"))
	if err != nil {
		return nil, err
	}
	content := stringArg(args, "content")
	if int64(len(content)) > s.cfg.MaxReadBytes {
		return nil, fmt.Errorf("content exceeds max_read_bytes")
	}
	old := ""
	if b, err := os.ReadFile(rp.Abs); err == nil {
		if looksBinary(b) {
			return nil, fmt.Errorf("existing file appears binary")
		}
		old = string(b)
	} else if !boolArg(args, "create") {
		return nil, err
	}
	diff := simpleUnifiedDiff(filepath.ToSlash(rp.Rel), old, content)
	return map[string]any{
		"workspace":   rp.Workspace.Name,
		"path":        rp.Rel,
		"would_write": true,
		"old_bytes":   len(old),
		"new_bytes":   len(content),
		"diff":        diff,
	}, nil
}

func statMap(rp resolvedPath, info os.FileInfo) map[string]any {
	return map[string]any{
		"workspace": rp.Workspace.Name,
		"path":      filepath.ToSlash(rp.Rel),
		"is_dir":    info.IsDir(),
		"mode":      info.Mode().String(),
		"size":      info.Size(),
		"modified":  info.ModTime().UTC().Format(time.RFC3339),
	}
}

func shouldSkipTreeEntry(name string) bool {
	switch name {
	case ".git", "node_modules", ".anigate", ".DS_Store":
		return true
	default:
		return false
	}
}

func simpleUnifiedDiff(path, oldText, newText string) string {
	if oldText == newText {
		return ""
	}
	oldLines := strings.SplitAfter(oldText, "\n")
	newLines := strings.SplitAfter(newText, "\n")
	var b strings.Builder
	b.WriteString("--- a/" + path + "\n")
	b.WriteString("+++ b/" + path + "\n")
	b.WriteString(fmt.Sprintf("@@ -1,%d +1,%d @@\n", len(oldLines), len(newLines)))
	for _, line := range oldLines {
		if line == "" {
			continue
		}
		b.WriteString("-" + strings.TrimSuffix(line, "\n") + "\n")
	}
	for _, line := range newLines {
		if line == "" {
			continue
		}
		b.WriteString("+" + strings.TrimSuffix(line, "\n") + "\n")
	}
	return b.String()
}
