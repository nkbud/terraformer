# Rundeck Provider

Terraformer supports importing resources from Rundeck.

## Authentication

### Environment Variables

* `RUNDECK_URL` - Rundeck server URL (required)
* `RUNDECK_TOKEN` - API token for authentication (optional, but either this or username/password is required)
* `RUNDECK_USERNAME` - Username for basic authentication (optional)
* `RUNDECK_PASSWORD` - Password for basic authentication (optional)
* `RUNDECK_API_VERSION` - API version to use (optional, defaults to "38")
* `RUNDECK_INSECURE_SSL` - Skip SSL certificate verification (optional, defaults to "false")

## Resources

### Supported Resources

- `rundeck_job` - Rundeck job definitions

## Usage

```bash
export RUNDECK_URL="https://rundeck.example.com"
export RUNDECK_TOKEN="your-api-token"

terraformer import rundeck --resources=jobs
```

Alternatively, using username/password authentication:

```bash
export RUNDECK_URL="https://rundeck.example.com"
export RUNDECK_USERNAME="your-username"
export RUNDECK_PASSWORD="your-password"

terraformer import rundeck --resources=jobs
```

### List Available Resources

```bash
terraformer import rundeck list
```

### Import Specific Jobs

```bash
terraformer import rundeck --resources=jobs --filter=job=job-id-1:job-id-2
```

### Output to Specific Directory

```bash
terraformer import rundeck --resources=jobs --path-output=./rundeck-terraform
```

## Example

```bash
export RUNDECK_URL="https://rundeck.example.com:4440"
export RUNDECK_TOKEN="your-api-token"
export RUNDECK_API_VERSION="38"

terraformer import rundeck --resources=jobs --verbose
```

This will create Terraform files for all Rundeck jobs accessible with the provided credentials.

## Provider Configuration

The generated Terraform configuration will include a `rundeck` provider block:

```hcl
provider "rundeck" {
  url         = "https://rundeck.example.com:4440"
  token       = "your-api-token"
  api_version = "38"
  insecure    = false
}
```

## Notes

- Jobs are organized by project in Rundeck, and the importer will discover all projects and their associated jobs
- Job resource names are constructed from project, group, and job name for uniqueness
- The importer requires either an API token or username/password for authentication
- The default API version is 38, which is compatible with most modern Rundeck installations