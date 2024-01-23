package project

import v1 "github.com/godtool/kubeone/service/model/v1"

type Project struct {
	v1.BaseModel `storm:"inline"`
	v1.Metadata  `storm:"inline"`
	Users        []string `json:"users"`
	ClusterRef   string   `json:"clusterRef" storm:"index"`
	ProjectName  string   `json:"projectName"`
	Certificate  []byte   `json:"certificate"`
}
