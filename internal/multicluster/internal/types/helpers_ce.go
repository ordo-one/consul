// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !consulent

package types

import (
	"fmt"
	"github.com/hashicorp/consul/internal/resource"
	multiclusterv1alpha1 "github.com/hashicorp/consul/proto-public/pbmulticluster/v1alpha1"
	"github.com/hashicorp/consul/proto-public/pbresource"
	"github.com/hashicorp/go-multierror"
)

func validateExportedServicesConsumer(consumer *multiclusterv1alpha1.ExportedServicesConsumer, merr error, indx int) error {
	switch consumer.GetConsumerTenancy().(type) {
	case *multiclusterv1alpha1.ExportedServicesConsumer_Partition:
		{
			if consumer.GetPartition() != resource.DefaultPartitionName {
				merr = multierror.Append(merr, resource.ErrInvalidListElement{
					Name:    "Partition",
					Index:   indx,
					Wrapped: fmt.Errorf("can only be set in Enterprise"),
				})
			}
		}
	case *multiclusterv1alpha1.ExportedServicesConsumer_Peer:
		{
			if consumer.GetPeer() != "" {
				merr = multierror.Append(merr, resource.ErrInvalidListElement{
					Name:    "Peer",
					Index:   indx,
					Wrapped: fmt.Errorf("can only be set in Enterprise"),
				})
			}
		}
	case *multiclusterv1alpha1.ExportedServicesConsumer_SamenessGroup:
		{
			if consumer.GetSamenessGroup() != "" {
				merr = multierror.Append(merr, resource.ErrInvalidListElement{
					Name:    "Sameness Group",
					Index:   indx,
					Wrapped: fmt.Errorf("can only be set in Enterprise"),
				})
			}
		}
	}
	return merr
}

func ValidateExportedServicesEnterprise(_ *pbresource.Resource, exportedService *multiclusterv1alpha1.ExportedServices) error {
	var merr error

	for indx, consumer := range exportedService.Consumers {
		merr = validateExportedServicesConsumer(consumer, merr, indx)
	}

	return merr
}

func ValidateNamespaceExportedServicesEnterprise(_ *pbresource.Resource, exportedService *multiclusterv1alpha1.NamespaceExportedServices) error {
	var merr error

	for indx, consumer := range exportedService.Consumers {
		merr = validateExportedServicesConsumer(consumer, merr, indx)
	}

	return merr
}

func ValidatePartitionExportedServicesEnterprise(_ *pbresource.Resource, exportedService *multiclusterv1alpha1.PartitionExportedServices) error {
	var merr error

	for indx, consumer := range exportedService.Consumers {
		merr = validateExportedServicesConsumer(consumer, merr, indx)
	}

	return merr
}

func ValidateComputedExportedServicesEnterprise(_ *pbresource.Resource, computedExportedServices *multiclusterv1alpha1.ComputedExportedServices) error {

	var merr error

	for indx, consumer := range computedExportedServices.GetConsumers() {
		for _, computedExportedServiceConsumer := range consumer.GetConsumers() {
			switch computedExportedServiceConsumer.GetConsumerTenancy().(type) {
			case *multiclusterv1alpha1.ComputedExportedServicesConsumer_Partition:
				{
					if computedExportedServiceConsumer.GetPartition() != resource.DefaultPartitionName {
						merr = multierror.Append(merr, resource.ErrInvalidListElement{
							Name:    "Partition",
							Index:   indx,
							Wrapped: fmt.Errorf("can only be set in Enterprise"),
						})
					}
				}
			case *multiclusterv1alpha1.ComputedExportedServicesConsumer_Peer:
				{
					if computedExportedServiceConsumer.GetPeer() != "" {
						merr = multierror.Append(merr, resource.ErrInvalidListElement{
							Name:    "Peer",
							Index:   indx,
							Wrapped: fmt.Errorf("can only be set in Enterprise"),
						})
					}
				}
			}
		}
	}

	return merr
}

func MutateComputedExportedServices(res *pbresource.Resource) error {
	var ces multiclusterv1alpha1.ComputedExportedServices

	if err := res.Data.UnmarshalTo(&ces); err != nil {
		return err
	}

	var changed bool

	for _, cesConsumer := range ces.GetConsumers() {
		for _, consumer := range cesConsumer.GetConsumers() {
			switch t := consumer.GetConsumerTenancy().(type) {
			case *multiclusterv1alpha1.ComputedExportedServicesConsumer_Partition:
				if t.Partition == "" {
					changed = true
					t.Partition = resource.DefaultPartitionName
				}
			}
		}
	}

	if !changed {
		return nil
	}

	return res.Data.MarshalFrom(&ces)
}

func MutateExportedServices(res *pbresource.Resource) error {
	var es multiclusterv1alpha1.ExportedServices

	if err := res.Data.UnmarshalTo(&es); err != nil {
		return err
	}

	var changed bool

	for _, consumer := range es.Consumers {
		changed = changed || updatePartitionIfNotSet(consumer)
	}

	if !changed {
		return nil
	}

	return res.Data.MarshalFrom(&es)
}

func MutateNamespaceExportedServices(res *pbresource.Resource) error {
	var nes multiclusterv1alpha1.NamespaceExportedServices

	if err := res.Data.UnmarshalTo(&nes); err != nil {
		return err
	}

	var changed bool

	for _, consumer := range nes.Consumers {
		changed = changed || updatePartitionIfNotSet(consumer)
	}

	if !changed {
		return nil
	}

	return res.Data.MarshalFrom(&nes)
}

func MutatePartitionExportedServices(res *pbresource.Resource) error {
	var pes multiclusterv1alpha1.PartitionExportedServices

	if err := res.Data.UnmarshalTo(&pes); err != nil {
		return err
	}

	var changed bool

	for _, consumer := range pes.Consumers {
		changed = changed || updatePartitionIfNotSet(consumer)
	}

	if !changed {
		return nil
	}

	return res.Data.MarshalFrom(&pes)
}

func updatePartitionIfNotSet(consumer *multiclusterv1alpha1.ExportedServicesConsumer) bool {
	var updated bool

	switch t := consumer.GetConsumerTenancy().(type) {
	case *multiclusterv1alpha1.ExportedServicesConsumer_Partition:
		if t.Partition == "" {
			updated = true
			t.Partition = resource.DefaultPartitionName
		}
	}
	return updated
}