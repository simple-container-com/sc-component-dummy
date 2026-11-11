// SPDX-License-Identifier: MIT
// Copyright (c) Simple Container

package dummy

import (
	"testing"

	. "github.com/onsi/gomega"

	scPlugin "github.com/simple-container-com/api/pkg/plugin"
	scpluginv1 "github.com/simple-container-com/api/pkg/plugin/proto/v1"
)

// PreDestroy for the dummy is a fire-and-forget log + return. A real
// component's test would assert on the cleanup side effects — for dummy,
// we just prove the hook is invokable with the Phase 1e v2 request shape.
func TestPreDestroy_FireAndForget(t *testing.T) {
	RegisterTestingT(t)

	pub := scPlugin.NewPublisher(0)
	req := &scpluginv1.PreDestroyRequest{
		Header:       &scpluginv1.RequestHeader{StackName: "smoke-dummy", Environment: "smoke"},
		ResourceType: "sc-dummy",
	}

	err := PreDestroy(t.Context(), pub, req)
	Expect(err).ToNot(HaveOccurred())
}
