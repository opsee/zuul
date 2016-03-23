package router

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type Router struct {
	ec2Client  *ec2.EC2
	instanceID string
	table      string
}

func NewRouter(table string) (*Router, error) {
	creds := credentials.NewChainCredentials(
		[]credentials.Provider{
			&credentials.EnvProvider{},
			&ec2rolecreds.EC2RoleProvider{ExpiryWindow: 5 * time.Minute},
		})
	metadataClient := ec2metadata.New(nil)
	instanceID, err := metadataClient.GetMetadata("instance-id")
	if err != nil {
		return nil, err
	}

	region, err := metadataClient.Region()
	if err != nil {
		return nil, err
	}
	config := &aws.Config{Credentials: creds, Region: aws.String(region)}

	router := &Router{
		ec2Client:  ec2.New(session.New(config)),
		instanceID: instanceID,
		table:      table,
	}
	return router, nil
}

func (r *Router) getInterfaceID() (string, error) {
	input := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{
			aws.String(r.instanceID),
		},
	}

	out, err := r.ec2Client.DescribeInstances(input)
	if err != nil {
		log.Println("Error performing DescribeInstances request.")
		return "", err
	}

	reservation := out.Reservations[0]
	if reservation == nil {
		return "", fmt.Errorf("No reservations for instance ID: %s", r.instanceID)
	}

	instance := reservation.Instances[0]
	if instance == nil {
		return "", fmt.Errorf("No instance found: %s", r.instanceID)
	}

	iface := instance.NetworkInterfaces[0]
	if iface == nil {
		return "", fmt.Errorf("No network interface found for instance: %s", r.instanceID)
	}

	ifaceID := iface.NetworkInterfaceId
	if iface == nil {
		return "", fmt.Errorf("No id found on interface for instance: %s", r.instanceID)
	}

	return *ifaceID, nil
}

func (r *Router) Announce(subnet string) error {
	netIface, err := r.getInterfaceID()
	if err != nil {
		return err
	}

	createInput := &ec2.CreateRouteInput{
		DestinationCidrBlock: aws.String(subnet),
		RouteTableId:         aws.String(r.table),
		NetworkInterfaceId:   aws.String(netIface),
	}
	if _, err := r.ec2Client.CreateRoute(createInput); err != nil {
		log.Println("Error performing CreateRoute request.")
		return err
	}

	return nil
}

func (r *Router) Delete(subnet string) error {
	deleteInput := &ec2.DeleteRouteInput{
		DestinationCidrBlock: aws.String(subnet),
		RouteTableId:         aws.String(r.table),
	}
	if _, err := r.ec2Client.DeleteRoute(deleteInput); err != nil {
		log.Println("Error performing DeleteRoute request.")
		return err
	}

	return nil
}
