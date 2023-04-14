package services

type DummyNode struct {
	name       string
	neighbours []*DummyNode
}

func (n *DummyNode) bfs(depth int) map[string]DummyNode {

	queue := []*DummyNode{n}
	visited := map[string]DummyNode{}

	for len(queue) > 0 || depth > 0 {
		level_size := len(queue)
		for i := 0; i < level_size; i++ {
			current := queue[0]
			queue = queue[1:]
			visited[current.name] = *current
			for _, nghr := range n.neighbours {
				if _, exists := visited[nghr.name]; !exists {
					queue = append(queue, nghr)
				}
			}
		}
		depth -= 1
	}
	return visited
}
