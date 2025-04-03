package router

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandlerUser(t *testing.T) {

	type want struct {
		statusCode int
		body       string
	}

	tests := []struct {
		name  string
		token string
		want  want
	}{

		{
			name: "Hasn't token header",
			want: want{
				statusCode: 404,
				body:       "",
			},
			token: "",
		},

		{
			name: "Has wrong token",
			want: want{
				statusCode: 403,
				body:       "",
			},
			token: "",
		},

		{
			name: "Accepted token header",
			want: want{
				statusCode: 200,
				body:       "User JSON",
			},
			token: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/user", nil)
			r.Header.Set("x-token", tt.token)

			w := httptest.NewRecorder()
			h := http.HandlerFunc(HandlerUser)
			h(w, r)

			res := w.Result()

			defer res.Body.Close()
			bodyContent, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}

			assert.Equal(t, tt.want.statusCode, res.StatusCode)
			assert.Equal(t, tt.want.body, bodyContent)
		})
	}
}
