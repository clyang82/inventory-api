package k8sclusters

import (
	"context"

	"gorm.io/gorm"

	"github.com/go-kratos/kratos/v2/log"

	biz "github.com/project-kessel/inventory-api/internal/biz/k8sclusters"
)

type k8sclustersRepo struct {
	DB  *gorm.DB
	Log *log.Helper
}

func New(g *gorm.DB, l *log.Helper) *k8sclustersRepo {
	return &k8sclustersRepo{
		DB:  g,
		Log: l,
	}
}

func (r *k8sclustersRepo) Save(context.Context, *biz.K8sCluster) (*biz.K8sCluster, error) {
	return nil, nil
}

func (r *k8sclustersRepo) Update(context.Context, *biz.K8sCluster) (*biz.K8sCluster, error) {
	return nil, nil
}

func (r *k8sclustersRepo) Delete(context.Context, int64) error {
	return nil
}

func (r *k8sclustersRepo) FindByID(context.Context, int64) (*biz.K8sCluster, error) {
	return nil, nil
}

func (r *k8sclustersRepo) ListAll(context.Context) ([]*biz.K8sCluster, error) {
	return nil, nil
}
