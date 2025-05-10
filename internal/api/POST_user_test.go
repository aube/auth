package api

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aube/auth/internal/model"
	"github.com/aube/auth/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandlerCreateUser(t *testing.T) {
	// mockStore := new(mocks.Store)
	mockStore := mocks.NewStore(t)

	// Define the behavior of User method on Store instance.
	// userRepoMock := new(mocks.IUserRepository)
	userRepoMock := mocks.NewIUserRepository(t)

	// When Create is called with specific parameters, it should return nil error
	// userRepoMock.On("Create", context.Background(), &mockStore.User{Username: "testuser", Password: "testpass"}).Return(nil)

	// mockStore.On("User").Return(userRepoMock)
	// type User struct {
	// 	ID                int    `json:"id"`
	// 	Email             string `json:"email"`
	// 	Password          string `json:"password,omitempty"`
	// 	EncryptedPassword string `json:"-"`
	// }

	user := model.User{
		ID:                123,
		Email:             "testuser",
		Password:          "testpass",
		EncryptedPassword: "",
	}
	userRepoMock.EXPECT().Create("foo").Return("bar", nil).Once()
	userRepoMock.Create(context.Background(), &user).Return(nil)

	// userRepoMock.On("Create", context.Background(), user).Return(nil)

	fmt.Println(user)
	return

	server := Server{store: mockStore}

	handler := http.HandlerFunc(server.HandlerCreateUser)

	req, err := http.NewRequest("POST", "/user", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code) // or any other assertion you need
}

func TestHandlerCreateUser1(t *testing.T) {

	type resp struct {
		statusCode int
		body       string
	}

	tests := []struct {
		name string
		json string
		resp resp
	}{

		/* {
			name: "Hasn't token header",
			resp: resp{
				statusCode: 404,
				body:       `{"email":"ololo","password":"123321"}`,
			},
			json: `{"email":"ololo","password":"123321"}`,
		},
		{
			name: "Has wrong token",
			resp: resp{
				statusCode: 403,
				body:       `{"email":"ololo@mail.fake","password":""}`,
			},
			json: `{"email":"ololo","password":"123321"}`,
		}, */

		{
			name: "User created",
			resp: resp{
				statusCode: 201,
				body:       `Ololo, World!`,
			},
			json: `{"email":"ololo","password":"123321"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(
				http.MethodPost,
				"/user",
				bytes.NewBufferString(tt.json),
			)

			w := httptest.NewRecorder()

			// Set content type header
			r.Header.Set("Content-Type", "application/json")

			// TODO?
			// h := http.HandlerFunc(Server.HandlerCreatedUser)
			// h(w, r)

			res := w.Result()

			defer res.Body.Close()
			bodyContent, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}

			assert.Equal(t, tt.resp.statusCode, res.StatusCode)
			assert.Equal(t, tt.resp.body, bodyContent)
		})
	}
}
