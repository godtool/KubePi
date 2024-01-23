package imagerepo

import V1ImageRepo "github.com/godtool/kubeone/service/model/v1/imagerepo"

type RepoConfig struct {
	V1ImageRepo.ImageRepo
	Page          int    `json:"page" validate:"required"`
	Limit         int    `json:"limit" validate:"required"`
	Search        string `json:"search"`
	ContinueToken string `json:"continueToken"`
}
