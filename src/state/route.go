package state

// The route object and methods are used to keep track of navigation in the app
// When the route is empty, a pop will close the app
// This also powers non-tui navigation, like running sudo rover configure

type route struct {
	path []string
}

func NewRoute() route {
	return route{
		path: []string{},
	}
}

func (r *route) Push(path string) {
	r.path = append(r.path, path)
}

func (r *route) Pop() string {
	if len(r.path) == 0 {
		return ""
	}
	path := r.path[len(r.path)-1]
	r.path = r.path[:len(r.path)-1]
	return path
}

// Replaces the top of the stack with the given path
func (c *route) Replace(path string) {
	c.Pop()
	c.Push(path)
}

func (r *route) Peek() string {
	if len(r.path) == 0 {
		return ""
	}
	return r.path[len(r.path)-1]
}

func (r *route) IsEmpty() bool {
	return len(r.path) == 0
}

func (r *route) Clear() {
	r.path = []string{}
}
