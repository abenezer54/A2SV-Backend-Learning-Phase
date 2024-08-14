package repositories

import (
	"context"
	"errors"
	"testing"
	"time"

	"task-manager-api/domains"
	"task-manager-api/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskRepositoryTestSuite struct {
	suite.Suite
	mockCollection   *mocks.CollectionInteface
	mockCursor       *mocks.CursorInterface
	mockSingleResult *mocks.SingleResultInterface
	repo             *TaskRepoMongo
}

func (suite *TaskRepositoryTestSuite) SetupTest() {
	suite.mockCollection = mocks.NewCollectionInteface(suite.T())
	suite.mockCursor = mocks.NewCursorInterface(suite.T())
	suite.mockSingleResult = mocks.NewSingleResultInterface(suite.T())

	suite.repo = NewTaskRepositoryMongo(suite.mockCollection)
}

// TestCreateTask - tests

func (suite *TaskRepositoryTestSuite) TestCreateTask_Sucess() {
	task := domains.NewTask("MyTitle", "Message from Abenezer Mulugeta", false, time.Now(), primitive.NewObjectID())

	insertResult := &mongo.InsertOneResult{
		InsertedID: task.ID,
	}

	suite.mockCollection.On("InsertOne", mock.Anything, task).Return(insertResult, nil).Once()

	result, err := suite.repo.CreateTask(context.TODO(), task)

	suite.NoError(err, "Expected no error from CreateTask")
	suite.Equal(task, result, "Expected returned task to be the same as the input task")
	suite.mockCollection.AssertExpectations(suite.T())
}

// FindTasksBycreator - tests
func (suite *TaskRepositoryTestSuite) TestFindTasksByCreator_Sucess() {
	creatorID := primitive.NewObjectID()
	task1 := domains.NewTask("task1", "this is description for task1", false, time.Now(), creatorID)
	task2 := domains.NewTask("task2", "this is description for task2", false, time.Now(), creatorID)

	cursor := suite.mockCursor

	// Setup expectations for the cursor's Next, Decode, and Err methods
	cursor.On("Next", mock.Anything).Return(true).Once()
	cursor.On("Decode", mock.Anything).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*domains.Task)
		*arg = *task1
	}).Return(nil).Once() // Add a return value of nil error
	cursor.On("Next", mock.Anything).Return(true).Once()
	cursor.On("Decode", mock.Anything).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*domains.Task)
		*arg = *task2
	}).Return(nil).Once() // Add a return value of nil error
	cursor.On("Next", mock.Anything).Return(false).Once()
	cursor.On("Err").Return(nil).Once()
	cursor.On("Close", mock.Anything).Return(nil).Once()

	// Setup expectation for the collection's Find method
	suite.mockCollection.On("Find", mock.Anything, bson.M{"creator_id": creatorID}).Return(cursor, nil).Once()

	// Call the method under test
	tasks, err := suite.repo.FindTasksByCreator(context.TODO(), creatorID)

	// Assertions
	suite.NoError(err, "Expected no error from FindTasksByCreator")
	suite.Len(tasks, 2, "Expected to find two tasks")
	suite.Equal(task1, tasks[0], "Expected first task to match")
	suite.Equal(task2, tasks[1], "Expected second task to match")

	// Verify that all expectations were met
	suite.mockCollection.AssertExpectations(suite.T())
	cursor.AssertExpectations(suite.T())
}

func (suite *TaskRepositoryTestSuite) TestFindTasksByCreator_IDNotFound() {
    creatorID := primitive.NewObjectID()

    cursor := suite.mockCursor

    // Setup expectations for the cursor's Next and Err methods
    cursor.On("Next", mock.Anything).Return(false).Once()
    cursor.On("Err").Return(nil).Once()
    cursor.On("Close", mock.Anything).Return(nil).Once()

    // Setup expectation for the collection's Find method
    suite.mockCollection.On("Find", mock.Anything, bson.M{"creator_id": creatorID}).Return(cursor, nil).Once()

    // Call the method under test
    tasks, err := suite.repo.FindTasksByCreator(context.TODO(), creatorID)

    // Assertions
    suite.NoError(err, "Expected no error from FindTasksByCreator")
    suite.Len(tasks, 0, "Expected no tasks to be found")

    // Verify that all expectations were met
    suite.mockCollection.AssertExpectations(suite.T())
    cursor.AssertExpectations(suite.T())
}

//FIndTaskByIDAndCreator - tests

func (suite *TaskRepositoryTestSuite) TestFindTaskByIDAndCreator_Success() {
    taskID := primitive.NewObjectID()
    creatorID := primitive.NewObjectID()

   expectedTask := domains.NewTask("MyTitle", "Message from Abenezer Mulugeta", false, time.Now(), primitive.NewObjectID())
   expectedTask.ID = taskID
   expectedTask.CreatorID = creatorID

    
    suite.mockSingleResult.On("Decode", &domains.Task{}).Run(func(args mock.Arguments) {
        arg := args.Get(0).(*domains.Task)
        *arg = *expectedTask
    }).Return(nil).Once()

    suite.mockCollection.On("FindOne", mock.Anything, bson.M{"_id": taskID, "creator_id": creatorID}).Return(suite.mockSingleResult).Once()

    // Call the method under test
    task, err := suite.repo.FindTaskByIDAndCreator(context.TODO(), taskID, creatorID)

    // Assertions
    suite.NoError(err)
    suite.NotNil(task)
    suite.Equal(expectedTask, task)

    // Verify that all expectations were met
    suite.mockCollection.AssertExpectations(suite.T())
    suite.mockSingleResult.AssertExpectations(suite.T())
}

func (suite *TaskRepositoryTestSuite) TestFindTaskByIDAndCreator_NotFound() {
    taskID := primitive.NewObjectID()
    creatorID := primitive.NewObjectID()

    
    suite.mockSingleResult.On("Decode", mock.Anything).Return(mongo.ErrNoDocuments).Once()

    suite.mockCollection.On("FindOne", mock.Anything, bson.M{"_id": taskID, "creator_id": creatorID}).Return(suite.mockSingleResult).Once()

    // Call the method under test
    task, err := suite.repo.FindTaskByIDAndCreator(context.TODO(), taskID, creatorID)

    // Assertions
    suite.NoError(err)
    suite.Nil(task)

    // Verify that all expectations were met
    suite.mockCollection.AssertExpectations(suite.T())
    suite.mockSingleResult.AssertExpectations(suite.T())
}

func (suite *TaskRepositoryTestSuite) TestFindTaskByIDAndCreator_DBError() {
    taskID := primitive.NewObjectID()
    creatorID := primitive.NewObjectID()

    mockError := errors.New("database error")

    suite.mockSingleResult.On("Decode", mock.Anything).Return(mockError).Once()

    suite.mockCollection.On("FindOne", mock.Anything, bson.M{"_id": taskID, "creator_id": creatorID}).Return(suite.mockSingleResult).Once()

    // Call the method under test
    task, err := suite.repo.FindTaskByIDAndCreator(context.TODO(), taskID, creatorID)

    // Assertions
    suite.Error(err)
    suite.Nil(task)
    suite.Equal(mockError, err)

    // Verify that all expectations were met
    suite.mockCollection.AssertExpectations(suite.T())
    suite.mockSingleResult.AssertExpectations(suite.T())
}

// GetTaskByID - tests
func (suite *TaskRepositoryTestSuite) TestGetTaskByID_Success() {
    taskID := primitive.NewObjectID()
    expectedTask := domains.NewTask("MyTitle", "Message from Abenezer Mulugeta", false, time.Now(), primitive.NewObjectID())
    expectedTask.ID = taskID

    suite.mockSingleResult.On("Decode", &domains.Task{}).Run(func(args mock.Arguments) {
        arg := args.Get(0).(*domains.Task)
        *arg = *expectedTask
    }).Return(nil).Once()

    suite.mockCollection.On("FindOne", mock.Anything, bson.M{"_id": taskID}).Return(suite.mockSingleResult).Once()

    // Call the method under test
    task, err := suite.repo.GetTaskByID(taskID.Hex())

    // Assertions
    suite.NoError(err, "Expected no error from GetTaskByID")
    suite.NotNil(task, "Expected task to be returned")
    suite.Equal(expectedTask, task, "Expected returned task to match the expected task")

    // Verify that all expectations were met
    suite.mockCollection.AssertExpectations(suite.T())
    suite.mockSingleResult.AssertExpectations(suite.T())
}

func (suite *TaskRepositoryTestSuite) TestGetTaskByID_NotFound() {
    taskID := primitive.NewObjectID()

    suite.mockSingleResult.On("Decode", mock.Anything).Return(mongo.ErrNoDocuments).Once()
    suite.mockCollection.On("FindOne", mock.Anything, bson.M{"_id": taskID}).Return(suite.mockSingleResult).Once()

    // Call the method under test
    task, err := suite.repo.GetTaskByID(taskID.Hex())

    // Assertions
    suite.Error(err, "Expected an error from GetTaskByID when the task is not found")
    suite.Nil(task, "Expected no task to be returned when not found")
    suite.EqualError(err, "task not found", "Expected 'task not found' error message")

    // Verify that all expectations were met
    suite.mockCollection.AssertExpectations(suite.T())
    suite.mockSingleResult.AssertExpectations(suite.T())
}
 
func (suite *TaskRepositoryTestSuite) TestGetTaskByID_DBError() {
    taskID := primitive.NewObjectID()
    mockError := errors.New("database error")

    suite.mockSingleResult.On("Decode", mock.Anything).Return(mockError).Once()
    suite.mockCollection.On("FindOne", mock.Anything, bson.M{"_id": taskID}).Return(suite.mockSingleResult).Once()

    // Call the method under test
    task, err := suite.repo.GetTaskByID(taskID.Hex())

    // Assertions
    suite.Error(err, "Expected an error from GetTaskByID on database error")
    suite.Nil(task, "Expected no task to be returned on database error")
    suite.Equal(mockError, err, "Expected the same error returned by the database operation")

    // Verify that all expectations were met
    suite.mockCollection.AssertExpectations(suite.T())
    suite.mockSingleResult.AssertExpectations(suite.T())
}

// UpdateTaskByCreatorID - tests
func (suite *TaskRepositoryTestSuite) TestUpdateTaskByCreatorID_Success() {
    task := domains.NewTask("MyTitle", "Message from Abenezer Mulugeta", false, time.Now(), primitive.NewObjectID())

    filter := bson.M{"_id": task.ID, "creator_id": task.CreatorID}
    update := bson.M{
        "$set": bson.M{
            "title":       task.Title,
            "description": task.Description,
            "completed":   task.Completed,
            "due_date":    task.DueDate,
        },
    }

    // Mock the UpdateOne method to return a successful result
    suite.mockCollection.On("UpdateOne", mock.Anything, filter, update).Return(&mongo.UpdateResult{MatchedCount: 1, ModifiedCount: 1}, nil).Once()

    // Call the method under test
    err := suite.repo.UpdateTaskByCreatorID(context.TODO(), task)

    // Assertions
    suite.NoError(err, "Expected no error from UpdateTaskByCreatorID")
    suite.mockCollection.AssertExpectations(suite.T())
}

func (suite *TaskRepositoryTestSuite) TestUpdateTaskByCreatorID_NotFound() {
    task := domains.NewTask("MyTitle", "Message from Abenezer Mulugeta", false, time.Now(), primitive.NewObjectID())

    filter := bson.M{"_id": task.ID, "creator_id": task.CreatorID}
    update := bson.M{
        "$set": bson.M{
            "title":       task.Title,
            "description": task.Description,
            "completed":   task.Completed,
            "due_date":    task.DueDate,
        },
    }

    // Mock the UpdateOne method to return a result indicating no document was matched
    suite.mockCollection.On("UpdateOne", mock.Anything, filter, update).Return(&mongo.UpdateResult{MatchedCount: 0, ModifiedCount: 0}, nil).Once()

    // Call the method under test
    err := suite.repo.UpdateTaskByCreatorID(context.TODO(), task)

    // Assertions
    suite.NoError(err, "Expected no error from UpdateTaskByCreatorID even if no document is found")
    suite.mockCollection.AssertExpectations(suite.T())
}

func (suite *TaskRepositoryTestSuite) TestUpdateTaskByCreatorID_DBError() {
    task := domains.NewTask("MyTitle", "Message from Abenezer Mulugeta", false, time.Now(), primitive.NewObjectID())

    filter := bson.M{"_id": task.ID, "creator_id": task.CreatorID}
    update := bson.M{
        "$set": bson.M{
            "title":       task.Title,
            "description": task.Description,
            "completed":   task.Completed,
            "due_date":    task.DueDate,
        },
    }

    // Mock the UpdateOne method to return an error
    mockError := errors.New("database error")
    suite.mockCollection.On("UpdateOne", mock.Anything, filter, update).Return(nil, mockError).Once()

    // Call the method under test
    err := suite.repo.UpdateTaskByCreatorID(context.TODO(), task)

    // Assertions
    suite.Error(err, "Expected an error from UpdateTaskByCreatorID due to database error")
    suite.Equal(mockError, err, "Expected the error to match the mock error")
    suite.mockCollection.AssertExpectations(suite.T())
}

// DeleteTaskByCreatorID - tests

func (suite *TaskRepositoryTestSuite) TestDeleteTaskByCreatorID_Success() {
    taskID := primitive.NewObjectID()
    creatorID := primitive.NewObjectID()

    filter := bson.M{"_id": taskID, "creator_id": creatorID}

    // Mock the DeleteOne method to return a successful result
    suite.mockCollection.On("DeleteOne", mock.Anything, filter).Return(&mongo.DeleteResult{DeletedCount: 1}, nil).Once()

    // Call the method under test
    err := suite.repo.DeleteTaskByCreatorID(context.TODO(), taskID, creatorID)

    // Assertions
    suite.NoError(err, "Expected no error from DeleteTaskByCreatorID")
    suite.mockCollection.AssertExpectations(suite.T())
}

func (suite *TaskRepositoryTestSuite) TestDeleteTaskByCreatorID_NotFound() {
    taskID := primitive.NewObjectID()
    creatorID := primitive.NewObjectID()

    filter := bson.M{"_id": taskID, "creator_id": creatorID}

    // Mock the DeleteOne method to return a result indicating no document was deleted
    suite.mockCollection.On("DeleteOne", mock.Anything, filter).Return(&mongo.DeleteResult{DeletedCount: 0}, nil).Once()

    // Call the method under test
    err := suite.repo.DeleteTaskByCreatorID(context.TODO(), taskID, creatorID)

    // Assertions
    suite.NoError(err, "Expected no error from DeleteTaskByCreatorID even if no document is found")
    suite.mockCollection.AssertExpectations(suite.T())
}

func (suite *TaskRepositoryTestSuite) TestDeleteTaskByCreatorID_DBError() {
    taskID := primitive.NewObjectID()
    creatorID := primitive.NewObjectID()

    filter := bson.M{"_id": taskID, "creator_id": creatorID}

    // Mock the DeleteOne method to return an error
    mockError := errors.New("database error")
    suite.mockCollection.On("DeleteOne", mock.Anything, filter).Return(nil, mockError).Once()

    // Call the method under test
    err := suite.repo.DeleteTaskByCreatorID(context.TODO(), taskID, creatorID)

    // Assertions
    suite.Error(err, "Expected an error from DeleteTaskByCreatorID due to database error")
    suite.Equal(mockError, err, "Expected the error to match the mock error")
    suite.mockCollection.AssertExpectations(suite.T())
}

func TestTaskRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(TaskRepositoryTestSuite))
}
