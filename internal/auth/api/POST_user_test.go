package api

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandlerCreateUser(t *testing.T) {

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
