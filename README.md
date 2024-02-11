# LemmeKnow
[![Backend](https://github.com/cse403-lemmeknow/lemmeknow/actions/workflows/backend.yml/badge.svg)](https://github.com/cse403-lemmeknow/lemmeknow/actions/workflows/backend.yml)
[![Frontend](https://github.com/cse403-lemmeknow/lemmeknow/actions/workflows/frontend.yml/badge.svg)](https://github.com/cse403-lemmeknow/lemmeknow/actions/workflows/frontend.yml)
[![Terraform](https://github.com/cse403-lemmeknow/lemmeknow/actions/workflows/terraform.yml/badge.svg)](https://github.com/cse403-lemmeknow/lemmeknow/actions/workflows/terraform.yml)
[![Living document](https://img.shields.io/badge/Google_Docs-Living_document-green)](https://docs.google.com/document/d/1d1Dfsa-rxboUDKKB_DPqz7EQ5tSm0tDeaKtedVD5LeU/edit?usp=sharing)

LemmeKnow is a group activity planning platform with an integrated calendar to show availability and an integrated polling system to help make decisions. It offers a centralized dashboard of who is attending an event, their status, and that of their assigned tasks.

## Goals
- Interactive web-based dashboard
- Calendar, shared between members, showing activities and time conflicts
- Polls, for reaching a consensus based on individual preferences
- Messaging system, for general discussion and coordination
- Reminders, to increase the chance of attendance

## Layout

- [`frontend/`](./frontend/) contains the frontend website
- [`backend/`](./backend/) contains the backend serverless function 
- [`terraform/`](./terraform/) contains infrastructure definitions

## Prerequisites

1. Use a Linux or MacOS environment. Alternatively, be prepared for additional troubleshooting beyond what is described below.
2. Install `go` v1.21 or higher ([instructions](https://go.dev/doc/install)).
3. Install `node` v18 or higher ([instructions](https://nodejs.org/en/learn/getting-started/how-to-install-nodejs)).
4. Install `make` ([instructions](https://www.gnu.org/software/make/)). Alternatively, run the corresponding command(s) from each `Makefile`.
5. Install `zip` ([Ubuntu instructions](https://www.mysoftkey.com/linux/how-to-do-zip-and-unzip-file-in-ubuntu-linux/), [MacOS instructions](https://formulae.brew.sh/formula/zip)). Alternatively, manually create any applicable zip archives.
6. Optionally, for containerized testing, install Docker and then `act` ([instructions](https://nektosact.com/installation/index.html)).
7. Optionally, for cloud deployment, install and configure AWS CLI and then Terraform ([instructions](./terraform/README.md#prerequisites)).

## Running

The first entry in each component's `Makefile` will run it locally, so simply type `make` in [`frontend/`](./frontend/) and/or [`backend/`](./backend/).

By default, the fronted will host http://localhost:5173/ and the backend will host http://localhost:8080. If you run both, the backend will reverse-proxy the frontend, so navigate to the latter URL only.

## Building

Each `Makefile` includes a `build` step, so simply type `make build` in [`frontend/`](`frontend/`) and/or [`backend/`](./backend/). The frontend will build to `frontend/build/*` and the backend will build to `backend/bin/bootstrap.zip`

## Testing

Each `Makefile` includes a `test` step, so simply type `make test` in [`frontend/`](`frontend/`) and/or [`backend/`](./backend/).

Some tests, notably DynamoDB unit tests, have dependencies best suited for installation in a CI workflow or Docker container. To run *all* tests locally, run `act` ([instructions](#prerequisites)) in this directory.

## Deployment

After building, you may deploy to the AWS cloud ([instructions](./terraform/README.md#provisioning-infrastructure)).