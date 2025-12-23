# Go-Wallet  🚀

A high-performance, thread-safe E-Wallet core service built with **Go (Golang)** and **PostgreSQL**.

This project demonstrates **Backend Engineering fundamentals** by avoiding ORMs and Frameworks to interact directly with the database using `database/sql`. It focuses on Data Consistency, ACID Transactions, and Clean Code architecture.

---

## 🛠 Tech Stack

- **Language:** Go (1.20+)
- **Database:** PostgreSQL
- **Driver:** `lib/pq`
- **Architecture:** Repository Pattern (Modular & Testable)

---

## 📂 Project Structure

```text
.
├── config/
│   └── database.go         # Singleton Database Connection (Postgres)
├── internal/
│   ├── entity/             # Data Structures (Structs)
│   └── repository/         # Data Access Layer (Raw SQL & Transactions)
├── main.go                 # Entry point & Wiring
├── go.mod                  # Dependency Manager
└── README.md               # Documentation
```

---

## 💾 Database Schema

The system uses a relational schema designed for data integrity.

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    balance INT DEFAULT 0 CHECK (balance >= 0) -- Constraint to prevent negative balance
);
```

---

## 🔥 Key Features & Implementation Details

### 1. Atomic Transfer (ACID)
The core feature is the `Transfer` method. It ensures that money is never lost during a transaction, even if the server crashes.

**Logic Flow:**
1.  `tx.Begin()`: Start transaction.
2.  **Debit Sender:** `UPDATE users SET balance = balance - amount WHERE id = sender AND balance >= amount`.
3.  **Check RowsAffected:** If 0, rollback (Insufficient funds).
4.  **Credit Receiver:** `UPDATE users SET balance = balance + amount WHERE id = receiver`.
5.  **Check RowsAffected:** If 0, rollback (Receiver not found).
6.  `tx.Commit()`: Permanent save.

### 2. Repository Pattern
Separating business logic from database queries allows for easier testing and maintenance.

---

## 🚀 How to Run

1.  **Clone the repository**
    ```bash
    git clone https://github.com/FarhanDH/e-wallet-core.git
    cd e-wallet
    ```

2.  **Setup Database**
    - Create a PostgreSQL database named `ewallet`.
    - Run the SQL creation script above.
    - Update `config/database.go` with your credentials.

3.  **Run the Application**
    ```bash
    go mod tidy
    go run main.go
    ```

---

## 🔜 Roadmap (Upcoming Features)

- [x] Database Connection (Postgres)
- [x] CRUD User (Raw SQL)
- [x] Atomic Transfer (Transaction)
- [x] REST API Implementation (using Gin)
- [ ] Unit Testing with Mocking
- [ ] Dockerization

---