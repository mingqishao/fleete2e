from azure.identity import DefaultAzureCredential
from azure.mgmt.containerservicefleet import ContainerServiceFleetMgmtClient
from azure.mgmt.containerservicefleet.v2023_08_15_preview.models import FleetPatch,ManagedServiceIdentity
import os
import logging

logging.basicConfig(level=logging.DEBUG)


sub_id = "26fe00f8-9173-4872-9134-bb1d2e00343a"
client = ContainerServiceFleetMgmtClient(credential=DefaultAzureCredential(), subscription_id=sub_id, api_version="2023-08-15-preview")
fleet = client.fleets
# res = fleet.get("minsha", "fleet2")
patch = FleetPatch(identity=ManagedServiceIdentity(type="SystemAssigned"))
# print(res)
print(patch)

try:
    poller = fleet.begin_update("minsha", "fleet2", properties=patch)
    
    response = poller.result() # result() returns the object of the corresponding resource
    print(response)
except Exception as e:
    print(e)
