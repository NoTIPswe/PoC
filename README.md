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

- **Tecnologia:** NATS JetStream.
- **Ruolo:**
  Agisce come unico punto di ingresso e buffer ad alte prestazioni. Sostituisce l'Ingestion Service custom, accettando direttamente le connessioni dai gateway e validandone l'identità tramite certificati statici. I messaggi validati vengono persistiti immediatamente nello stream JetStream per garantire la durabilità.
- **Perché:**
  Rimuove un hop architetturale superfluo, riducendo latenza e punti di fallimento. NATS gestisce nativamente la persistenza, garantendo che nessun dato venga perso se il database è sotto stress.

### 3. Data Consumer

- **Linguaggio:** Go.
- **Ruolo:**
  Worker interno che sottoscrive allo stream di NATS, preleva i messaggi crittografati in arrivo ed esegue operazioni di bulk write sul database per massimizzare il throughput.
- **Perché:**
  Separare la logica di ricezione dalla scrittura (pattern CQRS) permette di gestire i picchi di ingestion disaccoppiando i tempi di ricezione dai tempi di scrittura su disco.

### 4. Measures Database

- **Tecnologia:** PostgreSQL + TimescaleDB
- **Ruolo:**
Storage persistente delle serie temporali. I dati verranno salvati in **Hypertables** partizionate, includendo una colonna `tenant_id` per garantire la segregazione logica dei dati fin dal livello di storage.
- **Perché:**
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
- **Observability Avanzata:** È presente un setup base di Prometheus + Grafana per il monitoraggio di NATS, ma non sono implementati alerting, logging centralizzato o tracing distribuito.
- **Flusso Bidirezionale:** Il sistema gestirà solo l'ingestion (Gateway -> Cloud), escludendo l'invio di comandi ai dispositivi.

---

# Guida di Esecuzione - PoC

## Prerequisiti

- **Docker Desktop** installato e avviato sulla macchina
  - [Download Docker Desktop](https://www.docker.com/products/docker-desktop/)
- Porte libere: `4200` (Dashboard), `3000` (API), `3001` (Grafana), `9090` (Prometheus), `5432` (Database), `4222` (NATS)

## Esecuzione

### 1. Posizionati nella cartella infra

```bash
cd /path/to/PoC/infra
```

### 2. Avvia l'ambiente

**Modalità interattiva** (con log visibili):

```bash
docker compose up --build
```

**Modalità background** (terminale libero):

```bash
docker compose up --build -d
```

> ⏱️ **Nota**: Il primo avvio può richiedere 2-3 minuti per scaricare le immagini e buildare i container.

### 3. Verifica che i servizi siano avviati

Controlla lo stato dei container:

```bash
docker compose ps
```

Dovresti vedere tutti i servizi con stato `Up` o `running`.

### 4. Accedi alla Dashboard

Apri il browser e vai su: **http://localhost:4200**

## Credenziali di Accesso

Per fare il login nella Dashboard, usa uno dei seguenti **Tenant ID**:

- **Tenant 1**: `605e76a6-9812-4632-8418-43d99d9403d1`
- **Tenant 2**: `a66b9370-13f8-43e3-b097-f58c704f0f62`

Incolla uno di questi ID nel campo di login e clicca "Access Dashboard".

## Cosa Aspettarsi

Dopo il login, dovresti vedere:

- Una tabella con dati di telemetria cifrati e decifrati lato client
- Dati in arrivo ogni 5 secondi da diversi dispositivi IoT simulati
- Tipi di sensori: temperature, humidity, power meter, air quality, motion sensor

## Monitoring Stack (Prometheus + Grafana)

Il PoC include un setup di monitoraggio per NATS basato su Prometheus e Grafana.

### Accesso alle interfacce

- **Grafana**: http://localhost:3001
  - Username: `admin`
  - Password: `admin`
- **Prometheus**: http://localhost:9090
- **NATS Exporter Metrics**: http://localhost:7777/metrics

### Dashboard NATS

Una dashboard pre-configurata "NATS Server Monitoring" è automaticamente disponibile in Grafana con le seguenti metriche:

- **Message Rate**: Messaggi in entrata/uscita al secondo
- **Throughput**: Bytes in entrata/uscita al secondo
- **Active Connections**: Connessioni attive al server NATS
- **Subscriptions**: Numero di sottoscrizioni attive
- **Server Memory/CPU**: Utilizzo risorse del server NATS
- **Slow Consumers**: Contatore di consumer lenti

## Troubleshooting

### La dashboard non carica dati

1. Verifica che tutti i servizi siano running:

```bash
   docker compose ps
```

2. Controlla i log del simulatore:

```bash
   docker compose logs gateway-simulator
```

3. Controlla i log del data-consumer:

```bash
   docker compose logs data-consumer
```

### Errore "Cannot connect to the Docker daemon"

- Assicurati che Docker Desktop sia avviato

### Porta già in uso

Se una porta è già occupata, puoi modificarla nel file `docker-compose.yml`

## Arresto dell'Ambiente

**Se in modalità interattiva**: Premi `Ctrl+C` nel terminale

**Se in background**:

```bash
docker compose down
```

**Per rimuovere anche i dati (reset completo)**:

```bash
docker compose down -v
```
