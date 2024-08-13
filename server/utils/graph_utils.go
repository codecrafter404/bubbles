package utils

type GraphNode struct {
	Id        int
	DependsOn *int
}

func (n GraphNode) resolveDependency(items []GraphNode, deps []GraphNode) ([]GraphNode, bool) {
	for _, i := range items {
		if i.Id == *n.DependsOn {
			deps = append(deps, i)

			return deps, true
		}
	}
	//TODO:
}

// TODO:
func HasCycle(items []GraphNode, visited []int) bool {
	for _, item := range items {
		visited := []int{}
		for _, s := range items {
			if item.DependsOn == s.Id {

			}
		}
	}
	return false
}
