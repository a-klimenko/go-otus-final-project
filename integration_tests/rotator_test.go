package integration_test

import (
	"context"
	"github.com/a-klimenko/go-otus-final-project/internal/server/grpc/pb"
	"github.com/a-klimenko/go-otus-final-project/internal/storage"
	"github.com/google/uuid"
	"log"
	"net"
	"os"
	"testing"

	sqlstorage "github.com/a-klimenko/go-otus-final-project/internal/storage/sql"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/a-klimenko/go-otus-final-project/internal/app"
	"github.com/a-klimenko/go-otus-final-project/internal/logger"
	"github.com/stretchr/testify/suite"
)

type RotatorTestSuite struct {
	suite.Suite
	storage  *sqlstorage.Storage
	logFile  *os.File
	logger   *logger.Logger
	rotator  *app.App
	client   pb.RotatorClient
	bannerId uuid.UUID
	slotId   uuid.UUID
	groupId  uuid.UUID
}

func (suite *RotatorTestSuite) SetupTest() {
	suite.storage = sqlstorage.New()
	err := suite.storage.Connect()
	if err != nil {
		log.Fatal(err)
	}
	logFile, err := os.CreateTemp("", "test-logs.*.log")
	if err != nil {
		log.Fatal(err)
	}
	suite.logFile = logFile
	suite.logger = logger.New("info", suite.logFile)
	suite.rotator = app.New(suite.logger, suite.storage)
	host := net.JoinHostPort("rotator", "50051")
	conn, err := grpc.Dial(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	suite.client = pb.NewRotatorClient(conn)

	suite.bannerId = uuid.New()
	_, err = suite.storage.Db.ExecContext(
		context.Background(),
		"INSERT INTO banners (id, description) VALUES ($1, $2)",
		suite.bannerId,
		"test banner",
	)
	suite.NoError(err)

	suite.slotId = uuid.New()
	_, err = suite.storage.Db.ExecContext(
		context.Background(),
		"INSERT INTO slots (id, description) VALUES ($1, $2)",
		suite.slotId,
		"test slot",
	)
	suite.NoError(err)

	suite.groupId = uuid.New()
	_, err = suite.storage.Db.ExecContext(
		context.Background(),
		"INSERT INTO groups (id, description) VALUES ($1, $2)",
		suite.groupId,
		"test group",
	)
	suite.NoError(err)
}

func (suite *RotatorTestSuite) TearDownTest() {
	_, err := suite.storage.Db.ExecContext(
		context.Background(),
		"DELETE FROM rotations WHERE banner_id = $1",
		suite.bannerId,
	)
	suite.NoError(err)
	_, err = suite.storage.Db.ExecContext(
		context.Background(),
		"DELETE FROM banners WHERE id = $1",
		suite.bannerId,
	)
	suite.NoError(err)
	_, err = suite.storage.Db.ExecContext(
		context.Background(),
		"DELETE FROM slots WHERE id = $1",
		suite.slotId,
	)
	suite.NoError(err)
	_, err = suite.storage.Db.ExecContext(
		context.Background(),
		"DELETE FROM groups WHERE id = $1",
		suite.groupId,
	)
	suite.NoError(err)
	os.Remove(suite.logFile.Name())
	suite.storage.Close()
}

func (suite *RotatorTestSuite) TestAddBanner() {
	in := &pb.AddBannerRequest{
		BannerID: suite.bannerId.String(),
		SlotID:   suite.slotId.String(),
	}
	_, err := suite.client.AddBanner(context.Background(), in)
	suite.NoError(err)

	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM rotations WHERE banner_id=$1 AND slot_id=$2)`
	err = suite.storage.Db.QueryRowxContext(
		context.Background(), query, suite.bannerId, suite.slotId,
	).Scan(&exists)
	suite.NoError(err)
	suite.True(exists)
}

func (suite *RotatorTestSuite) TestRemoveBanner() {
	addRequest := &pb.AddBannerRequest{
		BannerID: suite.bannerId.String(),
		SlotID:   suite.slotId.String(),
	}
	_, err := suite.client.AddBanner(context.Background(), addRequest)
	suite.NoError(err)

	in := &pb.RemoveBannerRequest{
		SlotID:   suite.slotId.String(),
		BannerID: suite.bannerId.String(),
	}
	_, err = suite.client.RemoveBanner(context.Background(), in)
	suite.NoError(err)

	rotationsQuery := `
				SELECT id, banner_id, slot_id, group_id, clicks, shows, deleted_at
				FROM rotations 
				WHERE slot_id=$1 AND banner_id=$2
	`
	rows, err := suite.storage.Db.QueryxContext(
		context.Background(), rotationsQuery, suite.slotId, suite.bannerId,
	)
	suite.NoError(err)

	for rows.Next() {
		var rotation storage.Rotation
		err := rows.StructScan(&rotation)
		suite.NoError(err)
		suite.True(rotation.DeletedAt.Valid)
	}
}

func (suite *RotatorTestSuite) TestClickBanner() {
	addRequest := &pb.AddBannerRequest{
		BannerID: suite.bannerId.String(),
		SlotID:   suite.slotId.String(),
	}
	_, err := suite.client.AddBanner(context.Background(), addRequest)
	suite.NoError(err)

	in := &pb.ClickBannerRequest{
		SlotID:   suite.slotId.String(),
		BannerID: suite.bannerId.String(),
		GroupID:  suite.groupId.String(),
	}
	_, err = suite.client.ClickBanner(context.Background(), in)
	suite.NoError(err)

	var clicks int
	query := `SELECT clicks FROM rotations WHERE slot_id=$1 AND banner_id=$2 AND group_id=$3`
	err = suite.storage.Db.QueryRowxContext(
		context.Background(), query, suite.slotId, suite.bannerId, suite.groupId,
	).Scan(&clicks)
	suite.NoError(err)
	suite.Equal(1, clicks)
}

func (suite *RotatorTestSuite) TestChooseBanner() {
	addRequest := &pb.AddBannerRequest{
		BannerID: suite.bannerId.String(),
		SlotID:   suite.slotId.String(),
	}
	_, err := suite.client.AddBanner(context.Background(), addRequest)
	suite.NoError(err)

	in := &pb.ChooseBannerRequest{
		SlotID:  suite.slotId.String(),
		GroupID: suite.groupId.String(),
	}
	resp, err := suite.client.ChooseBanner(context.Background(), in)
	suite.NoError(err)
	suite.Equal(resp.BannerID, suite.bannerId.String())
}

func TestServerTestSuite(t *testing.T) {
	suite.Run(t, new(RotatorTestSuite))
}
