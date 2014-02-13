package bitverse

type Service struct {
	name     string
	observer ServiceObserver
	edgeNode *EdgeNode
}

func composeService(name string, observe ServiceObserver, edgeNode *EdgeNode) *Service {
	service := new(Service)
	service.name = name
	service.observer = observe
	service.edgeNode = edgeNode
	return service
}

func (service *Service) SendMsg(dst string, payload string) {
	service.edgeNode.send(dst, payload, service.name)
}

func (service *Service) Nop() {
}
