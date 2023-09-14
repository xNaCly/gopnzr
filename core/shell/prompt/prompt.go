package prompt

import (
	"gopnzr/core/shell/env"
	"gopnzr/core/shell/system"
	"os"
	"os/user"
	"strings"
)

const DEFAULT_PROMPT = "\033[94m'u\033[0m@\033[93m'h\033[0m \033[95m'd\033[0m \033[32m>\033[0m "

// contains all possible placeholders a prompt could contain
var prompt_placeholders = map[rune]string{
	'u': "",
	'h': "",
	'w': "",
	'd': "",
	// TODO: support git-branch (b) (either nothing or the branch name, see 'git branch')
	// TODO: support git-status (s) (either nothing or M for modified, see 'git status --short')
	// TODO: support time (t) (hh:mm:ss, 24hr)
	// TODO: support time (T) (hh:mm:ss, 12hr)
	// TODO: support date (D) (yyyy-mm-dd)
	// TODO: shell name (S)
}

// computes placeholder values that are known at startup, this decreases load
// on the main loop prompt computation
func PreComputePlaceholders() (e error) {
	u, e := user.Current()

	prompt_placeholders['u'] = u.Username
	pwd, _ := env.GetEnv("PWD")
	prompt_placeholders['w'] = pwd
	prompt_placeholders['d'] = system.Getdir()

	h, e := os.Hostname()
	prompt_placeholders['h'] = h
	return
}

// checks if custom prompt is set, returns either that prompt or the default
// prompt with placeholders replaced
func ComputePrompt() string {
	prompt := DEFAULT_PROMPT
	if val, ok := env.GetEnv("GPNZR_PROMPT"); ok {
		prompt = val
	}
	return replacePlaceholders(prompt)
}

// formats the working directory according to the configuration
func formatWd(path string) string {
	if !env.GetEnvBool("GPNZR_PROMPT_SHORT_PWD") {
		return path
	}
	b := strings.Builder{}
	var lc rune
	for _, c := range path {
		if lc == '/' {
			b.WriteRune(c)
			b.WriteRune('/')
		}
		lc = c
	}
	return b.String()
}

// replaces placeholders in the given prompt with the values in
// 'prompt_placeholders', works by detecting slashes and writing the
// 'prompt_placeholders' value of the placeholder into a string builder, which
// gets returned, this should be incredibly faster than calling strings.Replace
// on each placeholder
func replacePlaceholders(prompt string) string {
	prompt_placeholders['w'] = formatWd(prompt_placeholders['w'])
	b := strings.Builder{}
	placeHolderMode := false
	for _, c := range prompt {
		if c == '\'' {
			placeHolderMode = true
		} else if placeHolderMode {
			if t, ok := prompt_placeholders[c]; ok {
				b.WriteString(t)
			}
			placeHolderMode = false
		} else {
			b.WriteRune(c)
		}
	}
	return b.String()
}
