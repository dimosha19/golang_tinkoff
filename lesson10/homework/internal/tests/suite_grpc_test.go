package tests

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"homework10/internal/adapters/adrepo"
	"homework10/internal/adapters/userrepo"
	"homework10/internal/app"
	grpc2 "homework10/internal/ports/grpc"
	grpcPort "homework10/internal/ports/grpc"
	"net"
	"sync"
	"testing"
	"time"
)

func Client(t *testing.T) (*bufconn.Listener, *grpc.Server, context.Context, grpc2.AdServiceClient, context.CancelFunc, *grpc.ClientConn) {
	lis := bufconn.Listen(1024 * 1024)

	srv := grpc.NewServer()

	svc := grpcPort.NewService(app.NewApp(adrepo.New(), userrepo.New()))
	grpc2.RegisterAdServiceServer(srv, svc)
	go func() {
		assert.NoError(t, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	conn, _ := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer), grpc.WithTransportCredentials(insecure.NewCredentials()))

	client := grpc2.NewAdServiceClient(conn)
	return lis, srv, ctx, client, cancel, conn
}

type grpcSuite struct {
	suite.Suite
	client                grpc2.AdServiceClient
	userNick, title, text string
	userID                int64
	ctx                   context.Context
	cancel                context.CancelFunc
	conn                  *grpc.ClientConn
	lis                   *bufconn.Listener
	srv                   *grpc.Server
	mu                    sync.Mutex
}

func (suite *grpcSuite) SetupTest() {
	suite.userNick = "Oleg"
	suite.userID = 0
	suite.title = "title"
	suite.text = "text"
	suite.lis, suite.srv, suite.ctx, suite.client, suite.cancel, suite.conn = Client(&testing.T{})
}

func (suite *grpcSuite) TearDownTest() {
	suite.lis.Close()
	suite.srv.Stop()
	suite.cancel()
	suite.conn.Close()
}

func (suite *grpcSuite) TestGRPCCreateAd() {
	_, _ = suite.client.CreateUser(suite.ctx, &grpc2.CreateUserRequest{Nickname: suite.userNick})
	res, err := suite.client.CreateAd(suite.ctx, &grpc2.CreateAdRequest{Title: suite.title, Text: suite.text, UserId: suite.userID})

	suite.NoError(err, "client.Ad")
	suite.Equal("title", res.Title)
	suite.Equal("text", res.Text)
	suite.Equal(false, res.Published)
	suite.Equal(int64(0), res.Id)
	suite.Equal(int64(0), res.AuthorId)
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(grpcSuite))
}
