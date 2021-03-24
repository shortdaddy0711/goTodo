package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/shortdaddy0711/goTodo/model"
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
	var todo model.Todo
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

	resp, err = http.Get(ts.URL + "/todos")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	todos := []*model.Todo{}
	err = json.NewDecoder(resp.Body).Decode(&todos)
	assert.NoError(err)
	assert.Equal(len(todos), 2)
	for _, td := range todos {
		if td.ID == id1 {
			assert.Equal("Test todo1", td.Name)
		}
		if td.ID == id2 {
			assert.Equal("Test todo2", td.Name)
		} else {
			assert.Error(fmt.Errorf("testID should be id1 or id2"))
		}
	}

	resp, err = http.Get(ts.URL + "/complete/" + id1.String() + "?complete=true")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	resp, err = http.Get(ts.URL + "/todos")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	todos = []*model.Todo{}
	err = json.NewDecoder(resp.Body).Decode(&todos)
	assert.NoError(err)
	assert.Equal(len(todos), 2)
	for _, td := range todos {
		if td.ID == id1 {
			assert.True(td.Completed)
		}
	}

	req, _ := http.NewRequest("DELETE", ts.URL+"/todos/"+id1.String(), nil)
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	resp, err = http.Get(ts.URL + "/todos")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	todos = []*model.Todo{}
	err = json.NewDecoder(resp.Body).Decode(&todos)
	assert.NoError(err)
	assert.Equal(len(todos), 1)
	for _, td := range todos {
		assert.Equal(td.ID, id2)
	}

}
