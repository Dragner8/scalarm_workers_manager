package main

type IInfrastructureFacade interface {
	StatusCheck() ([]string, error)
	ExtractSiMFiles(*SMRecord)
	ResourceStatus([]string, *SMRecord) (string, error)
	PrepareResource(string, string) (string, error)
}

func NewInfrastructureFacades() map[string]IInfrastructureFacade {
	return map[string]IInfrastructureFacade{
		"qsub": QsubFacade{PLGridFacade{Name: "qsub"}},
		"qcg":  QcgFacade{PLGridFacade{Name: "qcg"}},
		//"private_machine": PrivateMachineFacade{},
	}
}
