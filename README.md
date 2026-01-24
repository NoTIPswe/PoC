# Proof of Concept - Progetto SWE UNIPD

L'obiettivo primario di questo PoC non è realizzare una "versione in miniatura" del prodotto finale, bensì validare la fattibilità tecnica e l'integrazione delle tecnologie critiche selezionate per l'architettura.
Vogliamo dimostrare che lo stack tecnologico scelto è un *good fit* per i requisiti di performance e scalabilità, e che il flusso dei dati (dall'ingestion alla visualizzazione sicura) è sostenibile. Ci concentreremo sulla "Hot Path" dei dati, tralasciando feature accessorie o infrastrutture di contorno che potrebbero variare nell'MVP finale.

## Cosa c'è (Perimetro del PoC)

### 1. Gateway Simulator

**Linguaggio:** Go
**Ruolo:**
Il simulatore genera traffico sintetico emulando molteplici dispositivi connessi simultaneamente. Si occupa di generare payload di telemetria, cifrarli alla fonte (mockando il comportamento E2EE con chiavi statiche) e inviarli tramite protocollo MQTTS.
**Perché:**
Go e le sue goroutines ci permettono di simulare un alto numero di connessioni concorrenti con un basso footprint di risorse, stressando il sistema come farebbe una flotta reale.

### 2. Ingestion Service

**Linguaggio:** Go
**Ruolo:**
È l'unico punto di ingresso esposto ai gateway. Non decifra i dati (che rimangono opachi), ma si occupa di validare il protocollo e inoltrare il messaggio sulla coda NATS nel topic corretto.
**Perché:**
Protegge il core del sistema e normalizza i dati in ingresso. Go garantisce la velocità necessaria per non diventare un collo di bottiglia durante i picchi di ingestion.

### 3. Data Queue

**Tecnologia:** NATS JetStream
**Ruolo:**
Agisce come buffer ad alte prestazioni tra l'Ingestion e la scrittura su DB. Garantisce che nessun dato venga perso se il database è sotto stress o se il Consumer è momentaneamente offline.
**Perché:**
Essenziale per disaccoppiare la velocità di ricezione (alta) dai tempi di scrittura su disco. JetStream è stato scelto per la sua leggerezza e facilità operativa rispetto ad alternative più pesanti come Kafka.

### 4. Data Consumer

**Linguaggio:** Go
**Ruolo:**
Worker che sottoscrive allo stream di NATS, preleva i messaggi in arrivo ed esegue operazioni di bulk write sul database per massimizzare il throughput.
**Perché:**
Separare la logica di scrittura da quella di lettura e ingestion permette di scalare le componenti indipendentemente (pattern CQRS).

### 5. Measures Database

**Tecnologia:** PostgreSQL + TimescaleDB
**Ruolo:**
Storage persistente delle serie temporali. I dati verranno salvati in **Hypertables** partizionate, includendo una colonna `tenant_id` per garantire la segregazione logica dei dati fin dal livello di storage.
**Perché:**
TimescaleDB offre le performance di un TSDB verticale mantenendo l'affidabilità e l'ecosistema SQL di PostgreSQL.

### 6. Data API

**Linguaggio:** TypeScript - NestJS
**Ruolo:**
Espone endpoint **Read-Only** per la dashboard. Il suo compito principale è interrogare il database applicando filtri rigorosi per l'isolamento dei tenant (es. `WHERE tenant_id = X`), basandosi su un'autenticazione mockata.
**Perché:**
NestJS offre una struttura modulare e manutenibile per la logica di business e di accesso ai dati.

### 7. Web Dashboard & Client SDK

**Tecnologia:** Angular + TypeScript SDK
**Ruolo:**
Interfaccia di visualizzazione per l'utente. La caratteristica chiave è che la logica di decifrazione avviene qui (lato client): l'applicazione riceve i payload criptati dalle API e utilizza una chiave (simulata/statica nel PoC) per renderli leggibili all'utente.
**Perché:**
Dimostra la fattibilità dell'approccio End-to-End Encryption (Zero Knowledge per il cloud), dove il dato in chiaro è visibile solo ai bordi del sistema (Gateway e Utente).

# Cosa non c'è

Per motivi di tempo e focalizzazione, il PoC esclude le parti non critiche per la validazione del flusso dati principale:

* **Autenticazione reale:** Gateway e Utenti useranno un sistema di *Mock Auth* (es. header statici) invece di integrare Keycloak.
* **Gestione chiavi dinamica:** Non ci sarà una CA interna o un Key Vault; useremo chiavi simmetriche statiche per dimostrare il concetto di E2EE.
* **API Esterne:** Non esporremo API per integratori terzi.
* **Segregazione multi-tenant "hard":** Useremo una segregazione logica (colonna nel DB), riservandoci di valutare schemi più complessi (es. RLS o schema-per-tenant) per l'MVP.
* **Testing & Code Quality:** Il codice sarà funzionale alla dimostrazione ("quick & dirty" dove necessario), senza copertura di test estensiva.
* **Observability:** Non implementeremo stack di monitoraggio (Prometheus/Grafana) né sistemi di logging centralizzato per auditing.
* **Sistema di Notifiche:** Nessun alert o motore di regole in questa fase.
* **Comandi ai Gateway:** Il flusso sarà monodirezionale (Gateway -> Cloud).
* **Interfaccia Simulazione:** La simulazione sarà configurata via codice/config file, senza una UI dedicata.

*Nota: Nel caso in cui il PoC richiedesse minor tempo rispetto al previsto, valuteremo l'introduzione anticipata di Keycloak o dell'Observability stack.*