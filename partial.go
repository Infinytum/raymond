package raymond

import (
	"fmt"
	"sync"
)

// Partial represents a Partial template
type Partial struct {
	Name   string
	Source string
	Tpl    *Template
}

// partials stores all global partials
var partials map[string]*Partial

var ResolvePartial func(view string) *Partial = func(view string) *Partial { return nil }

// protects global partials
var partialsMutex sync.RWMutex

func init() {
	partials = make(map[string]*Partial)
}

// NewPartial instanciates a new partial
func NewPartial(name string, source string, tpl *Template) *Partial {
	return &Partial{
		Name:   name,
		Source: source,
		Tpl:    tpl,
	}
}

// RegisterPartial registers a global partial. That partial will be available to all templates.
func RegisterPartial(name string, source string) {
	partialsMutex.Lock()
	defer partialsMutex.Unlock()

	if partials[name] != nil {
		panic(fmt.Errorf("Partial already registered: %s", name))
	}

	partials[name] = NewPartial(name, source, nil)
}

// RegisterPartials registers several global partials. Those partials will be available to all templates.
func RegisterPartials(partials map[string]string) {
	for name, p := range partials {
		RegisterPartial(name, p)
	}
}

// RegisterPartialTemplate registers a global partial with given parsed template. That partial will be available to all templates.
func RegisterPartialTemplate(name string, tpl *Template) {
	partialsMutex.Lock()
	defer partialsMutex.Unlock()

	if partials[name] != nil {
		panic(fmt.Errorf("Partial already registered: %s", name))
	}

	partials[name] = NewPartial(name, "", tpl)
}

// RemovePartial removes the partial registered under the given name. The partial will not be available globally anymore. This does not affect partials registered on a specific template.
func RemovePartial(name string) {
	partialsMutex.Lock()
	defer partialsMutex.Unlock()

	delete(partials, name)
}

// RemoveAllPartials removes all globally registered partials. This does not affect partials registered on a specific template.
func RemoveAllPartials() {
	partialsMutex.Lock()
	defer partialsMutex.Unlock()

	partials = make(map[string]*Partial)
}

// findPartial finds a registered global partial
func findPartial(name string) *Partial {
	partialsMutex.RLock()
	defer partialsMutex.RUnlock()

	if partial, ok := partials[name]; ok {
		return partial
	}

	return ResolvePartial(name)
}

// template returns parsed partial template
func (p *Partial) template() (*Template, error) {
	if p.Tpl == nil {
		var err error

		p.Tpl, err = Parse(p.Source)
		if err != nil {
			return nil, err
		}
	}

	return p.Tpl, nil
}
