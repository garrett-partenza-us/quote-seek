# QuoteSeek
## Retrieval Augmented Generation for Meditations by Marcus Aurelius

This repository contains an implementation of RAG for Meditations by Marcus Aurelius by leveraging
vector search and structured output from large language models. Users queries are embedded using a
pretrained FastText algorithm and distances are computed to all quotes in the book. Search results
are then passed through ChatGPT 4o-mini for interpretation and advice.

#### Directory Structure
## Repository Directory Structure

- **`ansible/`**  
  This directory contains Ansible playbooks for automating deployment processes.

- **`backend/`**  
  This directory houses the backend logic and GoLang API services for the application.

- **`data/`**  
  This directory contains the cleaned Meditations text.

- **`frontend/`**  
  This directory contains the frontend code and user interface.

- **`model/`**  
  This directory holds machine learning models or other related data.

- **`Dockerfile.backend`**  
  The Dockerfile for building the backend image for the application.

- **`Dockerfile.frontend`**  
  The Dockerfile for building the frontend image for the application.
