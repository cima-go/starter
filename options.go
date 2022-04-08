package starter

type Option func(s *Starter)

func FindConfInHome(name string) Option {
	return func(s *Starter) {
		s.confName = name
	}
}

func WithCustomFlags(flags Flags) Option {
	return func(s *Starter) {
		s.flags = flags
	}
}
