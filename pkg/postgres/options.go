package postgres

type Option func(*Server)

func WithPgPath(pgPath string) Option {
	return func(s *Server) {
		s.pgPath = pgPath
	}
}
