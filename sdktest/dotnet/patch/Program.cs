using System;
using System.Threading.Tasks;
using Azure;
using Azure.Core;
using Azure.Identity;
using Azure.ResourceManager.ContainerServiceFleet.Models;
using Azure.ResourceManager.Resources;
using Azure.ResourceManager;
using Azure.ResourceManager.ContainerService;
using Azure.ResourceManager.ContainerServiceFleet;
using Azure.ResourceManager.Models;
using Azure.Core.Diagnostics;
using System.Diagnostics.Tracing;
// using Microsoft.Azure.Management.Monitor.Models;

using AzureEventSourceListener listener = AzureEventSourceListener.CreateConsoleLogger(EventLevel.Verbose);
// using AzureEventSourceListener listener = AzureEventSourceListener.CreateTraceLogger();

// See https://aka.ms/new-console-template for more information
Console.WriteLine("Hello, World!");



// Generated from example definition: specification/containerservice/resource-manager/Microsoft.ContainerService/fleet/stable/2023-10-15/examples/Fleets_PatchTags.json
// this example is just showing the usage of "Fleets_Update" operation, for the dependent resources, they will have to be created separately.

// get your azure access token, for more details of how Azure SDK get your access token, please refer to https://learn.microsoft.com/en-us/dotnet/azure/sdk/authentication?tabs=command-line
TokenCredential cred = new DefaultAzureCredential();
// authenticate your client
ArmClient client = new ArmClient(cred);

// this example assumes you already have this ContainerServiceFleetResource created on azure
// for more information of creating ContainerServiceFleetResource, please refer to the document of ContainerServiceFleetResource
string subscriptionId = "26fe00f8-9173-4872-9134-bb1d2e00343a";
string resourceGroupName = "minsha";
string fleetName = "fleet";
ResourceIdentifier containerServiceFleetResourceId = ContainerServiceFleetResource.CreateResourceIdentifier(subscriptionId, resourceGroupName, fleetName);
ContainerServiceFleetResource containerServiceFleet = client.GetContainerServiceFleetResource(containerServiceFleetResourceId);

// invoke the operation
ContainerServiceFleetPatch patch = new ContainerServiceFleetPatch()
{
    Identity = new ManagedServiceIdentity(ManagedServiceIdentityType.SystemAssigned)
};
ArmOperation<ContainerServiceFleetResource> lro = await containerServiceFleet.UpdateAsync(WaitUntil.Completed, patch, ifMatch: null);
ContainerServiceFleetResource result = lro.Value;

// the variable result is a resource, you could call other operations on this instance as well
// but just for demo, we get its data from this resource instance
ContainerServiceFleetData resourceData = result.Data;
// for demo we just print out the id
Console.WriteLine($"Succeeded on id: {resourceData.Id}");