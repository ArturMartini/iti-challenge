package rest

import (
	"encoding/json"
	"errors"
	"github.com/arturmartini/iti-challenge/entities"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type integrationPasswordTest struct {
	suite.Suite
	server *httptest.Server
	client http.Client
}

type testCase struct {
	Uri           string
	Method        string
	Value         string
	ExpectedValue bool
}

type mockResponseWriter struct {
	code int
}

type errReader int

func (*mockResponseWriter) Header() http.Header {
	return http.Header{}
}

func (*mockResponseWriter) Write([]byte) (int, error) {
	return 200, nil
}

func (r *mockResponseWriter) WriteHeader(code int) {
	r.code = code
}

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}

func TestRunPassword(t *testing.T) {
	suite.Run(t, new(integrationPasswordTest))
}

func (r *integrationPasswordTest) SetupTest() {
	r.server = httptest.NewServer(router)
	r.client = http.Client{}
}

func (suite *integrationPasswordTest) TestIntegrationPasswordValidate() {
	testCases := []testCase{
		{
			Uri:           uri,
			Method:        http.MethodPost,
			Value:         "",
			ExpectedValue: false,
		},
		{
			Uri:           uri,
			Method:        http.MethodPost,
			Value:         "aa",
			ExpectedValue: false,
		},
		{
			Uri:           uri,
			Method:        http.MethodPost,
			Value:         "ab",
			ExpectedValue: false,
		},
		{
			Uri:           uri,
			Method:        http.MethodPost,
			Value:         "AAAbbbCc",
			ExpectedValue: false,
		},
		{
			Uri:           uri,
			Method:        http.MethodPost,
			Value:         "AbTp9!foo",
			ExpectedValue: false,
		},
		{
			Uri:           uri,
			Method:        http.MethodPost,
			Value:         "AbTp9!foA",
			ExpectedValue: false,
		},
		{
			Uri:           uri,
			Method:        http.MethodPost,
			Value:         "AbTp9 fok",
			ExpectedValue: false,
		},
		{
			Uri:           uri,
			Method:        http.MethodPost,
			Value:         "AbTp9!fok",
			ExpectedValue: true,
		},
	}

	validateTest(suite, testCases)
}

func (suite *integrationPasswordTest) TestFromJsonMethodUnmarshalError() {
	w := new(mockResponseWriter)
	r := &http.Request{
		Body: ioutil.NopCloser(strings.NewReader("")),
	}
	err := fromJson(w, r, nil)
	assert.Equal(suite.T(), "unexpected end of JSON input", err.Error())
	assert.Equal(suite.T(), http.StatusBadRequest, w.code)
}

func (suite *integrationPasswordTest) TestFromJsonMethodRequestError() {
	w := new(mockResponseWriter)
	r := httptest.NewRequest(http.MethodPost, "/something", errReader(0))
	err := fromJson(w, r, nil)
	assert.Equal(suite.T(), "test error", err.Error())
	assert.Equal(suite.T(), http.StatusBadRequest, w.code)
}

func (suite *integrationPasswordTest) TestHandleResponse() {
	w := new(mockResponseWriter)
	handleResponse(entities.Response{Valid: true}, w)
	assert.Equal(suite.T(), http.StatusOK, w.code)

	handleResponse(make(chan int), w)
	assert.Equal(suite.T(), http.StatusInternalServerError, w.code)
}

func request(client http.Client, url, method string, value interface{}) (*http.Response, error) {
	bytes, err := json.Marshal(value)
	if err != nil {
		log.WithError(err).Panic("Error when marshal value")
	}

	req, err := http.NewRequest(method, url, strings.NewReader(string(bytes)))
	if err != nil {
		log.WithError(err).Panic("Error when execute request http test")
	}
	return client.Do(req)
}

func response(response *http.Response, err error) (entities.Response, error) {
	if err != nil {
		return entities.Response{}, err
	}

	defer response.Body.Close()

	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return entities.Response{}, err
	}

	resp := entities.Response{}
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		return entities.Response{}, err
	}
	return resp, nil
}

func validateTest(suite *integrationPasswordTest, testCases []testCase) {
	for _, c := range testCases {
		resp, err := response(request(suite.client, suite.server.URL+c.Uri, c.Method, entities.Password{Value: c.Value}))
		assert.Nil(suite.T(), err)
		assert.Equal(suite.T(), c.ExpectedValue, resp.Valid)
	}
}
