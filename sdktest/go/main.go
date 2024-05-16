package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/log"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservicefleet/armcontainerservicefleet"
)

func main() {
	os.Setenv("AZURE_SDK_GO_LOGGING", "all")
	log.SetEvents(log.EventRequest, log.EventResponse, log.EventRetryPolicy, log.EventLRO)
	log.SetListener(func(cls log.Event, msg string) {
		// simple console logger, it writes to stderr in the following format:
		// [time-stamp] Event: message
		fmt.Fprintf(os.Stderr, "[%s] %s: %s\n", time.Now().Format(time.StampMicro), cls, msg)
	})
	ctx := context.Background()
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	client, err := armcontainerservicefleet.NewFleetsClient("26fe00f8-9173-4872-9134-bb1d2e00343a", cred, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("start patching")
	poller, err := client.BeginUpdate(ctx, "minsha", "fleet2", armcontainerservicefleet.FleetPatch{
		Identity: &armcontainerservicefleet.ManagedServiceIdentity{
			Type: to.Ptr(armcontainerservicefleet.ManagedServiceIdentityTypeSystemAssigned),
		},
	}, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	res, err := poller.PollUntilDone(ctx, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
}
