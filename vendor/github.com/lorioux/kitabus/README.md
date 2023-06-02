# Kitabus 

A tool to automate the google cloud resources export into terraform (infrastructure as code) replicating, side-by-side,  the GCloud resource hierarchy. Practitioners should leverage this tool for an accelerated, user-friendly resource configuration export, and repeatable infrastructure updating. Users who are unfamiliar to GCloud should start by navigating the GCloud console to provision resources, in user-friendly manner, then export all resources into terraform.

## How it works (step-by-step)
```bash
# Generate the resources catalogue from the org level
go run . -generate catalogue --scope "organizations/$ORG_ID" --kind ".*Project" --path "."
```
```bash
# Generate the organization resource hierarchy
go run . -generate orgtree --path "."
```
```bash
# Generate the resources' TF templates and import resource
go run . -generate modules --path "."
```

