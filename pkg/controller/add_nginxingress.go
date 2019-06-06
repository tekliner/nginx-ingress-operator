package controller

import (
	"github.com/tekliner/nginx-ingress-operator/pkg/controller/nginxingress"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, nginxingress.Add)
}
