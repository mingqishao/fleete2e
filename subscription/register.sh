#!/bin/bash

body=$(cat <<-END
{
   "properties":{
      "accountOwner":{
         "email":"shaomq@gmail.com",
         "puid":"0003BFFDE58B5244"
      },
      "additionalProperties":{
         "billingProperties":{
            "additionalStateInformation":{
               "blockNewResourceCreation":{
                  "value":false
               },
               "releaseNonDataRetentionResource":{
                  "value":false
               }
            },
            "billingAccount":{
               "id":"/providers/Microsoft.Billing/billingAccounts/65fad9d7-e54c-4a55-b0d5-e2b16b461b11"
            },
            "billingType":"Legacy",
            "channelType":"CustomerLed",
            "paymentType":"Paid",
            "tier":"Standard",
            "workloadType":"Production"
         },
         "resourceProviderProperties":"{\"resourceProviderNamespace\":\"Microsoft.ContainerService\"}"
      },
      "locationPlacementId":"Public_2014-09-01",
      "managedByTenants":[
         
      ],
      "quotaId":"PayAsYouGo_2014-09-01",
      "registeredFeatures":[
         
      ],
      "tenantId":"87c63cb4-bf39-43bb-ad51-f93648332628"
   },
   "registrationDate":"Mon, 08 Apr 2024 21:07:49 GMT",
   "state":"Registered"
}
END
)

echo $body

#  -H "x-ms-identity-url: http://msi-simulator.msi-simulator.svc.cluster.local" \

curl -v  -X PUT \
 -H "Content-Type: application/json" \
 -H "x-ms-home-tenant-id: 87c63cb4-bf39-43bb-ad51-f93648332628" \
 --data "$body" \
 "http://localhost:8080/subscriptions/c0eb2ba0-d276-4dcb-a083-fb13c5eeb0be?api-version=2.0"
