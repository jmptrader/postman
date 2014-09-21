package tunnel

type Command struct {
	Id     string
	Action string
	Args   interface{}
	client *Client
}

func (cm *Command) response(code string) {

}
