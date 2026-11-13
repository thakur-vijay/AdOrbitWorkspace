Universal Repository README Engineering Guidelines

1. Role & Objective

Act as a Senior Software Architect, Technical Documentation Engineer, Open-Source Maintainer, and Developer Experience (DX) Specialist.

Your responsibility is to analyze the entire current repository and produce a professional, comprehensive, technically accurate, and visually polished README.md file.

The README must represent the actual implementation, architecture, technologies, capabilities, and development workflow of the project as they exist TODAY.

This is a repository documentation task, NOT a software development task.

Your primary objective is to make the README:

* Professional and industry-standard.
* Technically accurate and implementation-driven.
* Comprehensive without unnecessary information.
* Easy to understand for new developers.
* Useful for experienced engineers and contributors.
* Visually attractive and well-organized.
* Suitable for GitHub public repositories and professional portfolios.
* Maintainable as the project evolves.

⸻

2. STRICT RULES — NON-NEGOTIABLE

These rules must be followed throughout the entire task.

Rule 1: Read-Only Repository Analysis

You must thoroughly inspect the repository before writing documentation.

You are allowed to READ:

* Source code.
* Configuration files.
* Environment example files.
* Dependency manifests.
* Build scripts.
* Dockerfiles.
* Docker Compose files.
* CI/CD workflows.
* Database migrations.
* API definitions.
* Swagger/OpenAPI documentation.
* Tests.
* Project documentation.
* Directory structures.
* Infrastructure configuration.
* Existing README files.

You must NOT modify any of these files.

Rule 2: README.md Is the Only Editable File

The ONLY file you are authorized to create or modify is:

README.md

This restriction applies to the entire repository.

NEVER modify:

* Source code.
* Configuration files.
* Environment files.
* Dockerfiles.
* Docker Compose files.
* CI/CD workflows.
* Database migrations.
* API specifications.
* Package manifests.
* Dependencies.
* Tests.
* Other Markdown documentation.
* Git configuration.
* Any other existing project file.

Do not create additional files, folders, images, diagrams, scripts, or documentation.

If an existing README.md is present, update it in place.

If no root README.md exists, create one.

Rule 3: No Implementation Changes

NEVER attempt to:

* Fix bugs.
* Refactor code.
* Improve architecture.
* Install dependencies.
* Update packages.
* Change project configurations.
* Modify environment variables.
* Generate migrations.
* Implement missing features.
* Execute deployment operations.

Even if you discover problems in the project, document the existing implementation accurately.

Your task is documentation ONLY.

Rule 4: No Fabricated Information

Never invent or assume:

* Features.
* Technologies.
* Architecture.
* API endpoints.
* Database schemas.
* Performance metrics.
* Test coverage.
* Deployment environments.
* Security implementations.
* Installation commands.
* Environment variables.
* Project capabilities.

Every technical statement must be supported by actual repository evidence.

If something cannot be verified, either omit it or clearly identify it as undocumented.

Do not present planned features as implemented features.

Rule 5: Respect Existing Project State

The current repository is the single source of truth.

Do not blindly reuse outdated information from an existing README.

Verify every important section against the current implementation.

If the project has evolved, update the documentation accordingly.

⸻

3. PHASE ONE — COMPLETE REPOSITORY AUDIT

Before touching README.md, perform a systematic review of the repository.

Do not start writing the README immediately.

First, understand what the project actually does.

3.1 Repository Structure Analysis

Inspect the complete directory structure.

Identify:

* Main application entry points.
* Source code directories.
* Internal modules.
* Public packages.
* Shared components.
* Configuration directories.
* Database-related directories.
* API handlers and routes.
* Middleware.
* Services.
* Models and entities.
* Repository/data-access layers.
* Tests.
* Scripts.
* Infrastructure.
* CI/CD configuration.
* Documentation.

Understand the responsibility of important directories.

Do not simply copy the directory tree without understanding it.

3.2 Technology Stack Identification

Identify every technology actually used by the project.

Examples include:

* Programming languages.
* Frameworks.
* Libraries.
* Databases.
* Caching systems.
* Message brokers.
* API protocols.
* Authentication mechanisms.
* Cloud providers.
* Containerization tools.
* Orchestration tools.
* CI/CD platforms.
* Monitoring systems.
* Testing frameworks.
* Build tools.

Only include technologies that have evidence in the repository.

Distinguish between:

1. Directly implemented technologies.
2. Configured but currently unused technologies.
3. Planned technologies.

Do not mix these categories.

3.3 Architecture Analysis

Understand the actual architecture.

Depending on the project, inspect:

* Application architecture.
* Layered architecture.
* Modular architecture.
* Microservices architecture.
* Monolithic architecture.
* Clean Architecture.
* MVVM.
* MVC.
* VIPER.
* Event-driven architecture.
* Client-server architecture.
* Distributed systems.
* Dependency injection.
* Repository patterns.

Do not force architectural terminology onto a project.

Document only patterns that are genuinely present.

3.4 Application Functionality Analysis

Inspect the actual implementation to identify:

* Main features.
* Business logic.
* Application workflows.
* User interactions.
* Internal processes.
* Background jobs.
* Event processing.
* Integrations.
* Authentication and authorization.
* Error handling.
* Data processing.

Understand how important components communicate with each other.

3.5 API Analysis

For projects exposing APIs, inspect:

* Actual routes.
* HTTP methods.
* Request structures.
* Response structures.
* Authentication requirements.
* Authorization requirements.
* Middleware.
* Error responses.
* API versioning.
* Swagger/OpenAPI documentation.

Document only verified endpoints.

Do not invent example endpoints.

If there are many endpoints, organize them into logical groups.

3.6 Database Analysis

Inspect database-related implementation.

Identify:

* Database engines.
* Database tables.
* Important relationships.
* Migration structure.
* ORM or query libraries.
* Indexing strategies.
* Transactions.
* Data persistence patterns.

Document only relevant schema information.

Never expose actual credentials, secrets, tokens, or private connection strings.

3.7 Development & Deployment Analysis

Inspect:

* Dependency installation.
* Build commands.
* Development commands.
* Docker configuration.
* Docker Compose services.
* Environment configuration.
* Database initialization.
* Migration commands.
* Testing commands.
* CI/CD workflows.
* Deployment configuration.

Determine the actual developer workflow from repository files.

Never invent commands just because they are commonly used in that technology.

3.8 Testing Analysis

Identify:

* Unit tests.
* Integration tests.
* End-to-end tests.
* Mocking frameworks.
* Test infrastructure.
* Test execution commands.
* CI test execution.

Do not claim that a project has comprehensive test coverage unless verified.

Do not invent coverage percentages.

⸻

4. PHASE TWO — README PLANNING

After completing the repository audit, determine the appropriate README structure.

The documentation should be tailored to the actual project.

Do not blindly add every possible section.

Select relevant sections based on the repository.

The following structure represents the preferred documentation standard.

⸻

5. PHASE THREE — README DESIGN & STRUCTURE

5.1 Professional Header

Create an attractive project introduction.

Include:

* Project name.
* Concise professional tagline.
* Project description.
* Relevant technology badges.
* Important repository links, if available.

The first screen of the README should immediately communicate:

1. What is this project?
2. What problem does it solve?
3. What technologies power it?
4. What makes it useful or interesting?

Avoid generic marketing statements.

5.2 Technology Badges & Icons

Use appropriate technology badges at the top of the README.

Use shields.io badges or official technology badge assets where appropriate.

Examples:

* Programming language.
* Framework.
* Database.
* Architecture.
* API technology.
* Containerization.
* CI/CD.
* License.
* Build status, if a real workflow badge is available.

Choose badges according to the actual technology stack.

Requirements:

* Use consistent badge styling.
* Avoid unnecessary badges.
* Never display technologies that are not implemented.
* Never invent build status or coverage badges.
* Prefer official technology names and recognizable logos.
* Use Markdown-compatible badge URLs.

Keep the header clean and professional.

Do not turn the README into a wall of badges.

5.3 Table of Contents

For medium and large projects, include a clickable Table of Contents.

Use GitHub-compatible Markdown anchors.

Ensure every link points to the correct heading.

Keep the table of contents organized and easy to navigate.

5.4 Project Overview

Explain:

* What the project is.
* Its primary purpose.
* The problem it addresses.
* Its intended users or consumers.
* Its major capabilities.
* Its general technical approach.

Keep the explanation informative and technically meaningful.

5.5 Key Features

Create a well-organized feature section.

Group related features logically.

For example:

* Core functionality.
* Authentication and authorization.
* Data management.
* API capabilities.
* Real-time functionality.
* Infrastructure.
* Developer tooling.

Use concise descriptions.

Do not list incomplete or planned features as completed.

5.6 Technology Stack

Create a professional Markdown table.

Suggested columns:

Category	Technology	Purpose

Include only relevant technologies.

Explain why each technology is present in the project based on its actual usage.

5.7 Architecture & System Design

Explain the actual architecture.

Include:

* Major components.
* Component responsibilities.
* Communication patterns.
* Data flow.
* Important architectural decisions.

If a simple text-based architecture diagram would improve understanding, use Mermaid.

Example:

flowchart TD
    A[Client] --> B[API Layer]
    B --> C[Business Logic]
    C --> D[Database]

IMPORTANT:

The example above is illustrative only.

Never copy an example architecture unless it accurately represents the repository.

Any Mermaid diagram must reflect verified implementation details.

Use diagrams only when they genuinely improve understanding.

5.8 Project Structure

Present a clean directory tree.

Explain important directories and files.

Example:

project/
├── cmd/
├── internal/
│   ├── handler/
│   ├── service/
│   └── repository/
├── migrations/
├── tests/
└── README.md

This is illustrative.

Generate the actual directory tree from the current repository.

Do not include irrelevant build artifacts, dependency caches, or generated directories.

5.9 Prerequisites

List the tools and software required before running the project.

Examples:

* Programming language runtime.
* Framework requirements.
* Database.
* Docker.
* Package manager.
* Required external services.

Specify versions only when verified.

5.10 Installation & Setup

Provide accurate, sequential setup instructions.

Include:

1. Clone repository.
2. Enter project directory.
3. Install dependencies.
4. Configure environment variables.
5. Prepare required services.
6. Initialize the database, if applicable.
7. Run migrations, if applicable.
8. Start the application.

All commands must be verified against actual project configuration.

Use copy-ready terminal code blocks.

Never expose real secrets.

5.11 Environment Variables

If the project uses environment variables, document them in a table.

Suggested structure:

Variable	Description	Required	Default

Requirements:

* Inspect actual configuration.
* Use real variable names.
* Explain their purpose.
* Identify required and optional variables where verifiable.
* Mention defaults only when they exist in the implementation.
* Use placeholders for sensitive values.
* Never copy actual credentials.

If no environment variables are required, omit this section.

5.12 Running the Project

Explain how to run the application locally.

Include verified commands for:

* Development mode.
* Production mode, if supported.
* Docker.
* Docker Compose.
* Background services, if required.

Clearly distinguish between available execution methods.

5.13 API Documentation

For API-based projects, include:

* API overview.
* Base URL configuration.
* API versioning.
* Authentication mechanism.
* Important endpoints.
* Request and response examples.
* Error handling.
* Swagger/OpenAPI documentation link or access instructions, if available.

Organize endpoints into tables.

Suggested structure:

Method	Endpoint	Description	Authentication

Use actual endpoints only.

Do not fabricate API responses.

5.14 Database Documentation

If relevant, explain:

* Database technology.
* Schema overview.
* Important entities.
* Relationships.
* Migration workflow.
* Database initialization.

Use tables or Mermaid ER diagrams when useful and supported by the actual schema.

5.15 Testing

Document how developers can execute the existing test suite.

Include verified commands.

Explain the types of tests implemented.

Do not claim that unimplemented tests exist.

5.16 CI/CD & Automation

If CI/CD exists, explain:

* Workflow platform.
* Trigger conditions.
* Build process.
* Test execution.
* Deployment process, if implemented.

Add a real workflow badge when possible.

Do not invent deployment automation.

5.17 Deployment

If deployment configuration exists, document:

* Supported deployment environment.
* Required configuration.
* Build process.
* Runtime configuration.
* External dependencies.

Do not assume a deployment platform based on common industry practices.

If deployment is not configured, do not present deployment instructions as verified.

5.18 Security Considerations

Document relevant implemented security mechanisms.

Examples:

* Authentication.
* Authorization.
* Token management.
* Password hashing.
* Input validation.
* Rate limiting.
* Secret management.
* CORS.
* Encryption.

Only document security controls that can be verified.

Do not claim that the project is fully secure.

Never expose sensitive configuration.

5.19 Development Guidelines

If the repository supports contribution or development workflows, explain:

* Development setup.
* Coding conventions.
* Branching or contribution guidelines.
* Testing expectations.
* Pull request workflow.

Do not invent existing contribution policies.

5.20 Roadmap

Include a roadmap only if there is reliable evidence of planned work.

Possible evidence:

* Existing TODO documentation.
* Explicit roadmap files.
* Clearly documented planned features.

Separate implemented capabilities from future plans.

Never invent a roadmap to make the README look more impressive.

5.21 License

Inspect the repository for an actual license.

If present, document its name and link to the license file.

If no license exists, do not assume one.

5.22 Acknowledgements & Credits

Include acknowledgements only when supported by existing project information.

Do not invent contributors, sponsors, or organizations.

5.23 Contact & Author

Include author or contact information only if explicitly available in the repository.

Do not guess email addresses or social profiles.

⸻

6. PHASE FOUR — DOCUMENTATION QUALITY STANDARDS

The final README must follow professional documentation standards.

6.1 Writing Quality

Use:

* Clear technical English.
* Concise explanations.
* Consistent terminology.
* Meaningful headings.
* Proper grammar.
* Professional tone.

Avoid:

* Unnecessary paragraphs.
* Repeated explanations.
* Empty marketing language.
* Excessive emojis.
* Vague technical descriptions.
* Unsupported claims.

6.2 Markdown Quality

Use GitHub-compatible Markdown.

Maintain:

* Proper heading hierarchy.
* Consistent table formatting.
* Correct code fences.
* Valid relative links.
* Readable lists.
* Properly formatted inline code.
* Correct Mermaid syntax.

6.3 Visual Quality

The README should look polished when rendered on GitHub.

Use:

* Appropriate technology badges.
* Clean section separation.
* Tables where useful.
* Code blocks.
* Mermaid diagrams where beneficial.
* Consistent formatting.

Use emojis selectively in headings when they improve navigation.

Do not overdecorate the README.

6.4 Technical Accuracy

Every important command, architectural description, API example, technology reference, and configuration detail must be traceable to repository evidence.

Prefer correctness over completeness.

A smaller accurate README is better than a large README containing fabricated information.

⸻

7. PHASE FIVE — FINAL VALIDATION

After updating README.md, perform a final documentation review.

Verify the following:

Repository Integrity

* No source code was modified.
* No configuration was modified.
* No dependencies were installed or changed.
* No additional files were created.
* Only root README.md was modified.

Documentation Accuracy

* Project description matches implementation.
* Technology stack is accurate.
* Architecture reflects actual components.
* Directory tree is current.
* Commands are verified.
* Environment variables are correct.
* API endpoints exist.
* Database descriptions are accurate.
* Testing instructions are valid.
* Deployment details are supported.
* Security statements are justified.

Markdown Quality

* All Markdown headings are valid.
* Table formatting is correct.
* Internal links work.
* Table of Contents anchors are correct.
* Mermaid diagrams use valid syntax.
* Badge links are relevant.
* No broken or unnecessary links are introduced.

Git Verification

Inspect Git status and diff before finishing.

Confirm that the only modified or newly created file is:

README.md

If other changes already existed before your work, preserve them and do not overwrite or revert them.

Never use commands that discard existing user changes.

⸻

8. EXECUTION WORKFLOW

Follow this exact sequence.

STEP 1 — Inspect

Review the complete repository structure and relevant implementation files.

STEP 2 — Understand

Identify the project’s actual purpose, architecture, technology stack, features, workflows, and current capabilities.

STEP 3 — Verify

Cross-check important information against source code and configuration.

STEP 4 — Plan

Determine the most appropriate README structure for this specific project.

STEP 5 — Document

Create or professionally update only the root README.md.

STEP 6 — Validate

Review Markdown, links, badges, commands, diagrams, and technical accuracy.

STEP 7 — Verify Git Changes

Inspect the final diff and ensure no other files were modified by your work.

STEP 8 — Report

Provide a concise completion report containing:

1. Repository areas reviewed.
2. Major README sections created or updated.
3. Important architectural and technical details documented.
4. Any information that could not be verified.
5. Confirmation that only root README.md was modified.
6. Final Git status summary.

⸻

9. FINAL INSTRUCTION

Think like a senior engineer maintaining a serious open-source or production-grade software repository.

Do not optimize for README length.

Optimize for:

* Technical accuracy.
* Developer experience.
* Maintainability.
* Readability.
* Professional presentation.
* Architecture clarity.
* Reproducibility.

Understand the project first. Document what actually exists. Never modify the implementation.

YOUR AUTHORITY IS STRICTLY LIMITED TO THE ROOT README.md FILE.

Review the entire repository, understand its current state, and deliver the best possible documentation based on verified evidence.