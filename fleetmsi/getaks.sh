#!/bin/bash

ID="${1}"

az rest --method get --url https://management.azure.com$ID?api-version=2023-06-01
