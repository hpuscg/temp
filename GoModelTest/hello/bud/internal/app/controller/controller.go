package controller

import (
	request "github.com/livebud/bud/framework/controller/controllerrt/request"
	response "github.com/livebud/bud/framework/controller/controllerrt/response"
	http "net/http"
	about "temp/GoModelTest/hello/controller/about"
	users "temp/GoModelTest/hello/controller/users"
)

// Controller struct
type Controller struct {
	About *AboutController
	Post  *PostController
	Users *UsersController
}

// Controller struct
type AboutController struct {
	Index *AboutIndexAction
}

// AboutIndexAction struct
type AboutIndexAction struct {
}

// Key is a unique identifier of this action
func (i *AboutIndexAction) Key() string {
	return "/about/index"
}

// Path is the default RESTful path to this action
func (i *AboutIndexAction) Path() string {
	return "/about"
}

// Method is the default RESTful method of this action
func (i *AboutIndexAction) Method() string {
	return "GET"
}

// ServeHTTP fn
func (i *AboutIndexAction) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	i.handler(r).ServeHTTP(w, r)
}

// Handler function
func (i *AboutIndexAction) handler(httpRequest *http.Request) http.Handler {
	aboutController, _nil := loadAboutController()
	if _nil != nil {
		return &response.Format{
			JSON: response.Status(500).Set("Content-Type", "application/json").JSON(map[string]string{"error": _nil.Error()}),
		}
	}
	fn := aboutController.Index
	// Call the controller
	in0, in1 := fn()
	if in1 != nil {
		return &response.Format{
			JSON: response.Status(500).Set("Content-Type", "application/json").JSON(map[string]string{"error": in1.Error()}),
		}
	}

	// Respond
	return &response.Format{
		JSON: response.JSON(in0),
	}
}

// Controller struct
type PostController struct {
}

// Controller struct
type UsersController struct {
	Index  *UsersIndexAction
	New    *UsersNewAction
	Create *UsersCreateAction
	Show   *UsersShowAction
	Update *UsersUpdateAction
	Delete *UsersDeleteAction
	Edit   *UsersEditAction
}

// UsersIndexAction struct
type UsersIndexAction struct {
}

// Key is a unique identifier of this action
func (i *UsersIndexAction) Key() string {
	return "/users/index"
}

// Path is the default RESTful path to this action
func (i *UsersIndexAction) Path() string {
	return "/users"
}

// Method is the default RESTful method of this action
func (i *UsersIndexAction) Method() string {
	return "GET"
}

// ServeHTTP fn
func (i *UsersIndexAction) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	i.handler(r).ServeHTTP(w, r)
}

// Handler function
func (i *UsersIndexAction) handler(httpRequest *http.Request) http.Handler {
	usersController, _nil := loadUsersController()
	if _nil != nil {
		return &response.Format{
			JSON: response.Status(500).Set("Content-Type", "application/json").JSON(map[string]string{"error": _nil.Error()}),
		}
	}
	fn := usersController.Index
	// Call the controller
	in0, in1 := fn()
	if in1 != nil {
		return &response.Format{
			JSON: response.Status(500).Set("Content-Type", "application/json").JSON(map[string]string{"error": in1.Error()}),
		}
	}

	// Respond
	return &response.Format{
		JSON: response.JSON(in0),
	}
}

// UsersNewAction struct
type UsersNewAction struct {
}

// Key is a unique identifier of this action
func (n *UsersNewAction) Key() string {
	return "/users/new"
}

// Path is the default RESTful path to this action
func (n *UsersNewAction) Path() string {
	return "/users/new"
}

// Method is the default RESTful method of this action
func (n *UsersNewAction) Method() string {
	return "GET"
}

// ServeHTTP fn
func (n *UsersNewAction) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	n.handler(r).ServeHTTP(w, r)
}

// Handler function
func (n *UsersNewAction) handler(httpRequest *http.Request) http.Handler {
	usersController, _nil := loadUsersController()
	if _nil != nil {
		return &response.Format{
			JSON: response.Status(500).Set("Content-Type", "application/json").JSON(map[string]string{"error": _nil.Error()}),
		}
	}
	fn := usersController.New
	// Call the controller
	fn()

	// Respond
	return &response.Format{
		JSON: response.Status(204),
	}
}

// UsersCreateAction struct
type UsersCreateAction struct {
}

// Key is a unique identifier of this action
func (c *UsersCreateAction) Key() string {
	return "/users/create"
}

// Path is the default RESTful path to this action
func (c *UsersCreateAction) Path() string {
	return "/users"
}

// Method is the default RESTful method of this action
func (c *UsersCreateAction) Method() string {
	return "POST"
}

// ServeHTTP fn
func (c *UsersCreateAction) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.handler(r).ServeHTTP(w, r)
}

// Handler function
func (c *UsersCreateAction) handler(httpRequest *http.Request) http.Handler {
	// Define the input struct
	var in struct {
		Name string `json:"name,omitempty"`
		Age  int    `json:"age,omitempty"`
	}
	// Unmarshal the request body
	if err := request.Unmarshal(httpRequest, &in); err != nil {
		return &response.Format{
			JSON: response.Status(400).Set("Content-Type", "application/json").JSON(map[string]string{"error": err.Error()}),
		}
	}
	usersController, _nil := loadUsersController()
	if _nil != nil {
		return &response.Format{
			JSON: response.Status(500).Set("Content-Type", "application/json").JSON(map[string]string{"error": _nil.Error()}),
		}
	}
	fn := usersController.Create
	// Call the controller
	in0, in1 := fn(
		in.Name,
		in.Age,
	)
	if in1 != nil {
		return &response.Format{
			JSON: response.Status(500).Set("Content-Type", "application/json").JSON(map[string]string{"error": in1.Error()}),
		}
	}

	// Respond
	return &response.Format{
		HTML: response.Status(302).Redirect(response.RedirectPath(httpRequest, "")),
		JSON: response.JSON(in0),
	}
}

// UsersShowAction struct
type UsersShowAction struct {
}

// Key is a unique identifier of this action
func (s *UsersShowAction) Key() string {
	return "/users/show"
}

// Path is the default RESTful path to this action
func (s *UsersShowAction) Path() string {
	return "/users/:id"
}

// Method is the default RESTful method of this action
func (s *UsersShowAction) Method() string {
	return "GET"
}

// ServeHTTP fn
func (s *UsersShowAction) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler(r).ServeHTTP(w, r)
}

// Handler function
func (s *UsersShowAction) handler(httpRequest *http.Request) http.Handler {
	// Define the input struct
	var in struct {
		ID int `json:"id,omitempty"`
	}
	// Unmarshal the request body
	if err := request.Unmarshal(httpRequest, &in); err != nil {
		return &response.Format{
			JSON: response.Status(400).Set("Content-Type", "application/json").JSON(map[string]string{"error": err.Error()}),
		}
	}
	usersController, _nil := loadUsersController()
	if _nil != nil {
		return &response.Format{
			JSON: response.Status(500).Set("Content-Type", "application/json").JSON(map[string]string{"error": _nil.Error()}),
		}
	}
	fn := usersController.Show
	// Call the controller
	in0, in1 := fn(
		in.ID,
	)
	if in1 != nil {
		return &response.Format{
			JSON: response.Status(500).Set("Content-Type", "application/json").JSON(map[string]string{"error": in1.Error()}),
		}
	}

	// Respond
	return &response.Format{
		JSON: response.JSON(in0),
	}
}

// UsersUpdateAction struct
type UsersUpdateAction struct {
}

// Key is a unique identifier of this action
func (u *UsersUpdateAction) Key() string {
	return "/users/update"
}

// Path is the default RESTful path to this action
func (u *UsersUpdateAction) Path() string {
	return "/users/:id"
}

// Method is the default RESTful method of this action
func (u *UsersUpdateAction) Method() string {
	return "PATCH"
}

// ServeHTTP fn
func (u *UsersUpdateAction) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	u.handler(r).ServeHTTP(w, r)
}

// Handler function
func (u *UsersUpdateAction) handler(httpRequest *http.Request) http.Handler {
	// Define the input struct
	var in struct {
		ID   int    `json:"id,omitempty"`
		Name string `json:"name,omitempty"`
		Age  int    `json:"age,omitempty"`
	}
	// Unmarshal the request body
	if err := request.Unmarshal(httpRequest, &in); err != nil {
		return &response.Format{
			JSON: response.Status(400).Set("Content-Type", "application/json").JSON(map[string]string{"error": err.Error()}),
		}
	}
	usersController, _nil := loadUsersController()
	if _nil != nil {
		return &response.Format{
			JSON: response.Status(500).Set("Content-Type", "application/json").JSON(map[string]string{"error": _nil.Error()}),
		}
	}
	fn := usersController.Update
	// Call the controller
	in0 := fn(
		in.ID,
		in.Name,
		in.Age,
	)
	if in0 != nil {
		return &response.Format{
			JSON: response.Status(500).Set("Content-Type", "application/json").JSON(map[string]string{"error": in0.Error()}),
		}
	}

	// Respond
	return &response.Format{
		HTML: response.Status(302).Redirect(response.RedirectPath(httpRequest, "")),
		JSON: response.Status(204),
	}
}

// UsersDeleteAction struct
type UsersDeleteAction struct {
}

// Key is a unique identifier of this action
func (d *UsersDeleteAction) Key() string {
	return "/users/delete"
}

// Path is the default RESTful path to this action
func (d *UsersDeleteAction) Path() string {
	return "/users/:id"
}

// Method is the default RESTful method of this action
func (d *UsersDeleteAction) Method() string {
	return "DELETE"
}

// ServeHTTP fn
func (d *UsersDeleteAction) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	d.handler(r).ServeHTTP(w, r)
}

// Handler function
func (d *UsersDeleteAction) handler(httpRequest *http.Request) http.Handler {
	// Define the input struct
	var in struct {
		ID int `json:"id,omitempty"`
	}
	// Unmarshal the request body
	if err := request.Unmarshal(httpRequest, &in); err != nil {
		return &response.Format{
			JSON: response.Status(400).Set("Content-Type", "application/json").JSON(map[string]string{"error": err.Error()}),
		}
	}
	usersController, _nil := loadUsersController()
	if _nil != nil {
		return &response.Format{
			JSON: response.Status(500).Set("Content-Type", "application/json").JSON(map[string]string{"error": _nil.Error()}),
		}
	}
	fn := usersController.Delete
	// Call the controller
	in0 := fn(
		in.ID,
	)
	if in0 != nil {
		return &response.Format{
			JSON: response.Status(500).Set("Content-Type", "application/json").JSON(map[string]string{"error": in0.Error()}),
		}
	}

	// Respond
	return &response.Format{
		HTML: response.Status(302).Redirect(response.RedirectPath(httpRequest, "")),
		JSON: response.Status(204),
	}
}

// UsersEditAction struct
type UsersEditAction struct {
}

// Key is a unique identifier of this action
func (e *UsersEditAction) Key() string {
	return "/users/edit"
}

// Path is the default RESTful path to this action
func (e *UsersEditAction) Path() string {
	return "/users/:id/edit"
}

// Method is the default RESTful method of this action
func (e *UsersEditAction) Method() string {
	return "GET"
}

// ServeHTTP fn
func (e *UsersEditAction) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e.handler(r).ServeHTTP(w, r)
}

// Handler function
func (e *UsersEditAction) handler(httpRequest *http.Request) http.Handler {
	// Define the input struct
	var in struct {
		ID int `json:"id,omitempty"`
	}
	// Unmarshal the request body
	if err := request.Unmarshal(httpRequest, &in); err != nil {
		return &response.Format{
			JSON: response.Status(400).Set("Content-Type", "application/json").JSON(map[string]string{"error": err.Error()}),
		}
	}
	usersController, _nil := loadUsersController()
	if _nil != nil {
		return &response.Format{
			JSON: response.Status(500).Set("Content-Type", "application/json").JSON(map[string]string{"error": _nil.Error()}),
		}
	}
	fn := usersController.Edit
	// Call the controller
	in0, in1 := fn(
		in.ID,
	)
	if in1 != nil {
		return &response.Format{
			JSON: response.Status(500).Set("Content-Type", "application/json").JSON(map[string]string{"error": in1.Error()}),
		}
	}

	// Respond
	return &response.Format{
		JSON: response.JSON(in0),
	}
}

func loadAboutController() (*about.Controller, error) {
	aboutController := &about.Controller{}
	return aboutController, nil
}

func loadUsersController() (*users.Controller, error) {
	usersController := &users.Controller{}
	return usersController, nil
}
