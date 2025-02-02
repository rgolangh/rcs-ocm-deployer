package webhooks

import (
	"context"
	"strings"

	rcsv1alpha1 "github.com/dana-team/container-app-operator/api/v1alpha1"

	"net/http"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

type cappValidator struct {
	Client  client.Client
	Decoder *admission.Decoder
	Log     logr.Logger
}

// +kubebuilder:webhook:path=/validate-v1-namespace,mutating=false,sideEffects=NoneOnDryRun,failurePolicy=fail,groups="core",resources=namespaces,verbs=delete,versions=v1,name=namespace.dana.io,admissionReviewVersions=v1;v1beta1

func (c *cappValidator) Handle(ctx context.Context, req admission.Request) admission.Response {
	log := c.Log.WithValues("webhook", "capp Webhook", "Name", req.Name)
	log.Info("webhook request received")
	capp := rcsv1alpha1.Capp{}
	if err := c.Decoder.DecodeRaw(req.OldObject, &capp); err != nil {
		log.Error(err, "could not decode capp object")
		return admission.Errored(http.StatusBadRequest, err)
	}
	return c.handle(ctx, req, capp)

}

func (c *cappValidator) handle(ctx context.Context, req admission.Request, capp rcsv1alpha1.Capp) admission.Response {
	if !isScaleMetricSupported(capp) {
		return admission.Denied(unSupportedScaleMetric + " " + capp.Spec.ScaleMetric + " the avilable options are " + strings.Join(SupportedScaleMetrics, ","))
	}
	if !isSiteClusterName(capp, c.Client, ctx) {
		return admission.Denied(unSupportedSite)
	}
	return admission.Allowed("")
}
