package app

type searchResult struct {
	score int
}

func search(move boardMove, results chan searchResult) {
	// TODO: do minimax search on 'boardMove'
	results <- searchResult{1}
}
