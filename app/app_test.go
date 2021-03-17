package app

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTodos(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(MakeHandler())
	defer ts.Close()
	v1 := url.Values{}
	v1.Set("name", "Test todo1")
	resp, err := http.PostForm(ts.URL+"/todos", v1)
	assert.NoError(err)
	defer resp.Body.Close()
	assert.Equal(http.StatusCreated, resp.StatusCode)
	var todo Todo
	err = json.NewDecoder(resp.Body).Decode(&todo)
	assert.NoError(err)
	assert.Equal(v1["name"][0], todo.Name)
	id1 := todo.ID

	v2 := url.Values{}
	v2.Set("name", "Test todo2")
	resp, err = http.PostForm(ts.URL+"/todos", v2)
	assert.NoError(err)
	defer resp.Body.Close()
	assert.Equal(http.StatusCreated, resp.StatusCode)
	err = json.NewDecoder(resp.Body).Decode(&todo)
	assert.NoError(err)
	assert.Equal(v2["name"][0], todo.Name)
	id2 := todo.ID

	resp, err = http.Get(ts.URL+"/todos")
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)
	todos := []*Todo{}
	err = json.NewDecoder(resp.Body).Decode(&todos)
	assert.NoError(err)
	assert.Equal(len(todos), 2)
}
