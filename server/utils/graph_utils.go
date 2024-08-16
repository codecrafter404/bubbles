package utils

type GraphNode struct {
	Id        int
	DependsOn *int
}

func (n GraphNode) ResolveDependency(items []GraphNode, deps []GraphNode) ([]GraphNode, bool) {
	if n.DependsOn == nil {
		return deps, true
	}
	for _, i := range items {
		if *n.DependsOn == i.Id { // found the dependency
			for _, x := range deps { // check if we already have the dependency
				if x.Id == *n.DependsOn {
					return deps, false
				}
			}
			deps = append(deps, i)
			res, successful := i.ResolveDependency(items, deps)
			if !successful {
				return deps, false
			}
			deps = res
			return deps, true
		}
	}
	return deps, false
}

// returns true if there is NO loop
func CheckDependencyLoop(items []GraphNode) bool {
	for _, x := range items {
		_, successful := x.ResolveDependency(items, []GraphNode{})
		if !successful {
			return false
		}
	}
	return true
}
