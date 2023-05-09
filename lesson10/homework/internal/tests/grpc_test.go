package tests

import (
	"context"
	"google.golang.org/grpc/credentials/insecure"
	"homework10/internal/adapters/userrepo"
	grpc2 "homework10/internal/ports/grpc"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"homework10/internal/adapters/adrepo"
	"homework10/internal/app"
	grpcPort "homework10/internal/ports/grpc"
)

func makeClient(t *testing.T) (context.Context, grpc2.AdServiceClient) {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})

	svc := grpcPort.NewService(app.NewApp(adrepo.New(), userrepo.New()))
	grpc2.RegisterAdServiceServer(srv, svc)

	go func() {
		assert.NoError(t, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err, "grpc.DialContext")

	t.Cleanup(func() {
		conn.Close()
	})

	client := grpc2.NewAdServiceClient(conn)
	return ctx, client
}

func TestGRRPCCreateUser(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})

	svc := grpcPort.NewService(app.NewApp(adrepo.New(), userrepo.New()))
	grpc2.RegisterAdServiceServer(srv, svc)

	go func() {
		assert.NoError(t, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err, "grpc.DialContext")

	t.Cleanup(func() {
		conn.Close()
	})

	client := grpc2.NewAdServiceClient(conn)
	res, err := client.CreateUser(ctx, &grpc2.CreateUserRequest{Nickname: "Oleg"})
	assert.NoError(t, err, "client.GetUser")

	assert.Equal(t, "Oleg", res.Name)
}

func TestGRPCCreateAd(t *testing.T) {
	ctx, client := makeClient(t)

	_, _ = client.CreateUser(ctx, &grpc2.CreateUserRequest{Nickname: "Oleg"})
	res, err := client.CreateAd(ctx, &grpc2.CreateAdRequest{Title: "title", Text: "text", UserId: 0})
	assert.NoError(t, err, "client.Ad")
	assert.Equal(t, "title", res.Title)
	assert.Equal(t, "text", res.Text)
	assert.Equal(t, false, res.Published)
	assert.Equal(t, int64(0), res.Id)
	assert.Equal(t, int64(0), res.AuthorId)
}

func TestGRPCChangeAdStatus(t *testing.T) {
	ctx, client := makeClient(t)

	_, _ = client.CreateUser(ctx, &grpc2.CreateUserRequest{Nickname: "Oleg"})
	_, _ = client.CreateAd(ctx, &grpc2.CreateAdRequest{Title: "title", Text: "text", UserId: 0})
	res, err := client.ChangeAdStatus(ctx, &grpc2.ChangeAdStatusRequest{AdId: 0, UserId: 0, Published: true})
	assert.NoError(t, err, "client.Ad")
	assert.Equal(t, "title", res.Title)
	assert.Equal(t, "text", res.Text)
	assert.Equal(t, true, res.Published)
	assert.Equal(t, int64(0), res.Id)
	assert.Equal(t, int64(0), res.AuthorId)
}

func TestGRPCChangeAd(t *testing.T) {
	ctx, client := makeClient(t)

	_, _ = client.CreateUser(ctx, &grpc2.CreateUserRequest{Nickname: "Oleg"})
	_, _ = client.CreateAd(ctx, &grpc2.CreateAdRequest{Title: "title", Text: "text", UserId: 0})
	res, err := client.UpdateAd(ctx, &grpc2.UpdateAdRequest{AdId: 0, UserId: 0, Title: "new title", Text: "new text"})
	assert.NoError(t, err, "client.Ad")
	assert.Equal(t, "new title", res.Title)
	assert.Equal(t, "new text", res.Text)
	assert.Equal(t, false, res.Published)
	assert.Equal(t, int64(0), res.Id)
	assert.Equal(t, int64(0), res.AuthorId)
}

func TestGRPCListAd(t *testing.T) {
	ctx, client := makeClient(t)

	_, _ = client.CreateUser(ctx, &grpc2.CreateUserRequest{Nickname: "Oleg"})
	_, _ = client.CreateAd(ctx, &grpc2.CreateAdRequest{Title: "title", Text: "text", UserId: 0})
	_, _ = client.ChangeAdStatus(ctx, &grpc2.ChangeAdStatusRequest{AdId: 0, UserId: 0, Published: true})
	_, _ = client.CreateAd(ctx, &grpc2.CreateAdRequest{Title: "title", Text: "text", UserId: 0})
	_, err := client.CreateAd(ctx, &grpc2.CreateAdRequest{Title: "title", Text: "text", UserId: 0})
	assert.NoError(t, err, "client.Ad")
	res, err := client.ListAds(ctx, &grpc2.GetAds{Pub: "true", Author: "0", Date: "all"})
	assert.NoError(t, err, "client.Ad")
	assert.Len(t, res.List, 1)
	assert.Equal(t, "title", res.List[0].Title)
	assert.Equal(t, "text", res.List[0].Text)
	assert.Equal(t, true, res.List[0].Published)
	assert.Equal(t, int64(0), res.List[0].Id)
	assert.Equal(t, int64(0), res.List[0].AuthorId)
}

func TestGRPCDeleteAd(t *testing.T) {
	ctx, client := makeClient(t)

	_, _ = client.CreateUser(ctx, &grpc2.CreateUserRequest{Nickname: "Oleg"})
	_, _ = client.CreateAd(ctx, &grpc2.CreateAdRequest{Title: "title", Text: "text", UserId: 0})
	res, err := client.DeleteAd(ctx, &grpc2.DeleteAdRequest{AdId: 0, AuthorId: 0})
	assert.NoError(t, err, "client.Ad")
	assert.Equal(t, "ad was successfully deleted", res.Success)
}

func TestGRPCDeleteUser(t *testing.T) {
	ctx, client := makeClient(t)

	_, _ = client.CreateUser(ctx, &grpc2.CreateUserRequest{Nickname: "Oleg"})
	res, err := client.DeleteUser(ctx, &grpc2.DeleteUserRequest{Id: 0})
	assert.NoError(t, err, "client.Ad")
	assert.Equal(t, "user was successfully deleted", res.Success)
}

func TestGRPCGetUser(t *testing.T) {
	ctx, client := makeClient(t)

	_, _ = client.CreateUser(ctx, &grpc2.CreateUserRequest{Nickname: "Oleg"})
	res, err := client.GetUser(ctx, &grpc2.GetUserRequest{Id: 0})
	assert.NoError(t, err, "client.Ad")
	assert.Equal(t, int64(0), res.Id)
	assert.Equal(t, "Oleg", res.Name)
	assert.Equal(t, "", res.Email)
}
