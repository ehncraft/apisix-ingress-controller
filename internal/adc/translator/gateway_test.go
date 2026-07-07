// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package translator

import (
	"context"
	"testing"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/apache/apisix-ingress-controller/internal/provider"
)

func TestTranslateSecret_Passthrough(t *testing.T) {
	t.Run("passthrough listener with nil certificateRefs is valid, not an error", func(t *testing.T) {
		tr := &Translator{Log: logr.Discard()}
		gateway := &gatewayv1.Gateway{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "default",
				Name:      "gw",
			},
			Spec: gatewayv1.GatewaySpec{
				Listeners: []gatewayv1.Listener{
					{
						Name: "tls-passthrough",
						TLS: &gatewayv1.GatewayTLSConfig{
							Mode: ptr.To(gatewayv1.TLSModePassthrough),
							// Per the Gateway API spec, Passthrough listeners MUST NOT
							// set certificateRefs: the Gateway never terminates TLS, so
							// there is no certificate for it to present.
							CertificateRefs: nil,
						},
					},
				},
			},
		}
		tctx := provider.NewDefaultTranslateContext(context.Background())

		ssl, err := tr.translateSecret(tctx, gateway.Spec.Listeners[0], gateway)
		require.NoError(t, err)
		assert.Empty(t, ssl)
	})
}
