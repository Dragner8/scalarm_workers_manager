package main

type IInfrastructureFacade interface {
	StatusCheck() ([]string, error)
	ExtractSiMFiles(sm_record *Sm_record)
	ResourceStatus(statusArray []string, jobID string) (string, error)
	PrepareResource(ids, command string) (string, error)
}

func NewInfrastructureFacades() map[string]IInfrastructureFacade {
	return map[string]IInfrastructureFacade{
		"qsub": QsubFacade{PLGridFacade{Name: "qsub"}},
		"qcg":  QcgFacade{PLGridFacade{Name: "qcg"}},
		//"private_machine": PrivateMachineFacade{},
	}
}
