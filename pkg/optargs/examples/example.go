package optargsexamples

import (
	"time"

	"github.com/zodimo/go-patterns/pkg/optargs"
)

type Preferences struct {
	Theme   string
	Timeout time.Duration
}

func WithTheme(theme string) func(*Preferences) {
	return func(p *Preferences) {
		p.Theme = theme
	}
}

func WithTimeout(timeout time.Duration) func(*Preferences) {
	return func(p *Preferences) {
		p.Timeout = timeout
	}
}

type Person struct {
	name       string
	preference Preferences
}

func DefaultPreferences() Preferences {
	return Preferences{
		Theme:   "Default",
		Timeout: 5 * time.Second,
	}
}

func NewPerson(
	name string,
	preference Preferences,
) Person {
	return Person{
		name:       name,
		preference: preference,
	}
}

func NewPerson2(name string, preferences ...func(*Preferences)) Person {
	// pattern for dealing with optional arguments
	opts := DefaultPreferences() // default or empty , Preferences{}
	for _, option := range preferences {
		if option != nil {
			option(&opts)
		}
	}

	return Person{
		name:       name,
		preference: opts,
	}
}

func NewPerson3(name string, preferences ...func(*Preferences)) Person {
	opts := optargs.HandleOptions(
		DefaultPreferences,
		preferences...,
	)

	return Person{
		name:       name,
		preference: opts,
	}
}

func Example_NoOptionalArgs() {
	_ = NewPerson(
		"Joe",
		Preferences{
			Theme:   "Default",
			Timeout: 5 * time.Second,
		},
	)
	// or
	_ = NewPerson(
		"Joe",
		DefaultPreferences(),
	)

}

func Example_OptionalArgs() {

	_ = NewPerson2("Joe")

	// or

	_ = NewPerson2("Joe", WithTheme("Special"), WithTimeout(10*time.Second))

}

func Example_OptionalArgs_WithHandler() {

	_ = NewPerson3("Joe")

	// or

	_ = NewPerson3("Joe", WithTheme("Special"), WithTimeout(10*time.Second))

}
