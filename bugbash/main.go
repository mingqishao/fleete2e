package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os/exec"
	"time"

	"bugbash/v20220702preview"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice"
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
	isPrivates   = []bool{true, false}
	skuTires     = []string{"Free", "Paid"}
	msiTypes     = []string{"SystemAssigned", "UserAssigned"}
	aadAuthNs    = []bool{true, false}
	aadAuthRBACs = []bool{true, false}
)

type aksOpts struct {
	location     string
	isPrivate    bool
	skuTire      string // Free / Paid
	msiType      string // SystemAssigned / UserAssigned
	userIdentity string
	aadAuthN     bool
	aadAuthRBAC  bool
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

func randomAKSOpts() aksOpts {
	r := aksOpts{
		location:    locations[random(len(locations))],
		isPrivate:   isPrivates[random(len(isPrivates))],
		skuTire:     skuTires[random(len(skuTires))],
		msiType:     msiTypes[random(len(msiTypes))],
		aadAuthN:    aadAuthNs[random(len(aadAuthNs))],
		aadAuthRBAC: aadAuthRBACs[random(len(aadAuthRBACs))],
	}
	if r.msiType == "UserAssigned" {
		r.userIdentity = ""
	}
	if !r.aadAuthN {
		r.aadAuthRBAC = false
	}
	return r
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

func createAKS(rg, name string, opts aksOpts) (*armcontainerservice.ManagedCluster, error) {
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
		},
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

func bCreateFleet() {
	name := "bugbash-fleet-2"
	l := randomAKSOpts().location
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
	for i := 1; i <= 20; i++ {
		aksName := fmt.Sprintf("%s-aks-%d", name, i)
		opts := randomAKSOpts()
		_, err := createAKS(name, aksName, opts)
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
		a = opts.toCSV(a)
		csv(a...)
	}
}

func main() {
	err := initialize()
	if err != nil {
		fmt.Println(err)
		return
	}

	// bCreateFleet()
	csv("aaa", "bbb")

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

	// mc, err := createAKS("minsha-rg-2", "minsha-aks-private-5", aksOpts{
	// 	location:     "eastus",
	// 	isPrivate:    true,
	// 	skuTire:      "Free",
	// 	msiType:      "UserAssigned",
	// 	userIdentity: "/subscriptions/3959ec86-5353-4b0c-b5d7-3877122861a0/resourcegroups/minsha-rg-2/providers/Microsoft.ManagedIdentity/userAssignedIdentities/msi",
	// 	aadAuthN:     false,
	// 	aadAuthRBAC:  false,
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
