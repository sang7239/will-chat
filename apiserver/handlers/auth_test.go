package handlers

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/will-slack/apiserver/models/users"
	"github.com/will-slack/apiserver/sessions"
)

type request struct {
	method  string
	handler http.HandlerFunc
	path    string
	body    interface{}
	expBody string
	expCode int
}

func NewContext() *Context {
	return &Context{
		SessionKey:   "secretkey",
		SessionStore: sessions.NewMemStore(time.Hour),
		UserStore:    users.NewMemStore(),
	}
}

func testCase(t *testing.T, c *request, wg *sync.WaitGroup) {
	defer wg.Done()
	var body io.Reader
	if b, ok := c.body.(string); ok {
		bodyStr := []byte(b)
		body = bytes.NewBuffer(bodyStr)
	} else if c.body == nil {
		body = nil
	}
	resRec := httptest.NewRecorder()
	handler := http.HandlerFunc(c.handler)
	req, err := http.NewRequest(c.method, c.path, body)
	if err != nil {
		t.Fatal(err)
	}
	handler.ServeHTTP(resRec, req)
	if resRec.Code != c.expCode {
		t.Errorf("handler returned wrong status code: expected `%d` but got `%d`\n", c.expCode, resRec.Code)
	}
	expected := []byte(c.expBody)
	if !bytes.Equal(resRec.Body.Bytes(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			resRec.Body.Bytes(), expected)
	}
}

func TestSessionHandler(t *testing.T) {
	ctx := NewContext()
	cases := []request{
		// valid user sign up
		{
			method:  "POST",
			handler: ctx.UserHandler,
			path:    "/v1/users",
			body: `{
						"email": "testsesh@gmail.com",
						"password": "password",
						"passwordConf": "password",
						"userName": "testsesh",
						"firstName": "testsesh",
						"lastName": "sesh"
					}`,
			expBody: `{
					"id":"5efd79be2456a833785a1b19",
					"email":"testsesh@gmail.com",
					"userName":"testsesh",
					"firstName":"testsesh",
					"lastName":"sesh",
					"photoURL":"https://www.gravatar.com/avatar/93b02540c167ddf9a1761d9f3701754e"
				}`,
			expCode: http.StatusOK,
		},
		// valid sign in
		{
			method:  "POST",
			handler: ctx.SessionHandler,
			path:    "/v1/sessions",
			body: `{
				"email":"testsesh@gmail.com",
				"password":"password"
			}`,
			expBody: ` {
				"id":"5efd79be2456a833785a1b19",
				"email":"testsesh@gmail.com",
				"userName":"testsesh",
				"firstName":"testsesh",
				"lastName":"sesh",
				"photoURL":"https://www.gravatar.com/avatar/93b02540c167ddf9a1761d9f3701754e"
			}
			`,
			expCode: http.StatusOK,
		},
		// invalid sign in
		{
			method:  "POST",
			handler: ctx.SessionHandler,
			path:    "/v1/sessions",
			body: `{
				"email":"testsesh@gmail.com",
				"password":"wrongpassword"
			}`,
			expBody: "User Not Authorized",
			expCode: http.StatusUnauthorized,
		},
	}
	wg := sync.WaitGroup{}
	for _, c := range cases {
		wg.Add(1)
		go testCase(t, &c, &wg)
	}
	wg.Wait()
}

func TestSessionMineHandler(t *testing.T) {
	ctx := NewContext()
	handler := http.HandlerFunc(ctx.UserHandler)
	resRec := httptest.NewRecorder()
	body := `{
		"email": "signout@gmail.com",
		"password": "password",
		"passwordConf": "password",
		"userName": "signout",
		"firstName": "signout",
		"lastName": "out"
	}`
	bodyStr := []byte(body)
	req, err := http.NewRequest("POST", "/v1/users", bytes.NewBuffer(bodyStr))
	if err != nil {
		t.Fatal(err)
	}
	handler.ServeHTTP(resRec, req)

	if status := resRec.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: expected `%d` but got `%d`\n", http.StatusOK, resRec.Code)
	}
	auth := resRec.Header().Get("Authorization")
	handler = http.HandlerFunc(ctx.SessionMineHandler)
	resRec = httptest.NewRecorder()
	req, err = http.NewRequest("DELETE", "/v1/sessions/mine", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Authorization", auth)
	handler.ServeHTTP(resRec, req)
	if status := resRec.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: expected `%d` but got `%d`\n", http.StatusOK, resRec.Code)
	}
	if resRec.Body.String() != "signed out\n" {
		t.Errorf("handler returned unexpected body: got %v want %v",
			resRec.Body.String(), "signed out")
	}
}

func TestUserMeHandler(t *testing.T) {
	ctx := NewContext()
	handler := http.HandlerFunc(ctx.UserHandler)
	resRec := httptest.NewRecorder()
	body := `{
		"email": "userme@gmail.com",
		"password": "password",
		"passwordConf": "password",
		"userName": "userme",
		"firstName": "user",
		"lastName": "me"
	}`
	bodyStr := []byte(body)
	req, err := http.NewRequest("POST", "/v1/users", bytes.NewBuffer(bodyStr))
	if err != nil {
		t.Fatal(err)
	}
	handler.ServeHTTP(resRec, req)

	if status := resRec.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: expected `%d` but got `%d`\n", http.StatusOK, resRec.Code)
	}
	auth := resRec.Header().Get("Authorization")
	handler = http.HandlerFunc(ctx.UserMeHandler)
	resRec = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/v1/users/me", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Authorization", auth)
	handler.ServeHTTP(resRec, req)
	if status := resRec.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: expected `%d` but got `%d`\n", http.StatusOK, resRec.Code)
	}
}

func TestUserHandler(t *testing.T) {
	ctx := NewContext()
	cases := []request{
		// valid user sign up
		{
			method:  "POST",
			handler: ctx.UserHandler,
			path:    "/v1/users",
			body: `{
						"email": "test@gmail.com",
						"password": "password",
						"passwordConf": "password",
						"userName": "test",
						"firstName": "user1",
						"lastName": "one"
					}`,
			expBody: `{
						"id":"5efd16192456a81d1de4dd8b",
						"email":"test@gmail.com",
						"userName":"test",
						"firstName":"user",
						"lastName":"one",
						"photoURL":"https://www.gravatar.com/avatar/1aedb8d9dc4751e229a335e371db8058"
				}`,
			expCode: http.StatusOK,
		},
		{
			method:  "POST",
			handler: ctx.UserHandler,
			path:    "/v1/users",
			body: `{
						"email": "test2@gmail.com",
						"password": "password",
						"passwordConf": "password",
						"userName": "test2",
						"firstName": "user2",
						"lastName": "two"
				}`,
			expBody: `
				{
					"id":"5efd17302456a81d1de4dd8c",
					"email":"test2@gmail.com",
					"userName":"test2",
					"firstName":"user2",
					"lastName":"two",
					"photoURL":"https://www.gravatar.com/avatar/3c4f419e8cd958690d0d14b3b89380d3"
				}`,
			expCode: http.StatusOK,
		},

		// invalid sign-up (Invalid Email)
		{
			method:  "POST",
			handler: ctx.UserHandler,
			path:    "/v1/users/",
			body: `{
				"email": "notvalidemail",
				"password": "password",
				"passwordConf": "password",
				"userName": "test3",
				"firstName": "user3",
				"lastName": "three"
			}`,
			expBody: "Unable to validate user",
			expCode: http.StatusBadRequest,
		},

		// invalid sign-up (Invalid username)
		{
			method:  "POST",
			handler: ctx.UserHandler,
			path:    "/v1/users/",
			body: `{
				"email": "test4@gmail.com",
				"password": "password",
				"passwordConf": "password",
				"userName": "",
				"firstName": "user4",
				"lastName": "four"
			}`,
			expBody: "Unable to validate user",
			expCode: http.StatusBadRequest,
		},
		// invalid sign-up (Duplicated User Email)
		{
			method:  "POST",
			handler: ctx.UserHandler,
			path:    "/v1/users/",
			body: `{
				"email": "test@gmail.com",
				"password": "password",
				"passwordConf": "password",
				"userName": "test",
				"firstName": "user1",
				"lastName": "one"
			}`,
			expBody: "Invalid Email. Already exists",
			expCode: http.StatusBadRequest,
		},
		// invalid sign-up (Duplicated username)
		{
			method:  "POST",
			handler: ctx.UserHandler,
			path:    "/v1/users/",
			body: `{
				"email": "test5@gmail.com",
				"password": "password",
				"passwordConf": "password",
				"userName": "test",
				"firstName": "user5",
				"lastName": "five"
			}`,
			expBody: "Invalid username. Already exists",
			expCode: http.StatusBadRequest,
		},
		// valid user retrieval
		{
			method:  "GET",
			handler: ctx.UserHandler,
			path:    "/v1/users/",
			body:    nil,
			expBody: `
			{
				"id":"5efd16192456a81d1de4dd8b",
				"email":"test@gmail.com",
				"userName":"test",
				"firstName":"user",
				"lastName":"one",
				"photoURL":"https://www.gravatar.com/avatar/1aedb8d9dc4751e229a335e371db8058"
			}
			{
				"id":"5efd17302456a81d1de4dd8c",
				"email":"test2@gmail.com",
				"userName":"test2",
				"firstName":"user2",
				"lastName":"two",
				"photoURL":"https://www.gravatar.com/avatar/3c4f419e8cd958690d0d14b3b89380d3"
			}
			`,
			expCode: http.StatusOK,
		},
	}
	wg := sync.WaitGroup{}
	for _, c := range cases {
		wg.Add(1)
		go testCase(t, &c, &wg)
	}
	wg.Wait()
}
