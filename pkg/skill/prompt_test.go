package skill

import (
	"strings"
	"testing"
)

func TestSkillsPrompt(t *testing.T) {
	t.Parallel()

	xml := SkillsPrompt(map[string]*Skill{
		"skill1": {Frontmatter: Frontmatter{Name: "skill1", Description: "desc1"}},
		"skill2": {Frontmatter: Frontmatter{Name: "skill2", Description: "desc<2>"}},
	})
	if xml == "" {
		t.Fatalf("expected non-empty xml")
	}
	if want := "<available_skills>"; !strings.Contains(xml, want) {
		t.Fatalf("missing %q", want)
	}
	if want := "desc&lt;2&gt;"; !strings.Contains(xml, want) {
		t.Fatalf("missing escaped description")
	}
}
