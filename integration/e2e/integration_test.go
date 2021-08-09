package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gavv/httpexpect/v2"
	"github.com/google/uuid"
	"github.com/zacscoding/echo-gorm-realworld-app/api/types"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
)

var env *Env

func TestMain(m *testing.M) {
	setup()
	ret := m.Run()
	if ret == 0 {
		teardown()
	}
	os.Exit(ret)
}

func setup() {
	env = &Env{env: make(map[string]interface{})}
	//----------------------------------------------
	// Setup users
	//----------------------------------------------
	userCounts := 3
	for i := 1; i <= userCounts; i++ {
		var (
			username = "user-" + uuid.NewString()[:8]
			email    = username + "@gmail.com"
			password = "pass" + uuid.NewString()[:8]
		)
		env.SetToEnv(fmt.Sprintf("usertest.user%d.username", i), username)
		env.SetToEnv(fmt.Sprintf("usertest.user%d.email", i), email)
		env.SetToEnv(fmt.Sprintf("usertest.user%d.password", i), password)

		body := new(SignUpRequest)
		body.User.Username = username
		body.User.Email = email
		body.User.Password = password
		b, _ := json.Marshal(body)

		req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/api/users", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			panic(err)
		}
		resp, err := http.DefaultClient.Do(req)
		log.Println("resp", resp.StatusCode, "err", err)
		if err != nil {
			panic(err)
		}
		if resp.StatusCode != http.StatusOK {
			panic("failed to setup users")
		}
		defer resp.Body.Close()

		var user types.UserResponse
		respBody, _ := ioutil.ReadAll(resp.Body)
		if err := json.Unmarshal(respBody, &user); err != nil {
			panic(err)
		}
		env.SetToEnv(fmt.Sprintf("usertest.user%d.token", i), user.User.Token)
	}
}

func teardown() {}

func newTester(t *testing.T) *Tester {
	return NewTester(t, &TesterParams{
		BaseURL:      "http://localhost:8080",
		CurlPrinter:  false,
		DebugPrinter: true,
	})
}

func assertError(e *httpexpect.Response, code int, msg string) {
	e.Status(code)
	e.JSON().Path("$.errors.body").String().Contains(msg)
}
