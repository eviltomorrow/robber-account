package server

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"net"
	"time"

	"github.com/eviltomorrow/robber-account/internal/middleware"
	"github.com/eviltomorrow/robber-account/internal/model"
	"github.com/eviltomorrow/robber-account/internal/service"
	"github.com/eviltomorrow/robber-account/pkg/pb"
	"github.com/eviltomorrow/robber-core/pkg/grpclb"
	"github.com/eviltomorrow/robber-core/pkg/mysql"
	"github.com/eviltomorrow/robber-core/pkg/system"
	"github.com/eviltomorrow/robber-core/pkg/zlog"
	"github.com/eviltomorrow/robber-core/pkg/znet"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var (
	Host           = "0.0.0.0"
	Port           = 27323
	Endpoints      = []string{}
	RevokeEtcdConn func() error
	Key            = "grpclb/service/account"

	server *grpc.Server
)

type GRPC struct {
	pb.UnimplementedAccountServer
}

func (g *GRPC) Version(ctx context.Context, _ *emptypb.Empty) (*wrapperspb.StringValue, error) {
	var buf bytes.Buffer
	buf.WriteString("Server: \r\n")
	buf.WriteString(fmt.Sprintf("   Robber-account Version (Current): %s\r\n", system.MainVersion))
	buf.WriteString(fmt.Sprintf("   Go Version: %v\r\n", system.GoVersion))
	buf.WriteString(fmt.Sprintf("   Go OS/Arch: %v\r\n", system.GoOSArch))
	buf.WriteString(fmt.Sprintf("   Git Sha: %v\r\n", system.GitSha))
	buf.WriteString(fmt.Sprintf("   Git Tag: %v\r\n", system.GitTag))
	buf.WriteString(fmt.Sprintf("   Git Branch: %v\r\n", system.GitBranch))
	buf.WriteString(fmt.Sprintf("   Build Time: %v\r\n", system.BuildTime))
	buf.WriteString(fmt.Sprintf("   HostName: %v\r\n", system.HostName))
	buf.WriteString(fmt.Sprintf("   IP: %v\r\n", system.IP))
	buf.WriteString(fmt.Sprintf("   Running Time: %v\r\n", system.RunningTime()))
	return &wrapperspb.StringValue{Value: buf.String()}, nil
}

// Create(context.Context, *User) (*wrapperspb.StringValue, error)
// Destroy(context.Context, *wrapperspb.StringValue) (*emptypb.Empty, error)
// List(*emptypb.Empty, Account_ListServer) error
// Find(context.Context, *wrapperspb.StringValue) (*User, error)

func (g *GRPC) Create(ctx context.Context, req *pb.User) (*wrapperspb.StringValue, error) {
	if req.Email == "" {
		return nil, fmt.Errorf("invalid email parameter")
	}
	if req.Phone == "" {
		return nil, fmt.Errorf("invalid phone parameter")
	}

	var user = &model.User{
		NickName: sql.NullString{String: req.NickName},
		Email:    req.Email,
		Phone:    req.Phone,
	}

	uuid, err := service.CreateUser(user)
	if err != nil {
		return nil, err
	}
	return &wrapperspb.StringValue{Value: uuid}, nil
}

func (g *GRPC) Destroy(ctx context.Context, req *wrapperspb.StringValue) (*emptypb.Empty, error) {
	if err := service.RemoveUser(req.Value); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (g *GRPC) List(_ *emptypb.Empty, resp pb.Account_ListServer) error {
	var (
		offset  int64 = 0
		limit   int64 = 50
		timeout       = 10 * time.Second
	)

	for {
		users, err := model.UserWithSelectRange(mysql.DB, offset, limit, timeout)
		if err != nil {
			return err
		}

		for _, user := range users {
			if err := resp.Send(&pb.User{
				Uuid:              user.UUID,
				NickName:          user.NickName.String,
				Email:             user.Email,
				Phone:             user.Phone,
				RegisterTimestamp: user.CreateTimestamp.Format("2006-01-02"),
			}); err != nil {
				return err
			}
		}

		if int64(len(users)) < limit {
			break
		}
		offset += limit
	}
	return nil
}

func (g *GRPC) Find(ctx context.Context, req *wrapperspb.StringValue) (*pb.User, error) {
	user, err := model.UserWithSelectOneByUUID(mysql.DB, req.Value, 10*time.Second)
	if err != nil {
		return nil, err
	}
	return &pb.User{
		Uuid:              user.UUID,
		NickName:          user.NickName.String,
		Email:             user.Email,
		Phone:             user.Phone,
		RegisterTimestamp: user.CreateTimestamp.Format("2006-01-02"),
	}, nil
}

func StartupGRPC() error {
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", Host, Port))
	if err != nil {
		return err
	}

	server = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.UnaryServerRecoveryInterceptor,
			middleware.UnaryServerLogInterceptor,
		),
		grpc.ChainStreamInterceptor(
			middleware.StreamServerRecoveryInterceptor,
			middleware.StreamServerLogInterceptor,
		),
	)

	reflection.Register(server)
	pb.RegisterAccountServer(server, &GRPC{})

	localIp, err := znet.GetLocalIP2()
	if err != nil {
		return fmt.Errorf("get local ip failure, nest error: %v", err)
	}

	close, err := grpclb.Register(Key, localIp, Port, Endpoints, 10)
	if err != nil {
		return fmt.Errorf("register service to etcd failure, nest error: %v", err)
	}
	RevokeEtcdConn = func() error {
		close()
		return nil
	}

	go func() {
		if err := server.Serve(listen); err != nil {
			zlog.Fatal("GRPC Server startup failure", zap.Error(err))
		}
	}()
	return nil
}

func ShutdownGRPC() error {
	if server == nil {
		return nil
	}
	server.Stop()
	return nil
}
