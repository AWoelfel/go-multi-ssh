package config

import (
	"fmt"
)

// Clients returns a prefiltered set of target clients based on the Configuration.
func (c *Configuration) Clients() ([]SshClientConfig, error) {

	allHosts, err := ReadClients(c.IndexFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read index file (%w)", err)
	}

	tagFilter := c.FilterTagsPredicate()

	//no filter set; use all hosts
	if tagFilter == nil {
		return allHosts, nil
	}

	//filter set; use filtered hosts
	var filteredHosts []SshClientConfig

	for i := 0; i < len(allHosts); i++ {
		if !tagFilter(allHosts[i].Tags) {
			continue
		}
		filteredHosts = append(filteredHosts, allHosts[i])
	}

	return filteredHosts, nil
}
