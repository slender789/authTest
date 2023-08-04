package test

import (
	"context"
	"testing"

	lib "git.phi-1.com/libauthorization"
	mocks "git.phi-1.com/libauthorization/mock"
	"github.com/golang/mock/gomock"
	"google.golang.org/grpc/metadata"
)

const (
	user            = "1"
	users_profile   = "1"
	request_profile = "1"
)

func buildGrpcContext(user, profile string) context.Context {
	md := metadata.New(map[string]string{
		"user_uuid":  user,
		"profile_id": profile,
	})
	return metadata.NewIncomingContext(context.Background(), md)
}

func TestViewRule(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	clientMock := mocks.NewMockPermissionApi(mockCtrl)
	clientMock.EXPECT().Check("profile:"+users_profile, string(lib.VIEW), "profile_info:"+request_profile).Return(nil).Times(1)

	viewRule := &lib.ViewRule{
		Profile: request_profile,
	}

	subj := &lib.RequestOwner{
		User:         user,
		UsersProfile: users_profile,
	}

	res := viewRule.Authorize(clientMock, *subj)
	if res != nil {
		t.Errorf("Wrong")
	}
}

func TestModify(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	clientMock := mocks.NewMockPermissionApi(mockCtrl)
	clientMock.EXPECT().Check("profile:"+users_profile, string(lib.MODIFY), "profile_info:"+request_profile).Return(nil).Times(1)

	modifyRule := &lib.ModifyRule{
		Profile: request_profile,
	}
	subj := &lib.RequestOwner{
		User:         user,
		UsersProfile: users_profile,
	}

	res := modifyRule.Authorize(clientMock, *subj)
	if res != nil {
		t.Errorf("Wrong")
	}
}

func TestMembershipRule(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	clientMock := mocks.NewMockPermissionApi(mockCtrl)
	clientMock.EXPECT().Check("user:"+user, string(lib.MEMBERSHIP), "profile:"+users_profile).Return(nil).Times(1)

	subjectRule := &lib.MembershipRule{}

	subj := &lib.RequestOwner{
		User:         user,
		UsersProfile: users_profile,
	}

	res := subjectRule.Authorize(clientMock, *subj)
	if res != nil {
		t.Errorf("Wrong")
	}
}
