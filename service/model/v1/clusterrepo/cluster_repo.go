package clusterrepo

import v1 "github.com/godtool/kubeone/service/model/v1"

type ClusterRepo struct {
	v1.BaseModel `storm:"inline"`
	v1.Metadata  `storm:"inline"`
	Cluster      string `json:"cluster"`
	Repo         string `json:"repo"`
}
