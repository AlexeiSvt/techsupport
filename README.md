# User Profiling & Trust Engine

-Neo4j storage is currently being implemented 

**User Profiling & Trust Engine** is a backend system designed to analyze user data, calculate trust scores, and support security-focused workflows such as identity verification and fraud detection.

The system processes both user-provided and historical data, transforming it into a structured trust score while identifying anomalies and suspicious behavior patterns.

---

## Overview

At its core, the engine evaluates how consistent and reliable user data is by comparing input against historical records.

The scoring model:
- Aggregates multiple trust signals (devices, location, account data)
- Applies weighted logic
- Penalizes suspicious activity (e.g., brute-force attempts, abnormal IP behavior)

---

## Core Functionality

### Trust Score Calculation

The final score is based on a weighted model:

- Matching user input with historical data increases trust
- Suspicious signals decrease trust
- The score is normalized and adjusted with penalties

### Data Processing

The system collects and analyzes:

- Account creation date  
- First known device  
- Registration country and city  
- Full device history  

### Risk Detection

Built-in detection mechanisms:

- Suspicious IP analysis  
- Brute-force attempt penalties  
- Inconsistent or conflicting user data  

---

## Project Structure

The project currently consists of four main directories:

### `/core`

The central engine of the system.

- Calculates the final trust score  
- Aggregates and processes user input  
- Applies penalties for suspicious behavior  
- Implements thread-safe logic  

---

### `/logger`

Responsible for system-wide logging.

- Tracks system events  
- Supports debugging and auditing  

---

### `/sysinfo`

Collects system-level information.

- Machine ID  
- CPU cores  
- RAM  
- Operating system  
- Device name and hardware details  

---

### `/storage`

Handles data persistence.

- Database interaction  
- Storage of user data and logs  

---

## Architecture

The project follows **Clean Architecture principles**:

- Clear separation of concerns  
- Decoupled business logic  
- Scalable and maintainable structure  

Key aspects:

- Internal logic is encapsulated in the `internal` directory  
- Shared packages are properly abstracted  
- Strict layering (services, repositories, models)  

---

## Engineering Principles

- **SOLID principles**  
- **KISS (Keep It Simple, Stupid)**  
- Thread-safe design in critical components  
- Fully documented codebase  

---

## Design Patterns & Testing

- **Singleton** — shared instances  
- **Strategy** — flexible scoring logic  
- **TDD (Test-Driven Development)** — core logic is covered with tests  

---

## Tech Stack

- **Language:** Go (Golang)  
- **Architecture:** Clean Architecture  
- **Concurrency:** Native goroutines (thread-safe)  
- **Database:** PostgreSQL  
- **IP Intelligence:** IPinfo API  
- **Documentation:** GoDoc  
- **OS:** Linux  

---

## Security & Privacy

- Sensitive data is processed internally  
- Built-in risk detection mechanisms  
- Safe behavior under concurrent load  
- Designed to minimize data exposure  

---

## Future Improvements

- AI-based anomaly detection  
- Advanced behavioral analytics  
- Implement Websockets for real-time communication

