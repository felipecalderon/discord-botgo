package store

import (
	"fmt"
	"math/rand"
	"sync"
)

type ImageStore struct {
	images []string
	rng    *rand.Rand
	mu     sync.RWMutex
}

func NewImageStore(rng *rand.Rand) *ImageStore {
	return &ImageStore{
		images: make([]string, 0),
		rng:    rng,
	}
}

func (s *ImageStore) AddImage(url string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.images = append(s.images, url)
}

func (s *ImageStore) GetRandomImage() (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.images) == 0 {
		return "", fmt.Errorf("no hay imágenes disponibles")
	}

	return s.images[s.rng.Intn(len(s.images))], nil
}
