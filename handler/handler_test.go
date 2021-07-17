package handler

import (
	"bytes"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	errors "golang.org/x/xerrors"
	storagemock "url-shortener/repository/mock"
)

type HandlerSuite struct {
	suite.Suite
}

func TestHandlerSuite(t *testing.T) {
	suite.Run(t, new(HandlerSuite))
}

func (s *HandlerSuite) SetupSuite() {
	rand.Seed(time.Now().Unix())
}

func (s *HandlerSuite) Test_getClientIP() {
	t := s.T()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		r *http.Request
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test getClientIP if header with X-Forwarded-For",
			args: args{
				r: &http.Request{
					Header: http.Header{
						"X-Forwarded-For": []string{"127.0.0.1"},
					},
				},
			},
			want: "127.0.0.1",
		},
		{
			name: "Test getClientIP if header with X-Real-Ip",
			args: args{
				r: &http.Request{
					Header: http.Header{
						"X-Real-Ip": []string{"198.126.12.12"},
					},
				},
			},
			want: "198.126.12.12",
		},
		{
			name: "Test getClientIP if header with RemoteAddr",
			args: args{
				r: &http.Request{
					RemoteAddr: "",
				},
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getClientIP(tt.args.r); got != tt.want {
				t.Errorf("getClientIP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (s *HandlerSuite) Test_handler_checkValidIP() {
	t := s.T()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		storage *storagemock.MockStorage
	}
	type args struct {
		r *http.Request
	}
	type mockArgs struct {
		clientIP string
		count    int
		err      error
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		want     bool
		wantErr  bool
		mockArgs mockArgs
	}{
		{
			name: "Test checkValidIP is valid",
			fields: fields{
				storage: storagemock.NewMockStorage(ctrl),
			},
			args: args{
				r: &http.Request{
					Header: http.Header{
						"X-Real-Ip": []string{"198.126.12.12"},
					},
				},
			},
			want:    true,
			wantErr: false,
			mockArgs: mockArgs{
				clientIP: "198.126.12.12",
				count:    1,
				err:      nil,
			},
		},
		{
			name: "Test checkValidIP is invalid",
			fields: fields{
				storage: storagemock.NewMockStorage(ctrl),
			},
			args: args{
				r: &http.Request{
					Header: http.Header{
						"X-Real-Ip": []string{"198.126.12.12"},
					},
				},
			},
			want:    false,
			wantErr: false,
			mockArgs: mockArgs{
				clientIP: "198.126.12.12",
				count:    100,
				err:      nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.storage.EXPECT().LoadClientIP(tt.mockArgs.clientIP).Return(tt.mockArgs.count, tt.mockArgs.err)

			h := &handler{
				storage: tt.fields.storage,
			}
			got, err := h.checkValidIP(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkValidIP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("checkValidIP() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func (s *HandlerSuite) Test_handler_UploadURL() {
	t := s.T()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		storage *storagemock.MockStorage
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	type mockArgs struct {
		clientIP        string
		count           int
		err             error
		urlID           string
		url             string
		expireAt        string
		isExpireAtValid bool
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		mockArgs mockArgs
	}{
		{
			name: "Test UploadURL is success",
			fields: fields{
				storage: storagemock.NewMockStorage(ctrl),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: &http.Request{
					Header: http.Header{
						"X-Real-Ip": []string{"198.126.12.12"},
					},
					Body: ioutil.NopCloser(bytes.NewReader([]byte(`{"url":"https://www.google.com","expireAt":"2021-07-18T16:58:30+08:00"}`))),
				},
			},
			mockArgs: mockArgs{
				clientIP:        "198.126.12.12",
				count:           1,
				err:             nil,
				urlID:           "aHR0cHM6Ly93d3cuZ29vZ2xlLmNvbTIwMjEtMDctMThUMTY6NTg6MzArMDg6MDA=",
				url:             "https://www.google.com",
				expireAt:        "2021-07-18T16:58:30+08:00",
				isExpireAtValid: true,
			},
		},
		{
			name: "Test UploadURL is fail because invalid expiredAt",
			fields: fields{
				storage: storagemock.NewMockStorage(ctrl),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: &http.Request{
					Header: http.Header{
						"X-Real-Ip": []string{"198.126.12.12"},
					},
					Body: ioutil.NopCloser(bytes.NewReader([]byte(`{"url":"https://www.google.com","expireAt":"2019-07-17T16:58:30+08:00"}`))),
				},
			},
			mockArgs: mockArgs{
				clientIP:        "198.126.12.12",
				count:           1,
				err:             nil,
				urlID:           "",
				url:             "https://www.google.com",
				expireAt:        "2019-07-17T16:58:30+08:00",
				isExpireAtValid: false,
			},
		},
	}
	for _, tt := range tests {
		if tt.mockArgs.isExpireAtValid {
			tt.fields.storage.EXPECT().LoadClientIP(tt.mockArgs.clientIP).Return(tt.mockArgs.count, tt.mockArgs.err)
			tt.fields.storage.EXPECT().Save(tt.mockArgs.urlID, tt.mockArgs.url, tt.mockArgs.expireAt).Return(tt.mockArgs.err)
			tt.fields.storage.EXPECT().SaveClientIP(getClientIP(tt.args.r), gomock.Any()).Return(tt.mockArgs.err)
		}

		t.Run(tt.name, func(t *testing.T) {
			h := &handler{
				storage: tt.fields.storage,
			}
			h.UploadURL(tt.args.w, tt.args.r)
		})
	}
}

func (s *HandlerSuite) Test_handler_RedirectURL() {
	t := s.T()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		storage *storagemock.MockStorage
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	type mockArgs struct {
		err          error
		urlID        string
		url          string
		isUrlIDValid bool
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		mockArgs mockArgs
	}{
		{
			name: "Test RedirectURL is fail",
			fields: fields{
				storage: storagemock.NewMockStorage(ctrl),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: &http.Request{
					Header: http.Header{
						"X-Real-Ip": []string{"198.126.12.12"},
					},
				},
			},
			mockArgs: mockArgs{
				err:          errors.New("get fail"),
				urlID:        "",
				url:          "",
				isUrlIDValid: false,
			},
		},
	}
	for _, tt := range tests {
		if tt.mockArgs.isUrlIDValid {
			tt.fields.storage.EXPECT().Load(tt.args.r).Return(tt.mockArgs.url, tt.mockArgs.err)
		}

		t.Run(tt.name, func(t *testing.T) {
			h := &handler{
				storage: tt.fields.storage,
			}
			h.RedirectURL(tt.args.w, tt.args.r)
		})
	}
}
