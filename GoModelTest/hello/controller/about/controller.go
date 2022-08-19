package about

type Controller struct {
	// Dependencies...
}

// About struct
type About struct {
	// Fields...
	Name string
	Age  int
}

func (c *Controller) Index() ([]*About, error) {
	return []*About{
		{
			Name: "string",
			Age:  12,
		},
	}, nil
}
