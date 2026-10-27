package service

import "context"

func (s *Service) CreateKey(ctx context.Context, name string, public bool, values []string) error {
	return s.strg.CreateKey(ctx, name, public, values)
}
