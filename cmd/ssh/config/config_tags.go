package config

// FilterTagsPredicate returns a predicate to test whether a set of tags shall be accepted based on all tags (include/exclude) stored in the Configuration.
// It will return nil when there are no tags configured hence all tags arr allowed
func (c *Configuration) FilterTagsPredicate() func([]string) bool {

	makeKeys := func(strings []string) map[string]bool {
		allTags := make(map[string]bool)

		for _, s := range strings {
			allTags[s] = true
		}

		return allTags
	}

	if len(c.IncludeTags) == 0 && len(c.ExcludeTags) == 0 {
		return nil
	}

	return func(tags []string) bool {

		hostTags := makeKeys(tags)

		for _, s := range c.IncludeTags {
			if _, found := hostTags[s]; !found {
				return false
			}
		}

		for _, s := range c.ExcludeTags {
			if _, found := hostTags[s]; found {
				return false
			}
		}

		return true

	}

}
