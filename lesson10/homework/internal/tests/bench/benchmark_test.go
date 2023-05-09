package bench

import (
	"context"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"homework10/internal/adapters/adrepo"
	"homework10/internal/adapters/userrepo"
	"homework10/internal/app"
	grpc2 "homework10/internal/ports/grpc"
	grpcPort "homework10/internal/ports/grpc"
	"homework10/internal/tests"
	"net"
	"testing"
	"time"
)

func makeClient(b *testing.B) (context.Context, grpc2.AdServiceClient) {
	lis := bufconn.Listen(1024 * 1024)
	b.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	b.Cleanup(func() {
		srv.Stop()
	})

	svc := grpcPort.NewService(app.NewApp(adrepo.New(), userrepo.New()))
	grpc2.RegisterAdServiceServer(srv, svc)

	go func() {
		assert.NoError(b, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	b.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(b, err, "grpc.DialContext")

	b.Cleanup(func() {
		conn.Close()
	})

	client := grpc2.NewAdServiceClient(conn)
	return ctx, client
}

func createUser(client *tests.TestClient) tests.UserResponse {
	res, _ := client.CreateUser("dimosha", "dmitriy@mail.ru")
	return res
}

func createUserGRPC(ctx context.Context, client grpc2.AdServiceClient) {
	_, _ = client.CreateUser(ctx, &grpc2.CreateUserRequest{Nickname: "Oleg"})
}

func BenchmarkREST(b *testing.B) {
	client := tests.GetTestClient()
	for i := 0; i < b.N; i++ {
		_ = createUser(client)
	}
}

func BenchmarkGRPC(b *testing.B) {
	ctx, client := makeClient(b)
	for i := 0; i < b.N; i++ {
		createUserGRPC(ctx, client)
	}
}
