package account_test

import (
	"encoding/json"
	"go-react-app/internal/dal"
	"go-react-app/internal/handlers/account"
	"go-react-app/internal/util/logger"
	"io"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

func mockLogger() *logger.Logger {
	return logger.NewWithWriter(&io.Discard)
}

func getFakeUsers() []*dal.User {
	return append(
		[]*dal.User{},
		&dal.User{
			Username:     "Alice",
			PasswordHash: "123",
			CreatedDate:  time.Now(),
		},
		&dal.User{
			Username:     "Bob",
			PasswordHash: "asd",
			CreatedDate:  time.Now(),
		},
		&dal.User{
			Username:     "Charlie",
			PasswordHash: "!@#",
			CreatedDate:  time.Now(),
		},
	)
}

type mockUserRepository struct {
	users []*dal.User
}

func (r *mockUserRepository) GetUserByUsername(username string) (*dal.User, error) {
	for _, item := range r.users {
		if item.Username == username {
			return item, nil
		}
	}
	return nil, errors.New("Not found User")
}

type mockSessionCreator struct{}

func (r *mockSessionCreator) CreateSession(userId int64) (uuid.UUID, error) {
	return uuid.Nil, nil
}

func Test_login_ServeHTTP(t *testing.T) {
	type fields struct {
		log      *logger.Logger
		users    account.UserByUsernameGetter
		sessions account.SessionCreator
	}
	type args struct {
		json string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Test valid input",
			fields: fields{
				log:      mockLogger(),
				users:    &mockUserRepository{users: getFakeUsers()},
				sessions: &mockSessionCreator{},
			},
			args: args{
				json: "{ \"username\": \"Alice\", \"password\": \"123\" }",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rw := httptest.NewRecorder()
			r := httptest.NewRequest(
				"POST", "http://example.com/api/account",
				strings.NewReader(tt.args.json),
			)

			h := account.NewLogin(
				tt.fields.log,
				tt.fields.users,
				tt.fields.sessions,
			)

			h.ServeHTTP(rw, r)
			defer rw.Result().Body.Close()

			payload, err := io.ReadAll(rw.Result().Body)
			if err != nil {
				t.Error(err)
				t.FailNow()
			}
			t.Log(string(payload))
			response := &account.LoginResponse{}
			err = json.Unmarshal(payload, response)
			if err != nil {
				t.Error(err)
				t.FailNow()
			}

			if response.SessionToken != uuid.Nil {
				t.FailNow()
			}
		})
	}
}
