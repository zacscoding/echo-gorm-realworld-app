package e2e

import (
	"github.com/gavv/httpexpect/v2"
	"github.com/google/uuid"
	"net/http"
	"testing"
)

// Tests "POST /api/users/login"
func TestUser_SignIn(t *testing.T) {
	t.Run("SignIn Success", func(t *testing.T) {
		tester := newTester(t)
		req := new(SignInRequest)
		req.User.Email = env.GetFromEnvString("usertest.user1.email")
		req.User.Password = env.GetFromEnvString("usertest.user1.password")

		e := tester.POST("/api/users/login").WithJSON(req).Expect().Status(http.StatusOK)

		resp := e.JSON()
		resp.Schema(UserJsonSchema)
		resp.Path("$.user.email").Equal(env.GetFromEnvString("usertest.user1.email"))
		resp.Path("$.user.username").Equal(env.GetFromEnvString("usertest.user1.username"))
		resp.Path("$.user.bio").Equal("")
		resp.Path("$.user.image").Equal("")
	})
	t.Run("SignIn Fail", func(t *testing.T) {
		cases := []struct {
			Name     string
			Email    string
			Password string
			// Expected
			Code int
			Msg  string
		}{
			{
				Name: "required email",
				Code: http.StatusUnprocessableEntity,
				Msg:  "Email validation error. reason: required",
			}, {
				Name:  "invalid email format",
				Email: "not email pattern",
				Code:  http.StatusUnprocessableEntity,
				Msg:   "Email validation error. reason: email",
			}, {
				Name:  "required password",
				Email: env.GetFromEnvString("usertest.user1.email"),
				Code:  http.StatusUnprocessableEntity,
				Msg:   "Password validation error. reason: required",
			}, {
				Name:     "mismatch password",
				Email:    env.GetFromEnvString("usertest.user1.email"),
				Password: "mismatch password",
				Code:     http.StatusUnprocessableEntity,
				Msg:      "password mismatch",
			}, {
				Name:     "user not found",
				Email:    uuid.NewString()[:5] + "@gmail.com",
				Password: "password",
				Code:     http.StatusNotFound,
				Msg:      "not found",
			},
		}

		for _, tc := range cases {
			t.Run(tc.Name, func(t *testing.T) {
				t.Parallel()
				tester := newTester(t)
				req := new(SignInRequest)
				req.User.Email = tc.Email
				req.User.Password = tc.Password

				e := tester.POST("/api/users/login").WithJSON(req).Expect()

				assertError(e, tc.Code, tc.Msg)
			})
		}
	})
}

// Tests POST /api/users
func TestUser_SignUp(t *testing.T) {
	t.Run("SignUp Success", func(t *testing.T) {
		tester := newTester(t)
		req := new(SignUpRequest)
		req.User.Email = uuid.NewString()[:8] + "@gmail.com"
		req.User.Password = uuid.NewString()[:4]
		req.User.Username = uuid.NewString()[:4]

		e := tester.POST("/api/users").WithJSON(req).Expect().Status(http.StatusOK)

		resp := e.JSON()
		resp.Schema(UserJsonSchema)
		u := resp.Path("$.user").Object()
		u.Value("email").Equal(req.User.Email)
		u.Value("username").Equal(req.User.Username)
		u.Value("token").NotNull()
		u.Value("bio").Equal("")
		u.Value("image").Equal("")
	})
	t.Run("SignUp Fail", func(t *testing.T) {
		cases := []struct {
			Name     string
			Username string
			Email    string
			Password string
			// Expected
			Code int
			Msg  string
		}{
			{
				Name:     "required email",
				Username: uuid.NewString()[:4],
				Password: uuid.NewString()[:4],
				Code:     http.StatusUnprocessableEntity,
				Msg:      "Email validation error. reason: required",
			}, {
				Name:     "invalid email format",
				Email:    "not email pattern",
				Username: uuid.NewString()[:4],
				Password: uuid.NewString()[:4],
				Code:     http.StatusUnprocessableEntity,
				Msg:      "Email validation error. reason: email",
			}, {
				Name:     "required password",
				Username: uuid.NewString()[:4],
				Email:    uuid.NewString()[:8] + "@gmail.com",
				Code:     http.StatusUnprocessableEntity,
				Msg:      "Password validation error. reason: required",
			}, {
				Name:     "duplicate email",
				Username: uuid.NewString()[:4],
				Email:    env.GetFromEnvString("usertest.user1.email"),
				Password: uuid.NewString()[:4],
				Code:     http.StatusUnprocessableEntity,
				Msg:      "duplicate email",
			},
		}

		for _, tc := range cases {
			t.Run(tc.Name, func(t *testing.T) {
				t.Parallel()
				tester := newTester(t)
				req := new(SignUpRequest)
				req.User.Email = tc.Email
				req.User.Password = tc.Password
				req.User.Username = tc.Username

				e := tester.POST("/api/users").WithJSON(req).Expect()

				assertError(e, tc.Code, tc.Msg)
			})
		}
	})
}

// Tests GET /api/user
func TestUser_CurrentUser(t *testing.T) {
	t.Run("CurrentUser Success", func(t *testing.T) {
		tester := newTester(t)

		e := tester.GetWithAuthToken("/api/user", env.GetFromEnvString("usertest.user1.token")).Expect().Status(http.StatusOK)

		resp := e.JSON()
		resp.Schema(UserJsonSchema)
		u := resp.Path("$.user").Object()
		u.Value("email").Equal(env.GetFromEnvString("usertest.user1.email"))
		u.Value("username").Equal(env.GetFromEnvString("usertest.user1.username"))
		u.Value("token").NotNull()
		u.Value("bio").Equal("")
		u.Value("image").Equal("")
	})
	t.Run("CurrentUser Fail", func(t *testing.T) {
		cases := []struct {
			Name      string
			AuthToken string
			// expected
			Code int
			Msg  string
		}{
			{
				Name:      "required token",
				AuthToken: "",
				Code:      http.StatusUnauthorized,
				Msg:       "auth required",
			}, {
				Name:      "invalid token",
				AuthToken: "invalid token",
				Code:      http.StatusUnauthorized,
				Msg:       "auth required",
			},
		}

		for _, tc := range cases {
			t.Run(tc.Name, func(t *testing.T) {
				t.Parallel()
				tester := newTester(t)

				e := tester.GetWithAuthToken("/api/user", tc.AuthToken).Expect()

				assertError(e, tc.Code, tc.Msg)
			})
		}
	})
}

// Tests PUT /api/user
func TestUser_Update(t *testing.T) {
	t.Run("Update Success", func(t *testing.T) {
		tester := newTester(t)
		req := new(UpdateUserRequest)
		req.User.Email = "user-" + uuid.NewString()[:8] + "@gmail.com"
		req.User.Password = uuid.NewString()[:4]
		req.User.Bio = uuid.NewString()
		req.User.Image = uuid.NewString()

		e := tester.PutWithAuthToken("/api/user", env.GetFromEnvString("usertest.user3.token")).WithJSON(req).Expect().Status(http.StatusOK)

		resp := e.JSON()
		resp.Schema(UserJsonSchema)
		u := resp.Path("$.user").Object()
		u.Value("email").Equal(req.User.Email)
		u.Value("username").Equal(env.GetFromEnvString("usertest.user3.username"))
		u.Value("token").NotNull()
		u.Value("bio").Equal(req.User.Bio)
		u.Value("image").Equal(req.User.Image)

		env.SetToEnv("usertest.user3.email", req.User.Email)
		env.SetToEnv("usertest.user3.password", req.User.Password)
	})
	t.Run("Update Fail", func(t *testing.T) {
		cases := []struct {
			Name      string
			AuthToken string
			Email     string
			Username  string
			Password  string
			Bio       string
			Image     string
			// expected
			Code int
			Msg  string
		}{
			{
				Name:      "invalid email format",
				AuthToken: env.GetFromEnvString("usertest.user3.token"),
				Email:     "invalid email pattern",
				Code:      http.StatusUnprocessableEntity,
				Msg:       "Email validation error. reason: email",
			},
		}

		for _, tc := range cases {
			t.Run(tc.Name, func(t *testing.T) {
				t.Parallel()
				tester := newTester(t)
				req := new(UpdateUserRequest)
				req.User.Email = tc.Email
				req.User.Username = tc.Username
				req.User.Password = tc.Password
				req.User.Bio = tc.Bio
				req.User.Image = tc.Image

				e := tester.PutWithAuthToken("/api/user", env.GetFromEnvString("usertest.user3.token")).WithJSON(req).Expect()

				assertError(e, tc.Code, tc.Msg)
			})
		}
	})
}

// Tests GET /api/profiles/:username
func TestUser_Profile(t *testing.T) {
	t.Run("Profile Success", func(t *testing.T) {
		tester := newTester(t)
		token := env.GetFromEnvString("usertest.user1.token")
		username := env.GetFromEnvString("usertest.user3.username")

		e := tester.GetWithAuthToken("/api/profiles/{username}", token, username).Expect().Status(http.StatusOK)

		resp := e.JSON()
		resp.Schema(ProfileJsonSchema)
		p := resp.Path("$.profile").Object()
		p.Value("username").Equal(env.GetFromEnvString("usertest.user3.username"))
		p.ContainsKey("image")
		p.ContainsKey("bio")
		p.ContainsKey("following")
	})
	t.Run("Profile Fail", func(t *testing.T) {
		cases := []struct {
			Name      string
			AuthToken string
			Username  string
			// expected
			Code int
			Msg  string
		}{
			{
				Name:      "Not Found",
				AuthToken: "",
				Username:  "notfounduser",
				Code:      http.StatusNotFound,
				Msg:       "not found",
			},
		}

		for _, tc := range cases {
			t.Run(tc.Name, func(t *testing.T) {
				t.Parallel()
				tester := newTester(t)

				var e *httpexpect.Response
				if tc.AuthToken != "" {
					e = tester.GetWithAuthToken("/api/profiles/{username}", tc.AuthToken, tc.Username).Expect()
				} else {
					e = tester.GET("/api/profiles/{username}", tc.Username).Expect()
				}
				assertError(e, tc.Code, tc.Msg)
			})
		}
	})
}

func TestUser_FollowUnFollow(t *testing.T) {
	//----------------------------------------------
	// Tests POST /api/profiles/:username/follow
	//----------------------------------------------
	t.Run("Follow Success", func(t *testing.T) {
		tester := newTester(t)
		user1Token := env.GetFromEnvString("usertest.user1.token")
		user2Name := env.GetFromEnvString("usertest.user2.username")

		e := tester.PostWithAuthToken("/api/profiles/{username}/follow", user1Token, user2Name).Expect().Status(http.StatusOK)

		resp := e.JSON()
		resp.Schema(ProfileJsonSchema)
		p := resp.Path("$.profile").Object()
		p.Value("username").Equal(env.GetFromEnvString("usertest.user2.username"))
		p.Value("following").Equal(true)
		p.ContainsKey("image")
		p.ContainsKey("bio")
	})
	t.Run("Follow Fail", func(t *testing.T) {
		cases := []struct {
			Name      string
			AuthToken string
			Username  string
			// expected
			Code int
			Msg  string
		}{
			{
				Name:      "Not Found",
				AuthToken: env.GetFromEnvString("usertest.user1.token"),
				Username:  "notfounduser",
				Code:      http.StatusNotFound,
				Msg:       "not found",
			},
		}

		for _, tc := range cases {
			t.Run(tc.Name, func(t *testing.T) {
				t.Parallel()
				tester := newTester(t)

				e := tester.PostWithAuthToken("/api/profiles/{username}/follow", tc.AuthToken, tc.Username).Expect()

				assertError(e, tc.Code, tc.Msg)
			})
		}
	})

	//----------------------------------------------
	// Tests DELETE /api/profiles/:username/follow
	//----------------------------------------------
	t.Run("UnFollow Success", func(t *testing.T) {
		tester := newTester(t)
		user1Token := env.GetFromEnvString("usertest.user1.token")
		user2Name := env.GetFromEnvString("usertest.user2.username")

		e := tester.DeleteWithAuthToken("/api/profiles/{username}/follow", user1Token, user2Name).Expect().Status(http.StatusOK)

		resp := e.JSON()
		resp.Schema(ProfileJsonSchema)
		p := resp.Path("$.profile").Object()
		p.Value("username").Equal(env.GetFromEnvString("usertest.user2.username"))
		p.Value("following").Equal(false)
		p.ContainsKey("image")
		p.ContainsKey("bio")
	})
	t.Run("UnFollow Fail", func(t *testing.T) {
		cases := []struct {
			Name      string
			AuthToken string
			Username  string
			// expected
			Code int
			Msg  string
		}{
			{
				Name:      "Not Found",
				AuthToken: env.GetFromEnvString("usertest.user1.token"),
				Username:  "notfounduser",
				Code:      http.StatusNotFound,
				Msg:       "not found",
			},
		}

		for _, tc := range cases {
			t.Run(tc.Name, func(t *testing.T) {
				t.Parallel()
				tester := newTester(t)

				e := tester.DeleteWithAuthToken("/api/profiles/{username}/follow", tc.AuthToken, tc.Username).Expect()

				assertError(e, tc.Code, tc.Msg)
			})
		}
	})
}

type (
	SignUpRequest struct {
		User struct {
			Username string `json:"username,omitempty"`
			Email    string `json:"email,omitempty"`
			Password string `json:"password,omitempty"`
		} `json:"user,omitempty"`
	}

	SignInRequest struct {
		User struct {
			Email    string `json:"email,omitempty"`
			Password string `json:"password,omitempty"`
		} `json:"user,omitempty"`
	}

	UpdateUserRequest struct {
		User struct {
			Email    string `json:"email,omitempty"`
			Username string `json:"username,omitempty"`
			Password string `json:"password,omitempty"`
			Bio      string `json:"bio,omitempty"`
			Image    string `json:"image,omitempty"`
		} `json:"user,omitempty"`
	}
)
