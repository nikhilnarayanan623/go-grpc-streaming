package usecase

import (
	"context"
	"errors"
	"io"
	"stream-service/pkg/mock/mock_file"
	"stream-service/pkg/mock/mock_repo"
	"stream-service/pkg/models/request"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUploadFileDetails(t *testing.T) {

	testCases := map[string]struct {
		buildStub         func(mockRepo *mock_repo.MockStreamRepository)
		input             request.FileDetails
		isExpectingOutput bool
		expectedError     error
	}{
		"db_error_should_return_error": {
			input: request.FileDetails{
				Name:        "Name",
				ContentType: "new-type",
			},
			buildStub: func(mockRepo *mock_repo.MockStreamRepository) {
				mockRepo.EXPECT().SaveFileDetails(gomock.Any(), gomock.Any()).Times(1).
					Return(errors.New("db error"))
			},
			isExpectingOutput: false,
			expectedError:     errors.New("db error"),
		},
		"successful_to_upload_return_id": {
			buildStub: func(mockRepo *mock_repo.MockStreamRepository) {
				mockRepo.EXPECT().SaveFileDetails(gomock.Any(), gomock.Any()).Times(1).
					Return(nil)
			},
			input: request.FileDetails{
				Name:        "File_Name",
				ContentType: "content_type",
			},
			isExpectingOutput: true,
			expectedError:     nil,
		},
	}

	for name, test := range testCases {
		test := test
		t.Run(name, func(t *testing.T) {

			t.Parallel()

			ctl := gomock.NewController(t)
			repo := mock_repo.NewMockStreamRepository(ctl)

			test.buildStub(repo)
			usecase := NewStreamUseCase(repo, nil)

			out, err := usecase.UploadFileDetails(context.TODO(), test.input)

			if test.isExpectingOutput {
				assert.NotEmpty(t, out, "expecting output of random generated file id")
			} else {
				assert.Empty(t, out, "not expecting out")
			}

			if test.expectedError == nil {
				assert.NoError(t, err)
			} else {
				// in usecase error is wrapping so checking error is containing the expected error
				assert.ErrorContains(t, err, test.expectedError.Error(), "should contain this error string")
			}

		})
	}
}

func TestUploadFileAsStream(t *testing.T) {

	testCases := map[string]struct {
		buildStub func(t *testing.T, mockFileHandler *mock_file.MockHandler)
		input     string // initial file information
		// data      []byte              // data which will stream
		// running this function concurrently to get data
		sendAndCheck  func(t *testing.T, cancel context.CancelFunc, dataChan chan<- []byte, errChan chan error, expError error)
		expectedError error
	}{
		"failed_to_create_dir_should_return_error": {
			input: "file_id",
			buildStub: func(t *testing.T, mockFileHandler *mock_file.MockHandler) {

				folderPath := generateFolderPath("file_id")
				// expect call mkdir and returning an error
				mockFileHandler.EXPECT().MkdirAll(folderPath, gomock.Any()).Times(1).
					Return(errors.New("failed to create directory"))
			},
			expectedError: errors.New("failed to create directory"),
			sendAndCheck: func(t *testing.T, cancel context.CancelFunc, dataChan chan<- []byte,
				errChan chan error, expError error) {

				// wait for the error no data send
				err := <-errChan
				if expError != nil {
					assert.ErrorContains(t, err, expError.Error())
				} else {
					assert.NoError(t, err)
				}
			},
		},
		"failed_to_create_file_should_return_error": {
			input: "file_id",
			buildStub: func(t *testing.T, mockFileHandler *mock_file.MockHandler) {

				folderPath := generateFolderPath("file_id")
				// expect call mkdir and returning no error
				mockFileHandler.EXPECT().MkdirAll(folderPath, gomock.Any()).Times(1).
					Return(nil)

				filePath := generateFilePath(folderPath, "file_id")
				// expecting call to create and returning nil file with an error
				mockFileHandler.EXPECT().Create(filePath).Times(1).
					Return(nil, errors.New("create file error"))
			},
			expectedError: errors.New("create file error"),
			sendAndCheck: func(t *testing.T, cancel context.CancelFunc, dataChan chan<- []byte,
				errChan chan error, expError error) {

				// wait for the error no data send
				err := <-errChan
				if expError != nil {
					assert.ErrorContains(t, err, expError.Error())
				} else {
					assert.NoError(t, err)
				}
			},
		},
		"file_write_error_should_return_error": {
			input: "file_id",
			buildStub: func(t *testing.T, mockFileHandler *mock_file.MockHandler) {
				folderPath := generateFolderPath("file_id")
				// expect call mkdir and returning no error
				mockFileHandler.EXPECT().MkdirAll(folderPath, gomock.Any()).Times(1).
					Return(nil)

				// create a mock file to send data
				ctl := gomock.NewController(t)
				mockFile := mock_file.NewMockFile(ctl)

				// always expecting a call for closing file
				mockFile.EXPECT().Close().Times(1).Return(nil)

				// expect a call to mock file write and returning an error
				mockFile.EXPECT().Write(gomock.Any()).Times(1).
					Return(0, errors.New("failed to write data on file"))

				filePath := generateFilePath(folderPath, "file_id")

				// expecting call to create and returning no mock file with error
				mockFileHandler.EXPECT().Create(filePath).Times(1).
					Return(mockFile, nil)

			},
			expectedError: errors.New("failed to write data on file"),
			sendAndCheck: func(t *testing.T, cancel context.CancelFunc, dataChan chan<- []byte,
				errChan chan error, expError error) {

				// send one data through channel and wait for the error
				dataChan <- []byte("first data")
				err := <-errChan
				if expError != nil {
					assert.ErrorContains(t, err, expError.Error())
				} else {
					assert.NoError(t, err)
				}
			},
		},
		"no_data_upload_will_function_return_after_max_function_wait": {
			input: "file_id",
			buildStub: func(t *testing.T, mockFileHandler *mock_file.MockHandler) {
				folderPath := generateFolderPath("file_id")
				// expect call mkdir and returning no error
				mockFileHandler.EXPECT().MkdirAll(folderPath, gomock.Any()).Times(1).
					Return(nil)

				// create a mock file to send data
				ctl := gomock.NewController(t)
				mockFile := mock_file.NewMockFile(ctl)

				// always expecting a call for closing file
				mockFile.EXPECT().Close().Times(1).Return(nil)

				filePath := generateFilePath(folderPath, "file_id")

				// expecting call to create and returning no mock file with error
				mockFileHandler.EXPECT().Create(filePath).Times(1).
					Return(mockFile, nil)

			},
			expectedError: nil,

			sendAndCheck: func(t *testing.T, cancel context.CancelFunc, dataChan chan<- []byte,
				errChan chan error, expError error) {
				/**
				*  not sending data or expecting an error
				*  considering the client who calling this function
				*  not sending data will return with in max time out on usecase
				* if its not return with in 30s(test default timeout ) will fire an error by test

				**/
			},
		},
		"client_side_error_send_on_chan_should_return": {
			input: "file_id",
			buildStub: func(t *testing.T, mockFileHandler *mock_file.MockHandler) {
				folderPath := generateFolderPath("file_id")
				// expect call mkdir and returning no error
				mockFileHandler.EXPECT().MkdirAll(folderPath, gomock.Any()).Times(1).
					Return(nil)

				// create a mock file to send data
				ctl := gomock.NewController(t)
				mockFile := mock_file.NewMockFile(ctl)

				// always expecting a call for closing file
				mockFile.EXPECT().Close().Times(1).Return(nil)

				filePath := generateFilePath(folderPath, "file_id")

				// expecting call to create and returning no mock file with error
				mockFileHandler.EXPECT().Create(filePath).Times(1).
					Return(mockFile, nil)

			},
			expectedError: errors.New("failed to write data on file"),

			sendAndCheck: func(t *testing.T, cancel context.CancelFunc, dataChan chan<- []byte,
				errChan chan error, expError error) {
				/**
				* sending error from client through errChan should return function
				* if its not return with in 30s(test default timeout ) will fire an error by test
				**/
				errChan <- errors.New("error on client side to send data")
			},
		},

		"cancel_on_context_should_return_function": {
			input: "file_id",
			buildStub: func(t *testing.T, mockFileHandler *mock_file.MockHandler) {
				folderPath := generateFolderPath("file_id")
				// expect call mkdir and returning no error
				mockFileHandler.EXPECT().MkdirAll(folderPath, gomock.Any()).Times(1).
					Return(nil)

				// create a mock file to send data
				ctl := gomock.NewController(t)
				mockFile := mock_file.NewMockFile(ctl)

				// always expecting a call for closing file
				mockFile.EXPECT().Close().Times(1).Return(nil)

				filePath := generateFilePath(folderPath, "file_id")

				// expecting call to create and returning no mock file with error
				mockFileHandler.EXPECT().Create(filePath).Times(1).
					Return(mockFile, nil)

			},
			expectedError: nil,

			sendAndCheck: func(t *testing.T, cancel context.CancelFunc, dataChan chan<- []byte,
				errChan chan error, expError error) {
				/**
				* calling ctx cancel func should return usecase function
				* if its not return with in 30s(test default timeout ) will fire an error by test
				**/
				cancel()
			},
		},

		"successful_send_5_data_should_write_data": {
			input: "file_id",
			buildStub: func(t *testing.T, mockFileHandler *mock_file.MockHandler) {
				folderPath := generateFolderPath("file_id")
				// expect call mkdir and returning no error
				mockFileHandler.EXPECT().MkdirAll(folderPath, gomock.Any()).Times(1).
					Return(nil)

				// create a mock file to send data
				ctl := gomock.NewController(t)
				mockFile := mock_file.NewMockFile(ctl)

				// always expecting a call for closing file
				mockFile.EXPECT().Close().Times(1).Return(nil)

				// expect a call to mock file write and returning an error
				mockFile.EXPECT().Write(gomock.Any()).Times(5).
					Return(0, nil)

				filePath := generateFilePath(folderPath, "file_id")

				// expecting call to create and returning no mock file with error
				mockFileHandler.EXPECT().Create(filePath).Times(1).
					Return(mockFile, nil)

			},
			expectedError: nil,

			sendAndCheck: func(t *testing.T, cancel context.CancelFunc, dataChan chan<- []byte,
				errChan chan error, expError error) {

				// send five data
				for i := 1; i <= 5; i++ {
					dataChan <- []byte("data")
				}
				// send EOF to notify stream completed
				errChan <- io.EOF
			},
		},
	}

	for name, test := range testCases {

		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// create mocks
			ctl := gomock.NewController(t)
			fileHandler := mock_file.NewMockHandler(ctl)

			test.buildStub(t, fileHandler)
			streamUseCase := NewStreamUseCase(nil, fileHandler)

			ctx, cancel := context.WithCancel(context.Background())
			dataChan, errChan := make(chan []byte), make(chan error)

			// run the send and check func on separate goroutine to check error and response
			go test.sendAndCheck(t, cancel, dataChan, errChan, test.expectedError)

			// running this function direct on this test helps to identify orphaned(a state )
			streamUseCase.UploadFileAsStream(ctx, test.input, dataChan, errChan)
		})
	}
}
