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

- `frontend/` contains the frontend website (see [./frontend/README.md](./frontend/README.md) for instructions)
- `backend/` contains the backend microservice (see [./backend/README.md](./backend/README.md) for instructions and an API reference)
- `terraform/` contains infrastructure definitions (see [./terraform/README.md](./terraform/README.md) for instructions)
