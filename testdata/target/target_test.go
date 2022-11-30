package target

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kazdevl/tgen/testdata/target/repository"
	"github.com/kazdevl/tgen/testdata/target/thirdparty"
	"github.com/stretchr/testify/assert"
)

func TestSampleService_GetSampleName(t *testing.T) {
	type fields struct {
		SampleRepository func(ctrl *gomock.Controller) repository.IFSampleRepository
		SampleClient     func(ctrl *gomock.Controller) thirdparty.IFSampleClient
	}
	type args struct {
		i int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr error
	}{
		{
			name: "正常",
			fields: fields{
				SampleRepository: func(ctrl *gomock.Controller) IFSampleRepository {
					mock := repository.NewMockIFSampleRepository(ctrl)
					// TODO embed expected args and return values
					mock.EXPECT().GetName(nil).Return(nil, nil)
					return mock
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			s := &SampleService{
				SampleRepository: tt.fields.SampleRepository(ctrl),
				SampleClient:     tt.fields.SampleClient(ctrl),
			}
			got, err := s.GetSampleName(tt.args.i)
			assert.True(t, errors.Is(tt.wantErr, err))
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSampleService_UpdateToRandomName(t *testing.T) {
	type fields struct {
		SampleRepository func(ctrl *gomock.Controller) repository.IFSampleRepository
		SampleClient     func(ctrl *gomock.Controller) thirdparty.IFSampleClient
	}
	type args struct {
		i int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "異常: 29行目のif文",
			fields: fields{
				SampleClient: func(ctrl *gomock.Controller) IFSampleClient {
					mock := thirdparty.NewMockIFSampleClient(ctrl)
					// TODO embed expected args and return values
					mock.EXPECT().GenrateRandomName().Return(nil)
					return mock
				},
				SampleRepository: func(ctrl *gomock.Controller) IFSampleRepository {
					mock := repository.NewMockIFSampleRepository(ctrl)
					// TODO embed expected args and return values
					mock.EXPECT().GetLastSaveTime(nil).Return(nil)
					return mock
				},
			},
		},
		{
			name: "正常",
			fields: fields{
				SampleClient: func(ctrl *gomock.Controller) IFSampleClient {
					mock := thirdparty.NewMockIFSampleClient(ctrl)
					// TODO embed expected args and return values
					mock.EXPECT().GenrateRandomName().Return(nil)
					return mock
				},
				SampleRepository: func(ctrl *gomock.Controller) IFSampleRepository {
					mock := repository.NewMockIFSampleRepository(ctrl)
					// TODO embed expected args and return values
					mock.EXPECT().GetLastSaveTime(nil).Return(nil)
					mock.EXPECT().Update(nil, nil).Return(nil)
					return mock
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			s := &SampleService{
				SampleRepository: tt.fields.SampleRepository(ctrl),
				SampleClient:     tt.fields.SampleClient(ctrl),
			}
			assert.True(t, errors.Is(tt.wantErr, s.UpdateToRandomName(tt.args.i)))
		})
	}
}

func TestSampleService_isValid(t *testing.T) {
	type fields struct {
		SampleRepository func(ctrl *gomock.Controller) repository.IFSampleRepository
		SampleClient     func(ctrl *gomock.Controller) thirdparty.IFSampleClient
	}
	type args struct {
		i          int
		updateName string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "異常: 44行目のif文",
			fields: fields{
				SampleRepository: func(ctrl *gomock.Controller) IFSampleRepository {
					mock := repository.NewMockIFSampleRepository(ctrl)
					// TODO embed expected args and return values
					mock.EXPECT().GetLastSaveTime(nil).Return(nil)
					return mock
				},
			},
		},
		{
			name: "正常",
			fields: fields{
				SampleRepository: func(ctrl *gomock.Controller) IFSampleRepository {
					mock := repository.NewMockIFSampleRepository(ctrl)
					// TODO embed expected args and return values
					mock.EXPECT().GetLastSaveTime(nil).Return(nil)
					return mock
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			s := &SampleService{
				SampleRepository: tt.fields.SampleRepository(ctrl),
				SampleClient:     tt.fields.SampleClient(ctrl),
			}
			assert.Equal(t, tt.want, s.isValid(tt.args.i, tt.args.updateName))
		})
	}
}

func TestSampleService_isUpdatable(t *testing.T) {
	type fields struct {
		SampleRepository func(ctrl *gomock.Controller) repository.IFSampleRepository
		SampleClient     func(ctrl *gomock.Controller) thirdparty.IFSampleClient
	}
	type args struct {
		i    int
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "異常: 57行目のif文",
			fields: fields{
				SampleRepository: func(ctrl *gomock.Controller) IFSampleRepository {
					mock := repository.NewMockIFSampleRepository(ctrl)
					// TODO embed expected args and return values
					mock.EXPECT().GetLastSaveTime(nil).Return(nil)
					return mock
				},
			},
		},
		{
			name: "正常",
			fields: fields{
				SampleRepository: func(ctrl *gomock.Controller) IFSampleRepository {
					mock := repository.NewMockIFSampleRepository(ctrl)
					// TODO embed expected args and return values
					mock.EXPECT().GetLastSaveTime(nil).Return(nil)
					return mock
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			s := &SampleService{
				SampleRepository: tt.fields.SampleRepository(ctrl),
				SampleClient:     tt.fields.SampleClient(ctrl),
			}
			assert.True(t, errors.Is(tt.wantErr, s.isUpdatable(tt.args.i, tt.args.name)))
		})
	}
}
