package anigate

import "fmt"

type ProductLine string

const (
	ProductLineMini ProductLine = "mini"
	ProductLineMax  ProductLine = "max"
)

func normalizeProductLine(line ProductLine) (ProductLine, error) {
	if line == "" {
		return ProductLineMax, nil
	}
	switch line {
	case ProductLineMini, ProductLineMax:
		return line, nil
	default:
		return "", fmt.Errorf("unknown product line %q", line)
	}
}

var miniToolNames = map[string]struct{}{
	"policy.info":         {},
	"sys.info":            {},
	"fs.list":             {},
	"fs.read":             {},
	"fs.stat":             {},
	"fs.tree":             {},
	"file.search":         {},
	"fs.write_preview":    {},
	"git.status":          {},
	"git.diff":            {},
	"git.log":             {},
	"git.show":            {},
	"artifact.list":       {},
	"artifact.read_range": {},
	"artifact.search":     {},
	"artifact.stats":      {},
	"context.health":      {},
	"handoff.create":      {},
	"handoff.resume":      {},
	"handoff.search":      {},
	"handoff.digest":      {},
}

func toolNamesForProduct(line ProductLine) []string {
	tools := allTools()
	names := make([]string, 0, len(tools))
	for _, tool := range tools {
		if line == ProductLineMini {
			if _, ok := miniToolNames[tool.Name]; !ok {
				continue
			}
		}
		names = append(names, tool.Name)
	}
	return names
}

func knownToolName(name string) bool {
	for _, tool := range allTools() {
		if tool.Name == name {
			return true
		}
	}
	return false
}

func (s *Service) toolAllowedForProduct(name string) bool {
	if s.productLine != ProductLineMini {
		return true
	}
	_, ok := miniToolNames[name]
	return ok
}

func (s *Service) requireToolForProduct(name string) error {
	if s.toolAllowedForProduct(name) {
		return nil
	}
	if !knownToolName(name) {
		return nil
	}
	return errString("tool %q is not available in %s product line", name, s.productLine)
}

func productLinesInfo() map[string]any {
	return map[string]any{
		"mini": map[string]any{
			"description": "Safe preview gateway with read, search, diff, artifact, context, and handoff tools only.",
			"entrypoint":  "anigate-mini",
			"tools":       toolNamesForProduct(ProductLineMini),
		},
		"max": map[string]any{
			"description":  "Controlled Linux workbench with execution, mutation, agent, task, and publish tools.",
			"entrypoint":   "anigate-max",
			"legacy_alias": "anigate",
			"tools":        toolNamesForProduct(ProductLineMax),
		},
		"enforcement": "Product line filters tools/list and gates tools/call before dispatch; workspace profile and read_only policy remain a second layer.",
	}
}
