"""
With this script, the artifacts created by goreleaser are pushed to the configured TFC private provider registry
NOTE: This script does not create the private provider registry, only uploads new versions to an existing one
see: https://developer.hashicorp.com/terraform/cloud-docs/registry/publish-providers#create-the-provider for details
"""
import json
import os
import requests
import glob
import logging

logging.basicConfig(
    format="%(asctime)s - %(message)s", datefmt="%d-%b-%y %H:%M:%S", level=logging.INFO
)

tfc_org_url = "https://app.terraform.io/api/v2/organizations"
service = "private"
provider = "awx"

namespace = os.getenv("TFC_PROVIDER_NAMESPACE")
org_name = os.getenv("TFC_ORG_NAME")
token = os.getenv("TFC_TOKEN")
gpg_key_id = os.getenv("TFC_GPG_KEY_ID")
github_ref = os.getenv("GITHUB_REF_NAME")

new_version = github_ref.replace("v", "")

os.chdir("dist/")

version_url = "{}/{}/registry-providers/{}/{}/{}/versions".format(
    tfc_org_url, org_name, service, namespace, provider
)

logging.info(f"VERSION URL: {version_url}")

headers = {
    "Authorization": "Bearer " + token,
    "Content-Type": "application/vnd.api+json",
}

resp = requests.get(version_url, headers=headers)

if resp.status_code != 200:
    raise ValueError("Failed to download versions file")

data = json.loads(resp.content)

"""
Create the endpoint for the new provider version
"""
upload_versions_dict = {
    "data": {
        "type": "registry-provider-versions",
        "attributes": {
            "version": new_version,
            "key-id": gpg_key_id,
            "protocols": ["5.0"],
        },
    }
}

json_payload = json.dumps(upload_versions_dict)

version_create_response = requests.post(version_url, data=json_payload, headers=headers)

if version_create_response.status_code != 201:
    raise ValueError(
        f"Unable to create new Provider version: received {version_create_response.status_code}\n{json.dumps(json.loads(version_create_response.content), indent=4)}"
    )

version_create_response_dict = json.loads(version_create_response.content)

logging.info(
    f"TFC PROVIDER VERSION CREATION: {json.dumps(version_create_response_dict, indent=4)}"
)

shasums_upload_url = version_create_response_dict["data"]["links"]["shasums-upload"]
shasums_sig_upload_url = version_create_response_dict["data"]["links"][
    "shasums-sig-upload"
]

shasums_file = "terraform-provider-{}_{}_SHA256SUMS".format(provider, new_version)
shasums_sig_file = "terraform-provider-{}_{}_SHA256SUMS.sig".format(
    provider, new_version
)

shasum_file_url_dict = {}
shasum_file_url_dict[shasums_file] = shasums_upload_url
shasum_file_url_dict[shasums_sig_file] = shasums_sig_upload_url

"""
Upload the SHASUM and SHASUM.sig files
"""
for sha_file, sha_url in shasum_file_url_dict.items():
    with open(sha_file, "rb") as ssum:
        shasum_data = ssum.read()

    sha_upload_response = requests.put(sha_url, data=shasum_data)

    if sha_upload_response.status_code != 200:
        raise ValueError(
            f"Unable to upload {sha_file} to {sha_url}: received {sha_upload_response.status_code}\n{json.dumps(json.loads(sha_upload_response.content), indent=4)}"
        )

    logging.info(f"SHASUM UPLOAD: {sha_file} -- {sha_upload_response.status_code}")

upload_platform_dict = {
    "data": {"type": "registry-provider-version-platforms", "attributes": {}}
}
shasum_dict = {}

with open(glob.glob("*SHA256SUMS")[0], "r") as ssum:
    for line in ssum:
        shasum_dict[line.split()[1]] = line.split()[0]

logging.info(f"ARTIFACT SHASUM MAP: {json.dumps(shasum_dict, indent=4)}")

"""
Upload the platform artifacts into the Registry
"""
for artifact in glob.glob("*.zip"):
    attr_dict = upload_platform_dict["data"]["attributes"]
    os_type = artifact.split("_")[2]
    arch = os.path.splitext(artifact)[0].split("_")[3]
    shasum = shasum_dict[artifact]

    attr_dict["os"] = os_type
    attr_dict["arch"] = arch
    attr_dict["shasum"] = shasum
    attr_dict["filename"] = artifact

    logging.info(
        f"TFC PROVIDER PLATFORM MAP: {json.dumps(upload_platform_dict, indent=4)}"
    )

    json_payload = json.dumps(upload_platform_dict)
    platform_url = "{}/{}/platforms".format(version_url, new_version)

    platform_upload_response = requests.post(
        platform_url, data=json_payload, headers=headers
    )

    if platform_upload_response.status_code != 201:
        raise ValueError(
            f"Unable to create platforms: received {platform_upload_response.status_code}\n{json.dumps(json.loads(platform_upload_response.content), indent=4)}"
        )

    platform_upload_response_dict = json.loads(platform_upload_response.content)

    logging.info(
        f"TFC PROVIDER PLATFORM CREATION: {json.dumps(platform_upload_response_dict, indent=4)}"
    )

    artifact_binary_upload_url = platform_upload_response_dict["data"]["links"][
        "provider-binary-upload"
    ]

    logging.info(f"Uploading {artifact} to {artifact_binary_upload_url}")

    with open(artifact, "rb") as bin_fobj:
        up_bin_data = bin_fobj.read()

    artifact_upload_response = requests.put(
        artifact_binary_upload_url, data=up_bin_data
    )

    if artifact_upload_response.status_code != 200:
        raise ValueError(
            f"Unable to upload {artifact} to designated Upload URL: received {artifact_upload_response.status_code}\n{json.dumps(json.loads(artifact_upload_response.content), indent=4)}"
        )
