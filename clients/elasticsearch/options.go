package elasticsearch

import (
	"fmt"
	es6 "github.com/elastic/go-elasticsearch/v6"
	es7 "github.com/elastic/go-elasticsearch/v7"
	es8 "github.com/elastic/go-elasticsearch/v8"
	"io/ioutil"
)

func setupOptions(c *Client, options map[string]interface{}) error {
	// Setup ES version and client
	if val, ok := options["version"].(string); ok {
		if val == "6" {
			c.version = "6"
			config, err := setupEs6Config(options)
			if err != nil {
				return fmt.Errorf("issue building es6 client config: %v", err)
			}
			c.es6Client, err = es6.NewClient(config)
			if err != nil {
				return fmt.Errorf("issue setting up es6 client: %v", err)
			}
		} else if val == "7" {
			c.version = "7"
			config, err := setupEs7Config(options)
			if err != nil {
				return fmt.Errorf("issue building es7 client config: %v", err)
			}
			c.es7Client, err = es7.NewClient(config)
			if err != nil {
				return fmt.Errorf("issue setting up es7 client: %v", err)
			}
		} else if val == "8" {
			c.version = "8"
			config, err := setupEs8Config(options)
			if err != nil {
				return fmt.Errorf("issue building es8 client config: %v", err)
			}
			c.es8Client, err = es8.NewClient(config)
			if err != nil {
				return fmt.Errorf("issue setting up es8 client: %v", err)
			}
		} else {
			options["version"] = "6"
			c.version = "6"
			config, err := setupEs6Config(options)
			if err != nil {
				return err
			}
			c.es6Client, err = es6.NewClient(config)
			if err != nil {
				return err
			}
		}
	} else {
		options["version"] = "6"
		c.version = "6"
		config, err := setupEs6Config(options)
		if err != nil {
			return err
		}
		c.es6Client, err = es6.NewClient(config)
		if err != nil {
			return err
		}
	}

	if val, ok := options["index"].(string); ok {
		c.index = val
	} else {
		return fmt.Errorf("missing elasticsearch index")
	}

	if val, ok := options["query-file"].(string); ok {
		q, err := ioutil.ReadFile(val)
		if err != nil {
			return fmt.Errorf("issue getting query file: %v", err)
		}
		c.query = q
	} else {
		return fmt.Errorf("missing elasticsearch query")
	}

	return nil
}

func setupEs6Config(options map[string]interface{}) (es6.Config, error) {
	config := es6.Config{}
	if val, ok := options["ca-cert"].(string); ok && val != "" {
		cert, err := ioutil.ReadFile(val)
		if err != nil {
			return config, fmt.Errorf("issue reading ca cert file for es: %v", err)
		}
		config.CACert = cert
	}

	if val, ok := options["addresses"].([]string); ok {
		config.Addresses = val
	} else {
		return config, fmt.Errorf("missing elasticsearch host addresses")
	}

	if val, ok := options["username"].(string); ok {
		config.Username = val
	}

	if val, ok := options["password"].(string); ok {
		config.Password = val
	}

	return config, nil
}

func setupEs7Config(options map[string]interface{}) (es7.Config, error) {
	config := es7.Config{}
	if val, ok := options["ca-cert"].(string); ok && val != "" {
		cert, err := ioutil.ReadFile(val)
		if err != nil {
			return config, fmt.Errorf("issue reading ca cert file for es: %v", err)
		}
		config.CACert = cert
	}

	if val, ok := options["addresses"].([]string); ok {
		config.Addresses = val
	} else {
		return config, fmt.Errorf("missing elasticsearch host addresses")
	}

	if val, ok := options["username"].(string); ok {
		config.Username = val
	}

	if val, ok := options["password"].(string); ok {
		config.Password = val
	}

	return config, nil
}

func setupEs8Config(options map[string]interface{}) (es8.Config, error) {
	config := es8.Config{}
	if val, ok := options["ca-cert"].(string); ok && val != "" {
		cert, err := ioutil.ReadFile(val)
		if err != nil {
			return config, fmt.Errorf("issue reading ca cert file for es: %v", err)
		}
		config.CACert = cert
	}

	if val, ok := options["addresses"].([]string); ok {
		config.Addresses = val
	} else {
		return config, fmt.Errorf("missing elasticsearch host addresses")
	}

	if val, ok := options["username"].(string); ok {
		config.Username = val
	}

	if val, ok := options["password"].(string); ok {
		config.Password = val
	}

	return config, nil
}