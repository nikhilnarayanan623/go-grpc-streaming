package service

import (
	"errors"
	"stream-service/pkg/mock/mock_service"
	"stream-service/pkg/mock/mock_usecase"
	"stream-service/pkg/models/request"
	"stream-service/pkg/pb"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestUpload(t *testing.T) {

	testCases := map[string]struct {
		buildStub func(mockStream *mock_service.MockStreamService_UploadServer,
			mockUsecase *mock_usecase.MockStreamUseCase)
		expectedStatusCode codes.Code
	}{
		"error_on_receive_stream_should_return_invalid_argument_code": {

			buildStub: func(mockStream *mock_service.MockStreamService_UploadServer,
				mockUsecase *mock_usecase.MockStreamUseCase) {
				mockStream.EXPECT().Recv().Times(1).
					Return(nil, errors.New("error_on_receive_stream"))
			},
			expectedStatusCode: codes.InvalidArgument,
		},
		"empty_file_info_should_return_invalid_argument_code": {
			// input: request.FileDetails{},
			buildStub: func(mockStream *mock_service.MockStreamService_UploadServer,
				mockUsecase *mock_usecase.MockStreamUseCase) {
				// returning an empty file metadata
				mockStream.EXPECT().Recv().Times(1).
					Return(&pb.UploadRequest{}, nil)
			},
			expectedStatusCode: codes.InvalidArgument,
		},
		"error_to_upload_file_details_should_return_internal_error_code": {
			buildStub: func(mockStream *mock_service.MockStreamService_UploadServer,
				mockUsecase *mock_usecase.MockStreamUseCase) {
				// return file details on first request
				mockStream.EXPECT().Recv().Times(1).
					Return(&pb.UploadRequest{
						File: &pb.UploadRequest_Info{
							Info: &pb.FileMetaData{
								Name:        "fileName",
								ContentType: "content-type",
							},
						},
					}, nil)

					// expecting the upload file with same data
				mockUsecase.EXPECT().UploadFileDetails(gomock.Any(), request.FileDetails{
					Name:        "fileName",
					ContentType: "content-type",
				}).Return("", errors.New("failed_to_upload_file_details"))
			},
			expectedStatusCode: codes.Internal,
		},
	}

	for name, test := range testCases {

		test := test
		t.Run(name, func(t *testing.T) {

			t.Parallel()

			ctl := gomock.NewController(t)
			// create mock upload stream server(grpc) and usecase
			uploadStreamServer := mock_service.NewMockStreamService_UploadServer(ctl)
			mockUsecase := mock_usecase.NewMockStreamUseCase(ctl)

			streamSrv := NewStreamService(mockUsecase)

			// call build stub with upload stream server and mock usecase
			test.buildStub(uploadStreamServer, mockUsecase)

			// call the actual upload with the mock upload stream server
			err := streamSrv.Upload(uploadStreamServer)
			// covert the error into grpc code and check
			actualCode := status.Code(err)
			assert.Equal(t, test.expectedStatusCode, actualCode)
		})
	}

}
