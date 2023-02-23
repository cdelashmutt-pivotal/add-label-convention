/*
Copyright 2020 VMware Inc.

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

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"

	"github.com/vmware-tanzu/cartographer-conventions/webhook"
)

func conventionHandler(template *corev1.PodTemplateSpec, images []webhook.ImageConfig) ([]string, error) {
	if template.ObjectMeta.Labels == nil {
		template.ObjectMeta.Labels = map[string]string{}
	}

	template.ObjectMeta.Labels["business-unit"] = "gbs"
	return []string{"add-bu-label"}, nil
}

func main() {
	ctx := context.Background()
	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}

	zapLog, err := zap.NewProductionConfig().Build()
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	logger := zapr.NewLogger(zapLog)
	ctx = logr.NewContext(ctx, logger)

	http.HandleFunc("/", webhook.ConventionHandler(ctx, conventionHandler))
	log.Fatal(webhook.NewConventionServer(ctx, fmt.Sprintf(":%s", port)))
}
