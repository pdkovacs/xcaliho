package xcalistores3_test

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"testing"
	xcalistores3 "xcalistore-s3"

	"github.com/stretchr/testify/suite"
)

type s3StoreTest struct {
	suite.Suite
	ctx   context.Context
	store *xcalistores3.DrawingStore
}

func TestXCaliStoreS3(t *testing.T) {
	suite.Run(t, &s3StoreTest{})
}

func (s *s3StoreTest) SetupSuite() {
	fmt.Println("SetupSuite running...")
	s.ctx = context.Background()
	var err error
	s.store, err = xcalistores3.NewStore(s.ctx, "test-xcalidrawings")
	s.NoError(err)
}

func (s *s3StoreTest) TestGetAllowedCredentials() {
	creds, getCredsErr := s.store.GetAllowedCredentials(s.ctx)
	s.NoError(getCredsErr)
	s.Equal("qwer\n", creds)
}

func (s *s3StoreTest) TestCreateListSessions() {
	_, createSessErr := s.store.CreateSession(s.ctx)
	s.NoError(createSessErr)
	sessId1, createSessErr1 := s.store.CreateSession(s.ctx)
	s.NoError(createSessErr1)

	sessionList, listErr := s.store.ListSessions(s.ctx)
	s.NoError(listErr)

	s.Equal(1, len(sessionList))
	s.Contains(sessionList, sessId1)
}

func (s *s3StoreTest) TestListDrawings() {
	inputTitles := []string{
		"some title",
		"some other title",
	}
	inputContents := []string{
		"some content",
		"some other content",
	}

	for index := range inputTitles {
		contentReader := strings.NewReader(inputContents[index])
		errPut := s.store.PutDrawing(s.ctx, inputTitles[index], contentReader)
		s.NoError(errPut)
	}

	outputTitles, getTitlesErr := s.store.ListDrawingTitles(s.ctx)
	s.NoError(getTitlesErr)
	s.Equal(2, len(outputTitles))

	slices.Sort(inputTitles)
	slices.Sort(outputTitles)
	s.Equal(inputTitles, outputTitles)
}

func (s *s3StoreTest) TestGetDrawing() {
	inputTitles := []string{
		"some title",
		"some other title",
	}
	inputContents := []string{
		"some content",
		"some other content",
	}

	for index := range inputTitles {
		contentReader := strings.NewReader(inputContents[index])
		errPut := s.store.PutDrawing(s.ctx, inputTitles[index], contentReader)
		s.NoError(errPut)
	}

	content1, getContentErr1 := s.store.GetDrawing(s.ctx, inputTitles[1])
	s.NoError(getContentErr1)
	s.Equal(inputContents[1], content1)
}
