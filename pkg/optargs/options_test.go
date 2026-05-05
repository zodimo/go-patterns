package optargs

import (
	"testing"
	"time"
)

type testOptions struct {
	name    string
	value   int
	enabled bool
}

func defaultTestOptions() testOptions {
	return testOptions{
		name:    "default",
		value:   0,
		enabled: false,
	}
}

func withName(name string) func(*testOptions) {
	return func(o *testOptions) {
		o.name = name
	}
}

func withValue(value int) func(*testOptions) {
	return func(o *testOptions) {
		o.value = value
	}
}

func withEnabled(enabled bool) func(*testOptions) {
	return func(o *testOptions) {
		o.enabled = enabled
	}
}

func TestHandleOptions(t *testing.T) {
	tests := []struct {
		name           string
		options        []func(*testOptions)
		expectedName   string
		expectedValue  int
		expectedEnabled bool
	}{
		{
			name:           "no options uses defaults",
			options:        nil,
			expectedName:   "default",
			expectedValue:  0,
			expectedEnabled: false,
		},
		{
			name:           "single option applied",
			options:        []func(*testOptions){withName("custom")},
			expectedName:   "custom",
			expectedValue:  0,
			expectedEnabled: false,
		},
		{
			name:           "multiple options applied",
			options:        []func(*testOptions){withName("custom"), withValue(42), withEnabled(true)},
			expectedName:   "custom",
			expectedValue:  42,
			expectedEnabled: true,
		},
		{
			name:           "nil options are skipped",
			options:        []func(*testOptions){nil, withName("after-nil"), nil},
			expectedName:   "after-nil",
			expectedValue:  0,
			expectedEnabled: false,
		},
		{
			name:           "all nil options uses defaults",
			options:        []func(*testOptions){nil, nil, nil},
			expectedName:   "default",
			expectedValue:  0,
			expectedEnabled: false,
		},
		{
			name:           "empty options slice uses defaults",
			options:        []func(*testOptions){},
			expectedName:   "default",
			expectedValue:  0,
			expectedEnabled: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HandleOptions(defaultTestOptions, tt.options...)

			if result.name != tt.expectedName {
				t.Errorf("expected name %q, got %q", tt.expectedName, result.name)
			}
			if result.value != tt.expectedValue {
				t.Errorf("expected value %d, got %d", tt.expectedValue, result.value)
			}
			if result.enabled != tt.expectedEnabled {
				t.Errorf("expected enabled %v, got %v", tt.expectedEnabled, result.enabled)
			}
		})
	}
}

func TestNewOptionsHandlerContext(t *testing.T) {
	tests := []struct {
		name           string
		defaultsFactory func() testOptions
		options        []func(*testOptions)
		expectedName   string
		expectedValue  int
	}{
		{
			name:           "context with defaults only",
			defaultsFactory: defaultTestOptions,
			options:        nil,
			expectedName:   "default",
			expectedValue:  0,
		},
		{
			name:           "context with options",
			defaultsFactory: defaultTestOptions,
			options:        []func(*testOptions){withName("contextual"), withValue(100)},
			expectedName:   "contextual",
			expectedValue:  100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := NewOptionsHandlerContext(tt.defaultsFactory, tt.options...)

			if ctx.defaultsFactory == nil {
				t.Error("expected defaultsFactory to be set")
			}
			if len(ctx.options) != len(tt.options) {
				t.Errorf("expected %d options, got %d", len(tt.options), len(ctx.options))
			}
		})
	}
}

func TestHandleOptionsFromContext(t *testing.T) {
	tests := []struct {
		name           string
		ctx            OptionsHandlerContext[testOptions]
		expectedName   string
		expectedValue  int
		expectedEnabled bool
	}{
		{
			name: "context with defaults factory",
			ctx: NewOptionsHandlerContext(
				defaultTestOptions,
				withName("from-context"),
				withValue(999),
				withEnabled(true),
			),
			expectedName:   "from-context",
			expectedValue:  999,
			expectedEnabled: true,
		},
		{
			name:           "context with no options",
			ctx:            NewOptionsHandlerContext(defaultTestOptions),
			expectedName:   "default",
			expectedValue:  0,
			expectedEnabled: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HandleOptionsFromContext(tt.ctx)

			if result.name != tt.expectedName {
				t.Errorf("expected name %q, got %q", tt.expectedName, result.name)
			}
			if result.value != tt.expectedValue {
				t.Errorf("expected value %d, got %d", tt.expectedValue, result.value)
			}
			if result.enabled != tt.expectedEnabled {
				t.Errorf("expected enabled %v, got %v", tt.expectedEnabled, result.enabled)
			}
		})
	}
}

func TestHandleOptionsRealWorldScenario(t *testing.T) {
	type Config struct {
		Host    string
		Port    int
		Timeout time.Duration
		Debug   bool
	}

	defaultConfig := func() Config {
		return Config{
			Host:    "localhost",
			Port:    8080,
			Timeout: 30 * time.Second,
			Debug:   false,
		}
	}

	withHost := func(host string) func(*Config) {
		return func(c *Config) {
			c.Host = host
		}
	}

	withPort := func(port int) func(*Config) {
		return func(c *Config) {
			c.Port = port
		}
	}

	withTimeout := func(timeout time.Duration) func(*Config) {
		return func(c *Config) {
			c.Timeout = timeout
		}
	}

	withDebug := func(debug bool) func(*Config) {
		return func(c *Config) {
			c.Debug = debug
		}
	}

	t.Run("default configuration", func(t *testing.T) {
		config := HandleOptions(defaultConfig)

		if config.Host != "localhost" {
			t.Errorf("expected host localhost, got %s", config.Host)
		}
		if config.Port != 8080 {
			t.Errorf("expected port 8080, got %d", config.Port)
		}
		if config.Timeout != 30*time.Second {
			t.Errorf("expected timeout 30s, got %v", config.Timeout)
		}
		if config.Debug != false {
			t.Errorf("expected debug false, got %v", config.Debug)
		}
	})

	t.Run("custom configuration", func(t *testing.T) {
		config := HandleOptions(
			defaultConfig,
			withHost("example.com"),
			withPort(3000),
			withTimeout(5*time.Minute),
			withDebug(true),
		)

		if config.Host != "example.com" {
			t.Errorf("expected host example.com, got %s", config.Host)
		}
		if config.Port != 3000 {
			t.Errorf("expected port 3000, got %d", config.Port)
		}
		if config.Timeout != 5*time.Minute {
			t.Errorf("expected timeout 5m, got %v", config.Timeout)
		}
		if config.Debug != true {
			t.Errorf("expected debug true, got %v", config.Debug)
		}
	})

	t.Run("partial configuration", func(t *testing.T) {
		config := HandleOptions(
			defaultConfig,
			withPort(9090),
		)

		if config.Host != "localhost" {
			t.Errorf("expected host localhost, got %s", config.Host)
		}
		if config.Port != 9090 {
			t.Errorf("expected port 9090, got %d", config.Port)
		}
		if config.Timeout != 30*time.Second {
			t.Errorf("expected timeout 30s, got %v", config.Timeout)
		}
		if config.Debug != false {
			t.Errorf("expected debug false, got %v", config.Debug)
		}
	})

	t.Run("mixed with nil options", func(t *testing.T) {
		config := HandleOptions(
			defaultConfig,
			nil,
			withHost("test.com"),
			nil,
			withDebug(true),
			nil,
		)

		if config.Host != "test.com" {
			t.Errorf("expected host test.com, got %s", config.Host)
		}
		if config.Debug != true {
			t.Errorf("expected debug true, got %v", config.Debug)
		}
	})
}

func TestHandleOptionsFromContextRealWorldScenario(t *testing.T) {
	type Preferences struct {
		Theme   string
		Timeout time.Duration
	}

	defaultPreferences := func() Preferences {
		return Preferences{
			Theme:   "light",
			Timeout: 5 * time.Second,
		}
	}

	withTheme := func(theme string) func(*Preferences) {
		return func(p *Preferences) {
			p.Theme = theme
		}
	}

	withTimeout := func(timeout time.Duration) func(*Preferences) {
		return func(p *Preferences) {
			p.Timeout = timeout
		}
	}

	t.Run("create person with context and default preferences", func(t *testing.T) {
		ctx := NewOptionsHandlerContext(defaultPreferences)
		prefs := HandleOptionsFromContext(ctx)

		if prefs.Theme != "light" {
			t.Errorf("expected theme light, got %s", prefs.Theme)
		}
		if prefs.Timeout != 5*time.Second {
			t.Errorf("expected timeout 5s, got %v", prefs.Timeout)
		}
	})

	t.Run("create person with context and custom preferences", func(t *testing.T) {
		ctx := NewOptionsHandlerContext(
			defaultPreferences,
			withTheme("dark"),
			withTimeout(10*time.Second),
		)
		prefs := HandleOptionsFromContext(ctx)

		if prefs.Theme != "dark" {
			t.Errorf("expected theme dark, got %s", prefs.Theme)
		}
		if prefs.Timeout != 10*time.Second {
			t.Errorf("expected timeout 10s, got %v", prefs.Timeout)
		}
	})
}

func BenchmarkHandleOptions(b *testing.B) {
	opts := []func(*testOptions){
		withName("benchmark"),
		withValue(100),
		withEnabled(true),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = HandleOptions(defaultTestOptions, opts...)
	}
}

func BenchmarkHandleOptionsWithNil(b *testing.B) {
	opts := []func(*testOptions){
		nil,
		withName("benchmark"),
		nil,
		withValue(100),
		nil,
		withEnabled(true),
		nil,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = HandleOptions(defaultTestOptions, opts...)
	}
}

func BenchmarkHandleOptionsFromContext(b *testing.B) {
	ctx := NewOptionsHandlerContext(
		defaultTestOptions,
		withName("benchmark"),
		withValue(100),
		withEnabled(true),
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = HandleOptionsFromContext(ctx)
	}
}

func BenchmarkHandleOptionsNoOptions(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = HandleOptions(defaultTestOptions)
	}
}

func TestHandleOptionsInto(t *testing.T) {
	tests := []struct {
		name            string
		initial         testOptions
		options         []func(*testOptions)
		expectedName    string
		expectedValue   int
		expectedEnabled bool
	}{
		{
			name:            "no options keeps initial values",
			initial:         testOptions{name: "initial", value: 10, enabled: true},
			options:         nil,
			expectedName:    "initial",
			expectedValue:   10,
			expectedEnabled: true,
		},
		{
			name:            "single option applied",
			initial:         testOptions{name: "initial", value: 10, enabled: true},
			options:         []func(*testOptions){withName("updated")},
			expectedName:    "updated",
			expectedValue:   10,
			expectedEnabled: true,
		},
		{
			name:            "multiple options applied",
			initial:         testOptions{name: "initial", value: 0, enabled: false},
			options:         []func(*testOptions){withName("updated"), withValue(42), withEnabled(true)},
			expectedName:    "updated",
			expectedValue:   42,
			expectedEnabled: true,
		},
		{
			name:            "nil options are skipped",
			initial:         testOptions{name: "initial", value: 5, enabled: false},
			options:         []func(*testOptions){nil, withName("after-nil"), nil},
			expectedName:    "after-nil",
			expectedValue:   5,
			expectedEnabled: false,
		},
		{
			name:            "nil target does nothing",
			initial:         testOptions{name: "unchanged", value: 99, enabled: true},
			options:         []func(*testOptions){withName("should-not-apply")},
			expectedName:    "unchanged",
			expectedValue:   99,
			expectedEnabled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var target *testOptions
			if tt.name != "nil target does nothing" {
				target = &testOptions{name: tt.initial.name, value: tt.initial.value, enabled: tt.initial.enabled}
			}

			HandleOptionsInto(target, tt.options...)

			if tt.name == "nil target does nothing" {
				return
			}

			if target.name != tt.expectedName {
				t.Errorf("expected name %q, got %q", tt.expectedName, target.name)
			}
			if target.value != tt.expectedValue {
				t.Errorf("expected value %d, got %d", tt.expectedValue, target.value)
			}
			if target.enabled != tt.expectedEnabled {
				t.Errorf("expected enabled %v, got %v", tt.expectedEnabled, target.enabled)
			}
		})
	}
}

func TestHandleOptionsIntoRealWorldScenario(t *testing.T) {
	type Config struct {
		Host    string
		Port    int
		Timeout time.Duration
		Debug   bool
	}

	withHost := func(host string) func(*Config) {
		return func(c *Config) {
			c.Host = host
		}
	}

	withPort := func(port int) func(*Config) {
		return func(c *Config) {
			c.Port = port
		}
	}

	t.Run("reusable config object", func(t *testing.T) {
		config := Config{
			Host:    "localhost",
			Port:    8080,
			Timeout: 30 * time.Second,
			Debug:   false,
		}

		HandleOptionsInto(&config,
			withHost("example.com"),
			withPort(3000),
		)

		if config.Host != "example.com" {
			t.Errorf("expected host example.com, got %s", config.Host)
		}
		if config.Port != 3000 {
			t.Errorf("expected port 3000, got %d", config.Port)
		}
		if config.Timeout != 30*time.Second {
			t.Errorf("expected timeout to remain 30s, got %v", config.Timeout)
		}
		if config.Debug != false {
			t.Errorf("expected debug to remain false, got %v", config.Debug)
		}
	})
}

func BenchmarkHandleOptionsInto(b *testing.B) {
	opts := []func(*testOptions){
		withName("benchmark"),
		withValue(100),
		withEnabled(true),
	}

	var target testOptions
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		target = testOptions{name: "default", value: 0, enabled: false}
		HandleOptionsInto(&target, opts...)
	}
}
