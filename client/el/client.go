package el

func New(host string) (*Client, error) {
	client := &Client{
		host: host,
	}

	return client, nil
}
