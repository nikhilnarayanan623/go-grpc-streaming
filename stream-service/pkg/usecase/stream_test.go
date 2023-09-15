package usecase

import (
	"context"
	"errors"
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
			usecase := NewStreamUseCase(repo)

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
