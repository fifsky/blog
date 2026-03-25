package skill

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSkill(t *testing.T) {
	// Create a temporary file
	dir := t.TempDir()
	filePath := filepath.Join(dir, "SKILL.md")

	content := `---
name: my-skill
description: This is a test skill
---
# Test Skill
This is the body of the skill.
`
	err := os.WriteFile(filePath, []byte(content), 0644)
	assert.NoError(t, err)

	skill, err := ParseSkillFile(filePath)
	assert.NoError(t, err)
	assert.Equal(t, "my-skill", skill.Name())
	assert.Equal(t, "This is a test skill", skill.Description())
	assert.Equal(t, "# Test Skill\nThis is the body of the skill.", skill.Content)
}

func TestParseSkillNoFrontmatter(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "SKILL.md")

	content := `# Test Skill
This is the body of the skill.
`
	err := os.WriteFile(filePath, []byte(content), 0644)
	assert.NoError(t, err)

	skill, err := ParseSkillFile(filePath)
	assert.NoError(t, err)
	assert.Equal(t, "", skill.Name())
	assert.Equal(t, "", skill.Description())
	assert.Equal(t, "# Test Skill\nThis is the body of the skill.", skill.Content)
}

func TestManagerLoad(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Skip("Could not get home directory")
	}

	skillsPath := filepath.Join(homeDir, ".agents", "skills")
	if _, err := os.Stat(skillsPath); os.IsNotExist(err) {
		t.Skip("Skills directory does not exist, skipping real directory test")
	}

	manager := NewManager(skillsPath)
	err = manager.Load()
	assert.NoError(t, err)

	skills := manager.GetSkills()
	assert.NotNil(t, skills)

	// Print loaded skills for debugging purposes
	t.Logf("Loaded %d skills from %s", len(skills), skillsPath)
	for name, skill := range skills {
		t.Logf("Skill: %s, Description: %s", name, skill.Description())
		assert.NotEmpty(t, skill.Name())
		// Description might be empty for some skills, but usually not
	}
}
