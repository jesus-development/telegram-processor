package processor

type ProcessorOption func(*messageProcessor)

func WithMessagesRepository(repository MessagesRepository) ProcessorOption {
	return func(p *messageProcessor) {
		p.MessagesRepository = repository
	}
}

func WithEmbeddingService(service EmbeddingService) ProcessorOption {
	return func(p *messageProcessor) {
		p.EmbeddingService = service
	}
}
