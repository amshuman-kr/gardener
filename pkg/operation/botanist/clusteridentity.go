// Copyright (c) 2021 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package botanist

import (
	"context"
	"fmt"

	"github.com/gardener/gardener/pkg/extensions"
	"github.com/gardener/gardener/pkg/operation/botanist/component/clusteridentity"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

// EnsureShootClusterIdentity ensures that Shoot's `status.clusterIdentity` field is set and updates the Cluster resource in
// the seed if necessary.
func (b *Botanist) EnsureShootClusterIdentity(ctx context.Context) error {
	if b.Shoot.Info.Status.ClusterIdentity == nil {
		clusterIdentity := fmt.Sprintf("%s-%s-%s", b.Shoot.SeedNamespace, b.Shoot.Info.Status.UID, b.GardenClusterIdentity)

		patch := client.MergeFrom(b.Shoot.Info.DeepCopy())
		b.Shoot.Info.Status.ClusterIdentity = &clusterIdentity

		if err := b.K8sGardenClient.Client().Status().Patch(ctx, b.Shoot.Info, patch); err != nil {
			return err
		}

		if err := extensions.SyncClusterResourceToSeed(ctx, b.K8sSeedClient.Client(), b.Shoot.SeedNamespace, b.Shoot.Info, nil, nil); err != nil {
			return err
		}
	}

	return nil
}

// DefaultClusterIdentity returns a deployer for the shoot's cluster-identity.
func (b *Botanist) DefaultClusterIdentity() clusteridentity.Interface {
	return clusteridentity.NewForShoot(b.K8sSeedClient.Client(), b.Shoot.SeedNamespace, "")
}

// DeployClusterIdentity deploys the shoot's cluster-identity.
func (b *Botanist) DeployClusterIdentity(ctx context.Context) error {
	if v := b.Shoot.Info.Status.ClusterIdentity; v != nil {
		b.Shoot.Components.SystemComponents.ClusterIdentity.SetIdentity(*v)
	}

	return b.Shoot.Components.SystemComponents.ClusterIdentity.Deploy(ctx)
}