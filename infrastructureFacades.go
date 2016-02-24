package main

type IInfrastructureFacade interface {
	StatusCheck() ([]string, error)
	SetId(*SMRecord, string)
	PrepareResource(string, string) (string, error)
	ExtractSiMFiles(*SMRecord) error
	ResourceStatus([]string, *SMRecord) (string, error)
}

func NewInfrastructureFacades() map[string]IInfrastructureFacade {
	return map[string]IInfrastructureFacade{
		"qsub":            QsubFacade{BashExecutor{}, PLGridFacade{}},
		"qcg":             QcgFacade{PLGridFacade{}},
		"private_machine": PrivateMachineFacade{},
		"slurm":           SlurmFacade{BashExecutor{}, PLGridFacade{}},
	}
}
