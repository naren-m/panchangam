package gateway

import (
	"context"

	ppb "github.com/naren-m/panchangam/proto"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type MockPanchangamClient struct {
	mock.Mock
}

func (m *MockPanchangamClient) Get(ctx context.Context, in *ppb.GetPanchangamRequest, opts ...grpc.CallOption) (*ppb.GetPanchangamResponse, error) {
	args := m.Called(ctx, in)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ppb.GetPanchangamResponse), args.Error(1)
}
