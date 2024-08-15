package utils

type GraphNode struct {
	Id        int
	DependsOn *int
}

func (n GraphNode) resolveDependency(items []GraphNode, deps []GraphNode) ([]GraphNode, bool) {
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
			res, successful := i.resolveDependency(items, deps)
			if !successful {
				return deps, false
			}
			deps = res
			return deps, true
		}
	}
	return deps, false
}
