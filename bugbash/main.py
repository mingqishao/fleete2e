import os

from azure.identity import DefaultAzureCredential
from azure.mgmt.containerservice import ContainerServiceClient
from azure.mgmt.resource import ResourceManagementClient

SUBSCRIPTION_ID = "3959ec86-5353-4b0c-b5d7-3877122861a0"

def create_group(name, location):
    resource_client = ResourceManagementClient(
        credential=DefaultAzureCredential(),
        subscription_id=SUBSCRIPTION_ID
    )
    rg = resource_client.resource_groups.create_or_update(
        name,
        {"location": location}
    )
    print("successfully create RG: %s" % rg.name)

def create_aks(rg, name, location):
    containerservice_client = ContainerServiceClient(
        credential=DefaultAzureCredential(),
        subscription_id=SUBSCRIPTION_ID
    )
    managed_clusters = containerservice_client.managed_clusters.begin_create_or_update(
        rg,
        name,
        {
            "identity": {
                "type": "SystemAssigned",
            },
            "dns_prefix": name,
            "agent_pool_profiles": [
                {
                    "name": "aksagent",
                    "count": 1,
                    "vm_size": "Standard_DS2_v2",
                    "max_pods": 110,
                    "min_count": 1,
                    "max_count": 100,
                    "os_type": "Linux",
                    "type": "VirtualMachineScaleSets",
                    "enable_auto_scaling": True,
                    "mode": "System",
                }
            ],
            "location": location 
        }).result()
    print("Create managed clusters:\n{}".format(managed_clusters))

    

def main():
    SUBSCRIPTION_ID = "3959ec86-5353-4b0c-b5d7-3877122861a0"
    # CLIENT_ID = os.environ.get("CLIENT_ID", None)
    # CLIENT_SECRET = os.environ.get("CLIENT_SECRET", None)
    GROUP_NAME = "minsha-test"
    MANAGED_CLUSTERS = "managed_aks"
    AZURE_LOCATION = "eastus"

    # Create client
    # # For other authentication approaches, please see: https://pypi.org/project/azure-identity/
    resource_client = ResourceManagementClient(
        credential=DefaultAzureCredential(),
        subscription_id=SUBSCRIPTION_ID
    )
    containerservice_client = ContainerServiceClient(
        credential=DefaultAzureCredential(),
        subscription_id=SUBSCRIPTION_ID
    )
    # - init depended client -
    # - end -

    # Create resource group
    resource_client.resource_groups.create_or_update(
        GROUP_NAME,
        {"location": AZURE_LOCATION}
    )

    # - init depended resources -
    # - end -

    # Create managed clusters
    managed_clusters = containerservice_client.managed_clusters.begin_create_or_update(
        GROUP_NAME,
        MANAGED_CLUSTERS,
        {
            "dns_prefix": "akspythonsdk",
            "agent_pool_profiles": [
                {
                    "name": "aksagent",
                    "count": 1,
                    "vm_size": "Standard_DS2_v2",
                    "max_pods": 110,
                    "min_count": 1,
                    "max_count": 100,
                    "os_type": "Linux",
                    "type": "VirtualMachineScaleSets",
                    "enable_auto_scaling": True,
                    "mode": "System",
                }
            ],
            "service_principal_profile": {
                "client_id": CLIENT_ID,
                "secret": CLIENT_SECRET
            },
            "location": AZURE_LOCATION
        }
    ).result()
    print("Create managed clusters:\n{}".format(managed_clusters))

    # Get managed clusters
    managed_clusters = containerservice_client.managed_clusters.get(
        GROUP_NAME,
        MANAGED_CLUSTERS
    )
    print("Get managed clusters:\n{}".format(managed_clusters))



if __name__ == "__main__":
    # main()
    # create_group("minsha-rg", "westus")
    create_aks("minsha-rg", "minsha-aks", "westus")