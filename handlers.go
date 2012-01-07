package admin

import (
	"launchpad.net/gobson/bson"
	"net/http"
	"path"
	"strings"
)

//Parse request grabs the paramaters out of the request URL for the collection
//and object id the handler will operate on.
func parseRequest(p string) (coll, id string) {
	chunks := strings.Split(path.Clean(p), "/")
	if chunks[0] == "." {
		return
	}
	coll = chunks[0]
	if len(chunks) >= 2 {
		id = chunks[1]
	}
	return
}

//Presents the detail view for an object in a collection
func (a *Admin) Detail(w http.ResponseWriter, req *http.Request) {
	coll, id := parseRequest(req.URL.Path)
	if coll == "" || id == "" {
		a.Renderer.NotFound(w, req)
		return
	}
	c, t := a.collFor(coll), a.newType(coll)

	//unknown collection
	if t == nil {
		a.Renderer.NotFound(w, req)
		return
	}

	//load into T
	if err := c.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(t); err != nil {
		a.Renderer.InternalError(w, req, err)
		return
	}

	a.Renderer.Detail(w, req, DetailContext{
		Object: t,
	})
}

//Presents the index page giving an overall view of the database
func (a *Admin) Index(w http.ResponseWriter, req *http.Request) {
	coll, id := parseRequest(req.URL.Path)
	if coll != "" || id != "" {
		a.Renderer.NotFound(w, req)
		return
	}

	a.Renderer.Index(w, req)
}

//Presents a list of objects in a collection matching filtering/sorting criteria
func (a *Admin) List(w http.ResponseWriter, req *http.Request) {
	coll, id := parseRequest(req.URL.Path)
	if coll == "" || id != "" {
		a.Renderer.NotFound(w, req)
	}

	//grab the data

	a.Renderer.List(w, req, ListContext{})
}

//Presents a handler that updates an object and shows the results of the update
func (a *Admin) Update(w http.ResponseWriter, req *http.Request) {
	coll, id := parseRequest(req.URL.Path)
	if coll == "" || id == "" {
		a.Renderer.NotFound(w, req)
		return
	}

	c, t := a.collFor(coll), a.newType(coll)

	//unknown collection
	if t == nil {
		a.Renderer.NotFound(w, req)
		return
	}

	//grab the data
	if err := c.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(t); err != nil {
		a.Renderer.InternalError(w, req, err)
		return
	}

	//attempt to update

	a.Renderer.Update(w, req, UpdateContext{
		Object: t,
	})
}

//Presents a handler that creates an object and shows the results of the create
func (a *Admin) Create(w http.ResponseWriter, req *http.Request) {
	coll, id := parseRequest(req.URL.Path)
	if coll == "" || id != "" {
		a.Renderer.NotFound(w, req)
	}

	//attempt to insert

	a.Renderer.Create(w, req, CreateContext{})
}
