package util

type ResponseReferences struct {
	Data []NodeReferences `json:"data"`
}

type NodeReferences struct {
	IsInfluential bool `json:"isInfluential"`
	CitedPaper    struct {
		RefPaperId string `json:"paperId"`
	} `json:"citedPaper"`
	CitingPaper struct {
		CitPaperId string `json:"paperId"`
	} `json:"citingPaper"`
}

type Author struct {
	AuthorId string `json:"authorId"`
	Name     string `json:"name"`
}

type Node struct {
	PaperId       string   `json:"paperId"`
	CitationCount int64    `json:"citationCount"`
	Year          int64    `json:"year"`
	Title         string   `json:"title"`
	Authors       []Author `json:"authors"`
	Reference     []Node   `json:"node"`
}
