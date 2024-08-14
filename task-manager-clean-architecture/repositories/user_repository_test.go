package repositories

import (
	"context"
	"errors"
	"task-manager-api/domains"
	"task-manager-api/mocks"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	mockCollection   *mocks.CollectionInteface
	mockSingleResult *mocks.SingleResultInterface
	repo             *UserRepositoryMongo
}

func (suite *UserRepositoryTestSuite) SetupTest() {
	suite.mockCollection = mocks.NewCollectionInteface(suite.T())
	suite.mockSingleResult = mocks.NewSingleResultInterface(suite.T())

	suite.repo = NewUserRepositoryMongo(suite.mockCollection)
}

// Test for CreateUser
func (suite *UserRepositoryTestSuite) TestCreateUser_Success() {
	user := &domains.User{
		Username: "abenezer",
		Password: "securepassword",
	}

	suite.mockCollection.On("InsertOne", mock.Anything, user).Return(nil, nil).Once()

	err := suite.repo.CreateUser(context.TODO(), user)

	suite.NoError(err, "Expected no error from CreateUser")
	suite.mockCollection.AssertExpectations(suite.T())
}

// Test for FindUserByUsername
func (suite *UserRepositoryTestSuite) TestFindUserByUsername_Success() {
	username := "abenezer"
	expectedUser := &domains.User{
		Username: username,
		Password: "securepassword",
	}

	suite.mockSingleResult.On("Decode", &domains.User{}).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*domains.User)
		*arg = *expectedUser
	}).Return(nil).Once()

	suite.mockCollection.On("FindOne", mock.Anything, bson.M{"username": username}).Return(suite.mockSingleResult).Once()

	user, err := suite.repo.FindUserByUsername(context.TODO(), username)

	suite.NoError(err, "Expected no error from FindUserByUsername")
	suite.Equal(expectedUser, user, "Expected user to be returned")
	suite.mockCollection.AssertExpectations(suite.T())
	suite.mockSingleResult.AssertExpectations(suite.T())
}

// Test for FindUserByUsername - Not Found
func (suite *UserRepositoryTestSuite) TestFindUserByUsername_NotFound() {
	username := "abenezer"

	suite.mockSingleResult.On("Decode", mock.Anything).Return(mongo.ErrNoDocuments).Once()
	suite.mockCollection.On("FindOne", mock.Anything, bson.M{"username": username}).Return(suite.mockSingleResult).Once()

	user, err := suite.repo.FindUserByUsername(context.TODO(), username)

	suite.NoError(err, "Expected no error when no user is found")
	suite.Nil(user, "Expected no user to be returned")
	suite.mockCollection.AssertExpectations(suite.T())
	suite.mockSingleResult.AssertExpectations(suite.T())
}

// Test for UserExists
func (suite *UserRepositoryTestSuite) TestUserExists_True() {
	username := "abenezer"
	suite.mockCollection.On("CountDocuments", mock.Anything, bson.M{"username": username}).Return(int64(1), nil).Once()

	exists, err := suite.repo.UserExists(context.TODO(), username)

	suite.NoError(err, "Expected no error from UserExists")
	suite.True(exists, "Expected user to exist")
	suite.mockCollection.AssertExpectations(suite.T())
}

// Test for UserExists - Not Found
func (suite *UserRepositoryTestSuite) TestUserExists_False() {
	username := "nonexistentuser"
	suite.mockCollection.On("CountDocuments", mock.Anything, bson.M{"username": username}).Return(int64(0), nil).Once()

	exists, err := suite.repo.UserExists(context.TODO(), username)

	suite.NoError(err, "Expected no error from UserExists")
	suite.False(exists, "Expected user to not exist")
	suite.mockCollection.AssertExpectations(suite.T())
}

// Test for UserExists - Database Error
func (suite *UserRepositoryTestSuite) TestUserExists_DBError() {
	username := "abenezer"
	mockError := errors.New("database error")

	suite.mockCollection.On("CountDocuments", mock.Anything, bson.M{"username": username}).Return(int64(0), mockError).Once()

	exists, err := suite.repo.UserExists(context.TODO(), username)

	suite.Error(err, "Expected error from UserExists when database error occurs")
	suite.False(exists, "Expected user to not exist on error")
	suite.Equal(mockError, err, "Expected the same error to be returned")
	suite.mockCollection.AssertExpectations(suite.T())
}

// Test Runner
func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}
