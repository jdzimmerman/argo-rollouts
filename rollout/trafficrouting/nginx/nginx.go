package nginx

import (
	"fmt"
	"github.com/argoproj/argo-rollouts/pkg/apis/rollouts/v1alpha1"
	logutil "github.com/argoproj/argo-rollouts/utils/log"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/record"
)

func NewNginxReconciler(r *v1alpha1.Rollout, client dynamic.Interface, recorder record.EventRecorder) *Reconciler {
	return &Reconciler{
		rollout: r,
		log:     logutil.WithRollout(r),

		client:   client,
		recorder: recorder,
	}
}
const Type = "nginx"

type Reconciler struct {
	rollout  *v1alpha1.Rollout
	log      *logrus.Entry
	client   dynamic.Interface
	recorder record.EventRecorder
}

func (r *Reconciler) Reconcile(desiredWeight int32) error {
	canarySvc := r.rollout.Spec.Strategy.Canary.CanaryService
	stableSvc := r.rollout.Spec.Strategy.Canary.StableService

	fmt.Println(desiredWeight)
	fmt.Println(canarySvc)
	fmt.Println(stableSvc)
	return nil

}

func (r *Reconciler) Type() string {
	return Type

}
func GetRolloutIngressName(rollout *v1alpha1.Rollout) string {

	return rollout.Spec.Strategy.Canary.TrafficRouting.Nginx.StableIngress

}

func SetCanaryIngressName(rollout *v1alpha1.Rollout) string {
	canaryName := fmt.Sprintf("%s-canary",GetRolloutIngressName(rollout))
	return canaryName
}

func (r *Reconciler) CreateCanaryIngress() {
	stableIngress := GetRolloutIngressName(r.rollout)
	canaryIngress := SetCanaryIngressName(r.rollout)
	client := r.client.Resource().Namespace(r.rollout.Namespace)
	client.Create(canaryIngress)

	fmt.Println(stableIngress, canaryIngress)
}

/*
this is the key logic from istio for routing to the specific services
func (r *Reconciler) generateVirtualServicePatches(httpRoutes []httpRoute, desiredWeight int64) virtualServicePatches {
	canarySvc := r.rollout.Spec.Strategy.Canary.CanaryService
	stableSvc := r.rollout.Spec.Strategy.Canary.StableService
	routes := map[string]bool{}
	for _, r := range r.rollout.Spec.Strategy.Canary.TrafficRouting.Istio.VirtualService.Routes {
		routes[r] = true
	}

	patches := virtualServicePatches{}
	for i := range httpRoutes {
		route := httpRoutes[i]
		if !routes[route.Name] {
			continue
		}
		for j := range route.Route {
			destination := httpRoutes[i].Route[j]
			host := destination.Destination.Host
			weight := destination.Weight
			if host == canarySvc && weight != desiredWeight {
				patch := virtualServicePatch{
					routeIndex:       i,
					destinationIndex: j,
					weight:           desiredWeight,
				}
				patches = append(patches, patch)
			}
			if host == stableSvc && weight != 100-desiredWeight {
				patch := virtualServicePatch{
					routeIndex:       i,
					destinationIndex: j,
					weight:           100 - desiredWeight,
				}
				patches = append(patches, patch)
			}
		}
	}
	return patches
}

 */
