# Proof of Concept - Progetto SWE UNIPD

L'obiettivo primario di questo PoC non è realizzare una "versione in miniatura" del prodotto finale, bensì validare la fattibilità tecnica e l'integrazione delle tecnologie critiche selezionate per l'architettura.
Vogliamo dimostrare che lo stack tecnologico scelto è un _good fit_ per i requisiti di performance e scalabilità, e che il flusso dei dati (dall'ingestion alla visualizzazione sicura) è sostenibile. Ci concentreremo sulla "Hot Path" dei dati, tralasciando feature accessorie o infrastrutture di contorno che potrebbero variare nell'MVP finale.

## Cosa c'è (Perimetro del PoC)

### 1. Gateway Simulator

- **Linguaggio:** Go.
- **Ruolo:**
  Il simulatore genera traffico sintetico emulando molteplici dispositivi connessi simultaneamente. Si occupa di generare payload di telemetria, cifrarli alla fonte (utilizzando chiavi statiche definite in configurazione) e inviarli tramite protocollo proprietario di **NATS**, autenticandosi tramite **certificati client (mTLS)**.
- **Perché:**
  Go e le sue goroutines permettono di simulare un alto numero di connessioni concorrenti con un basso footprint di risorse, stressando il sistema come farebbe una flotta reale.

### 2. Data Queue

- **Tecnologia:** NATS JetStream (Modulo MQTT abilitato).
- **Ruolo:**
  Agisce come unico punto di ingresso e buffer ad alte prestazioni. Sostituisce l'Ingestion Service custom, accettando direttamente le connessioni MQTTS dai gateway e validandone l'identità tramite certificati statici. I messaggi validati vengono persistiti immediatamente nello stream JetStream per garantire la durabilità.
- **Perché:**
  Rimuove un hop architetturale superfluo, riducendo latenza e punti di fallimento. NATS gestisce nativamente la persistenza e il protocollo MQTT, garantendo che nessun dato venga perso se il database è sotto stress.

### 3. Data Consumer

- **Linguaggio:** Go.
- **Ruolo:**
  Worker interno che sottoscrive allo stream di NATS, preleva i messaggi crittografati in arrivo ed esegue operazioni di bulk write sul database per massimizzare il throughput.
- **Perché:**
  Separare la logica di ricezione dalla scrittura (pattern CQRS) permette di gestire i picchi di ingestion disaccoppiando i tempi di ricezione dai tempi di scrittura su disco.

### 4. Measures Database

**Tecnologia:** PostgreSQL + TimescaleDB
**Ruolo:**
Storage persistente delle serie temporali. I dati verranno salvati in **Hypertables** partizionate, includendo una colonna `tenant_id` per garantire la segregazione logica dei dati fin dal livello di storage.
**Perché:**
TimescaleDB offre le performance di un TSDB verticale mantenendo l'affidabilità e l'ecosistema SQL di PostgreSQL.

### 5. Data API

- **Linguaggio:** TypeScript - NestJS.
- **Ruolo:**
  Espone endpoint **Read-Only** per la dashboard. Interroga il database applicando filtri rigorosi per l'isolamento dei tenant (es. `WHERE tenant_id = X`), basandosi su un'autenticazione mockata.
  Utilizzo di TypeORM per una gestione migliore delle Hypertables di TimescaleDB.
- **Perché:**
  NestJS offre una struttura modulare e manutenibile per la logica di business e di accesso ai dati. L'integrazione con TypeORM garantisce un perfetto equilibrio tra l'astrazione di un ORM e la necessità di eseguire query SQL ottimizzate per serie temporali.

### 6. Web Dashboard & Client SDK

- **Tecnologia:** Angular + TypeScript SDK.
- **Ruolo:**
  Interfaccia di visualizzazione per l'utente. La logica di decifrazione avviene lato client: l'applicazione riceve dalle API e utilizza la chiave statica corretta per rendere i dati leggibili.
- **Perché:**
  Dimostra la fattibilità dell'approccio End-to-End Encryption (Zero Knowledge), dove il dato in chiaro è visibile solo ai "bordi" del sistema (Gateway e Utente finale).

---

## Cosa non c'è

Per motivi di tempo e focalizzazione, il PoC esclude le parti non critiche per la validazione del flusso dati principale:

- **Autenticazione Utenti reale:** La Dashboard userà un sistema di _Mock Auth_ (es. header statici) invece di integrare Keycloak. I Gateway useranno invece mTLS reale con certificati statici.
- **Gestione PKI Dinamica:** Non ci sarà una CA interna attiva o un Key Vault; i certificati client e server sono trattati come **asset statici** generati una tantum.
- **API Esterne:** Non esporremo API per integratori terzi.
- **Segregazione multi-tenant "hard":** Useremo una segregazione logica (colonna nel DB), riservandoci di valutare schemi più complessi (es. RLS o schema-per-tenant) per l'MVP finale.
- **Testing & Code Quality:** Il codice sarà funzionale alla dimostrazione ("quick & dirty"), senza copertura di test estensiva.
- **Observability:** Non implementeremo stack di monitoraggio (Prometheus/Grafana) o logging centralizzato per auditing.
- **Flusso Bidirezionale:** Il sistema gestirà solo l'ingestion (Gateway -> Cloud), escludendo l'invio di comandi ai dispositivi.

---

## Guida all'Esecuzione

### Prerequisiti

Assicurarsi di aver installato sulla macchina **Docker Desktop**.

### Esecuzione

1. **Apri il terminale** nella cartella `infra` \\
   `cd /path/to/PoC/infra`

2. **Avvia l'ambiente:**\\
   `docker compose up --build`\\
   _Attendere che tutti i servizi siano started_\\

Se si vuole che il terminale sia lasciato libero e si vuole che tutta la struttura runni in background utilizzare:\\
`docker compose up --build -d`

3. **Accedi alla Dashboard**\\
   Apri il browser su: http://localhost:4200

### Credenziali e Accesso

Per fare il login nella Dashboard, usa uno dei **Tenant ID** configurati nel simulatore.

- **Tenant 1:** `605e76a6-9812-4632-8418-43d99d9403d1`
- **Tenant 2:** `a66b9370-13f8-43e3-b097-f58c704f0f62`
