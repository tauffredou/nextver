package provider

import "github.com/tauffredou/nextver/model"

type MockProvider struct {
}

func (p *MockProvider) GetLatestRelease() *model.Release {
	panic("implement me")
}

func (p *MockProvider) AddRelease(r model.Release) {

}

func (p *MockProvider) GetRelease(name string) (*model.Release, error) {
	return &model.Release{}, nil
}

func NewMockProvider() *MockProvider {
	return &MockProvider{}
}
