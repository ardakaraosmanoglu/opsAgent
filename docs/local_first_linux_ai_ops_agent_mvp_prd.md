# MVP PRD v1 — Local-First Linux AI Ops Agent

**Doküman amacı:**  
Bu doküman, projeyi geliştirecek developer’a doğrudan teslim edilmek üzere hazırlanmıştır. Ürün; Linux sunuculara tek komutla kurulan, lokal çalışan, kendi içinde web paneli ve SQLite veritabanı bulunan, AI destekli güvenli server monitoring ve operasyon agent’ıdır.

**MVP yaklaşımı:**  
Merkezi SaaS, PostgreSQL, Redis, NestJS, Vercel, Kubernetes gibi ek bağımlılıklar kullanılmayacak. İlk MVP tamamen kullanıcının kendi Linux sunucusunda lokal çalışacak.

---

## 1. Proje Özeti

Bu proje, Linux sunuculara tek komutla kurulan, lokal çalışan, AI destekli bir server monitoring ve operasyon asistanıdır.

Agent kullanıcının kendi Linux sunucusunda çalışır. Ekstra merkezi backend, harici database, Redis, Docker Compose veya ayrı servis kurulumu gerektirmez. Kurulumdan sonra kullanıcı lokal web paneline erişir, sunucusunun durumunu izler, sorunları inceler ve sistem değişikliği gerektiren işlemleri yalnızca açık onay vererek çalıştırır.

Ürünün temel prensibi:

```text
Read-only by default.
Write actions require explicit user approval.
```

Yani sistem varsayılan olarak sadece okuma ve analiz yapar. Dosya silme, paket güncelleme, servis restart etme, log temizleme gibi işlemler kullanıcı onayı olmadan asla yapılmaz.

---

## 2. MVP Hedefi

MVP’nin amacı, kullanıcıya şu deneyimi sunmaktır:

```text
1. Kullanıcı Linux sunucusuna SSH ile bağlanır.
2. Tek satırlık kurulum komutunu çalıştırır.
3. Agent systemd servisi olarak kurulur.
4. Lokal dashboard açılır.
5. Kullanıcı admin hesabı oluşturur.
6. AI API key girer veya AI olmadan devam eder.
7. İlk sistem taraması yapılır.
8. CPU, RAM, disk, servisler, portlar ve loglar izlenir.
9. Problem varsa alert oluşturulur.
10. Kullanıcı prompt veya hazır aksiyon ile analiz ister.
11. Agent çözüm planı ve çalıştırılacak komutları gösterir.
12. Kullanıcı onaylarsa komutlar çalıştırılır.
13. Sonuç ve audit log kaydedilir.
```

MVP’nin hedefi “her şeyi otomatik çözen AI sysadmin” değildir. Hedef:

> Kullanıcının Linux sunucusunu güvenli şekilde anlamasını, izlemesini ve onay kontrollü operasyon yapmasını sağlamak.

---

## 3. Ürün Konumlandırması

### Kısa Tanım

```text
Tek komutla Linux sunucunuza kurulan, lokal çalışan, AI destekli güvenli server monitoring ve operasyon agent’ı.
```

### Teknik Tanım

```text
Self-hosted Linux ops agent with local dashboard, SQLite storage, system metrics monitoring, AI-assisted planning, approval-based command execution, and audit logging.
```

### Ana Değer Önerisi

- Kullanıcı kendi sunucusuna kurar.
- Veriler kendi sunucusunda kalır.
- Ekstra database bağlantısı gerekmez.
- Tek binary/systemd service olarak çalışır.
- Sistem değişikliği yapan hiçbir işlem onaysız çalışmaz.
- Disk, CPU, RAM, servis, port ve log problemlerini tespit eder.
- Çözüm planı ve komutları kullanıcıya açıkça gösterir.

---

## 4. MVP Kapsamı

### 4.1 MVP’de Olacak Özellikler

#### Kurulum

- Tek komutla kurulum.
- Linux dağıtımı ve architecture kontrolü.
- Binary indirme.
- Gerekli klasörleri oluşturma.
- SQLite database oluşturma.
- systemd service oluşturma.
- Service başlatma.
- Terminalde dashboard erişim bilgisini gösterme.

#### Lokal Dashboard

- İlk kurulum wizard’ı.
- Admin hesabı oluşturma.
- AI provider/API key ayarı.
- Dashboard ana ekranı.
- Alerts ekranı.
- Assistant/prompt ekranı.
- Task approval ekranı.
- Audit log ekranı.
- Settings ekranı.

#### Monitoring

Agent şu bilgileri izlemeli:

- Hostname
- OS bilgisi
- Kernel version
- Uptime
- CPU usage
- RAM usage
- Disk usage
- Load average
- Top CPU processes
- Top memory processes
- Open/listening ports
- systemd services status
- `/var/log` altında büyük log dosyaları
- `journalctl --disk-usage`
- Docker varsa basic container listesi

#### Alert Sistemi

MVP threshold tabanlı çalışacak.

İlk kurallar:

- Disk usage > %85 warning
- Disk usage > %95 critical
- RAM usage > %90 warning
- CPU usage > %90 ve 5 dakika devam ederse warning
- Yeni açılan listening port tespit edilirse info/warning
- `/var/log` altında büyük dosya tespit edilirse warning
- Journal log kullanımı belirli eşiği aşarsa warning

#### Assistant

Kullanıcı panelden doğal dilde soru/görev yazabilecek.

Örnekler:

```text
Disk neden doluyor?
En çok RAM kullanan servis ne?
Sunucuyu güncelle.
Nginx loglarını kontrol et.
Açık portları göster.
Gereksiz büyük log dosyalarını bul.
```

#### Plan ve Approval

Sistem değişikliği gerektiren her işlem için önce plan oluşturulacak.

Plan içinde şunlar olacak:

- Kullanıcının isteği
- Sistemden toplanan veri
- Önerilen çözüm
- Çalıştırılacak komutlar
- Risk seviyesi
- Olası etkiler
- Onay butonu
- Reddet butonu

#### Command Execution

Onaylanan komutlar çalıştırılacak.

Çalıştırma sonrası:

- Komut çıktısı gösterilecek.
- Exit code gösterilecek.
- Başarılı/başarısız durumu gösterilecek.
- Audit log’a kaydedilecek.

#### Audit Log

Şunlar kayıt altına alınacak:

- Kim giriş yaptı?
- Hangi prompt gönderildi?
- Hangi plan üretildi?
- Hangi komutlar önerildi?
- Hangi komut onaylandı?
- Hangi komut çalıştırıldı?
- Çıktısı neydi?
- Exit code neydi?
- Ne zaman çalıştı?

---

### 4.2 MVP’de Olmayacak Özellikler

Bunlar bilinçli olarak dışarıda bırakılacak:

- Merkezi SaaS panel
- PostgreSQL
- Redis
- BullMQ
- NestJS backend
- Vercel deployment
- Kubernetes
- Telegram bot
- Slack/Discord entegrasyonu
- Multi-server merkezi yönetim
- Multi-user team workspace
- Role-based access control
- Otomatik self-healing
- Tam ML tabanlı anomaly detection
- Kubernetes cluster monitoring
- SIEM/security ürünü gibi davranma
- Firewall otomasyonu
- Root-level sınırsız komut çalıştırma

---

## 5. Önerilen Teknoloji Stack

### 5.1 Ana Stack

| Katman | Teknoloji |
|---|---|
| Agent/backend | Go |
| Database | SQLite |
| UI | React + Vite, Go binary içine embedded |
| Service manager | systemd |
| Config | YAML |
| Installer | Bash install script |
| AI provider | OpenAI-compatible API başlangıç |
| Local API | Go HTTP server |
| Auth | Local session/JWT |
| Logging | File log + SQLite audit |
| Deployment | Single Linux binary |

### 5.2 Neden Go?

Bu proje için Go seçilme sebepleri:

- Tek binary dağıtım kolaylığı.
- Linux agent geliştirmeye uygunluk.
- Düşük kaynak tüketimi.
- Kolay öğrenme eğrisi.
- Rust’a göre daha hızlı MVP geliştirme.
- Web server, SQLite, background worker ve system command execution için yeterli ekosistem.
- Bakımı daha kolay ve okunabilir kod yapısı.
- systemd service olarak çalıştırmaya uygunluk.

Rust teknik olarak güçlü bir seçenek olsa da MVP için geliştirme hızını düşürür. Bu ürünün ilk hedefi en mükemmel sistem diliyle yazılmak değil, güvenli ve stabil bir local-first MVP çıkarmaktır.

---

## 6. Sistem Mimarisi

### 6.1 Genel Mimari

```text
┌───────────────────────────────────────┐
│          User Linux Server             │
│                                       │
│  ┌─────────────────────────────────┐  │
│  │          opsagent binary         │  │
│  │                                 │  │
│  │  - Local Web Server              │  │
│  │  - Embedded UI                   │  │
│  │  - Local REST API                │  │
│  │  - Monitoring Engine             │  │
│  │  - Alert Engine                  │  │
│  │  - Assistant / Planner           │  │
│  │  - Policy Engine                 │  │
│  │  - Command Executor              │  │
│  │  - Audit Logger                  │  │
│  └───────────────┬─────────────────┘  │
│                  │                    │
│           SQLite Database             │
│    /var/lib/opsagent/opsagent.db      │
│                  │                    │
│           systemd service             │
│    /etc/systemd/system/opsagent.service│
│                                       │
└───────────────────────────────────────┘
```

### 6.2 Çalışma Modeli

Agent tek bir Go binary olarak çalışır.

Bu binary içinde:

- Web paneli serve eder.
- Local API sağlar.
- Metrik toplar.
- Alert üretir.
- AI planlama yapar.
- Komut policy kontrolü yapar.
- Onaylanan komutları çalıştırır.
- SQLite’a veri yazar.
- Audit log tutar.

---

## 7. Dosya Sistemi Yapısı

Kurulum sonrası dosya yapısı şöyle olmalı:

```text
/usr/local/bin/opsagent
/etc/opsagent/config.yaml
/var/lib/opsagent/opsagent.db
/var/log/opsagent/opsagent.log
/etc/systemd/system/opsagent.service
```

| Path | Amaç |
|---|---|
| `/usr/local/bin/opsagent` | Ana binary |
| `/etc/opsagent/config.yaml` | Config dosyası |
| `/var/lib/opsagent/opsagent.db` | SQLite database |
| `/var/log/opsagent/opsagent.log` | Agent logları |
| `/etc/systemd/system/opsagent.service` | systemd service |

---

## 8. Config Taslağı

```yaml
app:
  name: "OpsAgent"
  environment: "production"

server:
  bind_address: "127.0.0.1"
  port: 8787

database:
  path: "/var/lib/opsagent/opsagent.db"

logging:
  level: "info"
  file: "/var/log/opsagent/opsagent.log"

security:
  require_approval_for_write_actions: true
  command_timeout_seconds: 120
  max_output_size_kb: 1024

ai:
  enabled: false
  provider: ""
  api_key: ""
  model: ""

monitoring:
  interval_seconds: 30
  disk_warning_threshold: 85
  disk_critical_threshold: 95
  memory_warning_threshold: 90
  cpu_warning_threshold: 90
  cpu_warning_duration_seconds: 300

dashboard:
  session_secret: "generated-random-secret"
```

API key config dosyasında tutulursa dosya izinleri sıkı olmalı:

```bash
chmod 600 /etc/opsagent/config.yaml
chown root:root /etc/opsagent/config.yaml
```

Daha güvenli alternatif olarak API key encrypted şekilde SQLite içinde tutulabilir. MVP’de config dosyası yeterli olabilir ama permission zorunlu olmalı.

---

## 9. Kurulum Akışı

### 9.1 Kullanıcı Komutu

```bash
curl -fsSL https://opsagent.dev/install.sh | sudo bash
```

### 9.2 Installer’ın Yapacağı İşler

```text
1. Root yetkisi kontrol edilir.
2. OS kontrol edilir.
3. Architecture kontrol edilir.
4. Desteklenen sistem mi kontrol edilir.
5. En son opsagent binary indirilir.
6. Binary /usr/local/bin/opsagent konumuna taşınır.
7. /etc/opsagent oluşturulur.
8. /var/lib/opsagent oluşturulur.
9. /var/log/opsagent oluşturulur.
10. Varsayılan config.yaml oluşturulur.
11. systemd service dosyası oluşturulur.
12. systemctl daemon-reload çalıştırılır.
13. systemctl enable opsagent çalıştırılır.
14. systemctl start opsagent çalıştırılır.
15. Service status kontrol edilir.
16. Kullanıcıya panel erişim bilgisi gösterilir.
```

### 9.3 Terminal Çıktısı

Kurulum başarılı olduğunda kullanıcı şunu görmeli:

```text
OpsAgent installed successfully.

Service:
active

Dashboard:
http://127.0.0.1:8787

To access the dashboard from your computer:

ssh -L 8787:localhost:8787 root@YOUR_SERVER_IP

Then open:

http://localhost:8787

Your server data stays on this machine.
No write operation will run without approval.
```

---

## 10. İlk Kullanım Akışı

### 10.1 İlk Panel Açılışı

Kullanıcı şu adrese girer:

```text
http://localhost:8787
```

### 10.2 Setup Wizard

#### Step 1 — Welcome

```text
Welcome to OpsAgent.

This dashboard runs locally on your own Linux server.
OpsAgent will not modify your system without your approval.
```

Buton:

```text
Start Setup
```

#### Step 2 — Admin Account

Form alanları:

```text
Username
Password
Confirm Password
```

Kurallar:

- İlk oluşturulan kullanıcı admin olur.
- Şifre hashlenerek saklanır.
- Plain password asla saklanmaz.

#### Step 3 — AI Provider

Form alanları:

```text
AI Provider
API Key
Model
```

MVP için seçenekler:

```text
OpenAI-compatible
Skip for now
```

AI key girilmezse:

- Monitoring çalışır.
- Alerts çalışır.
- Basic analysis çalışır.
- AI assisted planning devre dışı olur.

#### Step 4 — Access Mode Info

Varsayılan erişim:

```text
Localhost only
```

Mesaj:

```text
For security, the dashboard is available only from 127.0.0.1 by default.
Use SSH tunnel to access it safely.
```

#### Step 5 — First Scan

Agent ilk sistem taramasını yapar.

Toplanacak bilgiler:

- Hostname
- OS
- Kernel
- Uptime
- CPU
- RAM
- Disk
- Load average
- Top processes
- Open ports
- systemd services
- Large log files
- Journal disk usage

#### Step 6 — Dashboard Ready

Kullanıcı ana dashboard’a yönlendirilir.

---

## 11. Dashboard Ekranları

### 11.1 Dashboard

Amaç: Kullanıcı sunucunun genel durumunu görür.

Kartlar:

```text
CPU Usage
Memory Usage
Disk Usage
Load Average
Uptime
Open Ports
Active Services
Critical Alerts
```

Alt bölümler:

```text
Recent Alerts
Top CPU Processes
Top Memory Processes
Disk Overview
Quick Actions
```

Quick Actions:

```text
Why is my disk full?
Check high CPU usage
Show open ports
Analyze system logs
Check package updates
```

### 11.2 Alerts

Amaç: Sistem tarafından tespit edilen uyarıları göstermek.

Alert alanları:

```text
Severity
Type
Message
Created At
Status
```

Severity:

```text
info
warning
critical
```

Status:

```text
open
acknowledged
resolved
ignored
```

### 11.3 Alert Detail

Örnek disk alert ekranı:

```text
Disk usage reached 91%.

Current usage:
91%

Largest directories:
- /var/log: 12.4 GB
- /var/lib/docker: 5.6 GB
- /home/app/storage: 3.2 GB

Suggested actions:
- Analyze logs
- Create cleanup plan
- Ignore alert
```

### 11.4 Assistant

Prompt input:

```text
Ask OpsAgent anything about this server...
```

Örnek promptlar:

```text
Disk neden doluyor?
Sunucuyu güncelle.
Nginx hata veriyor mu?
En çok RAM kullanan process ne?
Açık portları kontrol et.
```

Assistant cevapları ikiye ayrılır:

1. **Read-only cevap:** Sadece analiz yapılır, onay gerekmez.
2. **Write action plan:** Sistem değişikliği gerekiyorsa plan ve approval oluşturulur.

### 11.5 Task Approval

Amaç: Kullanıcı sistem değişikliği yapacak işlemleri açıkça görüp onaylar.

Task plan içeriği:

```text
Task:
Sunucuyu güncelle

Detected system:
Ubuntu 22.04

Commands:
1. apt update
2. apt list --upgradable
3. apt upgrade -y

Risk:
Medium

Possible impact:
Some services may restart.
Backup/snapshot is recommended.

Approval required.
```

Butonlar:

```text
Approve and Run
Reject
Edit / Ask Again
```

### 11.6 Task Execution Result

Komut çalıştıktan sonra:

```text
Command:
journalctl --vacuum-time=7d

Status:
Success

Exit code:
0

Output:
Vacuuming done, freed 3.1G of archived journals.

Saved to audit log.
```

### 11.7 Audit Log

Filtreler:

```text
Date
Event Type
Risk Level
Status
```

Event types:

```text
login
setup_completed
metric_collected
alert_created
prompt_submitted
plan_generated
command_approved
command_executed
command_failed
settings_updated
```

### 11.8 Settings

Ayarlar:

```text
AI Provider
API Key
Model
Dashboard Bind Address
Dashboard Port
Monitoring Interval
Alert Thresholds
Access Mode
Update Agent
```

Access mode default:

```text
127.0.0.1 only
```

Public mode açılırsa uyarı gösterilmeli:

```text
Warning: Public dashboard access may expose your server management panel to the internet.
Use HTTPS, strong password, and firewall rules.
```

---

## 12. Veri Modeli

SQLite kullanılacak.

### 12.1 users

```sql
CREATE TABLE users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  username TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  role TEXT NOT NULL DEFAULT 'admin',
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL
);
```

### 12.2 settings

```sql
CREATE TABLE settings (
  key TEXT PRIMARY KEY,
  value TEXT NOT NULL,
  updated_at DATETIME NOT NULL
);
```

### 12.3 metrics

```sql
CREATE TABLE metrics (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  cpu_usage REAL,
  memory_usage REAL,
  disk_usage REAL,
  load_average_1 REAL,
  load_average_5 REAL,
  load_average_15 REAL,
  uptime_seconds INTEGER,
  created_at DATETIME NOT NULL
);
```

### 12.4 process_snapshots

```sql
CREATE TABLE process_snapshots (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  pid INTEGER,
  name TEXT,
  user TEXT,
  cpu_usage REAL,
  memory_usage REAL,
  command TEXT,
  created_at DATETIME NOT NULL
);
```

### 12.5 ports

```sql
CREATE TABLE ports (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  protocol TEXT,
  local_address TEXT,
  port INTEGER,
  process_name TEXT,
  pid INTEGER,
  first_seen_at DATETIME NOT NULL,
  last_seen_at DATETIME NOT NULL,
  status TEXT NOT NULL
);
```

### 12.6 services

```sql
CREATE TABLE services (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  status TEXT,
  active_state TEXT,
  sub_state TEXT,
  first_seen_at DATETIME NOT NULL,
  last_seen_at DATETIME NOT NULL
);
```

### 12.7 alerts

```sql
CREATE TABLE alerts (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  type TEXT NOT NULL,
  severity TEXT NOT NULL,
  title TEXT NOT NULL,
  message TEXT NOT NULL,
  status TEXT NOT NULL DEFAULT 'open',
  metadata_json TEXT,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL
);
```

### 12.8 tasks

```sql
CREATE TABLE tasks (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  source TEXT NOT NULL,
  prompt TEXT,
  status TEXT NOT NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL
);
```

`source` değerleri:

```text
assistant
alert
quick_action
system
```

`status` değerleri:

```text
draft
waiting_approval
approved
running
completed
failed
rejected
```

### 12.9 task_plans

```sql
CREATE TABLE task_plans (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  task_id INTEGER NOT NULL,
  summary TEXT NOT NULL,
  risk_level TEXT NOT NULL,
  requires_approval INTEGER NOT NULL DEFAULT 1,
  metadata_json TEXT,
  created_at DATETIME NOT NULL,
  FOREIGN KEY (task_id) REFERENCES tasks(id)
);
```

### 12.10 commands

```sql
CREATE TABLE commands (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  task_plan_id INTEGER NOT NULL,
  command TEXT NOT NULL,
  command_type TEXT NOT NULL,
  risk_level TEXT NOT NULL,
  requires_approval INTEGER NOT NULL,
  status TEXT NOT NULL DEFAULT 'pending',
  created_at DATETIME NOT NULL,
  FOREIGN KEY (task_plan_id) REFERENCES task_plans(id)
);
```

`command_type`:

```text
read
write
dangerous
blocked
```

### 12.11 command_executions

```sql
CREATE TABLE command_executions (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  command_id INTEGER NOT NULL,
  output TEXT,
  error_output TEXT,
  exit_code INTEGER,
  started_at DATETIME NOT NULL,
  finished_at DATETIME,
  status TEXT NOT NULL,
  FOREIGN KEY (command_id) REFERENCES commands(id)
);
```

### 12.12 audit_logs

```sql
CREATE TABLE audit_logs (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  event_type TEXT NOT NULL,
  actor TEXT,
  message TEXT NOT NULL,
  metadata_json TEXT,
  created_at DATETIME NOT NULL
);
```

---

## 13. Backend / Local API Taslağı

Go içinde local HTTP API olacak.

### 13.1 Auth

```http
POST /api/setup/admin
POST /api/auth/login
POST /api/auth/logout
GET  /api/auth/me
```

### 13.2 Dashboard

```http
GET /api/dashboard/summary
GET /api/metrics/latest
GET /api/metrics/history
GET /api/processes/top
GET /api/ports
GET /api/services
```

### 13.3 Alerts

```http
GET  /api/alerts
GET  /api/alerts/:id
POST /api/alerts/:id/acknowledge
POST /api/alerts/:id/resolve
POST /api/alerts/:id/ignore
```

### 13.4 Assistant

```http
POST /api/assistant/message
GET  /api/tasks/:id
GET  /api/tasks
```

Request:

```json
{
  "message": "Disk neden doluyor?"
}
```

Response read-only example:

```json
{
  "type": "answer",
  "summary": "Disk usage is high because /var/log is using 12.4 GB.",
  "requires_approval": false
}
```

Response plan example:

```json
{
  "type": "plan",
  "task_id": 12,
  "requires_approval": true
}
```

### 13.5 Task Approval

```http
GET  /api/tasks/:id/plan
POST /api/tasks/:id/approve
POST /api/tasks/:id/reject
POST /api/tasks/:id/run
GET  /api/tasks/:id/executions
```

### 13.6 Settings

```http
GET  /api/settings
PUT  /api/settings
POST /api/settings/ai
POST /api/settings/access-mode
```

### 13.7 System

```http
GET  /api/system/info
POST /api/system/scan
GET  /api/system/service-status
```

---

## 14. Modül Yapısı

Clean Architecture uygulanabilir. Fakat aşırı soyutlama yapılmamalı. Go tarafında “pragmatic clean architecture” daha doğru olur.

Önerilen klasör yapısı:

```text
opsagent/
├── cmd/
│   └── opsagent/
│       └── main.go
│
├── internal/
│   ├── app/
│   │   ├── bootstrap.go
│   │   └── scheduler.go
│   │
│   ├── config/
│   │   └── config.go
│   │
│   ├── domain/
│   │   ├── alert.go
│   │   ├── metric.go
│   │   ├── task.go
│   │   ├── command.go
│   │   └── audit.go
│   │
│   ├── usecase/
│   │   ├── collect_metrics.go
│   │   ├── evaluate_alerts.go
│   │   ├── create_task_plan.go
│   │   ├── approve_task.go
│   │   ├── execute_command.go
│   │   └── write_audit_log.go
│   │
│   ├── infrastructure/
│   │   ├── storage/
│   │   │   └── sqlite/
│   │   ├── system/
│   │   │   ├── metrics.go
│   │   │   ├── processes.go
│   │   │   ├── ports.go
│   │   │   ├── services.go
│   │   │   └── logs.go
│   │   ├── ai/
│   │   │   └── openai_compatible.go
│   │   └── executor/
│   │       └── command_runner.go
│   │
│   ├── policy/
│   │   ├── classifier.go
│   │   ├── whitelist.go
│   │   └── denylist.go
│   │
│   ├── web/
│   │   ├── api/
│   │   ├── middleware/
│   │   ├── auth/
│   │   └── server.go
│   │
│   └── ui/
│       └── embed.go
│
├── web/
│   └── dashboard/
│       ├── src/
│       ├── package.json
│       └── dist/
│
├── scripts/
│   ├── install.sh
│   └── uninstall.sh
│
├── migrations/
│   └── 001_init.sql
│
├── go.mod
└── README.md
```

---

## 15. Clean Architecture Yaklaşımı

Bu projede full enterprise clean architecture yerine sade ama katmanlı yapı önerilir.

### 15.1 Domain Layer

Saf modeller burada olur.

Örnek:

```text
Metric
Alert
Task
TaskPlan
Command
CommandExecution
AuditLog
```

Bu layer sistem komutu çalıştırmaz, HTTP bilmez, SQLite bilmez.

### 15.2 Usecase Layer

İş kuralları burada olur.

Örnek usecase’ler:

```text
CollectMetrics
EvaluateAlerts
CreateTaskPlan
ApproveTask
ExecuteApprovedCommand
SaveAuditLog
```

Bu layer domain mantığını yönetir.

### 15.3 Infrastructure Layer

Dış dünya ile iletişim burada olur.

Örnekler:

```text
SQLite repository
Linux command collector
AI provider client
Command runner
File logger
```

### 15.4 Web Layer

HTTP API ve UI serving burada olur.

```text
Routes
Handlers
Middleware
Session/Auth
Static UI serving
```

### 15.5 Policy Layer

Bu proje için ayrıca policy layer şart.

Görevleri:

- Komut read/write/dangerous sınıflandırması.
- Whitelist kontrolü.
- Denylist kontrolü.
- Risk seviyesi belirleme.
- Approval zorunluluğunu belirleme.

---

## 16. Command Policy Kuralları

Bu ürünün en kritik parçası command policy sistemidir.

### 16.1 Command Types

```text
read
write
dangerous
blocked
```

### 16.2 Read-only Komut Örnekleri

Bunlar onaysız çalışabilir:

```bash
df -h
free -m
uptime
ps aux
top -b -n 1
systemctl status nginx
journalctl -u nginx --since "1 hour ago"
ss -tulpn
cat /etc/os-release
apt list --upgradable
du -sh /var/log/*
```

Not: Bazı `du` komutları permission açısından root isteyebilir ama sistem değişikliği yapmaz.

### 16.3 Write Komut Örnekleri

Bunlar mutlaka onay ister:

```bash
apt update
apt upgrade -y
systemctl restart nginx
systemctl reload nginx
journalctl --vacuum-time=7d
docker restart container_name
```

### 16.4 Dangerous / Blocked Komutlar

MVP’de direkt engellenmeli:

```bash
rm -rf /
mkfs
dd if=
:(){ :|:& };:
chmod -R 777 /
chown -R /
userdel
passwd
iptables -F
ufw disable
shutdown
reboot
curl something | bash
wget something | bash
```

`reboot` ileride yüksek riskli onayla açılabilir ama MVP’de kapalı kalması daha sağlıklı.

### 16.5 Komut Validasyon Akışı

```text
1. Plan oluşturulur.
2. Her komut policy classifier’dan geçer.
3. Komut read/write/dangerous/blocked olarak sınıflandırılır.
4. Blocked komut varsa plan reddedilir.
5. Write komut varsa approval required true olur.
6. Kullanıcı onaylamadan execution başlamaz.
7. Execution öncesi tekrar policy check yapılır.
8. Output audit log’a kaydedilir.
```

---

## 17. AI Planner Davranışı

### 17.1 AI Ne Yapabilir?

AI şunları yapabilir:

- Kullanıcı promptunu yorumlar.
- Sistemden gelen verileri açıklar.
- Sorunun muhtemel nedenini belirtir.
- Çözüm planı önerir.
- Komut önerir.
- Risk seviyesini açıklar.
- Kullanıcıya anlaşılır rapor üretir.

### 17.2 AI Ne Yapamaz?

AI tek başına şunları yapamaz:

- Komut çalıştıramaz.
- Dosya silemez.
- Servis restart edemez.
- Paket güncelleyemez.
- Firewall değiştiremez.
- Sistem kullanıcısı oluşturamaz/silemez.
- Policy engine’i bypass edemez.

### 17.3 AI Plan Formatı

AI’dan structured JSON response alınmalı.

```json
{
  "summary": "Disk usage is high. Journal logs are using significant space.",
  "risk_level": "medium",
  "requires_approval": true,
  "commands": [
    {
      "command": "journalctl --disk-usage",
      "purpose": "Check current journal log usage",
      "expected_effect": "Read-only",
      "risk_level": "low"
    },
    {
      "command": "journalctl --vacuum-time=7d",
      "purpose": "Remove journal logs older than 7 days",
      "expected_effect": "Modifies system logs",
      "risk_level": "medium"
    }
  ],
  "user_explanation": "This will clean old system journal logs. It will not delete application files."
}
```

Policy engine bu JSON’daki her komutu ayrıca kontrol etmeli.

---

## 18. Monitoring Engine

### 18.1 Çalışma Aralığı

Varsayılan:

```text
Her 30 saniyede bir lightweight metric collection.
Her 5 dakikada bir deeper scan.
Her 30 dakikada bir log/port/service baseline check.
```

### 18.2 Lightweight Metrics

- CPU
- RAM
- Disk
- Load average
- Uptime

### 18.3 Deeper Scan

- Top processes
- Open ports
- Active services
- Large log files
- Journal size

---

## 19. Alert Engine

### 19.1 Threshold Rules

```text
Disk > 85% → warning
Disk > 95% → critical
Memory > 90% → warning
CPU > 90% for 5 min → warning
New port detected → info/warning
Large log file detected → warning
Journal logs > configurable threshold → warning
```

### 19.2 Alert Deduplication

Aynı alert tekrar tekrar oluşturulmamalı.

Örnek:

Disk %90 olduğunda her 30 saniyede yeni alert açmak yerine mevcut açık disk alert güncellenmeli.

Kural:

```text
Same type + same target + open status = update existing alert
```

---

## 20. Güvenlik Gereksinimleri

### 20.1 Dashboard Access

Varsayılan:

```text
127.0.0.1:8787
```

Yani public internete açık olmayacak.

Public mode kullanıcı tarafından açılabilir ama uyarı gösterilmeli.

### 20.2 Authentication

- İlk girişte admin oluşturulmalı.
- Şifre hashlenmeli.
- Session/JWT kullanılmalı.
- Session secret otomatik generate edilmeli.
- Brute-force için basit rate limit uygulanmalı.

### 20.3 System Modification Protection

- Write action için approval zorunlu.
- Blocked komutlar asla çalışmamalı.
- Execution öncesi tekrar policy check yapılmalı.
- Timeout olmalı.
- Output size limiti olmalı.
- Audit log zorunlu olmalı.

### 20.4 File Permissions

```bash
chmod 755 /usr/local/bin/opsagent
chmod 700 /etc/opsagent
chmod 600 /etc/opsagent/config.yaml
chmod 700 /var/lib/opsagent
chmod 600 /var/lib/opsagent/opsagent.db
chmod 700 /var/log/opsagent
```

---

## 21. UI/UX Gereksinimleri

### 21.1 Temel Tasarım İlkeleri

- Developer tool hissi vermeli.
- Gereksiz animasyon olmamalı.
- Hızlı ve net olmalı.
- Komutlar açıkça görünmeli.
- Risk seviyesi net olmalı.
- “Onay olmadan değişiklik yapılmaz” mesajı düzenli görünmeli.
- Kritik aksiyonlarda kullanıcıyı aceleye getirmemeli.

### 21.2 Önemli UX Mesajları

```text
Read-only analysis completed. No system changes were made.
```

```text
This action modifies your system and requires approval.
```

```text
The following commands will be executed only after your approval.
```

```text
Execution completed. Full output was saved to audit log.
```

```text
This command was blocked by the safety policy.
```

---

## 22. Geliştirme Fazları

### Faz 1 — Project Bootstrap

Amaç: Go project yapısı, config, logging, SQLite migration.

Görevler:

```text
- Go module oluştur.
- Klasör yapısını kur.
- Config loader yaz.
- Logger kur.
- SQLite bağlantısı kur.
- Migration sistemi ekle.
- Basic health endpoint yaz.
```

Kabul kriterleri:

```text
opsagent binary çalışmalı.
Config okunmalı.
SQLite DB oluşturulmalı.
GET /api/system/service-status çalışmalı.
```

### Faz 2 — Installer + systemd

Amaç: Tek komutla kurulum.

Görevler:

```text
- install.sh yaz.
- OS/architecture check ekle.
- Binary download mantığı ekle.
- Dizinleri oluştur.
- Config oluştur.
- systemd service oluştur.
- Service enable/start yap.
- Başarılı kurulum mesajını göster.
- uninstall.sh yaz.
```

Kabul kriterleri:

```text
curl install komutu sonrası service active olmalı.
Dashboard localhost:8787 üzerinde açılmalı.
uninstall script sistemi temizlemeli.
```

### Faz 3 — Local Web Server + Setup Wizard

Amaç: Web panelin açılması ve ilk admin kurulumu.

Görevler:

```text
- Go HTTP server kur.
- React/Vite UI oluştur.
- UI build dosyalarını Go binary içine embed et.
- Setup required state kontrolü ekle.
- Admin create endpoint yaz.
- Login/logout yap.
- Session/JWT sistemi kur.
```

Kabul kriterleri:

```text
İlk açılışta setup wizard görünmeli.
Admin oluşturulmalı.
Login sonrası dashboard’a girilmeli.
```

### Faz 4 — Monitoring Engine

Amaç: Temel sistem verilerini toplamak.

Görevler:

```text
- CPU usage collector
- Memory usage collector
- Disk usage collector
- Load average collector
- Uptime collector
- Top process collector
- Open ports collector
- systemd services collector
- Journal disk usage collector
- Large log files scanner
```

Kabul kriterleri:

```text
Dashboard’da CPU/RAM/disk görünmeli.
Latest metrics API çalışmalı.
Top processes listelenmeli.
Open ports listelenmeli.
```

### Faz 5 — Alert Engine

Amaç: Threshold tabanlı alert üretmek.

Görevler:

```text
- Disk threshold rule
- Memory threshold rule
- CPU duration rule
- Large log file rule
- Journal usage rule
- New port detected rule
- Alert deduplication
- Alert status update
```

Kabul kriterleri:

```text
Disk yüksek olduğunda alert oluşmalı.
Aynı alert spam üretmemeli.
Alert list/detail ekranı çalışmalı.
```

### Faz 6 — Assistant + Basic Planner

Amaç: Kullanıcının prompt gönderebilmesi.

Görevler:

```text
- Assistant message endpoint
- Prompt intent detection
- Read-only command plan templates
- AI provider config
- OpenAI-compatible client
- AI disabled fallback mode
- Structured plan response parser
```

Kabul kriterleri:

```text
Kullanıcı “Disk neden doluyor?” yazınca read-only analiz yapılmalı.
Kullanıcı “Sunucuyu güncelle” yazınca approval gereken plan oluşmalı.
```

### Faz 7 — Command Policy Engine

Amaç: Güvenli komut sınıflandırması.

Görevler:

```text
- Whitelist tanımla.
- Denylist tanımla.
- Command classifier yaz.
- Risk level belirle.
- Requires approval belirle.
- Blocked command handling ekle.
```

Kabul kriterleri:

```text
Read-only komutlar onaysız çalışabilmeli.
Write komutlar approval istemeli.
Blocked komutlar asla çalışmamalı.
```

### Faz 8 — Approval + Execution

Amaç: Onaylı komut çalıştırma.

Görevler:

```text
- Task plan ekranı
- Approve endpoint
- Reject endpoint
- Command runner
- Timeout
- Output capture
- Exit code capture
- Execution result screen
- Audit log integration
```

Kabul kriterleri:

```text
Kullanıcı onay vermeden write komut çalışmamalı.
Onay sonrası komut çalışmalı.
Output ve exit code görünmeli.
Audit log yazılmalı.
```

### Faz 9 — Settings + Hardening

Amaç: MVP’yi kullanılabilir hale getirme.

Görevler:

```text
- AI settings ekranı
- Monitoring thresholds settings
- Dashboard access mode settings
- Password change
- Basic rate limiting
- File permission kontrolleri
- Service status ekranı
- Agent update placeholder
```

Kabul kriterleri:

```text
Kullanıcı AI key değiştirebilmeli.
Threshold değiştirebilmeli.
Settings kaydedilmeli.
```

---

## 23. MVP Test Senaryoları

### Test 1 — Kurulum

```text
Given clean Ubuntu server
When install.sh çalıştırılır
Then opsagent binary kurulmalı
And systemd service active olmalı
And dashboard localhost:8787’de açılmalı
```

### Test 2 — İlk Setup

```text
Given dashboard ilk kez açıldı
When admin oluşturulur
Then kullanıcı login olabilmeli
And setup tamamlandı olarak işaretlenmeli
```

### Test 3 — Monitoring

```text
Given agent çalışıyor
When 30 saniye geçer
Then metrics tablosuna kayıt yazılmalı
And dashboard son metrikleri göstermeli
```

### Test 4 — Disk Alert

```text
Given disk usage threshold üstünde
When alert engine çalışır
Then disk warning alert oluşmalı
And aynı alert tekrar spamlenmemeli
```

### Test 5 — Read-only Assistant

```text
Given kullanıcı "En çok RAM kullanan process ne?" yazar
When assistant çalışır
Then read-only komutlar çalışmalı
And kullanıcıya sonuç gösterilmeli
And approval istenmemeli
```

### Test 6 — Write Action Approval

```text
Given kullanıcı "Sunucuyu güncelle" yazar
When plan oluşturulur
Then apt update / apt upgrade komutları gösterilmeli
And approval required olmalı
And kullanıcı onaylamadan komut çalışmamalı
```

### Test 7 — Blocked Command

```text
Given AI veya kullanıcı rm -rf / komutu içeren plan oluşturur
When policy engine çalışır
Then komut blocked olmalı
And execution mümkün olmamalı
And audit log yazılmalı
```

---

## 24. MVP Kabul Kriterleri

MVP tamamlandı sayılması için:

```text
1. Tek komutla kurulum çalışmalı.
2. systemd service stabil çalışmalı.
3. Local dashboard açılmalı.
4. İlk admin hesabı oluşturulmalı.
5. CPU/RAM/disk metrikleri görünmeli.
6. Process/port/service bilgileri listelenmeli.
7. Disk/RAM/CPU alertleri oluşmalı.
8. Assistant prompt almalı.
9. Read-only analiz yapabilmeli.
10. Write action için plan ve komut listesi göstermeli.
11. Kullanıcı onayı olmadan write command çalışmamalı.
12. Denylist komutlar bloklanmalı.
13. Komut çıktıları gösterilmeli.
14. Audit log tutulmalı.
15. Settings üzerinden AI key ve threshold ayarlanmalı.
```

---

## 25. Developer İçin Öncelikli Implementasyon Sırası

Developer projeye şu sırayla başlamalı:

```text
1. Go project skeleton
2. Config + logger + SQLite
3. systemd compatible binary
4. install.sh
5. local web server
6. embedded UI
7. setup wizard + auth
8. metrics collection
9. dashboard
10. alert engine
11. assistant endpoint
12. AI planner
13. command policy engine
14. approval workflow
15. command executor
16. audit log
17. settings
18. hardening
```

Bu sıra önemli. Önce agent’ın lokal ve stabil çalışması kanıtlanmalı. AI kısmı en başta yapılmamalı.

---

## 26. Riskler ve Önlemler

| Risk | Önlem |
|---|---|
| Kullanıcı production server’a agent kurmaktan çekinir | Local-first, read-only default, approval workflow, audit log vurgulanmalı |
| AI tehlikeli komut önerir | Policy engine zorunlu olmalı |
| Dashboard public internete açılır | Default localhost only olmalı |
| Komut execution sunucuyu bozabilir | Approval + whitelist/denylist + risk label |
| Log tarama sistemi yorar | Büyük dosyalarda limit, timeout ve interval uygulanmalı |
| SQLite şişebilir | Metrics retention policy eklenmeli |
| Agent crash olur | systemd restart policy kullanılmalı |
| API key sızabilir | Config permission 600, mümkünse encrypted storage |

---

## 27. Retention / Cleanup Policy

SQLite zamanla şişmemeli.

MVP için öneri:

```text
Metrics: 7 gün sakla
Process snapshots: 24 saat sakla
Alerts: kalıcı sakla
Tasks: kalıcı sakla
Command executions: kalıcı sakla
Audit logs: kalıcı sakla
```

Settings’ten ileride değiştirilebilir.

---

## 28. systemd Service Taslağı

```ini
[Unit]
Description=OpsAgent - Local Linux AI Ops Assistant
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/opsagent serve --config /etc/opsagent/config.yaml
Restart=always
RestartSec=5
User=root
WorkingDirectory=/var/lib/opsagent

[Install]
WantedBy=multi-user.target
```

Not: MVP’de root olarak çalışması pratik olabilir. Sonraki sürümde daha güvenli `opsagent` user + sudo allowlist modeli düşünülmeli.

---

## 29. İlk Versiyon İçin Örnek Quick Actions

Dashboard’da hazır butonlar:

```text
Analyze disk usage
Check high CPU
Check memory usage
Show open ports
Analyze system logs
Check package updates
Check nginx status
Check docker containers
```

Bu butonlar kullanıcıyı prompt yazmaya mecbur bırakmaz.

---

## 30. Nihai MVP Tanımı

Bu MVP tamamlandığında ürün şu hale gelmiş olmalı:

> Kullanıcı Linux sunucusuna tek komutla OpsAgent kurar. Agent lokal olarak çalışır, kendi SQLite database’ini kullanır, web panelini kendi içinde sunar. Kullanıcı CPU/RAM/disk durumunu, servisleri, portları ve alertleri görür. Prompt yazarak sistem hakkında soru sorabilir. Sistem değişikliği gerektiren durumlarda agent önce çözüm planı ve çalıştırılacak komutları gösterir. Kullanıcı onay verirse komutlar çalışır. Tüm işlem çıktıları ve audit log kayıt altında tutulur.

---

## 31. Mimari Karar

Bu proje için önerilen mimari:

```text
Go single binary
+ SQLite
+ embedded React/Vite dashboard
+ systemd
+ local-only dashboard by default
+ AI provider key from user
+ approval-based command execution
+ audit-first security model
```

Clean Architecture kullanılabilir ama aşırı katmanlı hale getirilmemeli. En doğru yaklaşım:

> Pragmatic Clean Architecture.

Yani domain/usecase/infrastructure/web ayrımı olacak, fakat basit işler için gereksiz abstraction yazılmayacak.

Bu ürünün başarısı teknoloji kalabalığında değil; **kurulum kolaylığı, güvenli execution modeli ve kullanıcıya açık plan gösterme deneyiminde** yatıyor.

---

# 32. Kritik Kod Yapıları ve Skeleton Örnekleri

Bu bölümdeki kodlar production-ready final kod değildir. Developer’a mimari yön vermek için kritik yapıların skeleton örnekleridir. Kodlar geliştirirken test, error handling, logging, security ve edge-case kontrolleri genişletilmelidir.

---

## 32.1 Go Project Entry Point — `cmd/opsagent/main.go`

Amaç: Binary’nin `serve` komutu ile config okuyup application bootstrap etmesi.

```go
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"opsagent/internal/app"
	"opsagent/internal/config"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: opsagent serve --config /etc/opsagent/config.yaml")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "serve":
		serveCmd := flag.NewFlagSet("serve", flag.ExitOnError)
		configPath := serveCmd.String("config", "/etc/opsagent/config.yaml", "Path to config file")
		_ = serveCmd.Parse(os.Args[2:])

		cfg, err := config.Load(*configPath)
		if err != nil {
			log.Fatalf("failed to load config: %v", err)
		}

		application, err := app.New(cfg)
		if err != nil {
			log.Fatalf("failed to bootstrap app: %v", err)
		}

		if err := application.Run(); err != nil {
			log.Fatalf("app stopped with error: %v", err)
		}

	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}
```

---

## 32.2 Config Struct — `internal/config/config.go`

Amaç: YAML config dosyasını typed struct olarak okumak.

```go
package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	App struct {
		Name        string `yaml:"name"`
		Environment string `yaml:"environment"`
	} `yaml:"app"`

	Server struct {
		BindAddress string `yaml:"bind_address"`
		Port        int    `yaml:"port"`
	} `yaml:"server"`

	Database struct {
		Path string `yaml:"path"`
	} `yaml:"database"`

	Logging struct {
		Level string `yaml:"level"`
		File  string `yaml:"file"`
	} `yaml:"logging"`

	Security struct {
		RequireApprovalForWriteActions bool `yaml:"require_approval_for_write_actions"`
		CommandTimeoutSeconds          int  `yaml:"command_timeout_seconds"`
		MaxOutputSizeKB                int  `yaml:"max_output_size_kb"`
	} `yaml:"security"`

	AI struct {
		Enabled  bool   `yaml:"enabled"`
		Provider string `yaml:"provider"`
		APIKey   string `yaml:"api_key"`
		Model    string `yaml:"model"`
	} `yaml:"ai"`

	Monitoring struct {
		IntervalSeconds          int     `yaml:"interval_seconds"`
		DiskWarningThreshold     float64 `yaml:"disk_warning_threshold"`
		DiskCriticalThreshold    float64 `yaml:"disk_critical_threshold"`
		MemoryWarningThreshold   float64 `yaml:"memory_warning_threshold"`
		CPUWarningThreshold      float64 `yaml:"cpu_warning_threshold"`
		CPUWarningDurationSecond int     `yaml:"cpu_warning_duration_seconds"`
	} `yaml:"monitoring"`

	Dashboard struct {
		SessionSecret string `yaml:"session_secret"`
	} `yaml:"dashboard"`
}

func Load(path string) (*Config, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(raw, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
```

---

## 32.3 Application Bootstrap — `internal/app/bootstrap.go`

Amaç: Storage, HTTP server, scheduler ve worker parçalarını tek yerden başlatmak.

```go
package app

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"opsagent/internal/config"
	"opsagent/internal/infrastructure/storage/sqlite"
	"opsagent/internal/web"
)

type App struct {
	cfg    *config.Config
	server *http.Server
	store  *sqlite.Store
}

func New(cfg *config.Config) (*App, error) {
	store, err := sqlite.Open(cfg.Database.Path)
	if err != nil {
		return nil, err
	}

	if err := store.Migrate(); err != nil {
		return nil, err
	}

	router := web.NewRouter(cfg, store)

	addr := fmt.Sprintf("%s:%d", cfg.Server.BindAddress, cfg.Server.Port)
	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	return &App{
		cfg:    cfg,
		server: server,
		store:  store,
	}, nil
}

func (a *App) Run() error {
	ctx := context.Background()

	go func() {
		if err := StartScheduler(ctx, a.cfg, a.store); err != nil {
			log.Printf("scheduler stopped: %v", err)
		}
	}()

	log.Printf("OpsAgent dashboard listening on %s", a.server.Addr)
	return a.server.ListenAndServe()
}
```

---

## 32.4 Scheduler — `internal/app/scheduler.go`

Amaç: Lightweight metrics, deeper scan ve alert evaluation süreçlerini periyodik çalıştırmak.

```go
package app

import (
	"context"
	"log"
	"time"

	"opsagent/internal/config"
	"opsagent/internal/infrastructure/storage/sqlite"
	"opsagent/internal/usecase"
)

func StartScheduler(ctx context.Context, cfg *config.Config, store *sqlite.Store) error {
	metricsTicker := time.NewTicker(time.Duration(cfg.Monitoring.IntervalSeconds) * time.Second)
	deepScanTicker := time.NewTicker(5 * time.Minute)
	cleanupTicker := time.NewTicker(1 * time.Hour)

	defer metricsTicker.Stop()
	defer deepScanTicker.Stop()
	defer cleanupTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-metricsTicker.C:
			if err := usecase.CollectMetrics(ctx, store); err != nil {
				log.Printf("collect metrics failed: %v", err)
			}

			if err := usecase.EvaluateAlerts(ctx, cfg, store); err != nil {
				log.Printf("evaluate alerts failed: %v", err)
			}

		case <-deepScanTicker.C:
			if err := usecase.RunDeepScan(ctx, store); err != nil {
				log.Printf("deep scan failed: %v", err)
			}

		case <-cleanupTicker.C:
			if err := usecase.CleanupOldMetrics(ctx, store); err != nil {
				log.Printf("cleanup old metrics failed: %v", err)
			}
		}
	}
}
```

---

## 32.5 SQLite Store Skeleton — `internal/infrastructure/storage/sqlite/store.go`

Amaç: SQLite connection, migration ve repository fonksiyonlarını merkezi yönetmek.

```go
package sqlite

import (
	"database/sql"
	_ "modernc.org/sqlite"
)

type Store struct {
	DB *sql.DB
}

func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1)

	if _, err := db.Exec(`PRAGMA journal_mode=WAL;`); err != nil {
		return nil, err
	}

	if _, err := db.Exec(`PRAGMA foreign_keys=ON;`); err != nil {
		return nil, err
	}

	return &Store{DB: db}, nil
}

func (s *Store) Migrate() error {
	_, err := s.DB.Exec(initSQL)
	return err
}
```

---

## 32.6 Migration Embed — `internal/infrastructure/storage/sqlite/migrations.go`

Amaç: İlk migration’ı kod içine gömmek veya dosyadan çalıştırmak.

```go
package sqlite

const initSQL = `
CREATE TABLE IF NOT EXISTS users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  username TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  role TEXT NOT NULL DEFAULT 'admin',
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS metrics (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  cpu_usage REAL,
  memory_usage REAL,
  disk_usage REAL,
  load_average_1 REAL,
  load_average_5 REAL,
  load_average_15 REAL,
  uptime_seconds INTEGER,
  created_at DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS alerts (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  type TEXT NOT NULL,
  severity TEXT NOT NULL,
  title TEXT NOT NULL,
  message TEXT NOT NULL,
  status TEXT NOT NULL DEFAULT 'open',
  metadata_json TEXT,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS audit_logs (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  event_type TEXT NOT NULL,
  actor TEXT,
  message TEXT NOT NULL,
  metadata_json TEXT,
  created_at DATETIME NOT NULL
);
`
```

Not: Final projede tüm tablolar migration dosyalarında versioned şekilde tutulmalı. Bu skeleton minimum başlangıç içindir.

---

## 32.7 Domain Model — `internal/domain/command.go`

Amaç: Komut, risk, approval ve execution tarafındaki temel domain tiplerini netleştirmek.

```go
package domain

type CommandType string
type RiskLevel string
type CommandStatus string

const (
	CommandTypeRead      CommandType = "read"
	CommandTypeWrite     CommandType = "write"
	CommandTypeDangerous CommandType = "dangerous"
	CommandTypeBlocked   CommandType = "blocked"

	RiskLow      RiskLevel = "low"
	RiskMedium   RiskLevel = "medium"
	RiskHigh     RiskLevel = "high"
	RiskCritical RiskLevel = "critical"

	CommandPending  CommandStatus = "pending"
	CommandApproved CommandStatus = "approved"
	CommandRunning  CommandStatus = "running"
	CommandSuccess  CommandStatus = "success"
	CommandFailed   CommandStatus = "failed"
	CommandBlocked  CommandStatus = "blocked"
)

type Command struct {
	ID               int64
	TaskPlanID       int64
	Command          string
	Type             CommandType
	RiskLevel        RiskLevel
	RequiresApproval bool
	Status           CommandStatus
}
```

---

## 32.8 Command Policy Classifier — `internal/policy/classifier.go`

Amaç: AI veya kullanıcı tarafından üretilen her komutu execution öncesi sınıflandırmak.

```go
package policy

import (
	"strings"

	"opsagent/internal/domain"
)

type Result struct {
	Type             domain.CommandType
	RiskLevel        domain.RiskLevel
	RequiresApproval bool
	Allowed          bool
	Reason           string
}

var blockedPatterns = []string{
	"rm -rf /",
	"mkfs",
	"dd if=",
	"chmod -R 777 /",
	"chown -R /",
	"userdel",
	"passwd",
	"iptables -F",
	"ufw disable",
	"shutdown",
	"reboot",
	"curl ",
	"wget ",
	"| bash",
}

var writePrefixes = []string{
	"apt update",
	"apt upgrade",
	"apt install",
	"systemctl restart",
	"systemctl reload",
	"journalctl --vacuum-time",
	"docker restart",
}

var readPrefixes = []string{
	"df ",
	"free ",
	"uptime",
	"ps ",
	"top ",
	"systemctl status",
	"journalctl ",
	"ss ",
	"cat /etc/os-release",
	"apt list",
	"du ",
	"docker ps",
	"docker logs",
}

func Classify(command string) Result {
	normalized := strings.TrimSpace(command)

	lower := strings.ToLower(normalized)

	for _, pattern := range blockedPatterns {
		if strings.Contains(lower, strings.ToLower(pattern)) {
			return Result{
				Type:             domain.CommandTypeBlocked,
				RiskLevel:        domain.RiskCritical,
				RequiresApproval: false,
				Allowed:          false,
				Reason:           "Command matched blocked safety policy",
			}
		}
	}

	for _, prefix := range writePrefixes {
		if strings.HasPrefix(lower, strings.ToLower(prefix)) {
			return Result{
				Type:             domain.CommandTypeWrite,
				RiskLevel:        domain.RiskMedium,
				RequiresApproval: true,
				Allowed:          true,
				Reason:           "Write operation requires approval",
			}
		}
	}

	for _, prefix := range readPrefixes {
		if strings.HasPrefix(lower, strings.ToLower(prefix)) {
			return Result{
				Type:             domain.CommandTypeRead,
				RiskLevel:        domain.RiskLow,
				RequiresApproval: false,
				Allowed:          true,
				Reason:           "Read-only command",
			}
		}
	}

	return Result{
		Type:             domain.CommandTypeDangerous,
		RiskLevel:        domain.RiskHigh,
		RequiresApproval: true,
		Allowed:          false,
		Reason:           "Command is not whitelisted",
	}
}
```

Önemli not: MVP’de bilinmeyen komutlar otomatik çalıştırılmamalıdır. Whitelist dışında kalan komutlar varsayılan olarak engellenmelidir.

---

## 32.9 Command Runner — `internal/infrastructure/executor/command_runner.go`

Amaç: Komutları timeout, output limit ve policy check ile güvenli çalıştırmak.

```go
package executor

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"time"
)

type Result struct {
	Stdout   string
	Stderr   string
	ExitCode int
	Duration time.Duration
}

type Runner struct {
	Timeout       time.Duration
	MaxOutputSize int
}

func NewRunner(timeoutSeconds int, maxOutputSizeKB int) *Runner {
	return &Runner{
		Timeout:       time.Duration(timeoutSeconds) * time.Second,
		MaxOutputSize: maxOutputSizeKB * 1024,
	}
}

func (r *Runner) Run(ctx context.Context, command string) (*Result, error) {
	start := time.Now()

	ctx, cancel := context.WithTimeout(ctx, r.Timeout)
	defer cancel()

	// Güvenlik için shell kullanımını minimumda tutmak gerekir.
	// MVP’de kontrollü command string için bash -lc kullanılabilir,
	// fakat policy engine execution öncesi mutlaka çalışmalıdır.
	cmd := exec.CommandContext(ctx, "bash", "-lc", command)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &limitedBuffer{
		Buffer: &stdout,
		Limit:  r.MaxOutputSize,
	}
	cmd.Stderr = &limitedBuffer{
		Buffer: &stderr,
		Limit:  r.MaxOutputSize,
	}

	err := cmd.Run()

	if ctx.Err() == context.DeadlineExceeded {
		return nil, errors.New("command timed out")
	}

	exitCode := 0
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		} else {
			return nil, err
		}
	}

	return &Result{
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		ExitCode: exitCode,
		Duration: time.Since(start),
	}, nil
}

type limitedBuffer struct {
	Buffer *bytes.Buffer
	Limit  int
}

func (b *limitedBuffer) Write(p []byte) (int, error) {
	if b.Buffer.Len()+len(p) > b.Limit {
		remaining := b.Limit - b.Buffer.Len()
		if remaining > 0 {
			b.Buffer.Write(p[:remaining])
		}
		return len(p), nil
	}

	return b.Buffer.Write(p)
}
```

Güvenlik notu: `bash -lc` esneklik sağlar ama risklidir. Bu yüzden execution öncesi policy check zorunlu olmalıdır. Daha güvenli V2’de komutlar binary + args olarak ayrıştırılıp shell kullanılmadan çalıştırılmalıdır.

---

## 32.10 Execute Approved Command Usecase — `internal/usecase/execute_command.go`

Amaç: Komut çalıştırmadan önce approval ve policy check doğrulamak.

```go
package usecase

import (
	"context"
	"errors"

	"opsagent/internal/config"
	"opsagent/internal/infrastructure/executor"
	"opsagent/internal/infrastructure/storage/sqlite"
	"opsagent/internal/policy"
)

func ExecuteApprovedCommand(ctx context.Context, cfg *config.Config, store *sqlite.Store, commandID int64) error {
	cmd, err := store.GetCommandByID(ctx, commandID)
	if err != nil {
		return err
	}

	policyResult := policy.Classify(cmd.Command)
	if !policyResult.Allowed {
		_ = store.MarkCommandBlocked(ctx, commandID, policyResult.Reason)
		_ = store.InsertAuditLog(ctx, "command_blocked", "system", policyResult.Reason, nil)
		return errors.New(policyResult.Reason)
	}

	if policyResult.RequiresApproval && cmd.Status != "approved" {
		return errors.New("command requires approval")
	}

	runner := executor.NewRunner(
		cfg.Security.CommandTimeoutSeconds,
		cfg.Security.MaxOutputSizeKB,
	)

	_ = store.MarkCommandRunning(ctx, commandID)

	result, err := runner.Run(ctx, cmd.Command)
	if err != nil {
		_ = store.SaveCommandExecution(ctx, commandID, "", err.Error(), -1, "failed")
		_ = store.MarkCommandFailed(ctx, commandID)
		return err
	}

	status := "success"
	if result.ExitCode != 0 {
		status = "failed"
	}

	if err := store.SaveCommandExecution(ctx, commandID, result.Stdout, result.Stderr, result.ExitCode, status); err != nil {
		return err
	}

	if status == "success" {
		_ = store.MarkCommandSuccess(ctx, commandID)
	} else {
		_ = store.MarkCommandFailed(ctx, commandID)
	}

	_ = store.InsertAuditLog(ctx, "command_executed", "admin", "Command executed", map[string]any{
		"command_id": commandID,
		"exit_code":  result.ExitCode,
		"status":     status,
	})

	return nil
}
```

---

## 32.11 Metrics Collector Skeleton — `internal/infrastructure/system/metrics.go`

Amaç: Linux sistem metriklerini toplamak.

```go
package system

import (
	"os"
	"strconv"
	"strings"
)

type Metrics struct {
	CPUUsage      float64
	MemoryUsage   float64
	DiskUsage     float64
	LoadAverage1  float64
	LoadAverage5  float64
	LoadAverage15 float64
	UptimeSeconds int64
}

func CollectMetrics() (*Metrics, error) {
	load1, load5, load15, _ := readLoadAverage()
	uptime, _ := readUptimeSeconds()

	// CPU, memory ve disk için MVP’de gopsutil gibi bir kütüphane kullanılabilir.
	// Alternatif olarak /proc ve syscall üzerinden manuel okunabilir.

	return &Metrics{
		CPUUsage:      0,
		MemoryUsage:   0,
		DiskUsage:     0,
		LoadAverage1:  load1,
		LoadAverage5:  load5,
		LoadAverage15: load15,
		UptimeSeconds: uptime,
	}, nil
}

func readLoadAverage() (float64, float64, float64, error) {
	raw, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		return 0, 0, 0, err
	}

	parts := strings.Fields(string(raw))
	if len(parts) < 3 {
		return 0, 0, 0, nil
	}

	l1, _ := strconv.ParseFloat(parts[0], 64)
	l5, _ := strconv.ParseFloat(parts[1], 64)
	l15, _ := strconv.ParseFloat(parts[2], 64)

	return l1, l5, l15, nil
}

func readUptimeSeconds() (int64, error) {
	raw, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return 0, err
	}

	parts := strings.Fields(string(raw))
	if len(parts) == 0 {
		return 0, nil
	}

	uptimeFloat, _ := strconv.ParseFloat(parts[0], 64)
	return int64(uptimeFloat), nil
}
```

Not: MVP’de sistem metrikleri için `gopsutil` kullanmak geliştirmeyi hızlandırır. Daha sonra ihtiyaç olursa dependency azaltılabilir.

---

## 32.12 Alert Evaluation Usecase — `internal/usecase/evaluate_alerts.go`

Amaç: Son metriklere göre threshold alert üretmek ve spam/deduplication kontrolü yapmak.

```go
package usecase

import (
	"context"
	"fmt"

	"opsagent/internal/config"
	"opsagent/internal/infrastructure/storage/sqlite"
)

func EvaluateAlerts(ctx context.Context, cfg *config.Config, store *sqlite.Store) error {
	metric, err := store.GetLatestMetric(ctx)
	if err != nil {
		return err
	}

	if metric.DiskUsage >= cfg.Monitoring.DiskCriticalThreshold {
		return store.UpsertOpenAlert(ctx, sqlite.AlertInput{
			Type:     "disk_usage",
			Severity: "critical",
			Title:    "Disk usage is critical",
			Message:  fmt.Sprintf("Disk usage reached %.1f%%", metric.DiskUsage),
		})
	}

	if metric.DiskUsage >= cfg.Monitoring.DiskWarningThreshold {
		return store.UpsertOpenAlert(ctx, sqlite.AlertInput{
			Type:     "disk_usage",
			Severity: "warning",
			Title:    "Disk usage is high",
			Message:  fmt.Sprintf("Disk usage reached %.1f%%", metric.DiskUsage),
		})
	}

	if metric.MemoryUsage >= cfg.Monitoring.MemoryWarningThreshold {
		return store.UpsertOpenAlert(ctx, sqlite.AlertInput{
			Type:     "memory_usage",
			Severity: "warning",
			Title:    "Memory usage is high",
			Message:  fmt.Sprintf("Memory usage reached %.1f%%", metric.MemoryUsage),
		})
	}

	return nil
}
```

---

## 32.13 Web Router Skeleton — `internal/web/server.go`

Amaç: API route’larını ve embedded UI serving’i başlatmak.

```go
package web

import (
	"net/http"

	"opsagent/internal/config"
	"opsagent/internal/infrastructure/storage/sqlite"
)

func NewRouter(cfg *config.Config, store *sqlite.Store) http.Handler {
	mux := http.NewServeMux()

	// System
	mux.HandleFunc("GET /api/system/service-status", handleServiceStatus())
	mux.HandleFunc("GET /api/dashboard/summary", handleDashboardSummary(store))

	// Auth
	mux.HandleFunc("POST /api/setup/admin", handleCreateAdmin(store))
	mux.HandleFunc("POST /api/auth/login", handleLogin(store, cfg))

	// Assistant
	mux.HandleFunc("POST /api/assistant/message", handleAssistantMessage(cfg, store))

	// Tasks
	mux.HandleFunc("GET /api/tasks", handleTaskList(store))
	mux.HandleFunc("POST /api/tasks/approve", handleApproveTask(store))
	mux.HandleFunc("POST /api/tasks/run", handleRunTask(cfg, store))

	// Static UI
	mux.Handle("/", ServeEmbeddedUI())

	return withSecurityHeaders(mux)
}

func withSecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "no-referrer")
		next.ServeHTTP(w, r)
	})
}
```

---

## 32.14 Embedded UI — `internal/ui/embed.go`

Amaç: React/Vite build dosyalarını Go binary içine gömmek.

```go
package ui

import "embed"

//go:embed dist/*
var Files embed.FS
```

Web serving örneği:

```go
package web

import (
	"io/fs"
	"net/http"

	"opsagent/internal/ui"
)

func ServeEmbeddedUI() http.Handler {
	dist, err := fs.Sub(ui.Files, "dist")
	if err != nil {
		panic(err)
	}

	fileServer := http.FileServer(http.FS(dist))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// SPA fallback.
		// Eğer dosya bulunamazsa index.html dönülebilir.
		fileServer.ServeHTTP(w, r)
	})
}
```

---

## 32.15 AI Plan Struct — `internal/infrastructure/ai/types.go`

Amaç: AI’dan structured JSON formatında plan almak.

```go
package ai

type Plan struct {
	Summary         string        `json:"summary"`
	RiskLevel       string        `json:"risk_level"`
	RequiresApproval bool        `json:"requires_approval"`
	Commands        []PlanCommand `json:"commands"`
	UserExplanation  string       `json:"user_explanation"`
}

type PlanCommand struct {
	Command        string `json:"command"`
	Purpose        string `json:"purpose"`
	ExpectedEffect string `json:"expected_effect"`
	RiskLevel      string `json:"risk_level"`
}
```

---

## 32.16 OpenAI-Compatible Client Skeleton — `internal/infrastructure/ai/openai_compatible.go`

Amaç: OpenAI-compatible API endpoint’ine sistem durumu ve kullanıcı promptunu gönderip structured plan almak.

```go
package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

type Client struct {
	APIKey  string
	BaseURL string
	Model   string
}

func (c *Client) CreatePlan(ctx context.Context, userPrompt string, systemContext map[string]any) (*Plan, error) {
	payload := map[string]any{
		"model": c.Model,
		"messages": []map[string]string{
			{
				"role": "system",
				"content": `You are OpsAgent planner. You must only return safe structured JSON.
Never execute commands.
Never bypass approval.
Every command must be explicit.
Dangerous destructive commands must not be suggested.`,
			},
			{
				"role": "user",
				"content": userPrompt,
			},
		},
		"response_format": map[string]string{
			"type": "json_object",
		},
	}

	raw, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/chat/completions", bytes.NewReader(raw))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode >= 300 {
		return nil, errors.New("ai provider returned non-success status")
	}

	// Skeleton:
	// Final implementation must parse provider response,
	// extract message.content,
	// unmarshal it into Plan struct,
	// then pass all commands through policy engine.

	return nil, errors.New("not implemented")
}
```

---

## 32.17 Assistant Message Flow — `internal/usecase/create_task_plan.go`

Amaç: Kullanıcı promptundan task/plan üretmek, read-only ve write aksiyon ayrımını yapmak.

```go
package usecase

import (
	"context"

	"opsagent/internal/config"
	"opsagent/internal/infrastructure/storage/sqlite"
	"opsagent/internal/policy"
)

func CreateTaskPlanFromPrompt(ctx context.Context, cfg *config.Config, store *sqlite.Store, prompt string) (int64, error) {
	taskID, err := store.CreateTask(ctx, "assistant", prompt, "draft")
	if err != nil {
		return 0, err
	}

	// MVP fallback:
	// AI kapalıysa bazı promptlar template bazlı cevaplanabilir.
	// Örnek: disk neden doluyor -> df/du/journalctl read-only planı.
	plan := buildFallbackPlan(prompt)

	requiresApproval := false

	for i := range plan.Commands {
		result := policy.Classify(plan.Commands[i].Command)

		if !result.Allowed {
			plan.Commands[i].CommandType = "blocked"
			requiresApproval = false
			continue
		}

		if result.RequiresApproval {
			requiresApproval = true
		}

		plan.Commands[i].CommandType = string(result.Type)
		plan.Commands[i].RiskLevel = string(result.RiskLevel)
		plan.Commands[i].RequiresApproval = result.RequiresApproval
	}

	status := "completed"
	if requiresApproval {
		status = "waiting_approval"
	}

	planID, err := store.CreateTaskPlan(ctx, taskID, plan.Summary, plan.RiskLevel, requiresApproval)
	if err != nil {
		return 0, err
	}

	for _, cmd := range plan.Commands {
		if err := store.CreateCommand(ctx, planID, cmd); err != nil {
			return 0, err
		}
	}

	if err := store.UpdateTaskStatus(ctx, taskID, status); err != nil {
		return 0, err
	}

	return taskID, nil
}
```

---

## 32.18 install.sh Skeleton — `scripts/install.sh`

Amaç: Tek komutla kurulum akışını sağlamak.

```bash
#!/usr/bin/env bash
set -euo pipefail

APP_NAME="opsagent"
INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="/etc/opsagent"
DATA_DIR="/var/lib/opsagent"
LOG_DIR="/var/log/opsagent"
SERVICE_FILE="/etc/systemd/system/opsagent.service"
BINARY_URL_BASE="https://opsagent.dev/releases/latest"

if [ "$EUID" -ne 0 ]; then
  echo "Please run as root or with sudo."
  exit 1
fi

echo "OpsAgent Installer"

OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$ARCH" in
  x86_64) ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *)
    echo "Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

if [ "$OS" != "linux" ]; then
  echo "Unsupported OS: $OS"
  exit 1
fi

echo "[1/7] Detected OS: $OS"
echo "[2/7] Detected architecture: $ARCH"

mkdir -p "$CONFIG_DIR" "$DATA_DIR" "$LOG_DIR"

echo "[3/7] Downloading OpsAgent binary..."
curl -fsSL "$BINARY_URL_BASE/opsagent-linux-$ARCH" -o "$INSTALL_DIR/$APP_NAME"
chmod 755 "$INSTALL_DIR/$APP_NAME"

echo "[4/7] Creating config..."
if [ ! -f "$CONFIG_DIR/config.yaml" ]; then
  SESSION_SECRET="$(openssl rand -hex 32 || date +%s | sha256sum | cut -d' ' -f1)"

  cat > "$CONFIG_DIR/config.yaml" <<EOF
app:
  name: "OpsAgent"
  environment: "production"

server:
  bind_address: "127.0.0.1"
  port: 8787

database:
  path: "/var/lib/opsagent/opsagent.db"

logging:
  level: "info"
  file: "/var/log/opsagent/opsagent.log"

security:
  require_approval_for_write_actions: true
  command_timeout_seconds: 120
  max_output_size_kb: 1024

ai:
  enabled: false
  provider: ""
  api_key: ""
  model: ""

monitoring:
  interval_seconds: 30
  disk_warning_threshold: 85
  disk_critical_threshold: 95
  memory_warning_threshold: 90
  cpu_warning_threshold: 90
  cpu_warning_duration_seconds: 300

dashboard:
  session_secret: "$SESSION_SECRET"
EOF
fi

chmod 700 "$CONFIG_DIR"
chmod 600 "$CONFIG_DIR/config.yaml"
chmod 700 "$DATA_DIR"
chmod 700 "$LOG_DIR"

echo "[5/7] Creating systemd service..."
cat > "$SERVICE_FILE" <<EOF
[Unit]
Description=OpsAgent - Local Linux AI Ops Assistant
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/opsagent serve --config /etc/opsagent/config.yaml
Restart=always
RestartSec=5
User=root
WorkingDirectory=/var/lib/opsagent

[Install]
WantedBy=multi-user.target
EOF

echo "[6/7] Starting service..."
systemctl daemon-reload
systemctl enable opsagent
systemctl restart opsagent

echo "[7/7] Checking service status..."
sleep 1
systemctl is-active --quiet opsagent

echo ""
echo "OpsAgent installed successfully."
echo ""
echo "Service:"
echo "active"
echo ""
echo "Dashboard:"
echo "http://127.0.0.1:8787"
echo ""
echo "To access the dashboard from your computer:"
echo ""
echo "ssh -L 8787:localhost:8787 root@YOUR_SERVER_IP"
echo ""
echo "Then open:"
echo ""
echo "http://localhost:8787"
echo ""
echo "Your server data stays on this machine."
echo "No write operation will run without approval."
```

---

## 32.19 uninstall.sh Skeleton — `scripts/uninstall.sh`

Amaç: Kullanıcının agent’ı temiz şekilde kaldırabilmesi.

```bash
#!/usr/bin/env bash
set -euo pipefail

if [ "$EUID" -ne 0 ]; then
  echo "Please run as root or with sudo."
  exit 1
fi

echo "Stopping OpsAgent..."
systemctl stop opsagent || true
systemctl disable opsagent || true

echo "Removing service..."
rm -f /etc/systemd/system/opsagent.service
systemctl daemon-reload

echo "Removing binary..."
rm -f /usr/local/bin/opsagent

read -p "Remove config, database and logs? This cannot be undone. [y/N] " confirm
if [[ "$confirm" == "y" || "$confirm" == "Y" ]]; then
  rm -rf /etc/opsagent
  rm -rf /var/lib/opsagent
  rm -rf /var/log/opsagent
  echo "All OpsAgent files removed."
else
  echo "Config, database and logs were kept."
fi

echo "OpsAgent uninstalled."
```

---

## 32.20 React/Vite UI Minimum Route Plan

Amaç: UI tarafının hangi sayfalardan oluşacağını netleştirmek.

```text
web/dashboard/src/
├── app/
│   ├── App.tsx
│   ├── router.tsx
│   └── api.ts
│
├── pages/
│   ├── SetupPage.tsx
│   ├── LoginPage.tsx
│   ├── DashboardPage.tsx
│   ├── AlertsPage.tsx
│   ├── AlertDetailPage.tsx
│   ├── AssistantPage.tsx
│   ├── TaskApprovalPage.tsx
│   ├── AuditLogPage.tsx
│   └── SettingsPage.tsx
│
├── components/
│   ├── MetricCard.tsx
│   ├── AlertBadge.tsx
│   ├── CommandList.tsx
│   ├── RiskBadge.tsx
│   ├── ApprovalBox.tsx
│   └── ShellOutput.tsx
│
└── main.tsx
```

---

## 32.21 UI API Client Skeleton — `web/dashboard/src/app/api.ts`

Amaç: Frontend’in local API ile konuşması.

```ts
const API_BASE = "/api";

async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...(options.headers || {}),
    },
    credentials: "include",
  });

  if (!res.ok) {
    const text = await res.text();
    throw new Error(text || `Request failed: ${res.status}`);
  }

  return res.json();
}

export const api = {
  getDashboardSummary: () => request("/dashboard/summary"),

  sendAssistantMessage: (message: string) =>
    request("/assistant/message", {
      method: "POST",
      body: JSON.stringify({ message }),
    }),

  getTaskPlan: (taskId: number) =>
    request(`/tasks/${taskId}/plan`),

  approveTask: (taskId: number) =>
    request(`/tasks/${taskId}/approve`, {
      method: "POST",
    }),

  rejectTask: (taskId: number) =>
    request(`/tasks/${taskId}/reject`, {
      method: "POST",
    }),
};
```

---

## 32.22 Task Approval UI Component Skeleton

Amaç: Kullanıcı komutları görmeden onay verememeli.

```tsx
type Command = {
  id: number;
  command: string;
  command_type: "read" | "write" | "dangerous" | "blocked";
  risk_level: "low" | "medium" | "high" | "critical";
  requires_approval: boolean;
};

type Props = {
  taskTitle: string;
  summary: string;
  riskLevel: string;
  commands: Command[];
  onApprove: () => void;
  onReject: () => void;
};

export function TaskApprovalBox({
  taskTitle,
  summary,
  riskLevel,
  commands,
  onApprove,
  onReject,
}: Props) {
  const hasBlockedCommand = commands.some((cmd) => cmd.command_type === "blocked");

  return (
    <section className="rounded-xl border p-6">
      <h2 className="text-xl font-semibold">{taskTitle}</h2>
      <p className="mt-2 text-sm text-muted-foreground">{summary}</p>

      <div className="mt-4">
        <strong>Risk:</strong> {riskLevel}
      </div>

      <div className="mt-6">
        <h3 className="font-medium">Commands to be executed</h3>

        <div className="mt-3 space-y-3">
          {commands.map((cmd) => (
            <div key={cmd.id} className="rounded-lg border bg-muted p-3">
              <div className="mb-2 text-xs">
                {cmd.command_type} · {cmd.risk_level}
              </div>
              <pre className="overflow-x-auto text-sm">{cmd.command}</pre>
            </div>
          ))}
        </div>
      </div>

      {hasBlockedCommand && (
        <div className="mt-4 rounded-lg border border-red-500 p-3 text-sm">
          This plan contains blocked commands and cannot be executed.
        </div>
      )}

      <div className="mt-6 flex gap-3">
        <button disabled={hasBlockedCommand} onClick={onApprove}>
          Approve and Run
        </button>
        <button onClick={onReject}>
          Reject
        </button>
      </div>
    </section>
  );
}
```

---

## 32.23 Build Pipeline Taslağı

Amaç: UI build edildikten sonra Go binary içine gömülüp release binary üretilmesi.

```bash
#!/usr/bin/env bash
set -euo pipefail

echo "Building dashboard UI..."
cd web/dashboard
npm install
npm run build
cd ../..

echo "Building Go binary..."
GOOS=linux GOARCH=amd64 go build -o dist/opsagent-linux-amd64 ./cmd/opsagent
GOOS=linux GOARCH=arm64 go build -o dist/opsagent-linux-arm64 ./cmd/opsagent

echo "Build completed."
```

---

## 32.24 Minimum README Kurulum Bölümü

Amaç: Developer’ın repo içinde kullanacağı temel dokümantasyon.

```md
# OpsAgent

Local-first Linux AI Ops Agent.

## Install

```bash
curl -fsSL https://opsagent.dev/install.sh | sudo bash
```

## Access Dashboard

```bash
ssh -L 8787:localhost:8787 root@YOUR_SERVER_IP
```

Open:

```text
http://localhost:8787
```

## Safety

OpsAgent is read-only by default. Any command that modifies the system requires explicit user approval.
```

---

# 33. Developer’a Net Uygulama Talimatı

Bu projede önce UI veya AI’dan başlanmamalı. İlk hedef çalışan local agent olmalı.

Öncelik sırası:

```text
1. Go binary çalışsın.
2. Config okusun.
3. SQLite oluştursun.
4. systemd ile servis olarak çalışsın.
5. Localhost panel endpoint’i açılsın.
6. Metrik toplasın.
7. Alert üretsin.
8. Setup/login çalışsın.
9. Assistant ve plan sistemi eklensin.
10. Policy engine zorunlu hale getirilsin.
11. Approval olmadan write command çalışmasın.
12. Audit log eksiksiz olsun.
```

AI tarafı, güvenlik ve policy engine olmadan asla command execution’a bağlanmamalıdır.

---

# 34. Final Not

Bu ürünün ana farkı, AI kullanması değil; AI’ı production Linux sunucuda **kontrollü, onaylı ve audit edilebilir** hale getirmesidir.

Bu yüzden MVP’de en çok dikkat edilmesi gereken üç şey:

```text
1. Tek komutla sorunsuz kurulum.
2. Local-first güvenli çalışma.
3. Approval-based command execution.
```

Bu üçü iyi yapılırsa ürünün temeli doğru kurulmuş olur.
