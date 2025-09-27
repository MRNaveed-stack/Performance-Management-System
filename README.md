# Performance Management System (PMS)

A comprehensive backend system developed to manage employee information, attendance tracking, and performance evaluation. The system is designed with scalability, security, and maintainability in mind, ensuring enterprise-level reliability.

---

##  Modules

### 1. Employee Basic Information
- Maintains complete employee records (personal details, job history, qualifications, etc.).
- Fully normalized PostgreSQL schema.
- CRUD operations implemented using stored procedures and functions.

### 2. Attendance
- Integrated with biometric devices using **webhooks** for real-time attendance tracking.
- Supports overtime calculation, leave management, and detailed reporting.
- Ensures secure and accurate record-keeping.

### 3. Performance Management
- Defines and manages **KPIs (Key Performance Indicators)**.
- Handles **performance cycles** for structured evaluations.
- Stores **performance scores** for employee assessments.
- End-to-end workflows controlled through database procedures and functions.
- Work in progress. Implemented during internship; some sections could not be fully tested due to constraints. 
- Core logic is functional, but legacy or uneditable parts may contain errors.

---

##  Technical Implementation

- **Backend:** Go (Golang) with Gin Framework  
- **Database:** PostgreSQL (fully normalized design)  
- **Database Layer:** Stored Procedures & Functions (ensuring data integrity and consistency)  
- **Deployment:** Deployed on a Linux server  
- **Testing:** APIs tested and verified using Postman  
- **Integration:** Biometric attendance device via Webhooks  

---

##  Key Highlights

- Clean and modular API architecture.  
- Secure database access via procedures instead of raw queries.  
- Scalable design approved by senior developers.  
- Actively used within the company for daily operations.  

---

##  Status

The system is fully functional, tested, and deployed. It is being actively used and further extended to meet company requirements.
