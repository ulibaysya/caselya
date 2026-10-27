package service

import (
	"context"
	"fmt"

	"github.com/ulibaysya/caselya/internal/modules/item"
)

func (s *Service) Pick(ctx context.Context, items []item.WithWeight, n int) ([]int, error) {
	totalWeight := int64(0)

	for _, i := range items {
		if i.Weight < 0 {
			return nil, fmt.Errorf("item %d has negative weight: %d", i.ID, i.Weight)
		}
		totalWeight += i.Weight
	}

	if totalWeight <= 0 {
		return nil, fmt.Errorf("non-positive weight: %d", totalWeight)
	}

	results := make([]int, 0, n)
	for range n {
		result, err := s.singlePick(items, totalWeight)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}

func (s *Service) singlePick(items []item.WithWeight, totalWeight int64) (int, error) {
	roll := s.rng.Int64() % totalWeight

	weightAccumulation := int64(0)
	for i, item := range items {
		weightAccumulation += item.Weight
		if roll < weightAccumulation {
			return i, nil
		}
	}
	return 0, fmt.Errorf("unable to pick result for roll %v, total weight %v", roll, totalWeight)
}
