SentinelCore:
AI-Driven Trust & Safety Intelligence System

SentinelCore is a high-performance, microservice-based backend engine designed to automate and enhance user identity verification and account recovery processes. 

By combining Golang’s concurrency, gRPC for low-latency communication, and Local LLMs (Gemma) for private data analysis, SentinelCore transforms raw logs into actionable security intelligence.

Key Features

1. Hybrid Scoring Engine (Truth-Verification)
Dynamic Trust Score: A proprietary algorithm calculating the "truthfulness" of user-provided data against historical records.
Visual Risk Indicators:

🟢 Green (90-100%): High Trust / Automated Approval.

🟡 Yellow (50-89%): Manual Review Required.

🔴 Red (0-49%): High Risk / Probable Fraud.

2. AI-Powered Intelligence
Automated Case Summaries: Instant extraction of critical user data (Devices, Registration Country and City, etc.) from unstructured support tickets.

Decision Support: AI-generated recommendations for support agents on whether to approve or deny account recovery based on cross-referenced data.

Auto-Response Generation: One-click localized response generation based on the final investigation results.

3. Anomaly & Behavioral Analytics
Geo-Jump Detection: Visualizing login patterns to identify suspicious location shifts (e.g., "Chisinau to Bucharest" jumps).

Session Intelligence: A deep-dive log viewer covering:
Auth Logs: Login timestamps and IP history.
Transaction History: Verification of billing records and payment methods.
Device Fingerprinting: Full history of hardware IDs used to access the account.
Contact Evolution: History of associated emails and phone numbers.

🛠 Tech Stack
Language: Golang (Core Engine & Microservices)

Communication: gRPC + Protocol Buffers (Binary Serialization)

AI Inference: Ollama + Gemma3 (Running locally for 100% Data Privacy)

Architecture: Modular Monolith transitioning to Microservices

OS: Linux (Native Performance)
   
Unlike cloud-based AI solutions, SentinelCore processes all PII (Personally Identifiable Information) locally. Data never leaves the internal infrastructure, ensuring compliance with strict data protection regulations and preventing leaks of sensitive user logs.
