# ratelimiter
Here is a complete, phased plan of action. This roadmap takes you from an empty directory to a fully automated, cloud-native portfolio piece designed perfectly to showcase a Software Engineer who commands deep cloud infrastructure expertise.

## Phase 1: Local Development & Go Learning

*Goal: Master Go syntax and concurrency through TDD without touching the cloud yet.*

1. **Initialize the Monorepo:** Set up your project folder structure to house both the Go application code and the infrastructure code (e.g., AWS CDK) side-by-side.
2. **Setup the AI Mentor:** Configure Cursor, Windsurf, or Claude with the strict Socratic prompt to ensure it guides your learning rather than writing the code.
3. **Write the Interfaces:** Define your `StateStore` and `Algorithm` interfaces.
4. **Test-Driven Implementation:**
* Build the `MemoryStore` using Go's `sync.Mutex`.
* Implement the Token Bucket, Sliding Window Counter, and Sliding Window Log using pure functions (`Decide`).


5. **Add Local Redis:** Spin up a local Redis instance using Docker (leveraging your standard container workflows) and implement the `RedisStore` using Lua scripts and Sorted Sets.

## Phase 2: System Design Capture

*Goal: Document the "why" behind your architecture.*

1. **Write ADRs:** Create a `docs/adr` folder. For every major decision (e.g., Mutex vs. Channels in Go, Lua vs. native Redis commands), write a concise Markdown Architectural Decision Record.
2. **Diagram the Architecture:** Use Eraser.io to map out two diagrams:
* **Logical Architecture:** How the Go interfaces interact.
* **Cloud Architecture:** The target AWS deployment (VPC, ECS, ElastiCache, ALB).



## Phase 3: Benchmarking & Load Testing

*Goal: Prove the mathematical trade-offs of your algorithms.*

1. **Go Benchmarks:** Write internal `go test -bench` scripts to compare the memory and CPU footprint of your three algorithms using the `MemoryStore`.
2. **Local API Testing:** Wrap your rate limiter in a lightweight HTTP server (e.g., using Go's standard `net/http`).
3. **k6 Load Generation:** Write a JavaScript k6 script to hammer your local endpoints. Capture the latency differences between the `MemoryStore` and `RedisStore`.

## Phase 4: Cloud Infrastructure (IaC)

*Goal: Translate the local application into a production-grade AWS environment.*

1. **Containerize:** Write a multi-stage `Dockerfile` to compile your Go binary into a minimal scratch image.
2. **Define the Infrastructure:** Write the AWS CDK stack to provision:
* A custom VPC with public and private subnets.
* An ECS Fargate cluster to run your Go container.
* A Multi-AZ ElastiCache for Redis cluster in the private subnets.
* An Application Load Balancer (ALB) to route traffic to Fargate.


3. **Inject Observability:** Add OpenTelemetry to your Go code to push your algorithm metrics directly to Amazon CloudWatch.

## Phase 5: CI/CD & Ephemeral FinOps

*Goal: Automate the deployment, validate it, and tear it down to maintain strict cost optimization.*

1. **Build the Pipeline:** Create a GitHub Actions workflow (`.github/workflows/deploy-and-test.yml`) with the following stages:
* **Lint & Test:** Run `go test`.
* **Deploy:** Run `cdk deploy --require-approval never`.
* **Validate:** Run the k6 load test against the newly provisioned ALB endpoint.
* **Destroy:** Run `cdk destroy --force` (ensuring this runs `if: always()` so you never leave orphaned infrastructure).


2. **Capture Artifacts:** Run the pipeline once. Take screenshots of the CloudWatch dashboards showing the ECS tasks scaling and handling the k6 load.
3. **Publish:** Add the architecture diagrams, the benchmark results, and the CloudWatch screenshots to your `README.md`.