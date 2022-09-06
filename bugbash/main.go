package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"bugbash/v20220702preview"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
)

var (
	subscription = "3959ec86-5353-4b0c-b5d7-3877122861a0"
	cred         *azidentity.DefaultAzureCredential
	rgClient     *armresources.ResourceGroupsClient
	locations    = []string{
		"australiaeast",
		"brazilsouth",
		"canadacentral",
		"centralindia",
		"centralus",
		"eastasia",
		// "eastus",
		"eastus2",
		"francecentral",
		"germanywestcentral",
		"japaneast",
		"koreacentral",
		"northeurope",
		"norwayeast",
		"southeastasia",
		"southafricanorth",
		"southcentralus",
		"swedencentral",
		"switzerlandnorth",
		"uksouth",
		"westeurope",
		"westus2",
		"westus3",
	}
	isPrivates     = []bool{true, false}
	skuTires       = []string{"Free", "Paid"}
	msiTypes       = []string{"SystemAssigned", "UserAssigned"}
	aadAuthNs      = []bool{true, false}
	aadAuthRBACs   = []bool{true, false}
	networkPlugins = []string{"azure", "kubenet"}
)

type aksOpts struct {
	location      string
	isPrivate     bool
	skuTire       string // Free / Paid
	msiType       string // SystemAssigned / UserAssigned
	userIdentity  string
	aadAuthN      bool
	aadAuthRBAC   bool
	networkPlugin string // azure / kubenet
	addressSpace  int
}

func (o aksOpts) toCSV(r []string) []string {
	r = append(r, o.location)
	if o.isPrivate {
		r = append(r, "Private")
	} else {
		r = append(r, "Public")
	}
	r = append(r, o.skuTire)
	r = append(r, o.msiType)
	if o.aadAuthN {
		r = append(r, "AAD Enable")
	} else {
		r = append(r, "Not AAD Enable")
	}
	if o.aadAuthRBAC {
		r = append(r, "AAD RBAC")
	} else {
		r = append(r, "Not AAD RBAC")
	}
	r = append(r, o.networkPlugin)
	if o.networkPlugin == "azure" {
		r = append(r, fmt.Sprintf("10.240.%d.0/24", o.addressSpace))
	} else {
		r = append(r, "")
	}
	return r
}

func isNotFound(err error) bool {
	var resErr *azcore.ResponseError
	return errors.As(err, &resErr) && resErr.StatusCode == http.StatusNotFound
}

func random(n int) int {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	return r1.Intn(n)
}

func randomAKSOpts() *aksOpts {
	r := aksOpts{
		location:      locations[random(len(locations))],
		isPrivate:     isPrivates[random(len(isPrivates))],
		skuTire:       skuTires[random(len(skuTires))],
		msiType:       msiTypes[random(len(msiTypes))],
		aadAuthN:      aadAuthNs[random(len(aadAuthNs))],
		aadAuthRBAC:   aadAuthRBACs[random(len(aadAuthRBACs))],
		networkPlugin: networkPlugins[random(len(networkPlugins))],
		addressSpace:  random(255),
	}
	if r.msiType == "UserAssigned" {
		r.userIdentity = ""
	}
	if !r.aadAuthN {
		r.aadAuthRBAC = false
	}
	return &r
}

func createVNET(rg, name, location, addressPrefix string) (*armnetwork.VirtualNetwork, error) {
	fmt.Println("creating VNET", name)
	ctx := context.Background()
	client, err := armnetwork.NewVirtualNetworksClient(subscription, cred, nil)
	if err != nil {
		return nil, err
	}
	res, err := client.Get(ctx, rg, name, nil)
	if err != nil && !isNotFound(err) {
		return nil, err
	}
	if err == nil {
		return &res.VirtualNetwork, nil
	}

	p, err := client.BeginCreateOrUpdate(ctx, rg, name, armnetwork.VirtualNetwork{
		Location: &location,
		Properties: &armnetwork.VirtualNetworkPropertiesFormat{
			AddressSpace: &armnetwork.AddressSpace{
				AddressPrefixes: []*string{&addressPrefix},
			},
		},
	}, nil)
	if err != nil {
		return nil, err
	}
	createRes, err := p.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &createRes.VirtualNetwork, nil
}

func createSubnet(rg, vnet, name, addressPrefix string) (string, error) {
	fmt.Println("creating vnet/subnet", vnet, "/", name)
	ctx := context.Background()
	client, err := armnetwork.NewVirtualNetworksClient(subscription, cred, nil)
	if err != nil {
		return "", err
	}
	res, err := client.Get(ctx, rg, vnet, nil)
	if err != nil {
		return "", err
	}
	for _, subnet := range res.Properties.Subnets {
		if *subnet.Name == name {
			return *subnet.ID, nil
		}
	}
	v := res.VirtualNetwork
	v.Properties.Subnets = append(v.Properties.Subnets, &armnetwork.Subnet{
		Name: &name,
		Properties: &armnetwork.SubnetPropertiesFormat{
			AddressPrefix: &addressPrefix,
		},
	})
	fmt.Sprintln("submit request for subnet")
	p, err := client.BeginCreateOrUpdate(ctx, rg, vnet, v, nil)
	if err != nil {
		return "", err
	}
	createRes, err := p.PollUntilDone(ctx, nil)
	if err != nil {
		return "", err
	}
	for _, subnet := range createRes.Properties.Subnets {
		if *subnet.Name == name {
			return *subnet.ID, nil
		}
	}
	return "", errors.New("created subnet but can't found it from result")
}

func createMSI(rg, name, location string) (string, error) {
	fmt.Println("creating MSI", name)
	cmd := exec.Command("az", "identity", "delete", "-g", rg, "-n", name)
	if err := cmd.Run(); err != nil {
		return "", nil
	}
	cmd = exec.Command("az", "identity", "create", "-l", location, "-g", rg, "-n", name)
	if err := cmd.Run(); err != nil {
		return "", nil
	}
	msi := fmt.Sprintf("/subscriptions/3959ec86-5353-4b0c-b5d7-3877122861a0/resourcegroups/%s/providers/Microsoft.ManagedIdentity/userAssignedIdentities/%s", rg, name)
	fmt.Println("created MSI", msi)
	return msi, nil
}

func deleteAKS(rg, name string) error {
	fmt.Println("start to delete AKS", name)
	ctx := context.Background()
	client, err := armcontainerservice.NewManagedClustersClient(subscription, cred, nil)
	if err != nil {
		return err
	}
	p, err := client.BeginDelete(ctx, rg, name, nil)
	if err != nil {
		return err
	}
	_, err = p.PollUntilDone(ctx, nil)
	if err != nil {
		return err
	}
	return nil
}

func resetAKSOpts(mc *armcontainerservice.ManagedCluster, opts *aksOpts) error {
	opts.location = *mc.Location
	if mc.Properties.AADProfile != nil && mc.Properties.AADProfile.Managed != nil && *mc.Properties.AADProfile.Managed {
		opts.aadAuthN = true
	} else {
		opts.aadAuthN = false
	}
	if mc.Properties.AADProfile != nil && mc.Properties.AADProfile.EnableAzureRBAC != nil && *mc.Properties.AADProfile.EnableAzureRBAC {
		opts.aadAuthRBAC = true
	} else {
		opts.aadAuthRBAC = false
	}

	if mc.Properties.APIServerAccessProfile != nil && mc.Properties.APIServerAccessProfile.EnablePrivateCluster != nil && *mc.Properties.APIServerAccessProfile.EnablePrivateCluster {
		opts.isPrivate = true
	} else {
		opts.isPrivate = false
	}
	opts.skuTire = string(*mc.SKU.Tier)
	opts.msiType = string(*mc.Identity.Type)
	opts.networkPlugin = string(*mc.Properties.NetworkProfile.NetworkPlugin)
	if opts.networkPlugin == "azure" {
		subnetId := *mc.Properties.AgentPoolProfiles[0].VnetSubnetID
		addr, err := strconv.Atoi(strings.Split(subnetId, "subnet-")[1])
		if err != nil {
			return err
		}
		opts.addressSpace = addr
	}
	return nil
}

func createAKS(rg, name string, opts *aksOpts) (*armcontainerservice.ManagedCluster, error) {
	fmt.Println("start to create AKS", name)
	ctx := context.Background()
	client, err := armcontainerservice.NewManagedClustersClient(subscription, cred, nil)
	if err != nil {
		return nil, err
	}

	getRes, err := client.Get(ctx, rg, name, nil)
	if err != nil && !isNotFound(err) {
		return nil, err
	}
	if err == nil {
		if *getRes.Properties.ProvisioningState == "Succeeded" {
			resetAKSOpts(&getRes.ManagedCluster, opts)
			return &getRes.ManagedCluster, nil
		} else {
			if err := deleteAKS(rg, name); err != nil {
				return nil, err
			}
		}
	}

	var userIds map[string]*armcontainerservice.ManagedServiceIdentityUserAssignedIdentitiesValue = nil
	if opts.msiType == "UserAssigned" {
		userIds = map[string]*armcontainerservice.ManagedServiceIdentityUserAssignedIdentitiesValue{}
		id, err := createMSI(rg, fmt.Sprintf("%s-msi", name), opts.location)
		if err != nil {
			return nil, err
		}
		userIds[id] = &armcontainerservice.ManagedServiceIdentityUserAssignedIdentitiesValue{}
	}
	mc := armcontainerservice.ManagedCluster{
		Location: &opts.location,
		SKU: &armcontainerservice.ManagedClusterSKU{
			Name: to.Ptr(armcontainerservice.ManagedClusterSKUNameBasic),
			Tier: to.Ptr(armcontainerservice.ManagedClusterSKUTier(opts.skuTire)),
		},
		Identity: &armcontainerservice.ManagedClusterIdentity{
			Type:                   to.Ptr(armcontainerservice.ResourceIdentityType(opts.msiType)),
			UserAssignedIdentities: userIds,
		},
		Properties: &armcontainerservice.ManagedClusterProperties{
			DNSPrefix: to.Ptr(name),
			AgentPoolProfiles: []*armcontainerservice.ManagedClusterAgentPoolProfile{
				{
					Name:   to.Ptr("agentpool"),
					Count:  to.Ptr(int32(1)),
					VMSize: to.Ptr("Standard_DS3_v2"),
					OSType: to.Ptr(armcontainerservice.OSTypeLinux),
					Type:   to.Ptr(armcontainerservice.AgentPoolTypeVirtualMachineScaleSets),
					Mode:   to.Ptr(armcontainerservice.AgentPoolModeSystem),
				},
			},
			AADProfile: &armcontainerservice.ManagedClusterAADProfile{
				Managed:         to.Ptr(opts.aadAuthN),
				EnableAzureRBAC: to.Ptr(opts.aadAuthRBAC),
			},
			APIServerAccessProfile: &armcontainerservice.ManagedClusterAPIServerAccessProfile{
				EnablePrivateCluster: &opts.isPrivate,
			},
			NetworkProfile: &armcontainerservice.NetworkProfile{
				NetworkPlugin: to.Ptr(armcontainerservice.NetworkPlugin(opts.networkPlugin)),
			},
		},
	}
	if opts.networkPlugin == "azure" {
		vnetName := fmt.Sprintf("vnet-%s", opts.location)
		vnet, err := createVNET(rg, vnetName, opts.location, "10.0.0.0/8")
		if err != nil {
			return nil, err
		}
		fmt.Println("created VNET", *vnet.ID)
		subnetName := fmt.Sprintf("subnet-%d", opts.addressSpace)
		subnetAddr := fmt.Sprintf("10.240.%d.0/24", opts.addressSpace)
		subnetId, err := createSubnet(rg, vnetName, subnetName, subnetAddr)
		if err != nil {
			return nil, err
		}
		fmt.Println("created subnet", subnetId)
		mc.Properties.AgentPoolProfiles[0].VnetSubnetID = &subnetId
	}
	if !opts.aadAuthN {
		mc.Properties.AADProfile = nil
	}
	p, err := client.BeginCreateOrUpdate(ctx, rg, name, mc, nil)
	if err != nil {
		return nil, err
	}
	res, err := p.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &res.ManagedCluster, nil
}

func createFleet(rg, name, location string) (*v20220702preview.Fleet, error) {
	fmt.Println("creating fleet", name)
	ctx := context.Background()
	client, err := v20220702preview.NewFleetsClient(subscription, cred, nil)
	if err != nil {
		return nil, err
	}
	f, err := client.Get(ctx, rg, name, nil)
	if err != nil && !isNotFound(err) {
		return nil, err
	}
	if err == nil && *f.Properties.ProvisioningState == "Succeeded" {
		return &f.Fleet, nil
	}

	p, err := client.BeginCreateOrUpdate(ctx, rg, name, v20220702preview.Fleet{
		Location: &location,
		Properties: &v20220702preview.FleetProperties{
			HubProfile: &v20220702preview.FleetHubProfile{
				DNSPrefix: to.Ptr(name),
			},
		},
	}, nil)
	if err != nil {
		return nil, err
	}
	res, err := p.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &res.Fleet, nil
}

func createFleetMember(rg, fleet, memberRg, memberName string) (*v20220702preview.FleetMember, error) {
	fmt.Println("start create member", memberName)
	ctx := context.Background()
	client, err := v20220702preview.NewFleetMembersClient(subscription, cred, nil)
	if err != nil {
		return nil, err
	}
	getRes, err := client.Get(ctx, rg, fleet, memberName, nil)
	if err != nil && !isNotFound(err) {
		return nil, err
	}
	if err == nil && *getRes.Properties.ProvisioningState == "Succeeded" {
		return &getRes.FleetMember, nil
	}
	p, err := client.BeginCreateOrUpdate(ctx, rg, fleet, memberName, v20220702preview.FleetMember{
		Properties: &v20220702preview.FleetMemberProperties{
			ClusterResourceID: to.Ptr(fmt.Sprintf("/subscriptions/%s/resourcegroups/%s/providers/Microsoft.ContainerService/managedClusters/%s", subscription, memberRg, memberName)),
		},
	}, nil)
	if err != nil {
		return nil, err
	}
	res, err := p.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &res.FleetMember, nil
}

func createGroup(name, location string) (*armresources.ResourceGroup, error) {
	ctx := context.Background()
	getRes, err := rgClient.Get(ctx, name, nil)
	if err != nil && !isNotFound(err) {
		return nil, err
	}
	if err == nil && *getRes.Properties.ProvisioningState == "Succeeded" {
		return &getRes.ResourceGroup, nil
	}
	res, err := rgClient.CreateOrUpdate(context.TODO(), name, armresources.ResourceGroup{
		Location: &location,
	}, nil)
	if err != nil {
		return nil, err
	}
	return &res.ResourceGroup, nil
}

func deleteGroup(name string) error {
	p, err := rgClient.BeginDelete(context.TODO(), name, nil)
	if err != nil {
		return err
	}
	_, err = p.PollUntilDone(context.TODO(), nil)
	if err != nil {
		return err
	}
	return nil
}

func initialize() error {
	var err error
	cred, err = azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return err
	}
	rgClient, err = armresources.NewResourceGroupsClient(subscription, cred, nil)
	if err != nil {
		return err
	}
	return nil
}

func csv(a ...string) {
	fmt.Print("csv:")
	for i, v := range a {
		if i > 0 {
			fmt.Print(",")
		}
		fmt.Print(v)
	}
	fmt.Print("\n")
}

type batchAKSOptions struct {
	prefix string
	num    int
}

func BatchCreateAKS(opts *batchAKSOptions) {
	prefix := opts.prefix
	rgLocation := randomAKSOpts().location
	_, err := createGroup(prefix, rgLocation)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("created resource group", prefix)
	for i := 1; i <= opts.num; i++ {
		aksName := fmt.Sprintf("%s-aks-%d", prefix, i)
		opts := randomAKSOpts()
		_, err := createAKS(prefix, aksName, opts)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("created aks", aksName)
		a := []string{aksName}
		a = opts.toCSV(a)
		csv(a...)
	}
}

type batchFleetOptions struct {
	name           string
	memberNum      int
	isPrivate      string
	network        string
	fleetLocation  string
	memberLocation string
}

func BatchCreateFleet(opts *batchFleetOptions) {
	name := opts.name
	l := randomAKSOpts().location
	if opts.fleetLocation != "" {
		l = opts.fleetLocation
	}
	fmt.Println(l)
	_, err := createGroup(name, l)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("created resource group", name)

	_, err = createFleet(name, name, l)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("created fleet", name)
	csv(name, l)
	for i := 1; i <= opts.memberNum; i++ {
		aksName := fmt.Sprintf("aks-member-%d", i)
		aksOpts := randomAKSOpts()
		if opts.isPrivate == "yes" || opts.isPrivate == "true" {
			aksOpts.isPrivate = true
		}
		if opts.isPrivate == "no" || opts.isPrivate == "false" {
			aksOpts.isPrivate = false
		}
		if opts.network != "" {
			aksOpts.networkPlugin = opts.network
		}
		if opts.memberLocation != "" {
			aksOpts.location = opts.memberLocation
		}
		_, err := createAKS(name, aksName, aksOpts)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("created aks", aksName)

		_, err = createFleetMember(name, name, name, aksName)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("created member", aksName)
		a := []string{aksName}
		a = aksOpts.toCSV(a)
		csv(a...)
	}
}

func runCmd() {
	fleet := flag.NewFlagSet("fleet", flag.ExitOnError)
	fleetOpts := &batchFleetOptions{}
	fleet.StringVar(&fleetOpts.name, "name", "", "fleet name")
	fleet.IntVar(&fleetOpts.memberNum, "num", 10, "member number")
	fleet.StringVar(&fleetOpts.isPrivate, "private", "", "if private API Server: yes/or")
	fleet.StringVar(&fleetOpts.network, "network", "", "network plugin: azure/kubenet")
	fleet.StringVar(&fleetOpts.memberLocation, "member-location", "", "member location")
	fleet.StringVar(&fleetOpts.fleetLocation, "fleet-location", "", "fleet location")

	aks := flag.NewFlagSet("aks", flag.ExitOnError)
	aksOpts := &batchAKSOptions{}
	aks.StringVar(&aksOpts.prefix, "prefix", "", "name prefix")
	aks.IntVar(&aksOpts.num, "num", 10, "aks number")

	cmd := os.Args[1]
	if cmd == "fleet" {
		err := fleet.Parse(os.Args[2:])
		if err != nil {
			fmt.Println(err)
			return
		}
		BatchCreateFleet(fleetOpts)
	}
	if cmd == "aks" {
		err := aks.Parse(os.Args[2:])
		if err != nil {
			fmt.Println(err)
			return
		}
		BatchCreateAKS(aksOpts)
	}
}

func main() {
	err := initialize()
	if err != nil {
		fmt.Println(err)
		return
	}

	runCmd()

	// bCreateFleet()
	// csv("aaa", "bbb")

	// rg, err := createGroup("minsha-rg-2", "westus")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(rg.Name)

	// f, err := createFleet("minsha-rg-2", "minsha-fleet-1", "eastus")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(f.Name)

	// mc, err := createAKS("minsha-rg-2", "minsha-aks-vnet-6", aksOpts{
	// 	location:  "westus2",
	// 	isPrivate: false,
	// 	skuTire:   "Free",
	// 	msiType:   "SystemAssigned",
	// 	// userIdentity: "/subscriptions/3959ec86-5353-4b0c-b5d7-3877122861a0/resourcegroups/minsha-rg-2/providers/Microsoft.ManagedIdentity/userAssignedIdentities/msi",
	// 	aadAuthN:      false,
	// 	aadAuthRBAC:   false,
	// 	networkPlugin: "azure",
	// 	addressSpace:  15,
	// })
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(mc.Properties.ProvisioningState)

	// member, err := createFleetMember("minsha-rg-2", "minsha-fleet-1", "minsha-rg-2", "minsha-aks-3")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(member.Properties.ProvisioningState)

	// err = deleteGroup("minsha-rg-2")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(rg)

	// id, err := createMSI("minsha-rg-2", "msi2", "westus")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(id)
}
