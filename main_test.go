package main

import (
	"os"
	"reflect"
	"testing"
)

func TestParseCustomLabels(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		expected []LabelConfig
	}{
		{
			name:     "empty environment variable",
			envValue: "",
			expected: []LabelConfig{},
		},
		{
			name:     "unset environment variable",
			envValue: "",
			expected: []LabelConfig{},
		},
		{
			name:     "single label",
			envValue: "environment=production",
			expected: []LabelConfig{
				{Key: "environment", Value: "production"},
			},
		},
		{
			name:     "multiple labels",
			envValue: "environment=production,team=platform,region=us-west",
			expected: []LabelConfig{
				{Key: "environment", Value: "production"},
				{Key: "team", Value: "platform"},
				{Key: "region", Value: "us-west"},
			},
		},
		{
			name:     "labels with spaces",
			envValue: "environment = production , team = platform",
			expected: []LabelConfig{
				{Key: "environment", Value: "production"},
				{Key: "team", Value: "platform"},
			},
		},
		{
			name:     "empty pairs are ignored",
			envValue: "environment=production,,team=platform",
			expected: []LabelConfig{
				{Key: "environment", Value: "production"},
				{Key: "team", Value: "platform"},
			},
		},
		{
			name:     "invalid format is ignored",
			envValue: "environment=production,invalid,team=platform",
			expected: []LabelConfig{
				{Key: "environment", Value: "production"},
				{Key: "team", Value: "platform"},
			},
		},
		{
			name:     "labels with special characters",
			envValue: "app.kubernetes.io/name=myapp,app.kubernetes.io/version=v1.0.0",
			expected: []LabelConfig{
				{Key: "app.kubernetes.io/name", Value: "myapp"},
				{Key: "app.kubernetes.io/version", Value: "v1.0.0"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variable
			if tt.envValue != "" {
				os.Setenv("CUSTOM_LABELS", tt.envValue)
			} else {
				os.Unsetenv("CUSTOM_LABELS")
			}

			// Clean up after test
			defer os.Unsetenv("CUSTOM_LABELS")

			result := parseCustomLabels()

			if len(result) != len(tt.expected) {
				t.Errorf("parseCustomLabels() returned %d items, want %d items", len(result), len(tt.expected))
			} else if len(result) > 0 && !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("parseCustomLabels() = %v (type: %T), want %v (type: %T)", result, result, tt.expected, tt.expected)
			}
		})
	}
}

func TestParseCustomLabelsWithInvalidFormats(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		expected []LabelConfig
	}{
		{
			name:     "key-only format",
			envValue: "key-only",
			expected: []LabelConfig{},
		},
		{
			name:     "value-only format",
			envValue: "=value-only",
			expected: []LabelConfig{
				{Key: "", Value: "value-only"},
			},
		},
		{
			name:     "key=value=extra format",
			envValue: "key=value=extra",
			expected: []LabelConfig{
				{Key: "key", Value: "value=extra"},
			},
		},
		{
			name:     "empty string",
			envValue: "",
			expected: []LabelConfig{},
		},
		{
			name:     "whitespace only",
			envValue: "   ",
			expected: []LabelConfig{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("CUSTOM_LABELS", tt.envValue)
			defer os.Unsetenv("CUSTOM_LABELS")

			result := parseCustomLabels()

			if len(result) != len(tt.expected) {
				t.Errorf("parseCustomLabels() returned %d items, want %d items", len(result), len(tt.expected))
			} else if len(result) > 0 && !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("parseCustomLabels() = %v (type: %T), want %v (type: %T)", result, result, tt.expected, tt.expected)
			}
		})
	}
}
