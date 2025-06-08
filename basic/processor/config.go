package processor

type Processor interface {
	Pre()
	Post()
	Process()
}
