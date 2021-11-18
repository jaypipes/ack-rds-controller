// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

// Code generated by ack-generate. DO NOT EDIT.

package db_security_group

import (
	"context"
	"reflect"
	"strings"

	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	ackcondition "github.com/aws-controllers-k8s/runtime/pkg/condition"
	ackerr "github.com/aws-controllers-k8s/runtime/pkg/errors"
	ackrtlog "github.com/aws-controllers-k8s/runtime/pkg/runtime/log"
	"github.com/aws/aws-sdk-go/aws"
	svcsdk "github.com/aws/aws-sdk-go/service/rds"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	svcapitypes "github.com/aws-controllers-k8s/rds-controller/apis/v1alpha1"
)

// Hack to avoid import errors during build...
var (
	_ = &metav1.Time{}
	_ = strings.ToLower("")
	_ = &aws.JSONValue{}
	_ = &svcsdk.RDS{}
	_ = &svcapitypes.DBSecurityGroup{}
	_ = ackv1alpha1.AWSAccountID("")
	_ = &ackerr.NotFound
	_ = &ackcondition.NotManagedMessage
	_ = &reflect.Value{}
)

// sdkFind returns SDK-specific information about a supplied resource
func (rm *resourceManager) sdkFind(
	ctx context.Context,
	r *resource,
) (latest *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkFind")
	defer exit(err)
	// If any required fields in the input shape are missing, AWS resource is
	// not created yet. Return NotFound here to indicate to callers that the
	// resource isn't yet created.
	if rm.requiredFieldsMissingFromReadManyInput(r) {
		return nil, ackerr.NotFound
	}

	input, err := rm.newListRequestPayload(r)
	if err != nil {
		return nil, err
	}
	var resp *svcsdk.DescribeDBSecurityGroupsOutput
	resp, err = rm.sdkapi.DescribeDBSecurityGroupsWithContext(ctx, input)
	rm.metrics.RecordAPICall("READ_MANY", "DescribeDBSecurityGroups", err)
	if err != nil {
		if awsErr, ok := ackerr.AWSError(err); ok && awsErr.Code() == "DBSecurityGroupNotFound" {
			return nil, ackerr.NotFound
		}
		return nil, err
	}

	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := r.ko.DeepCopy()

	found := false
	for _, elem := range resp.DBSecurityGroups {
		if elem.DBSecurityGroupArn != nil {
			if ko.Status.ACKResourceMetadata == nil {
				ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
			}
			tmpARN := ackv1alpha1.AWSResourceName(*elem.DBSecurityGroupArn)
			ko.Status.ACKResourceMetadata.ARN = &tmpARN
		}
		if elem.DBSecurityGroupName != nil {
			ko.Spec.Name = elem.DBSecurityGroupName
		} else {
			ko.Spec.Name = nil
		}
		if elem.EC2SecurityGroups != nil {
			f3 := []*svcapitypes.EC2SecurityGroup{}
			for _, f3iter := range elem.EC2SecurityGroups {
				f3elem := &svcapitypes.EC2SecurityGroup{}
				if f3iter.EC2SecurityGroupId != nil {
					f3elem.EC2SecurityGroupID = f3iter.EC2SecurityGroupId
				}
				if f3iter.EC2SecurityGroupName != nil {
					f3elem.EC2SecurityGroupName = f3iter.EC2SecurityGroupName
				}
				if f3iter.EC2SecurityGroupOwnerId != nil {
					f3elem.EC2SecurityGroupOwnerID = f3iter.EC2SecurityGroupOwnerId
				}
				if f3iter.Status != nil {
					f3elem.Status = f3iter.Status
				}
				f3 = append(f3, f3elem)
			}
			ko.Status.EC2SecurityGroups = f3
		} else {
			ko.Status.EC2SecurityGroups = nil
		}
		if elem.IPRanges != nil {
			f4 := []*svcapitypes.IPRange{}
			for _, f4iter := range elem.IPRanges {
				f4elem := &svcapitypes.IPRange{}
				if f4iter.CIDRIP != nil {
					f4elem.CIDRIP = f4iter.CIDRIP
				}
				if f4iter.Status != nil {
					f4elem.Status = f4iter.Status
				}
				f4 = append(f4, f4elem)
			}
			ko.Status.IPRanges = f4
		} else {
			ko.Status.IPRanges = nil
		}
		if elem.OwnerId != nil {
			ko.Status.OwnerID = elem.OwnerId
		} else {
			ko.Status.OwnerID = nil
		}
		if elem.VpcId != nil {
			ko.Status.VPCID = elem.VpcId
		} else {
			ko.Status.VPCID = nil
		}
		found = true
		break
	}
	if !found {
		return nil, ackerr.NotFound
	}

	rm.setStatusDefaults(ko)
	return &resource{ko}, nil
}

// requiredFieldsMissingFromReadManyInput returns true if there are any fields
// for the ReadMany Input shape that are required but not present in the
// resource's Spec or Status
func (rm *resourceManager) requiredFieldsMissingFromReadManyInput(
	r *resource,
) bool {
	return r.ko.Spec.Name == nil

}

// newListRequestPayload returns SDK-specific struct for the HTTP request
// payload of the List API call for the resource
func (rm *resourceManager) newListRequestPayload(
	r *resource,
) (*svcsdk.DescribeDBSecurityGroupsInput, error) {
	res := &svcsdk.DescribeDBSecurityGroupsInput{}

	if r.ko.Spec.Name != nil {
		res.SetDBSecurityGroupName(*r.ko.Spec.Name)
	}

	return res, nil
}

// sdkCreate creates the supplied resource in the backend AWS service API and
// returns a copy of the resource with resource fields (in both Spec and
// Status) filled in with values from the CREATE API operation's Output shape.
func (rm *resourceManager) sdkCreate(
	ctx context.Context,
	desired *resource,
) (created *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkCreate")
	defer exit(err)
	input, err := rm.newCreateRequestPayload(ctx, desired)
	if err != nil {
		return nil, err
	}

	var resp *svcsdk.CreateDBSecurityGroupOutput
	_ = resp
	resp, err = rm.sdkapi.CreateDBSecurityGroupWithContext(ctx, input)
	rm.metrics.RecordAPICall("CREATE", "CreateDBSecurityGroup", err)
	if err != nil {
		return nil, err
	}
	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := desired.ko.DeepCopy()

	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	if resp.DBSecurityGroup.DBSecurityGroupArn != nil {
		arn := ackv1alpha1.AWSResourceName(*resp.DBSecurityGroup.DBSecurityGroupArn)
		ko.Status.ACKResourceMetadata.ARN = &arn
	}
	if resp.DBSecurityGroup.DBSecurityGroupDescription != nil {
		ko.Spec.Description = resp.DBSecurityGroup.DBSecurityGroupDescription
	} else {
		ko.Spec.Description = nil
	}
	if resp.DBSecurityGroup.DBSecurityGroupName != nil {
		ko.Spec.Name = resp.DBSecurityGroup.DBSecurityGroupName
	} else {
		ko.Spec.Name = nil
	}
	if resp.DBSecurityGroup.EC2SecurityGroups != nil {
		f3 := []*svcapitypes.EC2SecurityGroup{}
		for _, f3iter := range resp.DBSecurityGroup.EC2SecurityGroups {
			f3elem := &svcapitypes.EC2SecurityGroup{}
			if f3iter.EC2SecurityGroupId != nil {
				f3elem.EC2SecurityGroupID = f3iter.EC2SecurityGroupId
			}
			if f3iter.EC2SecurityGroupName != nil {
				f3elem.EC2SecurityGroupName = f3iter.EC2SecurityGroupName
			}
			if f3iter.EC2SecurityGroupOwnerId != nil {
				f3elem.EC2SecurityGroupOwnerID = f3iter.EC2SecurityGroupOwnerId
			}
			if f3iter.Status != nil {
				f3elem.Status = f3iter.Status
			}
			f3 = append(f3, f3elem)
		}
		ko.Status.EC2SecurityGroups = f3
	} else {
		ko.Status.EC2SecurityGroups = nil
	}
	if resp.DBSecurityGroup.IPRanges != nil {
		f4 := []*svcapitypes.IPRange{}
		for _, f4iter := range resp.DBSecurityGroup.IPRanges {
			f4elem := &svcapitypes.IPRange{}
			if f4iter.CIDRIP != nil {
				f4elem.CIDRIP = f4iter.CIDRIP
			}
			if f4iter.Status != nil {
				f4elem.Status = f4iter.Status
			}
			f4 = append(f4, f4elem)
		}
		ko.Status.IPRanges = f4
	} else {
		ko.Status.IPRanges = nil
	}
	if resp.DBSecurityGroup.OwnerId != nil {
		ko.Status.OwnerID = resp.DBSecurityGroup.OwnerId
	} else {
		ko.Status.OwnerID = nil
	}
	if resp.DBSecurityGroup.VpcId != nil {
		ko.Status.VPCID = resp.DBSecurityGroup.VpcId
	} else {
		ko.Status.VPCID = nil
	}

	rm.setStatusDefaults(ko)
	return &resource{ko}, nil
}

// newCreateRequestPayload returns an SDK-specific struct for the HTTP request
// payload of the Create API call for the resource
func (rm *resourceManager) newCreateRequestPayload(
	ctx context.Context,
	r *resource,
) (*svcsdk.CreateDBSecurityGroupInput, error) {
	res := &svcsdk.CreateDBSecurityGroupInput{}

	if r.ko.Spec.Description != nil {
		res.SetDBSecurityGroupDescription(*r.ko.Spec.Description)
	}
	if r.ko.Spec.Name != nil {
		res.SetDBSecurityGroupName(*r.ko.Spec.Name)
	}
	if r.ko.Spec.Tags != nil {
		f2 := []*svcsdk.Tag{}
		for _, f2iter := range r.ko.Spec.Tags {
			f2elem := &svcsdk.Tag{}
			if f2iter.Key != nil {
				f2elem.SetKey(*f2iter.Key)
			}
			if f2iter.Value != nil {
				f2elem.SetValue(*f2iter.Value)
			}
			f2 = append(f2, f2elem)
		}
		res.SetTags(f2)
	}

	return res, nil
}

// sdkUpdate patches the supplied resource in the backend AWS service API and
// returns a new resource with updated fields.
func (rm *resourceManager) sdkUpdate(
	ctx context.Context,
	desired *resource,
	latest *resource,
	delta *ackcompare.Delta,
) (*resource, error) {
	// TODO(jaypipes): Figure this out...
	return nil, ackerr.NotImplemented
}

// sdkDelete deletes the supplied resource in the backend AWS service API
func (rm *resourceManager) sdkDelete(
	ctx context.Context,
	r *resource,
) (latest *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkDelete")
	defer exit(err)
	input, err := rm.newDeleteRequestPayload(r)
	if err != nil {
		return nil, err
	}
	var resp *svcsdk.DeleteDBSecurityGroupOutput
	_ = resp
	resp, err = rm.sdkapi.DeleteDBSecurityGroupWithContext(ctx, input)
	rm.metrics.RecordAPICall("DELETE", "DeleteDBSecurityGroup", err)
	return nil, err
}

// newDeleteRequestPayload returns an SDK-specific struct for the HTTP request
// payload of the Delete API call for the resource
func (rm *resourceManager) newDeleteRequestPayload(
	r *resource,
) (*svcsdk.DeleteDBSecurityGroupInput, error) {
	res := &svcsdk.DeleteDBSecurityGroupInput{}

	if r.ko.Spec.Name != nil {
		res.SetDBSecurityGroupName(*r.ko.Spec.Name)
	}

	return res, nil
}

// setStatusDefaults sets default properties into supplied custom resource
func (rm *resourceManager) setStatusDefaults(
	ko *svcapitypes.DBSecurityGroup,
) {
	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	if ko.Status.ACKResourceMetadata.OwnerAccountID == nil {
		ko.Status.ACKResourceMetadata.OwnerAccountID = &rm.awsAccountID
	}
	if ko.Status.Conditions == nil {
		ko.Status.Conditions = []*ackv1alpha1.Condition{}
	}
}

// updateConditions returns updated resource, true; if conditions were updated
// else it returns nil, false
func (rm *resourceManager) updateConditions(
	r *resource,
	onSuccess bool,
	err error,
) (*resource, bool) {
	ko := r.ko.DeepCopy()
	rm.setStatusDefaults(ko)

	// Terminal condition
	var terminalCondition *ackv1alpha1.Condition = nil
	var recoverableCondition *ackv1alpha1.Condition = nil
	var syncCondition *ackv1alpha1.Condition = nil
	for _, condition := range ko.Status.Conditions {
		if condition.Type == ackv1alpha1.ConditionTypeTerminal {
			terminalCondition = condition
		}
		if condition.Type == ackv1alpha1.ConditionTypeRecoverable {
			recoverableCondition = condition
		}
		if condition.Type == ackv1alpha1.ConditionTypeResourceSynced {
			syncCondition = condition
		}
	}

	if rm.terminalAWSError(err) || err == ackerr.SecretTypeNotSupported || err == ackerr.SecretNotFound {
		if terminalCondition == nil {
			terminalCondition = &ackv1alpha1.Condition{
				Type: ackv1alpha1.ConditionTypeTerminal,
			}
			ko.Status.Conditions = append(ko.Status.Conditions, terminalCondition)
		}
		var errorMessage = ""
		if err == ackerr.SecretTypeNotSupported || err == ackerr.SecretNotFound {
			errorMessage = err.Error()
		} else {
			awsErr, _ := ackerr.AWSError(err)
			errorMessage = awsErr.Error()
		}
		terminalCondition.Status = corev1.ConditionTrue
		terminalCondition.Message = &errorMessage
	} else {
		// Clear the terminal condition if no longer present
		if terminalCondition != nil {
			terminalCondition.Status = corev1.ConditionFalse
			terminalCondition.Message = nil
		}
		// Handling Recoverable Conditions
		if err != nil {
			if recoverableCondition == nil {
				// Add a new Condition containing a non-terminal error
				recoverableCondition = &ackv1alpha1.Condition{
					Type: ackv1alpha1.ConditionTypeRecoverable,
				}
				ko.Status.Conditions = append(ko.Status.Conditions, recoverableCondition)
			}
			recoverableCondition.Status = corev1.ConditionTrue
			awsErr, _ := ackerr.AWSError(err)
			errorMessage := err.Error()
			if awsErr != nil {
				errorMessage = awsErr.Error()
			}
			recoverableCondition.Message = &errorMessage
		} else if recoverableCondition != nil {
			recoverableCondition.Status = corev1.ConditionFalse
			recoverableCondition.Message = nil
		}
	}
	// Required to avoid the "declared but not used" error in the default case
	_ = syncCondition
	if terminalCondition != nil || recoverableCondition != nil || syncCondition != nil {
		return &resource{ko}, true // updated
	}
	return nil, false // not updated
}

// terminalAWSError returns awserr, true; if the supplied error is an aws Error type
// and if the exception indicates that it is a Terminal exception
// 'Terminal' exception are specified in generator configuration
func (rm *resourceManager) terminalAWSError(err error) bool {
	// No terminal_errors specified for this resource in generator config
	return false
}
