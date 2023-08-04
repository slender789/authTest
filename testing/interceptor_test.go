package test

import (
	"context"
	"testing"

	lib "git.phi-1.com/libauthorization"
	mocks "git.phi-1.com/libauthorization/mock"
	"github.com/golang/mock/gomock"
	"google.golang.org/grpc"
)

func TestInterceptor(t *testing.T) {
	const (
		user    = "1"
		profile = "1"
	)
	ctx := buildGrpcContext(user, profile)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	rulesMock := mocks.NewMockRules(mockCtrl)
	clientMock := mocks.NewMockPermissionApi(mockCtrl)

	subj := lib.RequestOwner{
		User:         user,
		UsersProfile: profile,
	}

	rulesMock.EXPECT().Authorize(clientMock, subj).Return(nil).Times(1)

	configMap := make(map[string]func(interface{}) []lib.Rules)
	configMap["CreatePost"] = func(req interface{}) []lib.Rules {
		rules := []lib.Rules{
			rulesMock,
		}
		return rules
	}

	config := lib.Config{
		Rules: configMap,
	}

	UnaryServerInterceptor := lib.AuthLibraryInterceptor(config, clientMock)

	var request struct{}
	info := &grpc.UnaryServerInfo{FullMethod: "CreatePost"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return request, nil
	}

	_, err := UnaryServerInterceptor(ctx, request, info, handler)
	if err != nil {
		t.Errorf("Wrong")
	}
}
