package anigate

import (
	"crypto/sha256"
	"encoding/hex"
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
	if err := s.workspaceAllows(workspace, "read"); err != nil {
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

func (s *Service) fileEditApply(args map[string]any) (map[string]any, error) {
	workspace := stringArg(args, "workspace")
	if err := s.workspaceAllows(workspace, "write"); err != nil {
		return nil, err
	}
	rp, err := s.policy.resolve(workspace, stringArg(args, "path"))
	if err != nil {
		return nil, err
	}
	if err := ensureWriteParentInsideWorkspace(rp); err != nil {
		return nil, err
	}
	content := stringArg(args, "content")
	if int64(len(content)) > s.cfg.MaxReadBytes {
		return nil, fmt.Errorf("content exceeds max_read_bytes")
	}
	create := boolArg(args, "create")
	old := ""
	oldBytes, err := os.ReadFile(rp.Abs)
	if err != nil {
		if !os.IsNotExist(err) || !create {
			return nil, err
		}
	} else {
		if looksBinary(oldBytes) {
			return nil, fmt.Errorf("existing file appears binary")
		}
		old = string(oldBytes)
	}
	if expected := stringArg(args, "expected_sha256"); expected != "" && expected != sha256Hex(oldBytes) {
		return nil, fmt.Errorf("expected_sha256 does not match current file")
	}
	if expected := stringArg(args, "expected_old_text"); expected != "" && !strings.Contains(old, expected) {
		return nil, fmt.Errorf("expected_old_text not found in current file")
	}
	diff := simpleUnifiedDiff(filepath.ToSlash(rp.Rel), old, content)
	if err := os.WriteFile(rp.Abs, []byte(content), 0o644); err != nil {
		return nil, err
	}
	s.events.Append(Event{Kind: "file_edited", Tool: "file.edit_apply", Workspace: rp.Workspace.Name, Path: rp.Rel, OK: true, Fields: map[string]any{"actor": "web_gpt_direct"}})
	return map[string]any{
		"workspace":     rp.Workspace.Name,
		"path":          rp.Rel,
		"written":       true,
		"old_bytes":     len(old),
		"new_bytes":     len(content),
		"before_sha256": sha256Hex(oldBytes),
		"after_sha256":  sha256Hex([]byte(content)),
		"diff":          diff,
		"actor":         "web_gpt_direct",
		"next":          []string{"task.finish_preview", "git.diff"},
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

func sha256Hex(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

func ensureWriteParentInsideWorkspace(rp resolvedPath) error {
	root, err := filepath.Abs(rp.Workspace.Path)
	if err != nil {
		return err
	}
	if realRoot, err := filepath.EvalSymlinks(root); err == nil {
		root = realRoot
	}
	parent := filepath.Dir(rp.Abs)
	realParent, err := filepath.EvalSymlinks(parent)
	if err != nil {
		return err
	}
	rel, err := filepath.Rel(root, realParent)
	if err != nil {
		return err
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return fmt.Errorf("path escapes workspace %q", rp.Workspace.Name)
	}
	return nil
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
