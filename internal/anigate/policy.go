package anigate

func (s *Service) policyInfo() (map[string]any, error) {
	workspaces := make([]map[string]any, 0, len(s.cfg.Workspaces))
	for _, ws := range s.cfg.Workspaces {
		workspaces = append(workspaces, map[string]any{
			"name":      ws.Name,
			"path":      ws.Path,
			"read_only": ws.ReadOnly,
			"profile":   ws.Profile,
		})
	}
	presets := make([]map[string]any, 0, len(s.cfg.Presets))
	for _, p := range s.cfg.Presets {
		args := make([]map[string]any, 0, len(p.Args))
		for _, arg := range p.Args {
			args = append(args, map[string]any{
				"name":      arg.Name,
				"type":      arg.Type,
				"required":  arg.Required,
				"enum":      arg.Enum,
				"max_len":   arg.MaxLen,
				"max_items": arg.MaxItems,
			})
		}
		presets = append(presets, map[string]any{
			"name":        p.Name,
			"description": p.Description,
			"workspace":   p.Workspace,
			"async":       p.Async,
			"args":        args,
		})
	}
	agents := make([]map[string]any, 0, len(s.cfg.Agents))
	for _, a := range s.cfg.Agents {
		agents = append(agents, map[string]any{
			"name":        a.Name,
			"description": a.Description,
			"provider":    a.Provider,
			"workspace":   a.Workspace,
			"timeout_sec": a.TimeoutSec,
		})
	}
	projects := make([]map[string]any, 0, len(s.cfg.Projects))
	for _, p := range s.cfg.Projects {
		projects = append(projects, map[string]any{
			"name":           p.Name,
			"description":    p.Description,
			"workspace":      p.Workspace,
			"path":           p.Path,
			"remote_url":     redactRemoteURL(p.RemoteURL),
			"default_branch": p.DefaultBranch,
			"provider":       p.Provider,
			"allow_push":     p.AllowPush,
			"allow_pr":       p.AllowPR,
			"default_agent":  p.DefaultAgent,
		})
	}
	return map[string]any{
		"server":                "AniGate",
		"version":               Version,
		"mode":                  "controlled-linux-gateway",
		"auth_required":         s.cfg.AuthToken != "",
		"no_arbitrary_shell":    true,
		"max_read_bytes":        s.cfg.MaxReadBytes,
		"max_search_file_bytes": s.cfg.MaxSearchFileBytes,
		"max_search_results":    s.cfg.MaxSearchResults,
		"max_job_log_bytes":     s.cfg.MaxJobLogBytes,
		"max_artifact_bytes":    s.cfg.MaxArtifactBytes,
		"isolated_home":         s.cfg.IsolatedHome,
		"env_allowlist":         s.cfg.EnvAllowlist,
		"workspaces":            workspaces,
		"presets":               presets,
		"agents":                agents,
		"projects":              projects,
		"tools":                 s.Tools(),
	}, nil
}

func (s *Service) workspaceAllows(workspaceName, need string) error {
	ws, err := s.policy.workspace(workspaceName)
	if err != nil {
		return err
	}
	switch need {
	case "read":
		return nil
	case "operate":
		if ws.Profile == "operator" || ws.Profile == "agent" {
			return nil
		}
	case "agent":
		if ws.Profile == "agent" {
			return nil
		}
	case "write":
		if !ws.ReadOnly && (ws.Profile == "operator" || ws.Profile == "agent") {
			return nil
		}
	}
	return permissionError(workspaceName, need)
}

func permissionError(workspaceName, need string) error {
	return errString("workspace %q does not allow %s", workspaceName, need)
}
