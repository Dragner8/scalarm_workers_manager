package main

type IInfrastructureFacade interface {
	StatusCheck() ([]string, error)
	PrepareResource(string, string) (string, error)
	ExtractSiMFiles(*SMRecord) error
	ResourceStatus([]string, *SMRecord) (string, error)
}

func NewInfrastructureFacades() map[string]IInfrastructureFacade {
	return map[string]IInfrastructureFacade{
		"qsub":            QsubFacade{PLGridFacade{}},
		"qcg":             QcgFacade{PLGridFacade{}},
		"private_machine": PrivateMachineFacade{},
	}
}
