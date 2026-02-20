package service

import (
	"context"

	pb "github.com/go-kratos/kratos-layout/api/version/v1"
	"github.com/go-kratos/kratos-layout/internal/version"
)

type VersionService struct {
	pb.UnimplementedVersionServer
}

func NewVersionService() *VersionService {
	return &VersionService{}
}

func (s *VersionService) GetVersion(ctx context.Context, req *pb.GetVersionRequest) (*pb.GetVersionReply, error) {
	return &pb.GetVersionReply{
		Version:   version.Version,
		Commit:    version.Commit,
		BuildTime: version.BuildTime,
	}, nil
}
