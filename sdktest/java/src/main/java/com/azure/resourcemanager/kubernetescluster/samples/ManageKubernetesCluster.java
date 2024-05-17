// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package com.azure.resourcemanager.kubernetescluster.samples;

import com.azure.core.credential.TokenCredential;
import com.azure.core.http.policy.HttpLogDetailLevel;
import com.azure.core.http.policy.HttpLogOptions;
import com.azure.core.management.AzureEnvironment;
import com.azure.identity.DefaultAzureCredentialBuilder;
import com.azure.resourcemanager.AzureResourceManager;
import com.azure.resourcemanager.containerservice.models.AgentPoolMode;
import com.azure.resourcemanager.containerservice.models.ContainerServiceVMSizeTypes;
import com.azure.resourcemanager.containerservice.models.KubernetesCluster;
import com.azure.resourcemanager.containerservicefleet.ContainerServiceFleetManager;
import com.azure.core.management.Region;
import com.azure.core.management.profile.AzureProfile;
import com.azure.resourcemanager.samples.Utils;
import com.azure.resourcemanager.containerservicefleet.models.Fleet;
import com.azure.resourcemanager.containerservicefleet.models.ManagedServiceIdentity;
import com.azure.resourcemanager.containerservicefleet.models.ManagedServiceIdentityType;

import java.util.Date;

/**
 * Azure Container Service (AKS) sample for managing a Kubernetes cluster.
 * - Create an Azure Container Service (AKS) with managed Kubernetes cluster
 * - Create a SSH private/public key
 * - Update the number of agent virtual machines in the Kubernetes cluster
 */
public class ManageKubernetesCluster {

    /**
     * Main entry point.
     *
     * @param args the parameters
     */
    public static void main(String[] args) {
        try {
            // =============================================================
            // Authenticate
            System.setProperty("AZURE_LOG_LEVEL", "VERBOSE");
            System.setProperty("AZURE_HTTP_LOG_DETAIL_LEVEL", "BASIC");
            final AzureProfile profile = new AzureProfile("72f988bf-86f1-41af-91ab-2d7cd011db47",
                    "26fe00f8-9173-4872-9134-bb1d2e00343a", AzureEnvironment.AZURE);
            final TokenCredential credential = new DefaultAzureCredentialBuilder()
                    .authorityHost(profile.getEnvironment().getActiveDirectoryEndpoint())
                    .build();

            HttpLogOptions options = new HttpLogOptions();
            options.setLogLevel(HttpLogDetailLevel.BODY_AND_HEADERS);
            options.setPrettyPrintBody(true);
            ContainerServiceFleetManager manager = ContainerServiceFleetManager.configure().withLogOptions(options)
                    .authenticate(credential, profile);
            Fleet fleet = manager.fleets()
                    .getByResourceGroupWithResponse("minsha", "fleet5", com.azure.core.util.Context.NONE).getValue();
            // manager.fleets()
            System.out.println(fleet);

            ManagedServiceIdentity identity = new ManagedServiceIdentity();
            identity.withType(ManagedServiceIdentityType.SYSTEM_ASSIGNED);
            Fleet.Update update = fleet.update().withIdentity(identity);
            update.apply();
        } catch (Exception e) {
            System.out.println(e.getMessage());
            e.printStackTrace();
        }
    }
}
