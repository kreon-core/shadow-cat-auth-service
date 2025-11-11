## SUCCESS CODES

| Code | Message                                            |
| ---- | -------------------------------------------------- |
| 5    | Resource deleted successfully                      |
| 4    | Resource updated successfully                      |
| 3    | Resource created successfully                      |
| 2    | Request accepted (for Asynchronous jobs only)      |
| 1    | Success                                            |
| 0    | Ignored (still success but no-content, no-process) |

---

## ERROR CODES

### Unspecified && Validation

| Code | Message                                                                 |
| ---- | ----------------------------------------------------------------------- |
| -1   | Unspecified error (due to panic / unknown exception / developer issues) |
| -2   | Invalid request                                                         |
| -90  | Zalo signature mismatch (does not match provided data)                  |

**Note:** More validation-related codes can be added later as needed.

---

### Authentication & Authorization

**Enumeration counted down from `-100`**

| Code | Message                          |
| ---- | -------------------------------- |
| -100 | Invalid signature or credentials |
| -101 | Missing authorization header     |
| -102 | Invalid access token             |
| -103 | Expired access token             |
| -104 | Invalid refresh token            |
| -105 | Access denied or token revoked   |
| -106 | Account suspended or banned      |

---

### User & Account

**Enumeration counted down from `-200`**

| Code | Message                              |
| ---- | ------------------------------------ |
| -200 | User or Player not found             |
| -201 | User or Player already exists        |
| -202 | Username or Email already registered |
| -203 | Too many login attempts              |

**Note:**
Rate limiting is handled via Redis key pattern
`auth::username/email/provider_uid` to prevent brute-force login attacks.

---

### Game Logics

**Enumeration counted down from `-300`**

| Code | Message          |
| ---- | ---------------- |
| -300 | Not enough coins |
| -301 | Not enough gems  |
| -302 | Energy depleted  |

---

### Data & Resource

**Enumeration counted down from `-400`**

| Code | Message                                        |
| ---- | ---------------------------------------------- |
| -400 | Resource not found                             |
| -401 | Resource already exists or uniqueness conflict |
| -402 | Database error                                 |
| -403 | Data integrity violation                       |
| -404 | Record not allowed to modify                   |

---

### Forward External API

**Enumeration counted down from `-500`**

| Code | Message                         |
| ---- | ------------------------------- |
| -500 | External API error              |
| -501 | External API timeout            |
| -502 | Third-party service unavailable |
