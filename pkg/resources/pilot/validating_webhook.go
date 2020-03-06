/*
Copyright 2019 Banzai Cloud.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package pilot

import (
	"github.com/banzaicloud/istio-operator/pkg/util"
	admissionv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func (r *Reconciler) webhooks() []admissionv1beta1.Webhook {
	if !util.PointerToBool(r.Config.Spec.Istiod.Enabled) {
		return nil
	}

	ignore := admissionv1beta1.Ignore
	se := admissionv1beta1.SideEffectClassNone
	return []admissionv1beta1.Webhook{
		{
			Name: "validation.istio.io",
			ClientConfig: admissionv1beta1.WebhookClientConfig{
				Service: &admissionv1beta1.ServiceReference{
					Name:      serviceNameIstiod,
					Namespace: r.Config.Namespace,
					Path:      util.StrPointer("/validate"),
				},
				CABundle: []byte{},
			},
			Rules: []admissionv1beta1.RuleWithOperations{
				{
					Operations: []admissionv1beta1.OperationType{
						admissionv1beta1.Create,
						admissionv1beta1.Update,
					},
					Rule: admissionv1beta1.Rule{
						APIGroups:   []string{"config.istio.io", "rbac.istio.io", "security.istio.io", "authentication.istio.io", "networking.istio.io"},
						APIVersions: []string{"*"},
						Resources:   []string{"*"},
					},
				},
			},
			FailurePolicy: &ignore,
			SideEffects:   &se,
		},
	}
}

func (r *Reconciler) validatingWebhook() runtime.Object {
	return &admissionv1beta1.ValidatingWebhookConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name:   validatingWebhookName,
			Labels: util.MergeStringMaps(istiodLabels, istiodLabelSelector),
		},
		Webhooks: r.webhooks(),
	}
}

func (r *Reconciler) validatingWebhookGalley() runtime.Object {
	return &admissionv1beta1.ValidatingWebhookConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name:   validatingWebhookNameGalley,
			Labels: util.MergeStringMaps(galleyLabels, galleyLabelSelector),
		},
		Webhooks: nil,
	}
}