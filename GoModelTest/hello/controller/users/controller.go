package users

type Controller struct {
	// Dependencies...
}

// User struct
type User struct {
	// Fields...
	Name string
	Age  int
}

// GET /users
func (c *Controller) Index() ([]*User, error) {
	return []*User{
		{
			Name: "string",
			Age:  12,
		},
	}, nil
}

// New user page
// GET /users/new
func (c *Controller) New() {}

// Create a new user
// POST /users
func (c *Controller) Create(name string, age int) (*User, error) {
	return &User{}, nil
}

// Show a user
// GET /users/:id
func (c *Controller) Show(id int) (*User, error) {
	return &User{}, nil
}

// Update a user
// PATCH /users/:id
func (c *Controller) Update(id int, name string, age int) error {
	return nil
}

// Delete a user
// DELETE /users/:id
func (c *Controller) Delete(id int) error {
	return nil
}

// Edit user page
// GET /users/:id/edit
func (c *Controller) Edit(id int) (*User, error) {
	return &User{}, nil
}
