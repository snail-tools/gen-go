package gengo

type Generator interface {
	Init(sw SnippetWriter)
	New() Generator
	Generate()
}
