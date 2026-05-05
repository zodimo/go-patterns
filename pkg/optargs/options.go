package optargs

type OptionsHandlerContext[T any] struct {
	defaultsFactory func() T
	options         []func(*T)
}

func NewOptionsHandlerContext[T any](
	defaultsFactory func() T,
	options ...func(*T),
) OptionsHandlerContext[T] {
	return OptionsHandlerContext[T]{
		defaultsFactory: defaultsFactory,
		options:         options,
	}
}

func HandleOptionsFromContext[T any](ctx OptionsHandlerContext[T]) T {
	return HandleOptions[T](
		ctx.defaultsFactory,
		ctx.options...,
	)
}

func HandleOptions[T any](
	defaultsFactory func() T,
	options ...func(*T),
) T {
	opts := defaultsFactory()
	for _, option := range options {
		if option != nil {
			option(&opts)
		}

	}
	return opts
}

// HandleOptionsInto applies options to an existing value without allocation.
// Use this when you want to avoid heap allocations by managing your own storage.
// target must be non-nil.
func HandleOptionsInto[T any](
	target *T,
	options ...func(*T),
) {
	if target == nil {
		return
	}
	for _, option := range options {
		if option != nil {
			option(target)
		}
	}
}
