package main

import (
	"fmt"

	"github.com/charmbracelet/x/ansi"
)

type Setting struct {
	name    string
	setting string
}

func (setting *Setting) View() string {
	spaceAvailable := 32 - ansi.StringWidth(setting.name)
	return fmt.Sprintf("   %s: %*s\n", setting.name, spaceAvailable, setting.setting)
}

var settingPositions []string = []string{
	"work_mem",
	"random_page_cost",
	"join_collapse_limit",
	"effective_cache_size",
	"max_parallel_workers_per_gather",
}

func (setting *Setting) FindPosition() int {
	for i, sPos := range settingPositions {
		if sPos == setting.name {
			return i
		}
	}
	return -1
}

func SettingCompare(a, b Setting) int {
	return a.FindPosition() - b.FindPosition()
}

func (setting Setting) Sql() string {
	return fmt.Sprintf("SET %s = '%s'", setting.name, setting.setting)
}
